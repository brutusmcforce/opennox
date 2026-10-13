package opennox

import (
	"image"
	"log/slog"
	"unsafe"

	"github.com/brutusmcforce/opennox/libs/console"
	"github.com/brutusmcforce/opennox/libs/noximage"
	"github.com/spf13/viper"

	"github.com/brutusmcforce/opennox/client"
	"github.com/brutusmcforce/opennox/client/noxrender"
	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/legacy"
	"github.com/brutusmcforce/opennox/legacy/client/audio/ail"
)

var (
	noxClient *Client
)

func init() {
	legacy.GetClient = func() legacy.Client {
		return noxClient
	}
}

const (
	configVideoFiltering = "video.filtering"
	configVideoStretch   = "video.stretch"
)

// gpuBackend is an optional GPU renderer (see client/gpu).
type gpuBackend interface {
	noxrender.GPU
	// ClearTransparent clears the current render target to a transparent color.
	ClearTransparent()
	// Composite draws one render target over another one, scaling it to fit.
	Composite(dst, src *noximage.Image16)
	// Floor draws floor tiles in given screen spans.
	Floor(p *noxrender.FloorParams, spans []image.Rectangle)
	// Zoom replaces the content of a render target with a part of it, scaled to fill the whole target.
	Zoom(img *noximage.Image16, src image.Rectangle)
	// Upload replaces the content of the render target associated with an image with the image pixels.
	Upload(img *noximage.Image16)
	// SetTargetSize sets the physical size of the render target associated with an image (zero resets it).
	// Drawing coordinates are still in the image coordinates, and are scaled to the physical size.
	SetTargetSize(img *noximage.Image16, sz image.Point)
	// Forget releases the render target of an image that won't be used anymore.
	Forget(img *noximage.Image16)
	Close()
}

func NewClient(log *slog.Logger, pr console.Printer, srv *Server) (*Client, error) {
	cl := client.NewClient(log, pr, srv.Server)
	c := &Client{
		Client: cl,
		srv:    srv,
		r:      NewNoxRender(cl.Render()),
	}
	c.Client.ExtClient = unsafe.Pointer(c)
	c.info.Init(c)
	c.guiAdv.Init(c)
	c.screenshots.Init(c)
	c.mapsend.Init(c)
	c.srv.Teams.OnCreateOrRemove(legacy.Sub_459CD0)
	c.ctrl = new(CtrlEventHandler)
	return c, nil
}

type Client struct {
	*client.Client
	srv                    *Server
	r                      *NoxRender
	ctrl                   *CtrlEventHandler
	gamepad                *gamepadSeat
	mapsend                clientMapDownload
	guiAdv                 guiAdvOptions
	guiFPS                 guiFPS
	interp                 renderInterp
	gpu                    gpuBackend // OpenGL renderer; nil before the window is created
	fps                    fpsOverlay
	gpuFloor               noxrender.FloorParams
	gpuFloorSpans          []image.Rectangle
	sel                    selCircle
	zoom                   worldZoom
	screenshots            screenshots
	tiles                  clientTileData
	info                   clientDrawInfo
	netPrevMouse           image.Point
	inSub4312C0            bool
	ticks805996            uint64
	inDraw1                bool
	dragndrapSpell         uint32
	dragndropSpellType     int
	dragndropItem          *client.Drawable
	pos1097204             image.Point
	pos1097212             image.Point
	flag3798728            bool
	sysCursor              sysCursorState
	val1096556             int
	val1096560             int
	savedHighResFloors     int
	savedHighResFrontWalls int
	dword_5d4594_1046604   int
	tsFullPRev             uint32
	tsUpperPrev            uint32
}

func (c *Client) Cli() *client.Client {
	return c.Client
}

func (c *Client) R2() legacy.Render2 {
	return c.r
}

func (c *Client) GetCtrlEvent() legacy.CtrlEventHandler {
	return c.ctrl
}

func (c *Client) Close() error {
	c.freeSeat()
	return nil
}

func (c *Client) SetStretch(v bool) {
	if c == nil {
		return
	}
	c.Client.SetStretch(v)
	viper.Set(configVideoStretch, v)
}

func (c *Client) setStretchIfNotSet(v bool) {
	if !viper.IsSet(configVideoStretch) {
		c.SetStretch(v)
	}
}

func (c *Client) ToggleStretch() {
	c.SetStretch(!c.GetStretch())
}

func (c *Client) ToggleFiltering() {
	if c == nil {
		return
	}
	val := c.Client.ToggleFiltering()
	viper.Set(configVideoFiltering, val)
	writeConfigLater()
}

func (c *Client) processInput() {
	if c == nil || c.Inp == nil {
		return
	}
	c.Inp.Tick()
	c.nox_client_processInput_4308A0()
	c.nox_xxx_cursorUpdate_46B740()
	c.GUI.ProcessKeys(c.Inp)
}

func (c *Client) setClientPlayerUnit(dr *client.Drawable) {
	legacy.LocalPlayerDrawable = dr
}

func (c *Client) ClientPlayerUnit() *client.Drawable {
	return legacy.LocalPlayerDrawable
}

func (c *Client) Update() bool {
	defer noxPerfmon.startProfileClient()()
	if !isDedicatedServer {
		c.inDraw1 = true
		if c.DrawFunc != nil && !c.DrawFunc() {
			if debugMainloop {
				gameLog.Println("call_nox_draw_unk1 exit")
			}
			return false
		}
		c.inDraw1 = false
	}
	if !c.srv.Update2() {
		if debugMainloop {
			gameLog.Println("server.Update2 exit")
		}
		return false
	}
	legacy.Sub_4519C0()
	c.sub4312C0()
	if !isDedicatedServer {
		legacy.Sub_495430()
	}
	if noxflags.HasGame(noxflags.GameHost) && continueMenuOrHost && !mainloopStopError {
		mainloopMaybeSwitchMapXXX()
	}
	c.drawAndPresent()
	return true
}

func (c *Client) sub4312C0() {
	ail.Serve()
	if c.inSub4312C0 {
		return
	}
	c.inSub4312C0 = true
	defer func() {
		c.inSub4312C0 = false
	}()
	ticks := platformTicks()

	if ticks-c.ticks805996 > 33 {
		legacy.Dialogs.Sub_44D3A0()
		legacy.MusicModule.Update()
		c.ticks805996 = ticks
	}
}
