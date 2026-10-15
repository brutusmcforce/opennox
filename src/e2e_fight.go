package opennox

import (
	"image"
	"math"
	"time"

	"github.com/brutusmcforce/opennox/client"
	"github.com/brutusmcforce/opennox/libs/client/keybind"
	"github.com/brutusmcforce/opennox/libs/client/seat"
	"github.com/brutusmcforce/opennox/libs/object"

	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/server"
)

// Multiplayer test actions: find another player and fight them. The host knows every player from the server
// state; a client only sees the players its server sent (those near it).

// e2ePlayerTarget returns the world positions of the local player and of the nearest other player.
func e2ePlayerTarget() (self, target image.Point, ok bool) {
	best := math.MaxFloat64
	consider := func(p image.Point) {
		d := math.Hypot(float64(p.X-self.X), float64(p.Y-self.Y))
		if d < best {
			best, target, ok = d, p, true
		}
	}
	if noxflags.HasGame(noxflags.GameHost) {
		host := noxServer.Players.ByInd(server.HostPlayerIndex)
		if host == nil || host.PlayerUnit == nil {
			return
		}
		own := host.PlayerUnit
		self = image.Pt(int(own.PosVec.X), int(own.PosVec.Y))
		for _, u := range noxServer.Players.ListUnits() {
			if u != own && !u.Flags().Has(object.FlagDestroyed) {
				consider(image.Pt(int(u.PosVec.X), int(u.PosVec.Y)))
			}
		}
		return
	}
	c := noxClient
	own := c.ClientPlayerUnit()
	if own == nil {
		return
	}
	self = own.PosVec
	vp := c.Viewport()
	c.Objs.EachInRect(vp.World.Inset(-200), func(dr *client.Drawable) {
		if dr != own && dr.Class().Has(object.ClassPlayer) {
			consider(dr.PosVec)
		}
	})
	return
}

// e2ePlayerLog logs every player's health and position from the server state on the host, and the other
// players a client sees (the ones its server sent).
func e2ePlayerLog() {
	if !noxflags.HasGame(noxflags.GameHost) {
		c := noxClient
		own := c.ClientPlayerUnit()
		if own == nil {
			return
		}
		n := 0
		c.Objs.EachInRect(c.Viewport().World.Inset(-200), func(dr *client.Drawable) {
			if dr != own && dr.Class().Has(object.ClassPlayer) {
				e2eLog.Printf("VISIBLE player at (%d,%d)", dr.PosVec.X, dr.PosVec.Y)
				n++
			}
		})
		e2eLog.Printf("VISIBLE players %d, self at (%d,%d)", n, own.PosVec.X, own.PosVec.Y)
		return
	}
	for _, u := range noxServer.Players.ListUnits() {
		name := ""
		if pl := u.ControllingPlayer(); pl != nil {
			name = pl.Name()
		}
		hp, max := 0, 0
		if u.HealthData != nil {
			hp, max = int(u.HealthData.Cur), int(u.HealthData.Max)
		}
		e2eLog.Printf("PLAYER %q health %d/%d at (%.0f,%.0f)", name, hp, max, u.PosVec.X, u.PosVec.Y)
	}
}

// Chase runs towards the nearest other player for dt frames, and stops within stop units of them.
func (sc *e2eScenario) Chase(dt time.Duration, stop int, name string) {
	if stop <= 0 {
		stop = 60
	}
	pressed := false
	for i := time.Duration(0); i < dt; i++ {
		n := ""
		if i == 0 {
			n = name
		}
		sc.add(1, n, func() {
			self, target, ok := e2ePlayerTarget()
			dx, dy := float64(target.X-self.X), float64(target.Y-self.Y)
			if !ok || math.Hypot(dx, dy) < float64(stop) {
				if pressed {
					e2eQueueInput(&seat.MouseButtonEvent{Button: seat.MouseButtonRight, Pressed: false})
					pressed = false
				}
				return
			}
			ang := math.Atan2(dx, -dy) / math.Pi
			e2eQueueInput(&seat.MouseMoveEvent{Pos: e2eAngToPos(ang, e2eRunDist)})
			if !pressed {
				e2eQueueInput(&seat.MouseButtonEvent{Button: seat.MouseButtonRight, Pressed: true})
				pressed = true
			}
		})
	}
	sc.add(1, "", func() {
		if pressed {
			e2eQueueInput(&seat.MouseButtonEvent{Button: seat.MouseButtonRight, Pressed: false})
		}
	})
}

// Attack aims at the nearest other player every 15 frames for dt frames, and attacks (slot 0: the Action
// button) or casts the spell in the given quick slot (1-5).
func (sc *e2eScenario) Attack(dt time.Duration, slot int, name string) {
	keys := []keybind.Key{keybind.KeyA, keybind.KeyS, keybind.KeyD, keybind.KeyF, keybind.KeyG}
	for i := time.Duration(0); i < dt; i += 15 {
		n := ""
		if i == 0 {
			n = name
		}
		sc.add(15, n, func() {
			_, target, ok := e2ePlayerTarget()
			if !ok {
				return
			}
			pos := noxClient.Viewport().ToScreenPos(target)
			e2eQueueInput(&seat.MouseMoveEvent{Pos: pos})
			if slot >= 1 && slot <= len(keys) {
				e2eQueueInput(&seat.KeyboardEvent{Key: keys[slot-1], Pressed: true})
			} else {
				e2eQueueInput(&seat.MouseButtonEvent{Button: seat.MouseButtonLeft, Pressed: true})
			}
		})
		sc.add(1, "", func() {
			if slot >= 1 && slot <= len(keys) {
				e2eQueueInput(&seat.KeyboardEvent{Key: keys[slot-1], Pressed: false})
			} else {
				e2eQueueInput(&seat.MouseButtonEvent{Button: seat.MouseButtonLeft, Pressed: false})
			}
		})
	}
}

// Players logs every player's health (see e2ePlayerLog).
func (sc *e2eScenario) Players(name string) {
	sc.add(0, name, e2ePlayerLog)
}

// Gather moves every other player next to the host (host only), like the "tp" cheat: the arena maps are mazes,
// so the test brings the players together instead of path finding.
func (sc *e2eScenario) Gather(name string) {
	sc.add(0, name, func() {
		if !noxflags.HasGame(noxflags.GameHost) {
			return
		}
		host := noxServer.Players.Host()
		if host == nil || host.PlayerUnit == nil {
			return
		}
		pos := host.PlayerUnit.Pos()
		pos.X += 40
		for _, u := range noxServer.Players.ListUnits() {
			if u != host.PlayerUnit {
				asObjectS(u).SetPos(pos)
				e2eLog.Printf("GATHER: moved a player next to the host at (%.0f,%.0f)", pos.X, pos.Y)
			}
		}
	})
}
