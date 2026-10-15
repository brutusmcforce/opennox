package lobby

import (
	"net"
	"net/netip"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func natTestSock(t testing.TB) *net.UDPConn {
	c, err := net.ListenUDP("udp4", &net.UDPAddr{IP: net.IPv4(127, 0, 0, 1)})
	require.NoError(t, err)
	t.Cleanup(func() { _ = c.Close() })
	return c
}

func natAddr(c *net.UDPConn) netip.AddrPort {
	a := c.LocalAddr().(*net.UDPAddr).AddrPort()
	return netip.AddrPortFrom(a.Addr().Unmap(), a.Port())
}

func natRecv(t testing.TB, c *net.UDPConn) ([]byte, netip.AddrPort) {
	t.Helper()
	buf := make([]byte, 1500)
	_ = c.SetReadDeadline(time.Now().Add(3 * time.Second))
	n, src, err := c.ReadFromUDPAddrPort(buf)
	require.NoError(t, err)
	return buf[:n], netip.AddrPortFrom(src.Addr().Unmap(), src.Port())
}

// natRecvOp skips other NAT messages (the lobby repeats some of them).
func natRecvOp(t testing.TB, c *net.UDPConn, want NATOp) ([]byte, netip.AddrPort) {
	t.Helper()
	for {
		b, src := natRecv(t, c)
		op, payload, ok := ParseNAT(b)
		require.True(t, ok)
		if op == want {
			return payload, src
		}
	}
}

func natFreeRange(t testing.TB, n int) (uint16, uint16) {
	// find a few free consecutive ports
	for first := 38600; first < 60000; first += 100 {
		ok := true
		for p := first; p < first+n; p++ {
			c, err := net.ListenUDP("udp4", &net.UDPAddr{Port: p})
			if err != nil {
				ok = false
				break
			}
			_ = c.Close()
		}
		if ok {
			return uint16(first), uint16(first + n - 1)
		}
	}
	t.Skip("no free ports")
	return 0, 0
}

func TestNATPunchAndRelay(t *testing.T) {
	first, last := natFreeRange(t, 2)
	svc, err := ListenNAT("127.0.0.1:0", NATConfig{RelayFirst: first, RelayLast: last})
	require.NoError(t, err)
	defer func() { _ = svc.Close() }()
	go func() { _ = svc.Serve() }()
	lobbyAddr := netip.AddrPortFrom(netip.AddrFrom4([4]byte{127, 0, 0, 1}), svc.Port())

	host := natTestSock(t)
	cli := natTestSock(t)
	hostAddr := natAddr(host)

	// a join for an unknown host does nothing; the host registers with a hello
	_, err = host.WriteToUDPAddrPort(NATHostHelloMsg(hostAddr.Port()), lobbyAddr)
	require.NoError(t, err)
	require.Eventually(t, func() bool {
		svc.mu.Lock()
		defer svc.mu.Unlock()
		return len(svc.hosts) == 1
	}, time.Second, 10*time.Millisecond)

	// hole punching: the host learns the client's address
	_, err = cli.WriteToUDPAddrPort(NATJoinMsg(hostAddr), lobbyAddr)
	require.NoError(t, err)
	payload, src := natRecvOp(t, host, NATPunch)
	require.Equal(t, lobbyAddr, src)
	caddr, err := ParseNATPunch(payload)
	require.NoError(t, err)
	require.Equal(t, natAddr(cli), caddr)

	// relay
	_, err = svc.OpenRelay(netip.MustParseAddr("127.0.0.1"), netip.MustParseAddrPort("127.0.0.1:1"))
	require.ErrorIs(t, err, ErrNATNoHost)
	port, err := svc.OpenRelay(netip.MustParseAddr("127.0.0.1"), hostAddr)
	require.NoError(t, err)
	port2, err := svc.OpenRelay(netip.MustParseAddr("127.0.0.1"), hostAddr)
	require.NoError(t, err)
	require.Equal(t, port, port2, "a retry gets the same session")

	payload, _ = natRecvOp(t, host, NATRelayOpen)
	rport, id, err := ParseNATRelayOpen(payload)
	require.NoError(t, err)
	require.Equal(t, port, rport)
	relayAddr := netip.AddrPortFrom(lobbyAddr.Addr(), port)

	// client packets before the host's hello are dropped
	_, err = cli.WriteToUDPAddrPort([]byte("early"), relayAddr)
	require.NoError(t, err)
	_, err = host.WriteToUDPAddrPort(NATRelayHelloMsg(id+1), relayAddr) // wrong session: ignored
	require.NoError(t, err)
	_, err = host.WriteToUDPAddrPort(NATRelayHelloMsg(id), relayAddr)
	require.NoError(t, err)
	time.Sleep(50 * time.Millisecond)

	_, err = cli.WriteToUDPAddrPort([]byte("hello"), relayAddr)
	require.NoError(t, err)
	for {
		b, src := natRecv(t, host)
		if _, _, ok := ParseNAT(b); ok {
			continue // a repeated relay open
		}
		require.Equal(t, "hello", string(b))
		require.Equal(t, relayAddr, src)
		break
	}
	_, err = host.WriteToUDPAddrPort([]byte("world"), relayAddr)
	require.NoError(t, err)
	b, src := natRecv(t, cli)
	require.Equal(t, "world", string(b))
	require.Equal(t, relayAddr, src)

	// limits
	svc.cfg.MaxRelays = 1
	_, err = svc.OpenRelay(netip.MustParseAddr("127.0.0.2"), hostAddr)
	require.ErrorIs(t, err, ErrNATRelayLimit)
}

func TestNATJoinRateLimit(t *testing.T) {
	svc, err := ListenNAT("127.0.0.1:0", NATConfig{})
	require.NoError(t, err)
	defer func() { _ = svc.Close() }()
	now := time.Now()
	host := netip.MustParseAddrPort("10.0.0.1:18590")
	svc.hosts[host] = natHost{addr: host, seen: now}
	client := netip.MustParseAddr("10.0.0.2")
	for i := 0; i < natJoinsPerIP; i++ {
		_, ok := svc.joinHost(client, host, now)
		require.True(t, ok)
	}
	_, ok := svc.joinHost(client, host, now)
	require.False(t, ok)
	_, ok = svc.joinHost(client, host, now.Add(natJoinWindow+time.Second))
	require.True(t, ok)
	// stale hosts are not used
	_, ok = svc.joinHost(netip.MustParseAddr("10.0.0.3"), host, now.Add(natHostTTL+time.Second))
	require.False(t, ok)
}

func TestNATMessages(t *testing.T) {
	addr := netip.MustParseAddrPort("1.2.3.4:18590")
	op, payload, ok := ParseNAT(NATJoinMsg(addr))
	require.True(t, ok)
	require.Equal(t, NATJoin, op)
	got, err := ParseNATJoin(payload)
	require.NoError(t, err)
	require.Equal(t, addr, got)

	_, _, ok = ParseNAT([]byte{0x3a, 0xf1, 0x00, 0x00}) // a game extension packet
	require.False(t, ok)
	_, _, ok = ParseNAT([]byte{0x3a, 0xf1})
	require.False(t, ok)
	_, _, err = ParseNATRelayOpen([]byte{1})
	require.Error(t, err)
}
