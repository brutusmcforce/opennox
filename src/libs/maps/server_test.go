package maps

import (
	"context"
	"io"
	"log/slog"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/shoenig/test/must"

	"github.com/brutusmcforce/opennox/libs/ifs"
	"github.com/brutusmcforce/opennox/libs/noxtest"
)

func TestIsAllowedFile(t *testing.T) {
	var cases = []struct {
		path string
		exp  bool
	}{
		{"example.map", true},
		{"example.nxz", false},
		{"example.rul", true},
		{"example.zip", false},
		{"example.tar", false},
		{"example.tar.gz", false},
		{"user.rul", false},
		{"some.lua", true},
		{"go.mod", true},
		{"go.sum", true},
		{"some.go", true},
		{"sub/some.go", true},
		{"vendor/sub/some.go", true},
		{"LICENSE", true},
		{"README.md", true},
		{"README.txt", true},
		{"some.other", false},
		{"some.json", true},
		{"some.yaml", true},
		{"some.yml", true},
		{"some.png", true},
		{"some.jpg", true},
		{"some.mp3", true},
		{"some.ogg", true},
		{".git/config", false},
		{".git/refs/heads/fake.go", false},
	}
	for _, c := range cases {
		c := c
		t.Run(c.path, func(t *testing.T) {
			got := IsAllowedFile(c.path)
			if got != c.exp {
				t.FailNow()
			}
		})
	}
}

func copyFile(t testing.TB, dst, src string) {
	s, err := ifs.Open(src)
	must.NoError(t, err)
	defer s.Close()
	d, err := os.Create(dst)
	must.NoError(t, err)
	defer d.Close()
	_, err = io.Copy(d, s)
	must.NoError(t, err)
	err = d.Close()
	must.NoError(t, err)
}

func TestMapServer(t *testing.T) {
	dpath := noxtest.DataPath(t, "maps")
	srcdir, err := os.MkdirTemp("", "opennox-map-test-src-*")
	must.NoError(t, err)
	dstdir, err := os.MkdirTemp("", "opennox-map-test-dst-*")
	must.NoError(t, err)
	t.Cleanup(func() {
		_ = os.RemoveAll(srcdir)
		_ = os.RemoveAll(dstdir)
	})
	const mname = "estate"
	for _, name := range []string{
		"estate", "!estate",
		"Estate", "!Estate",
	} {
		name := name
		t.Run(name, func(t *testing.T) {
			// Copy a well-known map to a source folder.
			err = os.MkdirAll(filepath.Join(srcdir, name), 0755)
			must.NoError(t, err)
			smpath := filepath.Join(srcdir, name, name+".map")
			copyFile(t, smpath, filepath.Join(dpath, mname, mname+".map"))

			// Start serving maps from source folder.
			srv := NewServer(slog.Default(), srcdir)
			hsrv := &http.Server{Handler: srv}
			l, err := net.Listen("tcp", ":0")
			must.NoError(t, err)
			t.Cleanup(func() {
				_ = l.Close()
			})
			t.Logf("listening on %q", l.Addr())
			go func() {
				err := hsrv.Serve(l)
				if err == http.ErrServerClosed {
					return
				}
				must.NoError(t, err)
			}()
			t.Cleanup(func() {
				_ = hsrv.Close()
			})

			ctx := context.Background()
			cli, err := NewClient(ctx, slog.Default(), l.Addr().String())
			must.NoError(t, err)
			defer cli.Close()

			err = cli.DownloadMap(ctx, dstdir, name)
			must.NoError(t, err)
		})
	}
}

// TestMapServerRoutes checks the routes and that a map name from the URL cannot escape the maps folder.
// It needs no game data: the maps folder is empty.
func TestMapServerRoutes(t *testing.T) {
	dir := t.TempDir()
	srv := NewServer(slog.Default(), dir)

	get := func(path string) *http.Response {
		t.Helper()
		req := httptest.NewRequest("GET", path, nil)
		rec := httptest.NewRecorder()
		srv.ServeHTTP(rec, req)
		return rec.Result()
	}

	resp := get("/api/v0/maps/")
	must.EqOp(t, http.StatusOK, resp.StatusCode)
	must.EqOp(t, "application/json", resp.Header.Get("Content-Type"))

	for _, path := range []string{
		"/api/v0/maps/missing",
		"/api/v0/maps/missing/download",
		"/api/v0/maps/missing/other",
	} {
		resp := get(path)
		must.EqOp(t, http.StatusNotFound, resp.StatusCode)
	}
}

// TestMapName checks the map names that are rejected. The router unescapes the path wildcard, so a request for
// "/api/v0/maps/..%2F..%2Fetc" arrives with the name "../../etc": it must not reach the file paths.
func TestMapName(t *testing.T) {
	for _, c := range []struct {
		name string
		exp  bool
	}{
		{"estate", true},
		{"!estate", true},
		{"Estate", true},
		{"so.me", true},
		{"", false},
		{".", false},
		{"..", false},
		{"../../etc", false},
		{"a/b", false},
		{`..\..\etc`, false},
		{"a\x00b", false},
	} {
		t.Run(c.name, func(t *testing.T) {
			r := httptest.NewRequest("GET", "/api/v0/maps/x", nil)
			r.SetPathValue("map", c.name)
			got, ok := mapName(r)
			must.EqOp(t, c.exp, ok)
			if ok {
				must.EqOp(t, c.name, got)
			}
		})
	}
}
