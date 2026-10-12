package netstr

import (
	"testing"
	"time"

	"github.com/brutusmcforce/opennox/libs/noxnet"
	"github.com/brutusmcforce/opennox/libs/noxnet/netmsg"
	"github.com/stretchr/testify/require"

	"github.com/brutusmcforce/opennox/common/ntype"
	"github.com/brutusmcforce/opennox/server/netlib"
)

func TestNetstr(t *testing.T) {
	var frame uint32 = 1
	server := NewStreams(func() uint32 { return frame })
	server.IsHost = func() bool { return true }
	server.GetMaxPlayers = func() int { return 10 }
	server.Xor = false
	received := make(chan []byte, 1)
	opts := &Options{
		Port: 18501, Max: 10, BufferSize: 2048,
		OnReceive: func(id netlib.StreamID, buf []byte) int {
			select {
			case received <- append([]byte(nil), buf...):
			default:
			}
			return len(buf)
		},
		OnJoin: func(p *noxnet.MsgServerTryJoin, full bool, add func(ntype.Player) bool) netmsg.Message {
			return &noxnet.MsgJoinOK{}
		},
		CheckPass: func(p *noxnet.MsgServerPass) netmsg.Message { return &noxnet.MsgJoinOK{} },
	}
	listener, err := server.Listen(opts)
	require.NoError(t, err)
	// Start only after binding succeeds; stop and join before closing shared state.
	stop, done := make(chan struct{}), make(chan struct{})
	go func() {
		defer close(done)
		ticker := time.NewTicker(time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-stop:
				return
			case <-ticker.C:
				server.Update()
				listener.RecvLoop(false)
				frame++
			}
		}
	}()
	defer func() { close(stop); <-done; listener.Close() }()

	client := NewStreams(nil)
	client.Xor = false
	conn, err := client.NewClient(&Options{Max: 10, BufferSize: 2048})
	require.NoError(t, err)
	defer conn.Close()
	payload := "Hello\x01\x02\x03"
	// Listen may select a free port after the preferred one; let the OS allocate
	// the client port and keep all traffic on loopback.
	err = conn.Dial("127.0.0.1", opts.Port, 0, &fakeOpts{Str: payload})
	require.NoError(t, err)
	var got []byte
	err = conn.DialWait(3*time.Second, func() {}, func() bool {
		select {
		case got = <-received:
			return true
		default:
			return false
		}
	})
	require.NoError(t, err)
	require.True(t, client.Responded)
	require.Equal(t, append([]byte{byte(netmsg.MSG_CLIENT_ACCEPT)}, []byte(payload)...), got)
}

type fakeOpts struct{ Str string }

func (opt *fakeOpts) MarshalBinary() ([]byte, error) { return []byte(opt.Str), nil }
