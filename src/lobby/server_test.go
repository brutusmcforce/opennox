package lobby_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	lobby "github.com/brutusmcforce/opennox/lobby"
)

func TestGetIP(t *testing.T) {
	api := lobby.NewServer(nil)
	req := httptest.NewRequest(http.MethodGet, "/api/v0/address", nil)
	rec := httptest.NewRecorder()
	api.ServeHTTP(rec, req)
	require.Equal(t, 200, rec.Result().StatusCode)
	require.Equal(t, `{"data":{"ip":"192.0.2.1"}}`+"\n", rec.Body.String())
}

func TestGetIPFailed(t *testing.T) {
	api := lobby.NewServer(nil)
	req := httptest.NewRequest(http.MethodGet, "/api/v0/address", nil)
	req.RemoteAddr = ""
	rec := httptest.NewRecorder()
	api.ServeHTTP(rec, req)
	require.Equal(t, http.StatusBadRequest, rec.Result().StatusCode)
}

func TestGetIPProxy(t *testing.T) {
	get := func(trust bool, remote string, hdr map[string]string) string {
		api := lobby.NewServer(nil)
		api.SetTrustProxy(trust)
		req := httptest.NewRequest(http.MethodGet, "/api/v0/address", nil)
		req.RemoteAddr = remote
		for k, v := range hdr {
			req.Header.Set(k, v)
		}
		rec := httptest.NewRecorder()
		api.ServeHTTP(rec, req)
		return rec.Body.String()
	}
	ip := func(v string) string { return `{"data":{"ip":"` + v + `"}}` + "\n" }
	// a local proxy
	require.Equal(t, ip("203.0.113.5"), get(true, "127.0.0.1:4000", map[string]string{"X-Real-IP": "203.0.113.5"}))
	require.Equal(t, ip("203.0.113.6"), get(true, "[::1]:4000", map[string]string{"X-Forwarded-For": "198.51.100.1, 203.0.113.6"}))
	require.Equal(t, ip("127.0.0.1"), get(true, "127.0.0.1:4000", map[string]string{"X-Real-IP": "garbage"}))
	// not trusted: a remote host, or the option is off
	require.Equal(t, ip("192.0.2.9"), get(true, "192.0.2.9:4000", map[string]string{"X-Real-IP": "203.0.113.5"}))
	require.Equal(t, ip("127.0.0.1"), get(false, "127.0.0.1:4000", map[string]string{"X-Real-IP": "203.0.113.5"}))
}
