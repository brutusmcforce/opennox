package discover

import (
	"testing"

	"github.com/shoenig/test/must"
)

// Configs from older versions have the old default lobby saved.
func TestLobbyServerDefault(t *testing.T) {
	old := LobbyServer
	t.Cleanup(func() { LobbyServer = old })
	for _, c := range []struct{ conf, exp string }{
		{"", DefaultLobbyServer},
		{oldLobbyServer, DefaultLobbyServer},
		{DefaultLobbyServer, DefaultLobbyServer},
		{"http://127.0.0.1:8088", "http://127.0.0.1:8088"},
	} {
		LobbyServer = c.conf
		must.EqOp(t, c.exp, lobbyServer())
	}
}
