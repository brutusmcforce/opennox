package lobby

import (
	"crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/netip"
	"sync"
	"time"
)

// NATConfig configures the NAT traversal service of the lobby.
type NATConfig struct {
	// RelayFirst and RelayLast are the UDP ports for relay sessions. No relay if RelayFirst is zero.
	RelayFirst, RelayLast uint16
	// MaxRelays limits the number of relay sessions, and MaxRelaysPerIP the sessions of one client address.
	MaxRelays, MaxRelaysPerIP int
	// RelayIdle closes a relay session after this time without traffic.
	RelayIdle time.Duration
	// RelayRate limits the bytes per second a relay session forwards (both directions together).
	RelayRate int
	// Log is used for logging, if set.
	Log *slog.Logger
}

func (c *NATConfig) setDefaults() {
	if c.RelayLast < c.RelayFirst {
		c.RelayLast = c.RelayFirst
	}
	if c.MaxRelays <= 0 {
		c.MaxRelays = int(c.RelayLast-c.RelayFirst) + 1
	}
	if c.MaxRelaysPerIP <= 0 {
		c.MaxRelaysPerIP = 4
	}
	if c.RelayIdle <= 0 {
		c.RelayIdle = time.Minute
	}
	if c.RelayRate <= 0 {
		// a game uses a few KB/s per player; the map download of the game itself is the largest burst
		c.RelayRate = 512 << 10
	}
	if c.Log == nil {
		c.Log = slog.Default()
	}
}

const (
	natHostTTL      = 3 * NATHelloInterval * time.Second
	natMaxHosts     = 10000
	natJoinsPerIP   = 30 // per natJoinWindow
	natJoinWindow   = 10 * time.Second
	natMaxJoinIPs   = 10000
	natRelayOpenTTL = 30 * time.Second // a session the host never said hello to
)

var (
	ErrNATNoHost      = errors.New("the game host is not known to the lobby's NAT service")
	ErrNATNoRelay     = errors.New("relay is not available")
	ErrNATRelayLimit  = errors.New("too many relay sessions")
	ErrNATNotIPv4     = errors.New("only IPv4 is supported")
	errNATServiceDone = errors.New("NAT service closed")
)

type natHost struct {
	addr netip.AddrPort // the public address of the host's game socket
	seen time.Time
}

type natJoins struct {
	start time.Time
	n     int
}

// NATService is the UDP side of the lobby: it helps games behind a NAT to connect (see NATOp).
type NATService struct {
	cfg  NATConfig
	conn *net.UDPConn
	port uint16

	mu     sync.Mutex
	closed bool
	hosts  map[netip.AddrPort]natHost // by the address the lobby lists: public IP and game port
	joins  map[netip.Addr]natJoins
	relays map[uint16]*natRelay
}

// ListenNAT starts the NAT service on a UDP address. Serve must be called to process messages.
func ListenNAT(addr string, cfg NATConfig) (*NATService, error) {
	cfg.setDefaults()
	ua, err := net.ResolveUDPAddr("udp4", addr)
	if err != nil {
		return nil, err
	}
	conn, err := net.ListenUDP("udp4", ua)
	if err != nil {
		return nil, err
	}
	return &NATService{
		cfg:    cfg,
		conn:   conn,
		port:   uint16(conn.LocalAddr().(*net.UDPAddr).Port),
		hosts:  make(map[netip.AddrPort]natHost),
		joins:  make(map[netip.Addr]natJoins),
		relays: make(map[uint16]*natRelay),
	}, nil
}

// Port returns the UDP port of the service.
func (n *NATService) Port() uint16 { return n.port }

// HasRelay reports if the service has relay ports.
func (n *NATService) HasRelay() bool { return n.cfg.RelayFirst != 0 }

// Close stops the service and all relay sessions.
func (n *NATService) Close() error {
	n.mu.Lock()
	n.closed = true
	relays := n.relays
	n.relays = make(map[uint16]*natRelay)
	n.mu.Unlock()
	for _, r := range relays {
		r.close()
	}
	return n.conn.Close()
}

// Serve processes messages until the service is closed.
func (n *NATService) Serve() error {
	buf := make([]byte, 1500)
	for {
		sz, src, err := n.conn.ReadFromUDPAddrPort(buf)
		if err != nil {
			n.mu.Lock()
			closed := n.closed
			n.mu.Unlock()
			if closed {
				return nil
			}
			var ne net.Error
			if errors.As(err, &ne) && ne.Timeout() {
				continue
			}
			return err
		}
		n.handle(buf[:sz], netip.AddrPortFrom(src.Addr().Unmap(), src.Port()))
	}
}

func (n *NATService) handle(b []byte, src netip.AddrPort) {
	if !src.Addr().Is4() {
		return
	}
	op, payload, ok := ParseNAT(b)
	if !ok {
		return
	}
	now := time.Now()
	switch op {
	case NATHostHello:
		port, err := ParseNATHostHello(payload)
		if err != nil || port == 0 {
			return
		}
		n.mu.Lock()
		defer n.mu.Unlock()
		key := netip.AddrPortFrom(src.Addr(), port)
		if _, ok := n.hosts[key]; !ok && len(n.hosts) >= natMaxHosts {
			n.gcLocked(now)
			if len(n.hosts) >= natMaxHosts {
				return
			}
		}
		if old, ok := n.hosts[key]; !ok || old.addr != src {
			n.cfg.Log.Info("NAT host", "game", key, "addr", src)
		}
		n.hosts[key] = natHost{addr: src, seen: now}
	case NATJoin:
		host, err := ParseNATJoin(payload)
		if err != nil {
			return
		}
		h, ok := n.joinHost(src.Addr(), host, now)
		if !ok {
			return
		}
		_, _ = n.conn.WriteToUDPAddrPort(NATPunchMsg(src), h.addr)
	}
}

// joinHost checks the rate limit of a client address, and returns the host for a join.
func (n *NATService) joinHost(client netip.Addr, host netip.AddrPort, now time.Time) (natHost, bool) {
	n.mu.Lock()
	defer n.mu.Unlock()
	j, ok := n.joins[client]
	if !ok && len(n.joins) >= natMaxJoinIPs {
		n.gcLocked(now)
		if len(n.joins) >= natMaxJoinIPs {
			return natHost{}, false
		}
	}
	if now.Sub(j.start) > natJoinWindow {
		j = natJoins{start: now}
	}
	j.n++
	n.joins[client] = j
	if j.n > natJoinsPerIP {
		return natHost{}, false
	}
	h, ok := n.hosts[netip.AddrPortFrom(host.Addr().Unmap(), host.Port())]
	if !ok || now.Sub(h.seen) > natHostTTL {
		return natHost{}, false
	}
	return h, true
}

func (n *NATService) gcLocked(now time.Time) {
	for k, h := range n.hosts {
		if now.Sub(h.seen) > natHostTTL {
			delete(n.hosts, k)
		}
	}
	for k, j := range n.joins {
		if now.Sub(j.start) > natJoinWindow {
			delete(n.joins, k)
		}
	}
}

// OpenRelay opens a relay session between a client address and a host, as the lobby lists it, and returns the relay
// port. A client that asks again for the same host gets the same session.
func (n *NATService) OpenRelay(client netip.Addr, host netip.AddrPort) (uint16, error) {
	client = client.Unmap()
	host = netip.AddrPortFrom(host.Addr().Unmap(), host.Port())
	if !client.Is4() || !host.Addr().Is4() {
		return 0, ErrNATNotIPv4
	}
	if !n.HasRelay() {
		return 0, ErrNATNoRelay
	}
	now := time.Now()
	n.mu.Lock()
	defer n.mu.Unlock()
	if n.closed {
		return 0, errNATServiceDone
	}
	h, ok := n.hosts[host]
	if !ok || now.Sub(h.seen) > natHostTTL {
		return 0, ErrNATNoHost
	}
	perIP := 0
	for _, r := range n.relays {
		if r.clientIP == client {
			if r.hostKey == host {
				// a retry: tell the host again, it may have missed it
				n.sendRelayOpen(r, h.addr)
				return r.port, nil
			}
			perIP++
		}
	}
	if perIP >= n.cfg.MaxRelaysPerIP || len(n.relays) >= n.cfg.MaxRelays {
		return 0, ErrNATRelayLimit
	}
	var conn *net.UDPConn
	port := n.cfg.RelayFirst
	for ; ; port++ {
		if _, used := n.relays[port]; !used {
			c, err := net.ListenUDP("udp4", &net.UDPAddr{Port: int(port)})
			if err == nil {
				conn = c
				break
			}
		}
		if port >= n.cfg.RelayLast {
			return 0, ErrNATRelayLimit
		}
	}
	var id [4]byte
	_, _ = rand.Read(id[:])
	r := &natRelay{
		n:        n,
		port:     port,
		id:       binary.LittleEndian.Uint32(id[:]),
		conn:     conn,
		clientIP: client,
		hostKey:  host,
		hostIP:   h.addr.Addr(),
		opened:   now,
		tokens:   float64(n.cfg.RelayRate),
		refill:   now,
	}
	n.relays[port] = r
	n.cfg.Log.Info("relay open", "port", port, "client", client, "host", host)
	go r.serve()
	n.sendRelayOpen(r, h.addr)
	return port, nil
}

func (n *NATService) sendRelayOpen(r *natRelay, host netip.AddrPort) {
	msg := NATRelayOpenMsg(r.port, r.id)
	go func() {
		// UDP: say it a few times
		for _, dt := range []time.Duration{0, 200 * time.Millisecond, 800 * time.Millisecond} {
			time.Sleep(dt)
			_, _ = n.conn.WriteToUDPAddrPort(msg, host)
		}
	}()
}

func (n *NATService) relayDone(r *natRelay) {
	n.mu.Lock()
	if n.relays[r.port] == r {
		delete(n.relays, r.port)
	}
	n.mu.Unlock()
}

// natRelay is a relay session: it forwards packets between a client and a host on its own port, so the host sees
// each relayed client as a separate address.
type natRelay struct {
	n        *NATService
	port     uint16
	id       uint32
	conn     *net.UDPConn
	clientIP netip.Addr
	hostKey  netip.AddrPort // the listed host address
	hostIP   netip.Addr     // the public IP the host's hellos come from
	opened   time.Time

	// only used by serve
	client netip.AddrPort // the last client socket that sent to the relay
	host   netip.AddrPort // the host socket that said hello
	last   time.Time
	tokens float64
	refill time.Time
	sent   int64
}

func (r *natRelay) close() { _ = r.conn.Close() }

func (r *natRelay) allow(sz int, now time.Time) bool {
	rate := float64(r.n.cfg.RelayRate)
	r.tokens += now.Sub(r.refill).Seconds() * rate
	r.refill = now
	if r.tokens > rate {
		r.tokens = rate
	}
	if r.tokens < float64(sz) {
		return false
	}
	r.tokens -= float64(sz)
	return true
}

func (r *natRelay) serve() {
	defer r.n.relayDone(r)
	defer r.close()
	buf := make([]byte, 2048)
	for {
		now := time.Now()
		if r.last.IsZero() {
			if now.Sub(r.opened) > natRelayOpenTTL {
				r.n.cfg.Log.Info("relay unused", "port", r.port)
				return
			}
		} else if now.Sub(r.last) > r.n.cfg.RelayIdle {
			r.n.cfg.Log.Info("relay idle", "port", r.port, "bytes", r.sent)
			return
		}
		_ = r.conn.SetReadDeadline(now.Add(time.Second))
		sz, src, err := r.conn.ReadFromUDPAddrPort(buf)
		if err != nil {
			var ne net.Error
			if errors.As(err, &ne) && ne.Timeout() {
				continue
			}
			return
		}
		src = netip.AddrPortFrom(src.Addr().Unmap(), src.Port())
		b := buf[:sz]
		now = time.Now()
		if op, payload, ok := ParseNAT(b); ok && op == NATRelayHello {
			if id, err := ParseNATRelayHello(payload); err == nil && id == r.id && src.Addr() == r.hostIP {
				if r.host != src {
					r.n.cfg.Log.Info("relay host", "port", r.port, "addr", src)
				}
				r.host = src
			}
			continue
		}
		var dst netip.AddrPort
		switch {
		case r.host.IsValid() && src == r.host:
			dst = r.client
		case src.Addr() == r.clientIP:
			r.client = src
			dst = r.host
		default:
			continue
		}
		if !dst.IsValid() || !r.allow(sz, now) {
			continue
		}
		if _, err := r.conn.WriteToUDPAddrPort(b, dst); err == nil {
			r.last = now
			r.sent += int64(sz)
		}
	}
}

// String implements fmt.Stringer.
func (r *natRelay) String() string {
	return fmt.Sprintf("relay(%d: %s <-> %s)", r.port, r.clientIP, r.hostKey)
}
