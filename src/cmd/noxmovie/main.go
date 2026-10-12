//go:build !server

package main

import (
	"errors"
	"flag"
	"fmt"
	"image"
	"log/slog"
	"os"

	"github.com/brutusmcforce/opennox/libs/client/seat/sdl"
	"github.com/brutusmcforce/opennox/libs/noximage"

	"github.com/brutusmcforce/opennox/client/gpu"
	"github.com/brutusmcforce/opennox/client/noxmovie"
	"github.com/brutusmcforce/opennox/client/render"
	"github.com/brutusmcforce/opennox/legacy/client/audio/ail"
	"github.com/brutusmcforce/opennox/legacy/common/alloc/handles"
)

var (
	fMovie = flag.String("i", "intro.vqa", "movie file to open")
)

func main() {
	flag.Parse()
	if err := run(*fMovie); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(fname string) error {
	// Handles needed for communication with C (currently used by audio backend).
	// TODO: move them to engine instead
	handles.Init()
	defer handles.Release()

	// Initialize SDL window and the OpenGL renderer.
	log := slog.Default()
	sz := image.Pt(640, 480)
	win, err := sdl.New(log, "Nox Movie Player", sz)
	if err != nil {
		return err
	}
	defer win.Close()

	scr, err := render.New(win)
	if err != nil {
		return err
	}
	g, err := gpu.New(log)
	if err != nil {
		return err
	}
	defer g.Close()
	scr.SetPresentFunc(g.Present)

	var frame *noximage.Image16
	display := func(img *noximage.Image16) {
		if frame == nil || frame.Rect != img.Rect {
			frame = noximage.NewImage16(img.Rect)
		}
		copy(frame.Pix, img.Pix)
		g.Upload(frame)
		scr.CopyBuffer(frame)
	}

	// Initialize audio.
	ail.Startup()

	drv := ail.WaveOutOpen()
	if e := ail.LastError(); e != "" {
		return errors.New(e)
	}

	plr, err := noxmovie.NewPlayer(fname, win, drv, display)
	if err != nil {
		return err
	}

	defer plr.Close()

	plr.Start()
	plr.Play()

	plr2, err := noxmovie.NewPlayer(fname, win, drv, display)
	if err != nil {
		return err
	}

	defer plr2.Close()

	plr2.Start()
	plr2.Play()

	ail.Shutdown()
	return nil
}
