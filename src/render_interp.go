//go:build !server

package opennox

import (
	"image"
	"time"

	"github.com/spf13/viper"

	"github.com/brutusmcforce/opennox/libs/env"
	"github.com/brutusmcforce/opennox/libs/platform"

	"github.com/brutusmcforce/opennox/client"
	"github.com/brutusmcforce/opennox/client/noxrender"
	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/legacy"
)

// High refresh rate rendering.
//
// The game logic (and network updates) runs at a fixed rate of 30 ticks per second.
// Originally, the game was drawn exactly once per tick.
//
// To make the game look smooth on high refresh rate displays, we draw additional frames between the ticks.
// These frames don't advance the game state. Instead, positions of moving objects and the camera
// are interpolated between the last two ticks. This means the world is displayed one tick behind the
// game state, which is a common trade-off for frame interpolation.
//
// Objects remember their previous position and the frame of the last movement (see Nox_xxx_updateSpritePosition_49AA90),
// so we only need to track the camera position here.

const (
	configVideoMaxFPS = "video.max_fps"

	// maxFPSDef is the default frame rate limit. Values at or below the tick rate disable interpolation.
	maxFPSDef = 144
	maxFPSMax = 1000

	// interpTeleportDist is a squared distance between two ticks above which the movement is considered a teleport.
	interpTeleportDist = 64 * 64
)

func init() {
	viper.SetDefault(configVideoMaxFPS, maxFPSDef)
	registerOnConfigRead(func() {
		noxClient.interp.setMaxFPS(viper.GetInt(configVideoMaxFPS))
	})
}

type interpSaved struct {
	dr  *client.Drawable
	cur image.Point // real position
	set image.Point // interpolated position
}

type renderInterp struct {
	maxFPS int

	tickAt    time.Duration    // time of the last tick when the world was drawn
	tickFrame uint32           // game frame of the last tick when the world was drawn
	nextFrame time.Duration    // scheduled time of the next presented frame
	hasTick   bool             // world was drawn on the last tick
	frameCost [2]time.Duration // estimated time needed to draw an additional frame, for each betweenTicksMode

	camPrev, camCur image.Point
	camOK           bool

	// renderOnly is set while drawing an additional frame between the ticks.
	renderOnly bool

	applied  bool
	saved    []interpSaved
	savedCam image.Rectangle
	setCam   image.Rectangle
}

func (r *renderInterp) setMaxFPS(v int) {
	if v < 0 {
		v = 0
	} else if v > maxFPSMax {
		v = maxFPSMax
	}
	if r.maxFPS != v {
		videoLog.Printf("max fps: %d", v)
	}
	r.maxFPS = v
}

// enabled checks if additional frames should be drawn between the ticks.
func (r *renderInterp) enabled() bool {
	if env.IsE2E() {
		// tests rely on deterministic frame timing
		return false
	}
	return r.maxFPS > int(noxServer.TickRate())
}

func (r *renderInterp) tickDur() time.Duration {
	rate := noxServer.TickRate()
	if rate == 0 {
		rate = 30
	}
	return time.Second / time.Duration(rate)
}

// alpha returns the interpolation factor between the previous and the current tick.
func (r *renderInterp) alpha() float64 {
	if !r.enabled() {
		return 1
	}
	a := float64(platform.Ticks()-r.tickAt) / float64(r.tickDur())
	if a < 0 {
		a = 0
	} else if a > 1 {
		a = 1
	}
	return a
}

func lerpPoint(a, b image.Point, t float64) image.Point {
	return image.Point{
		X: a.X + int(float64(b.X-a.X)*t+0.5),
		Y: a.Y + int(float64(b.Y-a.Y)*t+0.5),
	}
}

func interpIsTeleport(a, b image.Point) bool {
	d := b.Sub(a)
	return d.X*d.X+d.Y*d.Y > interpTeleportDist
}

// onTick must be called on each game tick, right before drawing the world.
func (c *Client) interpOnTick(vp *noxrender.Viewport) {
	r := &c.interp
	r.tickAt = platform.Ticks()
	r.tickFrame = c.srv.Frame()
	r.hasTick = true
	cam := vp.World.Min
	if r.camOK {
		r.camPrev = r.camCur
	} else {
		r.camPrev = cam
	}
	r.camCur = cam
	r.camOK = true
	if interpIsTeleport(r.camPrev, r.camCur) {
		r.camPrev = r.camCur
	}
	if r.tickFrame%32 == 0 && r.tickFrame > 2 {
		c.ForgetMovedByDraw(r.tickFrame - 2)
	}
}

// interpReset disables interpolation until the next tick.
func (c *Client) interpReset() {
	c.interp.hasTick = false
	c.interp.camOK = false
}

// interpApply temporarily moves the camera and objects to their interpolated positions.
// It must be followed by interpRestore.
func (c *Client) interpApply(vp *noxrender.Viewport) {
	r := &c.interp
	if r.applied || !r.enabled() || !r.camOK {
		return
	}
	t := r.alpha()
	if t >= 1 {
		return
	}
	r.applied = true
	r.saved = r.saved[:0]

	// camera
	r.savedCam = vp.World
	if vp.World.Min == r.camCur {
		off := lerpPoint(r.camPrev, r.camCur, t).Sub(r.camCur)
		vp.World = vp.World.Add(off)
	}
	r.setCam = vp.World

	// objects that moved on the last tick
	frame := r.tickFrame
	const margin = 256
	rect := image.Rectangle{
		Min: r.savedCam.Min.Sub(image.Pt(margin, margin)),
		Max: r.savedCam.Min.Add(vp.Size).Add(image.Pt(margin, margin+128)),
	}
	c.Objs.EachInRect(rect, func(dr *client.Drawable) {
		if dr.Field_5 != frame {
			return
		}
		if frame > 0 && c.MovedByDraw(dr, frame-1) {
			// moved by its own draw function; such effects only move on ticks
			return
		}
		prev := image.Point{X: int(int32(dr.Field_8)), Y: int(int32(dr.Field_9))}
		cur := dr.PosVec
		if prev == cur || interpIsTeleport(prev, cur) {
			return
		}
		p := lerpPoint(prev, cur, t)
		dr.PosVec = p
		r.saved = append(r.saved, interpSaved{dr: dr, cur: cur, set: p})
	})
}

// interpRestore restores the real positions after interpApply.
func (c *Client) interpRestore(vp *noxrender.Viewport) {
	r := &c.interp
	if !r.applied {
		return
	}
	r.applied = false
	for i := len(r.saved) - 1; i >= 0; i-- {
		s := &r.saved[i]
		// the object might have been moved while drawing; keep the new position in this case
		if s.dr.PosVec == s.set {
			s.dr.PosVec = s.cur
		}
		s.dr = nil
	}
	r.saved = r.saved[:0]
	if vp.World == r.setCam {
		vp.World = r.savedCam
	}
}

// presented must be called each time a frame is presented. It schedules the next frame.
func (r *renderInterp) presented(now time.Duration) {
	if !r.enabled() {
		return
	}
	period := r.framePeriod()
	// Frames follow a fixed schedule, so that a frame delayed by a game tick is caught up later.
	// Don't catch up if we are too far behind, it would only cause a burst of frames.
	if r.nextFrame == 0 || now-r.nextFrame > 2*period {
		r.nextFrame = now
	}
	// Frames presented early (e.g. while the game is paused) must not move the schedule into the future,
	// otherwise no additional frames would be drawn until the time catches up with it.
	r.nextFrame = min(r.nextFrame+period, now+period)
}

// framePeriod returns the minimal time between presented frames.
func (r *renderInterp) framePeriod() time.Duration {
	return time.Second / time.Duration(r.maxFPS)
}

type betweenTicksMode int

const (
	betweenTicksNone  = betweenTicksMode(iota)
	betweenTicksWorld // in game: draw the world with interpolated positions and the GUI
	betweenTicksMenu  // in menus: only draw the GUI
)

// betweenTicksMode checks if it's safe to draw additional frames now, and what they should include.
func (c *Client) betweenTicksMode() betweenTicksMode {
	r := &c.interp
	if !r.enabled() || isDedicatedServer || noxflags.HasEngine(noxflags.EngineNoRendering) {
		return betweenTicksNone
	}
	if nox_client_gui_flag_815132 != 0 {
		return betweenTicksMenu
	}
	if r.hasTick && c.DrawFunc != nil &&
		!noxflags.HasEngine(noxflags.EnginePause) &&
		c.ClientPlayerUnit() != nil && nox_client_isConnected() {
		return betweenTicksWorld
	}
	return betweenTicksNone
}

// waitAndRender waits for the next game tick. If enabled, it draws additional interpolated frames while waiting.
func (c *Client) waitAndRender(dt time.Duration) {
	r := &c.interp
	if dt <= 0 {
		// The game is behind its schedule, or paused (the game time doesn't advance). Don't wait for the tick,
		// but still respect the frame rate limit: while paused, a frame is drawn on each iteration.
		// When the game is catching up, frames are not presented, so the schedule is in the past and we don't wait.
		if r.enabled() {
			if wait := r.nextFrame - platform.Ticks(); wait > 0 && wait <= r.framePeriod() {
				c.srv.LoopSleep(wait)
			}
		}
		return
	}
	deadline := platform.Ticks() + dt
	for {
		mode := c.betweenTicksMode()
		if mode == betweenTicksNone {
			break
		}
		est := &r.frameCost[mode-betweenTicksWorld]
		// start drawing so that the frame is presented at the scheduled time
		start := max(r.nextFrame-*est, platform.Ticks())
		// don't draw if the frame won't be ready before the next tick
		if start+*est >= deadline {
			if r.nextFrame < deadline && *est > r.tickDur()/4 {
				// A frame is due, but the estimate says it won't fit. The estimate is only updated
				// when frames are drawn, so a single slow frame could disable additional frames forever.
				// Let a suspiciously high estimate decay, so that the frame cost is measured again.
				*est = *est * 7 / 8
			}
			break
		}
		if now := platform.Ticks(); start > now {
			c.srv.LoopSleep(start - now)
		}
		if c.betweenTicksMode() != mode {
			continue
		}
		t0 := platform.Ticks()
		c.drawBetweenTicks(mode)
		// keep a moving average of the frame cost; limit the effect of a single slow frame
		cost := platform.Ticks() - t0
		if *est == 0 {
			*est = cost
		} else {
			*est = (*est*7 + min(cost, 2**est)) / 8
		}
	}
	if now := platform.Ticks(); deadline > now {
		c.srv.LoopSleep(deadline - now)
	}
}

// drawBetweenTicks draws and presents an additional frame without advancing the game state.
func (c *Client) drawBetweenTicks(mode betweenTicksMode) {
	r := &c.interp
	r.renderOnly = true
	setDrawIsTick(false)
	c.Inp.Poll()
	c.Inp.SetLiveMouse(true)
	defer func() {
		c.Inp.SetLiveMouse(false)
		setDrawIsTick(true)
		r.renderOnly = false
	}()
	if mode == betweenTicksWorld {
		c.drawWorldFrame()
	} else {
		// nothing to interpolate in menus
		c.interpReset()
	}
	c.mainloopDrawAndPresent()
}

// drawWorldFrame draws the game world for an additional frame. See nox_xxx_clientDrawAll_436100_draw.
func (c *Client) drawWorldFrame() {
	c.updateZoom()
	c.beginWorldDraw()
	vp := c.Viewport()
	legacy.Sub_430B50(vp.Screen.Min.X, vp.Screen.Min.Y, vp.Screen.Max.X, vp.Screen.Max.Y)
	c.drawWorldInterp(vp)
	c.endWorldDraw()
	c.worldDrawn()
	legacy.Nox_xxx_drawMinimapAndLines_4738E0()
	c.drawWorldOverlays()
}

// drawWorldInterp draws the world with interpolated object positions.
func (c *Client) drawWorldInterp(vp *noxrender.Viewport) {
	c.interpApply(vp)
	defer c.interpRestore(vp)
	c.nox_xxx_drawAllMB_475810_draw(vp)
}
