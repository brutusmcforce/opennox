package netstr

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/require"
)

func queuedSeqs(ns *reliableSender) []byte {
	var out []byte
	for it := ns.queue; it != nil; it = it.next {
		out = append(out, it.seq)
	}
	slices.Sort(out)
	return out
}

func seqRange(from byte, n int) []byte {
	var out []byte
	for i := 0; i < n; i++ {
		out = append(out, from+byte(i))
	}
	slices.Sort(out)
	return out
}

// Acknowledgements near the sequence wrap must keep the packets sent after the wrap
// and drop the ones sent before it.
func TestReliableAckWrap(t *testing.T) {
	ns := newReliableSender(func(data []byte) (int, error) { return len(data), nil })
	ns.seq = 0xEE
	for i := 0; i < 0x20; i++ { // 0xEE..0x0D
		_, err := ns.Send(1, []byte{byte(i)})
		require.NoError(t, err)
	}
	ns.SetCurrent(0xF5)
	require.Equal(t, seqRange(0xF5, 0x19), queuedSeqs(ns)) // 0xF5..0x0D
	ns.SetCurrent(0x05)
	require.Equal(t, seqRange(0x05, 0x09), queuedSeqs(ns)) // 0x05..0x0D
	ns.SetCurrent(0x40)
	require.Empty(t, queuedSeqs(ns))
}
