package discover

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/netip"
	"net/url"

	"github.com/brutusmcforce/opennox/lobby"
)

// NATLobby is the NAT traversal service of the lobby (see lobby.NATOp).
type NATLobby struct {
	// UDP is the address of the service.
	UDP netip.AddrPort
	// Relay is set if the lobby can relay games.
	Relay bool
}

// RelayAddr returns the address of a relay session port.
func (n *NATLobby) RelayAddr(port uint16) netip.AddrPort {
	return netip.AddrPortFrom(n.UDP.Addr(), port)
}

// ResolveNAT finds the NAT traversal service of the lobby.
func ResolveNAT(ctx context.Context) (*NATLobby, error) {
	cli := newLobbyClient()
	u, err := url.Parse(cli.ServerURL())
	if err != nil {
		return nil, err
	}
	info, err := cli.NATInfo(ctx)
	if err != nil {
		return nil, err
	}
	// the service runs on the lobby's host; the game is IPv4 only
	ips, err := net.DefaultResolver.LookupNetIP(ctx, "ip4", u.Hostname())
	if err != nil {
		return nil, err
	}
	if len(ips) == 0 {
		return nil, fmt.Errorf("no IPv4 address for %q", u.Hostname())
	}
	if info.UDPPort == 0 {
		return nil, errors.New("lobby has no NAT service")
	}
	return &NATLobby{UDP: netip.AddrPortFrom(ips[0].Unmap(), info.UDPPort), Relay: info.Relay}, nil
}

// OpenRelay asks the lobby for a relay session to a game (as the lobby lists it), and returns the relay address.
func (n *NATLobby) OpenRelay(ctx context.Context, game netip.AddrPort) (netip.AddrPort, error) {
	port, err := newLobbyClient().OpenRelay(ctx, game.Addr().String(), int(game.Port()))
	if err != nil {
		return netip.AddrPort{}, err
	}
	return n.RelayAddr(port), nil
}

// IsNATPacket reports if a packet is a NAT traversal message.
func IsNATPacket(b []byte) bool {
	_, _, ok := lobby.ParseNAT(b)
	return ok
}
