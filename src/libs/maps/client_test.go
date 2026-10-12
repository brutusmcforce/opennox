package maps

import (
	"archive/zip"
	"bytes"
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/shoenig/test/must"
)

// serveZip starts a map server that answers every download with the given archive.
func serveZip(t *testing.T, files map[string][]byte) *Client {
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	for name, data := range files {
		w, err := zw.Create(name)
		must.NoError(t, err)
		_, err = w.Write(data)
		must.NoError(t, err)
	}
	must.NoError(t, zw.Close())
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "HEAD" {
			return
		}
		w.Header().Set("Content-Type", contentTypeZIP)
		_, _ = w.Write(buf.Bytes())
	}))
	t.Cleanup(srv.Close)
	cli, err := NewClient(context.Background(), slog.Default(), strings.TrimPrefix(srv.URL, "http://"))
	must.NoError(t, err)
	return cli
}

func TestDownloadMapZipSlip(t *testing.T) {
	root := t.TempDir()
	dest := filepath.Join(root, "maps")
	must.NoError(t, os.Mkdir(dest, 0o755))
	cli := serveZip(t, map[string][]byte{
		"evil.map":                []byte("map"),
		"x/../../evil.lua":        []byte("planted"),
		"x/../../../outside.lua":  []byte("planted"),
		"x/../../other/other.lua": []byte("planted"),
		"/abs.lua":                []byte("planted"),
	})
	err := cli.DownloadMap(context.Background(), dest, "evil")
	if err != nil {
		t.Fatal(err)
	}
	_, err = os.Stat(filepath.Join(dest, "evil", "evil.map"))
	must.NoError(t, err)
	for _, p := range []string{"maps/evil.lua", "outside.lua", "maps/other/other.lua", "abs.lua"} {
		_, err = os.Stat(filepath.Join(root, p))
		must.True(t, os.IsNotExist(err), must.Sprintf("%s was written", p))
	}
}

func TestDownloadMapZipBomb(t *testing.T) {
	dest := t.TempDir()
	// compresses to almost nothing, unpacks above the limit
	cli := serveZip(t, map[string][]byte{
		"bomb.map": make([]byte, mapUnpackLimit+1),
	})
	err := cli.DownloadMap(context.Background(), dest, "bomb")
	must.Error(t, err)
}
