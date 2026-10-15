package opennox

import (
	"context"
	"errors"
	"fmt"
	"image"
	"log/slog"
	"path"
	"path/filepath"
	"strings"

	"github.com/brutusmcforce/opennox/libs/client/keybind"
	"github.com/brutusmcforce/opennox/libs/datapath"
	noxlog "github.com/brutusmcforce/opennox/libs/log"
	"github.com/brutusmcforce/opennox/libs/maps"

	"github.com/brutusmcforce/opennox/client/gui"
	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/common/sound"
	"github.com/brutusmcforce/opennox/legacy/cnxz"
)

type clientMapDownload struct {
	log        *slog.Logger
	inProgress bool
	downloadOK bool
	cancel     func()
	done       <-chan error
	native     bool
	ndl        *maps.NativeDownloader
	nativePath string // the .nxz being downloaded
	win        *gui.Window
	iface      struct {
		setMapPath      func(path string)
		sendCancelMap   func() int
		sendReceivedMap func() int
	}
}

func (c *clientMapDownload) Init(cl *Client) {
	c.log = noxlog.WithSystem(cl.Log, "mapsend")
	c.ndl = maps.NewNativeDownloader(datapath.Maps())
	c.iface.setMapPath = cl.srv.nox_xxx_gameSetMapPath_409D70
	c.iface.sendCancelMap = cl.Nox_xxx_cliSendCancelMap_43CAB0
	c.iface.sendReceivedMap = cl.Nox_xxx_netMapReceived_43CA80
}

func (c *clientMapDownload) Downloading() bool {
	return c.inProgress
}

func (c *clientMapDownload) setDownloading(v bool) {
	c.inProgress = v
}

func (c *clientMapDownload) setDownloadOK(v bool) {
	c.downloadOK = v
}

func (c *clientMapDownload) deleteFile() {
	c.log.Info("delete downloaded map")
	c.ndl.CancelAndDelete()
	c.setDownloading(false)
	c.setDownloadOK(true)
}

func (c *clientMapDownload) cancelNative() {
	if c.ndl.DownloadOK() {
		c.iface.sendCancelMap()
		c.ndl.CancelAndCleanup()
	}
}

func (c *clientMapDownload) startNative(name string, sz uint) error {
	if name == "" {
		return errors.New("empty name")
	} else if sz <= 0 {
		return errors.New("invalid size")
	}
	if !strings.HasSuffix(name, maps.Ext) {
		name += maps.Ext
	}
	c.ndl.Reset()
	// the name comes from the server: only keep the file name, whatever separators it uses
	base := path.Base(strings.ReplaceAll(name, "\\", "/"))
	mname := strings.TrimSuffix(base, maps.Ext)
	if mname == "" || mname == "." || mname == ".." || strings.HasPrefix(mname, ".") {
		return fmt.Errorf("unexpected map name: %q", name)
	}
	c.iface.setMapPath(base)
	mdir := datapath.Maps(mname)
	fpath := filepath.Join(mdir, mname+".nxz")
	c.log.Info("downloading map to", "path", fpath)
	if err := c.ndl.Start(fpath, sz); err != nil {
		return err
	}
	c.nativePath = fpath
	return nil
}

// unpackNative replaces the .map with the downloaded .nxz. The game loads the .map and only unpacks the .nxz when
// there's no .map, so an old .map would be kept and fail the CRC check on every join (opennox#639).
func (c *clientMapDownload) unpackNative() error {
	zpath := c.nativePath
	mpath := strings.TrimSuffix(zpath, ".nxz") + maps.Ext
	return cnxz.DecompressFile(zpath, mpath)
}

func (c *clientMapDownload) onMapDownloadStart(name string, sz uint) {
	name = strings.TrimSuffix(strings.ToLower(name), maps.Ext)
	c.log.Info("download start", "proto", "udp", "name", name, "size", int(sz))
	c.native = true
	if err := c.startNative(name, sz); err != nil {
		c.log.Error("download start failed", "err", err)
		c.iface.sendCancelMap()
		c.setDownloading(false)
		c.setDownloadOK(false)
	}
}

func (c *clientMapDownload) onMapDownloadPart(ind uint, data []byte) {
	if len(data) == 0 {
		return
	}
	c.ndl.WritePart(ind, data)
	if c.ndl.Complete() {
		err := c.ndl.Finish()
		if err != nil {
			c.log.Error("cannot save the downloaded map", "err", err)
		} else if err = c.unpackNative(); err != nil {
			c.log.Error("cannot unpack the downloaded map", "err", err)
		}
		c.setDownloading(false)
		c.setDownloadOK(err == nil)
		c.setProgress(100)
		c.iface.sendReceivedMap()
	} else {
		c.setProgress(int(c.ndl.Progress() * 100.0))
	}
}

func (c *clientMapDownload) setProgress(val int) int {
	win2 := c.win.ChildByID(1603)
	return gui.EventRespInt(win2.Func94(gui.AsWindowEvent(0x4020, uintptr(val), 0)))
}

func (c *clientMapDownload) onMapDownloadAbort() {
	c.log.Info("download aborted by the server")
	c.deleteFile()
	c.setDownloading(false)
	c.setDownloadOK(false)
}

func (c *clientMapDownload) Cancel() {
	c.log.Info("map download cancelled")
	if c.native {
		c.cancelNative()
		return
	}
	c.setDownloading(false)
	if c.cancel != nil {
		c.cancel()
		c.cancel = nil
		if err := <-c.done; err != nil && err != context.Canceled {
			c.log.Error("map download failed", "err", err)
		}
		c.done = nil
	}
}

func (c *clientMapDownload) abortAndClose(path string) {
	c.CloseDialog()
	c.Cancel()
	c.iface.setMapPath(path)
	c.setDownloading(false)
	c.setDownloadOK(false)
}

func (c *clientMapDownload) CloseDialog() {
	if c.win != nil {
		c.win.Destroy()
	}
	c.win = nil
}

func (c *Client) ShowMapDownload() {
	c.mapsend.win = newWindowFromFile(c.GUI, "mapdnld.wnd", c.dialogProc)
	c.mapsend.win.SetAllFuncs(c.dialogEvents, nil, nil)
	mode := c.videoGetGameMode()
	x, y := mode.X, mode.Y
	max := videoGetMaxSize()
	if x > max.X {
		x = max.X
	}
	if y > max.Y {
		y = max.Y
	}
	c.mapsend.win.SetPos(image.Point{X: (x - 1024) / 2, Y: (y - 768) / 2})
	c.mapsend.win.ShowModal()
}

func (c *Client) dialogEvents(a1 *gui.Window, ev gui.WindowEvent) gui.WindowEventResp {
	ekey, ok := ev.(gui.WindowKeyPress)
	if !ok {
		return nil
	}
	if ekey.Key != keybind.KeyEsc {
		return nil
	}
	if ekey.Pressed {
		clientPlaySoundSpecial(sound.SoundPermanentFizzle, 100)
		c.mapsend.abortAndClose(memmap.String(0x5D4594, 1522936))
		noxflags.UnsetGame(noxflags.GameFlag21 | noxflags.GameFlag24)
	}
	return gui.RawEventResp(1)
}

func (c *Client) dialogProc(a1 *gui.Window, ev gui.WindowEvent) gui.WindowEventResp {
	switch ev := ev.(type) {
	case *WindowEvent0x4007:
		id := ev.Win.ID()
		clientPlaySoundSpecial(sound.SoundButtonPress, 100)
		if id != 1601 {
			return nil
		}
		c.mapsend.abortAndClose(memmap.String(0x5D4594, 1522940))
		noxflags.UnsetGame(noxflags.GameFlag21 | noxflags.GameFlag24)
		return gui.RawEventResp(1)
	}
	return nil
}
