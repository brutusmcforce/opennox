package opennox

import (
	"context"
	"net"
	"net/netip"
	"sync"
	"time"

	"github.com/brutusmcforce/opennox/libs/log"
	"github.com/brutusmcforce/opennox/lobby"

	"github.com/brutusmcforce/opennox/common/discover"
	"github.com/brutusmcforce/opennox/legacy"
)

// NAT traversal through the lobby (see lobby.NATOp): hosts tell the lobby their public address, clients ask the lobby
// to open the host's router for them (hole punching), and fall back to a relay of the lobby if that doesn't work.

var (
	natLog     = log.New("nat-punch")
	natEnabled = true
)

func init() {
	configBoolPtr("network.nat_punch", "NOX_NET_PUNCH", true, &natEnabled)
	legacy.BrowserJoinWaitHook = natClientJoinWait
}

const (
	natResolveRetry = 30 * time.Second
	natResolveTTL   = 10 * time.Minute
	natRelayAfter   = 3 * time.Second  // a join without an answer uses the relay after this time
	natResend       = time.Second      // the join request is resent this often while waiting
	natRelayHellos  = time.Minute      // the host says hello to a new relay session for this long
	natDialPunches  = 10 * time.Second // a connecting client asks for punches this long
	natDialWait     = 400 * time.Millisecond
)

// natResolver finds the NAT service of the lobby once, in the background.
var natResolver struct {
	mu   sync.Mutex
	lb   *discover.NATLobby
	at   time.Time
	busy bool
}

// natLobby returns the NAT service of the lobby, if it's known. It doesn't block: the first call starts looking for it.
func natLobby() *discover.NATLobby {
	if !natEnabled {
		return nil
	}
	r := &natResolver
	r.mu.Lock()
	defer r.mu.Unlock()
	ttl := natResolveTTL
	if r.lb == nil {
		ttl = natResolveRetry
	}
	if !r.busy && (r.at.IsZero() || time.Since(r.at) > ttl) {
		r.busy = true
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			lb, err := discover.ResolveNAT(ctx)
			if err != nil {
				natLog.Info("lobby has no NAT service", "err", err)
			} else {
				natLog.Info("lobby NAT service", "addr", lb.UDP, "relay", lb.Relay)
			}
			r.mu.Lock()
			defer r.mu.Unlock()
			r.busy = false
			r.at = time.Now()
			if lb != nil {
				r.lb = lb // keep the last one if the lobby doesn't answer for a while
			}
		}()
	}
	return r.lb
}

// natPublic reports if an address may be behind a NAT the lobby can help with. LAN games need no help.
func natPublic(addr netip.AddrPort) bool {
	ip := addr.Addr().Unmap()
	return ip.Is4() && addr.Port() != 0 && !ip.IsPrivate() && !ip.IsLoopback() && !ip.IsLinkLocalUnicast() &&
		!ip.IsUnspecified() && !ip.IsMulticast()
}

func natWrite(pc net.PacketConn, b []byte, addr netip.AddrPort) {
	if pc == nil || !addr.IsValid() {
		return
	}
	_, _ = pc.WriteTo(b, net.UDPAddrFromAddrPort(addr))
}

// natOnPacket handles the NAT messages that arrive on the game sockets.
func natOnPacket(pc net.PacketConn, b []byte, from netip.AddrPort) bool {
	op, payload, ok := lobby.ParseNAT(b)
	if !ok {
		return false
	}
	from = netip.AddrPortFrom(from.Addr().Unmap(), from.Port())
	switch op {
	case lobby.NATPunch:
		natHost.onPunch(pc, payload, from)
	case lobby.NATRelayOpen:
		natHost.onRelayOpen(pc, payload, from)
	}
	// NATPing only opens the routers: nothing to do
	return true
}

// Host side.

type natHostRelay struct {
	id    uint32
	until time.Time
	next  time.Time
}

type natHostService struct {
	mu      sync.Mutex
	active  bool
	cancel  func()
	pc      net.PacketConn
	lb      *discover.NATLobby
	relays  map[uint16]*natHostRelay
	punchAt time.Time
	punches int
}

var natHost natHostService

// natHostStart tells the lobby the public address of the game socket, so clients can ask for hole punching.
func (s *Server) natHostStart() {
	pc := s.NetStr.ListenConn()
	if !natEnabled || pc == nil {
		return
	}
	h := &natHost
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.active && h.pc == pc {
		return
	}
	if h.cancel != nil {
		h.cancel()
	}
	ctx, cancel := context.WithCancel(context.Background())
	h.active, h.cancel, h.pc, h.lb = true, cancel, pc, nil
	h.relays = make(map[uint16]*natHostRelay)
	go h.run(ctx, pc, uint16(s.ServerPort()))
}

func (s *Server) natHostStop() {
	h := &natHost
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.cancel != nil {
		h.cancel()
	}
	h.active, h.cancel, h.pc, h.lb, h.relays = false, nil, nil, nil, nil
}

func (h *natHostService) run(ctx context.Context, pc net.PacketConn, port uint16) {
	tick := time.NewTicker(time.Second)
	defer tick.Stop()
	var next time.Time
	for {
		select {
		case <-ctx.Done():
			return
		case <-tick.C:
		}
		lb := natLobby()
		if lb == nil {
			continue
		}
		now := time.Now()
		h.mu.Lock()
		if !h.active || h.pc != pc {
			h.mu.Unlock()
			return
		}
		if h.lb == nil {
			natLog.Info("announcing the game socket to the lobby", "lobby", lb.UDP, "port", port, "local", pc.LocalAddr())
		}
		h.lb = lb
		for rport, r := range h.relays {
			if now.After(r.until) {
				delete(h.relays, rport)
			} else if !now.Before(r.next) {
				r.next = now.Add(3 * time.Second)
				natWrite(pc, lobby.NATRelayHelloMsg(r.id), lb.RelayAddr(rport))
			}
		}
		h.mu.Unlock()
		if !now.Before(next) {
			next = now.Add(lobby.NATHelloInterval * time.Second)
			natWrite(pc, lobby.NATHostHelloMsg(port), lb.UDP)
		}
	}
}

func (h *natHostService) fromLobby(pc net.PacketConn, from netip.AddrPort) *discover.NATLobby {
	if !h.active || h.pc != pc || h.lb == nil || from != h.lb.UDP {
		return nil
	}
	return h.lb
}

func (h *natHostService) onPunch(pc net.PacketConn, payload []byte, from netip.AddrPort) {
	client, err := lobby.ParseNATPunch(payload)
	if err != nil {
		return
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.fromLobby(pc, from) == nil || !natPublic(client) {
		return
	}
	// the lobby limits joins, but keep the pings bounded anyway
	now := time.Now()
	if now.Sub(h.punchAt) > time.Second {
		h.punchAt, h.punches = now, 0
	}
	if h.punches++; h.punches > 50 {
		return
	}
	ping := lobby.NATPingMsg()
	natWrite(pc, ping, client)
	natWrite(pc, ping, client)
}

func (h *natHostService) onRelayOpen(pc net.PacketConn, payload []byte, from netip.AddrPort) {
	port, id, err := lobby.ParseNATRelayOpen(payload)
	if err != nil || port == 0 {
		return
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	lb := h.fromLobby(pc, from)
	if lb == nil {
		return
	}
	now := time.Now()
	if r := h.relays[port]; r != nil && r.id == id {
		return // the lobby repeats it
	}
	if len(h.relays) >= 64 {
		return
	}
	natLog.Info("relay session", "addr", lb.RelayAddr(port))
	h.relays[port] = &natHostRelay{id: id, until: now.Add(natRelayHellos), next: now.Add(500 * time.Millisecond)}
	natWrite(pc, lobby.NATRelayHelloMsg(id), lb.RelayAddr(port))
}

// Client side.

// natJoin is the state of the join request the client waits an answer for (only used by the game loop).
var natJoin struct {
	active   bool
	host     netip.AddrPort // the game, as the lobby lists it
	pkt      []byte         // the join request, to resend it
	start    time.Time
	lastSend time.Time
	relay    netip.AddrPort // set once the client uses a relay
	relayReq bool
	relayRes chan netip.AddrPort
}

// natClientJoinStart is called when the client sends its join request to a game.
func natClientJoinStart(target netip.AddrPort, pkt []byte) {
	j := &natJoin
	if j.active && j.relay.IsValid() && target == j.relay {
		j.pkt = append(j.pkt[:0], pkt...) // a new request through the relay, e.g. with a password
		return
	}
	j.active = false
	if !natPublic(target) {
		return
	}
	now := time.Now()
	j.active = true
	j.host = target
	j.pkt = append(j.pkt[:0], pkt...)
	j.start, j.lastSend = now, now
	j.relay = netip.AddrPort{}
	j.relayReq = false
	j.relayRes = make(chan netip.AddrPort, 1)
	if lb := natLobby(); lb != nil {
		natWrite(lobbyBroadcast, lobby.NATJoinMsg(target), lb.UDP)
	}
}

// natClientJoinWait runs on every frame while the client waits for the answer to its join request.
func natClientJoinWait() {
	j := &natJoin
	if !j.active || lobbyBroadcast == nil {
		return
	}
	now := time.Now()
	lb := natLobby()
	select {
	case addr := <-j.relayRes:
		natLog.Info("joining through the lobby relay", "game", j.host, "relay", addr)
		j.relay = addr
		legacy.SetSelectedServerAddr(addr)
		j.lastSend = now
		_, _ = sendJoinGame(addr, j.pkt)
		return
	default:
	}
	if now.Sub(j.lastSend) >= natResend {
		// the first requests may get lost until the host's router lets them through
		j.lastSend = now
		target := j.host
		if j.relay.IsValid() {
			target = j.relay
		} else if lb != nil {
			natWrite(lobbyBroadcast, lobby.NATJoinMsg(j.host), lb.UDP)
		}
		_, _ = sendJoinGame(target, j.pkt)
	}
	if !j.relay.IsValid() && !j.relayReq && now.Sub(j.start) >= natRelayAfter && lb != nil && lb.Relay {
		j.relayReq = true
		host, res := j.host, j.relayRes
		natLog.Info("no answer from the game, asking the lobby for a relay", "game", host)
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			addr, err := lb.OpenRelay(ctx, host)
			if err != nil {
				natLog.Info("no relay", "game", host, "err", err)
				return
			}
			res <- addr
		}()
	}
}

// natClientListedAddr returns the address the lobby lists for the selected game, which differs from the address the
// client uses with a relay.
func natClientListedAddr(addr netip.AddrPort) netip.AddrPort {
	if j := &natJoin; j.relay.IsValid() && addr == j.relay {
		return j.host
	}
	return addr
}

// natOnDial is called when the client opens the socket it connects to the game with: that socket needs its own hole.
func natOnDial(pc net.PacketConn, addr netip.AddrPort) {
	addr = netip.AddrPortFrom(addr.Addr().Unmap(), addr.Port())
	if j := &natJoin; j.relay.IsValid() && addr == j.relay {
		return
	}
	lb := natLobby()
	if lb == nil || !natPublic(addr) {
		return
	}
	msg := lobby.NATJoinMsg(addr)
	natWrite(pc, msg, lb.UDP)
	// Give the host's ping time to open its router before the first packet: otherwise it's dropped, and the client
	// only resends it a second later.
	time.Sleep(natDialWait)
	go func() {
		for i := time.Duration(0); i < natDialPunches; i += natResend {
			time.Sleep(natResend)
			if _, err := pc.WriteTo(msg, net.UDPAddrFromAddrPort(lb.UDP)); err != nil {
				return // closed
			}
		}
	}()
}
