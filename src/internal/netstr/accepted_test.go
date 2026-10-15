package netstr

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// Players leaving give their slot back to the listener: it counted every player that ever joined,
// so a server became full after enough joins and leaves (opennox#739).
func TestAcceptedReleased(t *testing.T) {
	g := NewStreams(func() uint32 { return 1 })
	lis := g.newStreamAt(0, -1, &Options{BufferSize: 1024})
	g.streams[0] = lis
	for i := 0; i < 20; i++ {
		ind := 1 + i%5
		c := g.newStreamAt(ind, ind, &Options{BufferSize: 1024})
		g.streams[ind] = c
		lis.accepted++
		c.listener = lis
		g.playerIDs[ind] = struct{}{}
		g.cntX++
		if i%2 == 0 {
			c.SendServerClose()
		} else {
			_ = c.Close()
		}
		require.Equal(t, 0, lis.accepted, "after %d", i)
		require.Equal(t, 0, g.cntX)
		require.Empty(t, g.playerIDs)
		c.SendServerClose() // twice: still released once
		require.Equal(t, 0, lis.accepted)
	}
}
