package netxfer

import (
	"bytes"
	"testing"

	"github.com/shoenig/test/must"

	"github.com/brutusmcforce/opennox/libs/noxnet/netmsg"
)

// sinkConn drops everything the receiver sends back.
type sinkConn struct{ id int }

func (*sinkConn) SendReliable(netmsg.Message) error   { return nil }
func (*sinkConn) SendUnreliable(netmsg.Message) error { return nil }

func newTestReceiver() (*Receiver[*sinkConn], *[]Data) {
	var x Receiver[*sinkConn]
	x.Reset(minStreams)
	var got []Data
	x.OnReceive(func(_ *sinkConn, p Data) {
		got = append(got, p)
	})
	return &x, &got
}

func TestReceiveHugeSize(t *testing.T) {
	x, _ := newTestReceiver()
	c := &sinkConn{}
	x.HandleStart(c, 0, &MsgStart{Size: 0xF0000000})
	must.EqOp(t, 0, x.active)
	must.EqOp(t, 0, x.allocated)

	// many transfers below the per-transfer limit must not add up without bound either
	for range minStreams {
		x.HandleStart(c, 0, &MsgStart{Size: maxRecvSize})
	}
	must.EqOp(t, maxRecvTotal/maxRecvSize, x.active)
}

func TestReceiveOtherConn(t *testing.T) {
	x, _ := newTestReceiver()
	c1, c2 := &sinkConn{1}, &sinkConn{2}
	x.HandleStart(c1, 0, &MsgStart{Size: 10})
	must.EqOp(t, 1, x.active)
	x.HandleCancel(c2, &MsgCancel{RecvID: 0})
	must.EqOp(t, 1, x.active)
	x.HandleCancel(c1, &MsgCancel{RecvID: 0})
	must.EqOp(t, 0, x.active)
	must.EqOp(t, 0, x.allocated)
}

func TestReceiveLargeOutOfOrder(t *testing.T) {
	x, got := newTestReceiver()
	c := &sinkConn{}
	size := maxDataSize + 3*blockSize
	x.HandleStart(c, 0, &MsgStart{Size: uint32(size)})
	exp := make([]byte, size)
	for i := range exp {
		exp[i] = byte(i / blockSize)
	}
	n := size / blockSize
	order := make([]int, 0, n)
	for i := 0; i < n-3; i++ {
		order = append(order, i)
	}
	order = append(order, n-1, n-2, n-3) // the last three arrive in reverse
	for _, i := range order {
		x.HandleData(c, 0, &MsgData{RecvID: 0, Chunk: Chunk(i + 1), Data: exp[i*blockSize : (i+1)*blockSize]})
	}
	must.Len(t, 1, *got)
	must.True(t, bytes.Equal(exp, (*got)[0].Data))
	must.EqOp(t, 0, x.active)
}

func TestReceiveLargeOverflow(t *testing.T) {
	x, got := newTestReceiver()
	c := &sinkConn{}
	size := maxDataSize
	x.HandleStart(c, 0, &MsgStart{Size: uint32(size)})
	chunk := make([]byte, 0xFFFF) // larger than the remaining space at the end
	for i := 0; len(*got) == 0 && i < size/len(chunk)+2; i++ {
		x.HandleData(c, 0, &MsgData{RecvID: 0, Chunk: Chunk(i + 1), Data: chunk})
	}
	must.Len(t, 1, *got)
}
