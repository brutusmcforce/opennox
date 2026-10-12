//go:build linux

package netstr

import (
	"net"
	"syscall"
	"testing"
	"time"
)

// Exercise the Linux ioctl contract on real loopback UDP sockets.

func listenLoopbackUDP(t *testing.T) *net.UDPConn {
	t.Helper()
	c, err := net.ListenUDP("udp4", &net.UDPAddr{IP: net.IPv4(127, 0, 0, 1)})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = c.Close() })
	return c
}

func dialUDPSender(t *testing.T, dst *net.UDPConn) *net.UDPConn {
	t.Helper()
	src, err := net.DialUDP("udp4", nil, dst.LocalAddr().(*net.UDPAddr))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = src.Close() })
	return src
}

func sendUDP(t *testing.T, src *net.UDPConn, data []byte) {
	t.Helper()
	if _, err := src.Write(data); err != nil {
		t.Fatal(err)
	}
}

func waitReadableCount(t *testing.T, count func() (int, error), want int) {
	t.Helper()
	deadline := time.Now().Add(time.Second)
	for time.Now().Before(deadline) {
		got, err := count()
		if err != nil {
			t.Fatalf("readability query: %v", err)
		}
		if got == want {
			return
		}
		time.Sleep(time.Millisecond)
	}
	got, err := count()
	t.Fatalf("readable byte count = %d, want %d (err %v)", got, want, err)
}

func readDatagram(t *testing.T, pc *net.UDPConn, want []byte) {
	t.Helper()
	if err := pc.SetReadDeadline(time.Now().Add(time.Second)); err != nil {
		t.Fatal(err)
	}
	buf := make([]byte, 64)
	n, _, err := pc.ReadFromUDP(buf)
	if err != nil {
		t.Fatalf("ReadFromUDP: %v", err)
	}
	if n != len(want) || string(buf[:n]) != string(want) {
		t.Fatalf("read datagram %q (%d bytes), want %q (%d bytes)", buf[:n], n, want, len(want))
	}
}

func queryNetCanRead(t *testing.T, c *net.UDPConn) (uint32, syscall.Errno) {
	t.Helper()
	rc, err := c.SyscallConn()
	if err != nil {
		t.Fatal(err)
	}
	var n uint32
	var errno syscall.Errno
	if err := rc.Control(func(fd uintptr) { n, errno = netCanRead(fd) }); err != nil {
		t.Fatal(err)
	}
	return n, errno
}

func directCountQuery(t *testing.T, c *net.UDPConn) (int, error) {
	t.Helper()
	n, errno := queryNetCanRead(t, c)
	if errno != 0 {
		return int(n), errno
	}
	return int(n), nil
}

func connCountQuery(c *net.UDPConn) (int, error) {
	return canReadConn(false, nil, c)
}

func waitQueuedDatagram(t *testing.T, c *net.UDPConn) {
	t.Helper()
	rc, err := c.SyscallConn()
	if err != nil {
		t.Fatal(err)
	}
	deadline := time.Now().Add(time.Second)
	for time.Now().Before(deadline) {
		var peekErr error
		if err := rc.Control(func(fd uintptr) {
			_, _, peekErr = syscall.Recvfrom(int(fd), nil, syscall.MSG_PEEK|syscall.MSG_DONTWAIT)
		}); err != nil {
			t.Fatalf("RawConn.Control: %v", err)
		}
		if peekErr == nil {
			return
		}
		errno, ok := peekErr.(syscall.Errno)
		if !ok || (errno != syscall.EAGAIN && errno != syscall.EWOULDBLOCK) {
			t.Fatalf("non-consuming datagram probe: %v", peekErr)
		}
		time.Sleep(time.Millisecond)
	}
	t.Fatal("timed out waiting for queued UDP datagram")
}

func TestNetCanReadUDPQueue(t *testing.T) {
	pc := listenLoopbackUDP(t)
	if n, errno := queryNetCanRead(t, pc); n != 0 || errno != 0 {
		t.Fatalf("empty socket: count=%d errno=%v, want 0/0", n, errno)
	}

	first := []byte("first!")
	second := []byte("second-datagram")
	sender := dialUDPSender(t, pc)
	sendUDP(t, sender, first)
	waitReadableCount(t, func() (int, error) { return directCountQuery(t, pc) }, len(first))
	sendUDP(t, sender, second)
	waitQueuedDatagram(t, pc)                                                                // MSG_PEEK confirms the queue is nonempty without consuming its head.
	waitReadableCount(t, func() (int, error) { return directCountQuery(t, pc) }, len(first)) // Reports the next datagram, not a sum.
	readDatagram(t, pc, first)
	waitReadableCount(t, func() (int, error) { return directCountQuery(t, pc) }, len(second))
	readDatagram(t, pc, second)
	waitReadableCount(t, func() (int, error) { return directCountQuery(t, pc) }, 0)
}

func TestNetCanReadZeroLengthUDPDatagram(t *testing.T) {
	pc := listenLoopbackUDP(t)
	sender := dialUDPSender(t, pc)
	sendUDP(t, sender, nil)
	waitQueuedDatagram(t, pc)
	if n, errno := queryNetCanRead(t, pc); n != 0 || errno != 0 {
		t.Fatalf("queued zero-length datagram: count=%d errno=%v, want 0/0", n, errno)
	}
	readDatagram(t, pc, nil)
	waitReadableCount(t, func() (int, error) { return directCountQuery(t, pc) }, 0)
}

func TestNetCanReadInvalidFD(t *testing.T) {
	_, errno := netCanRead(^uintptr(0))
	if errno != syscall.EBADF {
		t.Fatalf("invalid fd errno = %v, want %v", errno, syscall.EBADF)
	}
}

func TestCanReadConnUDPQueue(t *testing.T) {
	pc := listenLoopbackUDP(t)
	if got, err := canReadConn(false, nil, pc); got != 0 || err != nil {
		t.Fatalf("empty socket: count=%d err=%v, want 0/nil", got, err)
	}

	first := []byte("abcde")
	second := []byte("a longer queued datagram")
	sender := dialUDPSender(t, pc)
	sendUDP(t, sender, first)
	waitReadableCount(t, func() (int, error) { return connCountQuery(pc) }, len(first))
	sendUDP(t, sender, second)
	waitQueuedDatagram(t, pc)
	waitReadableCount(t, func() (int, error) { return connCountQuery(pc) }, len(first))
	readDatagram(t, pc, first)
	waitReadableCount(t, func() (int, error) { return connCountQuery(pc) }, len(second))
	readDatagram(t, pc, second)
	if got, err := canReadConn(false, nil, pc); got != 0 || err != nil {
		t.Fatalf("drained socket: count=%d err=%v, want 0/nil", got, err)
	}
}

func TestCanReadConnZeroLengthUDPDatagram(t *testing.T) {
	pc := listenLoopbackUDP(t)
	sender := dialUDPSender(t, pc)
	sendUDP(t, sender, nil)
	waitQueuedDatagram(t, pc)
	if got, err := canReadConn(false, nil, pc); got != 0 || err != nil {
		t.Fatalf("queued zero-length datagram: count=%d err=%v, want 0/nil", got, err)
	}
	readDatagram(t, pc, nil)
}
