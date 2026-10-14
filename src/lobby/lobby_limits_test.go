package lobby

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

// Registrations with broken fields are refused: they end up in the game's server list.
func TestLobbyValidate(t *testing.T) {
	ctx := context.Background()
	l := NewLobby()
	base := Game{Name: "test", Address: "1.1.1.1"}
	testSetDefault(&base)
	for name, mod := range map[string]func(g *Game){
		"newline in name":  func(g *Game) { g.Name = "a\nb" },
		"escape in map":    func(g *Game) { g.Map = "a\x1b[31m" },
		"invalid utf8":     func(g *Game) { g.Name = "a\xff" },
		"cur over max":     func(g *Game) { g.Players.Cur = 33 },
		"max too big":      func(g *Game) { g.Players.Max = 100000 },
		"list over max":    func(g *Game) { g.Players.Max = 1; g.Players.List = make([]PlayerInfo, 2) },
		"bad player name":  func(g *Game) { g.Players.List = []PlayerInfo{{Name: "x\x00"}} },
		"negative res":     func(g *Game) { g.Res.Width = -1 },
		"negative stage":   func(g *Game) { g.Quest = &QuestInfo{Stage: -1} },
		"long player name": func(g *Game) { g.Players.List = []PlayerInfo{{Name: string(make([]byte, maxFieldLen+1))}} },
	} {
		g := base
		mod(&g)
		require.Error(t, l.RegisterGame(ctx, &g), name)
	}
	g := base
	g.Name = "Ünicode Server ⚔"
	g.Players = PlayersInfo{Cur: 1, Max: 32, List: []PlayerInfo{{Name: "Jack", Class: "conjurer"}}}
	require.NoError(t, l.RegisterGame(ctx, &g))
}

// One address can't fill the lobby, and the stored game doesn't share memory with the request.
func TestLobbyLimits(t *testing.T) {
	ctx := context.Background()
	l := NewLobby()
	g := Game{Name: "test", Address: "1.1.1.1"}
	testSetDefault(&g)
	for i := 0; i < maxGamesPerAddr; i++ {
		g.Port = 20000 + i
		require.NoError(t, l.RegisterGame(ctx, &g))
	}
	g.Port = 30000
	require.Error(t, l.RegisterGame(ctx, &g))
	// updating a registered game is still allowed
	g.Port = 20000
	require.NoError(t, l.RegisterGame(ctx, &g))
	// another address is not affected
	g.Address = "2.2.2.2"
	g.Players.List = []PlayerInfo{{Name: "a"}}
	require.NoError(t, l.RegisterGame(ctx, &g))
	g.Players.List[0].Name = "changed"
	list, err := l.ListGames(ctx)
	require.NoError(t, err)
	for _, v := range list {
		if v.Address == "2.2.2.2" {
			require.Equal(t, "a", v.Players.List[0].Name)
		}
	}
	// the total limit
	l2 := NewLobby()
	for i := 0; i < maxGames; i++ {
		g := Game{Name: "t", Address: fmt.Sprintf("10.%d.%d.%d", i>>16&255, i>>8&255, i&255)}
		testSetDefault(&g)
		require.NoError(t, l2.RegisterGame(ctx, &g))
	}
	g = Game{Name: "t", Address: "3.3.3.3"}
	testSetDefault(&g)
	require.Error(t, l2.RegisterGame(ctx, &g))
}
