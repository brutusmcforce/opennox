package opennoxcontrol

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

type testGame struct{ cmds []string }

func (g *testGame) GameInfo() (Info, error)     { return Info{}, nil }
func (g *testGame) ListMaps() ([]Map, error)    { return nil, nil }
func (g *testGame) ChangeMap(name string) error { return nil }
func (g *testGame) Command(cmd string) error {
	g.cmds = append(g.cmds, cmd)
	return nil
}

func TestCommandCSRF(t *testing.T) {
	g := &testGame{}
	cp := NewControlPanel(g, &Options{AllowCommands: true})
	send := func(r *http.Request) {
		cp.ServeHTTP(httptest.NewRecorder(), r)
	}
	// a GET (e.g. an <img> tag on another site) must not run anything
	send(httptest.NewRequest("GET", "/cmd/?data=evil", nil))
	// neither must a form posted from another site
	r := httptest.NewRequest("POST", "/cmd/", strings.NewReader(url.Values{"data": {"evil"}}.Encode()))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Set("Sec-Fetch-Site", "cross-site")
	send(r)
	// the panel's own form works
	r = httptest.NewRequest("POST", "/cmd/", strings.NewReader(url.Values{"data": {"ok"}}.Encode()))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Set("Sec-Fetch-Site", "same-origin")
	send(r)
	if len(g.cmds) != 1 || g.cmds[0] != "ok" {
		t.Fatalf("unexpected commands: %q", g.cmds)
	}
}
