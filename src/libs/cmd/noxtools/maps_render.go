package main

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/image/draw"

	"github.com/brutusmcforce/opennox/libs/datapath"
	"github.com/brutusmcforce/opennox/libs/ifs"
	"github.com/brutusmcforce/opennox/libs/maps"
	"github.com/brutusmcforce/opennox/libs/maps/maprender"
	"github.com/brutusmcforce/opennox/libs/maps/minidraw"
)

func init() {
	cmd := &cobra.Command{
		Use:   "render map...",
		Short: "Renders Nox maps to PNG images, the way the game draws them",
		Long: `Renders Nox maps to PNG images, the way the game draws them: the floor, the walls and the objects, with
their colors and lights. A map is a map directory or file, or the name of a map in the game's maps directory.`,
		Example: `  noxtools map render con01a
  noxtools map render -o out --scale 0.25 ~/games/Nox/maps/*/
  noxtools map render --minimap estate`,
	}
	mapCmd.AddCommand(cmd)
	fData := cmd.Flags().String("data", "", "Nox data directory (found automatically if not set)")
	fOut := cmd.Flags().StringP("out", "o", ".", "output directory, or the output file when rendering one map")
	fScale := cmd.Flags().Float64("scale", 1, "scale of the image")
	fNoFloor := cmd.Flags().Bool("no-floor", false, "do not draw the floor")
	fNoWalls := cmd.Flags().Bool("no-walls", false, "do not draw the walls")
	fNoObjects := cmd.Flags().Bool("no-objects", false, "do not draw the objects")
	fNoLight := cmd.Flags().Bool("no-light", false, "draw everything fully lit")
	fMinimap := cmd.Flags().Bool("minimap", false, "draw the walls as the game's minimap does instead")
	cmd.SilenceUsage = true
	cmd.SilenceErrors = true // main prints it
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("expected at least one map")
		}
		if *fData != "" {
			datapath.SetData(*fData)
		}
		opts := &mapRenderOpts{
			Out:     *fOut,
			Single:  len(args) == 1,
			Scale:   *fScale,
			Minimap: *fMinimap,
			Render: maprender.Options{
				NoFloor: *fNoFloor, NoWalls: *fNoWalls, NoObjects: *fNoObjects, NoLight: *fNoLight,
			},
		}
		return cmdMapRender(cmd, args, opts)
	}
}

type mapRenderOpts struct {
	Out     string
	Single  bool // Out is the file name
	Scale   float64
	Minimap bool
	Render  maprender.Options
}

// mapDir returns the directory of a map given as a directory, a map file or a map name.
func mapDir(arg string) (string, error) {
	if fi, err := os.Stat(arg); err == nil {
		if fi.IsDir() {
			return filepath.Clean(arg), nil
		}
		return filepath.Dir(arg), nil
	}
	if !strings.ContainsAny(arg, `/\`) {
		// map directories differ in case
		dir := ifs.Normalize(datapath.Maps(strings.TrimSuffix(arg, maps.Ext)))
		if fi, err := os.Stat(dir); err == nil && fi.IsDir() {
			return dir, nil
		}
	}
	return "", fmt.Errorf("map not found: %q", arg)
}

func cmdMapRender(cmd *cobra.Command, args []string, opts *mapRenderOpts) error {
	var r *maprender.Renderer
	if !opts.Minimap {
		var err error
		r, err = maprender.NewRenderer(datapath.Data())
		if err != nil {
			return fmt.Errorf("cannot load the game data (set it with --data): %w", err)
		}
		defer func() { _ = r.Close() }()
	}
	failed := 0
	for _, arg := range args {
		dir, err := mapDir(arg)
		if err == nil {
			err = renderMap(cmd, r, dir, opts)
		}
		if err != nil {
			if len(args) == 1 {
				return err
			}
			failed++
			cmd.PrintErrln(err)
		}
	}
	if failed != 0 {
		return fmt.Errorf("%d of %d maps failed", failed, len(args))
	}
	return nil
}

func renderMap(cmd *cobra.Command, r *maprender.Renderer, dir string, opts *mapRenderOpts) error {
	name := filepath.Base(dir)
	m, err := maps.ReadMap(dir)
	if err != nil {
		return fmt.Errorf("%s: %w", name, err)
	}
	var img image.Image
	if opts.Minimap {
		img = minidraw.MinimapRGBA(m, color.Black)
	} else {
		rgba, err := r.DrawMap(m, &opts.Render)
		if rgba == nil {
			return fmt.Errorf("%s: %w", name, err)
		} else if err != nil {
			// a partial image
			cmd.PrintErrf("%s: %v\n", name, err)
		}
		img = rgba
	}
	if opts.Scale > 0 && opts.Scale != 1 {
		b := img.Bounds()
		w, h := max(1, int(float64(b.Dx())*opts.Scale)), max(1, int(float64(b.Dy())*opts.Scale))
		dst := image.NewRGBA(image.Rect(0, 0, w, h))
		draw.CatmullRom.Scale(dst, dst.Rect, img, b, draw.Src, nil)
		img = dst
	}
	out := opts.Out
	if !opts.Single || strings.HasSuffix(out, string(filepath.Separator)) || isDir(out) {
		if err := os.MkdirAll(out, 0755); err != nil {
			return err
		}
		out = filepath.Join(out, strings.ToLower(name)+".png")
	} else if err := os.MkdirAll(filepath.Dir(out), 0755); err != nil {
		return err
	}
	f, err := os.Create(out)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()
	if err := png.Encode(f, img); err != nil {
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	cmd.Println(out)
	return nil
}

func isDir(path string) bool {
	fi, err := os.Stat(path)
	return err == nil && fi.IsDir()
}
