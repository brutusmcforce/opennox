package opennox

import (
	"context"
	"fmt"
	"image/png"
	"os"
	"path/filepath"
	"time"

	"github.com/brutusmcforce/opennox/libs/console"
	"github.com/brutusmcforce/opennox/libs/datapath"
)

func init() {
	noxConsole.Register(&console.Command{
		Token:  "image",
		HelpID: "imagehelp",
		Flags:  console.ClientServer | console.NoDedicated,
		Func:   cmdScreenshot,
	})
}

func cmdScreenshot(ctx context.Context, c *console.Console, tokens []string) bool {
	makeScreenshot()
	return true
}

// screenshotDir is where screenshots are saved: the user's directory, not the game's (which may be read-only).
func screenshotDir() string {
	return datapath.User("screenshots")
}

// createScreenshotFile creates a new file for a screenshot taken at a time.
func createScreenshotFile(dir string, t time.Time) (*os.File, error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}
	base := filepath.Join(dir, "opennox-"+t.Format("20060102-150405"))
	for i := 1; ; i++ {
		name := base + ".png"
		if i > 1 {
			name = fmt.Sprintf("%s-%d.png", base, i)
		}
		f, err := os.OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
		if !os.IsExist(err) {
			return f, err
		}
	}
}

func makeScreenshot() {
	// Screenshot will wait for the next clean frame, so must run in a goroutine
	go func() {
		log := noxClient.Log
		img, err := noxClient.Screenshot(context.Background())
		if err != nil {
			log.Error("cannot take screenshot", "err", err)
			return
		}
		f, err := createScreenshotFile(screenshotDir(), time.Now())
		if err != nil {
			log.Error("cannot save screenshot", "err", err)
			return
		}
		defer func() { _ = f.Close() }()
		if err = png.Encode(f, img); err != nil {
			log.Error("cannot encode screenshot", "err", err)
			return
		}
		if err = f.Close(); err != nil {
			log.Error("cannot save screenshot", "err", err)
			return
		}
		log.Info("screenshot saved", "path", f.Name())
		// tell the player, on the game's thread
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		noxServer.QueueInLoop(ctx, func() {
			// also printed to the console
			nox_xxx_printCentered_445490("Screenshot saved: " + filepath.Base(f.Name()))
		})
	}()
}
