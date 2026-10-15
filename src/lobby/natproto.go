package lobby

import (
	"encoding/binary"
	"errors"
	"net/netip"
)

// NAT traversal messages.
//
// Games behind a NAT can't be reached directly unless the router forwards the game port. The lobby helps in two ways:
//
//   - Hole punching: a host sends NATHostHello from its game socket to the lobby's UDP port every NATHelloInterval,
//     so the lobby knows the host's public address. A client sends NATJoin from the socket it connects with. The lobby
//     then sends NATPunch with the client's public address to the host, which sends NATPing to the client from its
//     game socket. Both routers have then seen outgoing packets to the other side, and let the game packets through.
//   - Relay: if that fails (e.g. a symmetric NAT on both sides), the client asks for a relay over HTTP. The lobby
//     opens a UDP port for the session and sends NATRelayOpen to the host, which sends NATRelayHello to that port.
//     The relay then forwards packets between the client and the host. The client uses the relay address as the
//     server address.
//
// All messages use the extension packet framing of the game (NATPacketCode, then the operation code), which the game
// keeps out of the Nox protocol, so they can be sent on the game sockets.

// NATPacketCode starts every extension packet of the game.
const NATPacketCode = 0xF13A

// NATOp is an operation code of a NAT traversal message.
type NATOp uint16

const (
	NATHostHello  = NATOp(0x100) // host -> lobby: version u8, game port u16
	NATPunch      = NATOp(0x101) // lobby -> host: client IPv4 [4]byte, client port u16
	NATPing       = NATOp(0x102) // host -> client: no payload; only opens the host's router for the client
	NATRelayOpen  = NATOp(0x103) // lobby -> host: relay port u16, session u32
	NATRelayHello = NATOp(0x104) // host -> lobby relay port: session u32
	NATJoin       = NATOp(0x110) // client -> lobby: version u8, host IPv4 [4]byte, host port u16
)

// NATVersion is the version of the NAT traversal messages.
const NATVersion = 1

// NATHelloInterval is how often a host sends NATHostHello. Routers usually drop idle UDP mappings after 30 s or more.
const NATHelloInterval = 15

var errNATShort = errors.New("short NAT message")

// AppendNAT appends a NAT message with the given payload.
func AppendNAT(b []byte, op NATOp, payload ...byte) []byte {
	b = binary.LittleEndian.AppendUint16(b, NATPacketCode)
	b = binary.LittleEndian.AppendUint16(b, uint16(op))
	return append(b, payload...)
}

// ParseNAT returns the operation and the payload of a NAT message.
func ParseNAT(b []byte) (NATOp, []byte, bool) {
	if len(b) < 4 || binary.LittleEndian.Uint16(b) != NATPacketCode {
		return 0, nil, false
	}
	op := NATOp(binary.LittleEndian.Uint16(b[2:]))
	if op < NATHostHello || op > NATJoin {
		return 0, nil, false
	}
	return op, b[4:], true
}

func appendAddr4(b []byte, addr netip.AddrPort) []byte {
	ip := addr.Addr().Unmap().As4()
	b = append(b, ip[:]...)
	return binary.LittleEndian.AppendUint16(b, addr.Port())
}

func parseAddr4(b []byte) (netip.AddrPort, error) {
	if len(b) < 6 {
		return netip.AddrPort{}, errNATShort
	}
	return netip.AddrPortFrom(netip.AddrFrom4([4]byte(b[:4])), binary.LittleEndian.Uint16(b[4:])), nil
}

// NATHostHelloMsg returns a NATHostHello message.
func NATHostHelloMsg(gamePort uint16) []byte {
	return AppendNAT(nil, NATHostHello, NATVersion, byte(gamePort), byte(gamePort>>8))
}

// ParseNATHostHello returns the game port of a NATHostHello payload.
func ParseNATHostHello(b []byte) (uint16, error) {
	if len(b) < 3 {
		return 0, errNATShort
	}
	return binary.LittleEndian.Uint16(b[1:]), nil
}

// NATPunchMsg returns a NATPunch message for a client address.
func NATPunchMsg(client netip.AddrPort) []byte {
	return appendAddr4(AppendNAT(nil, NATPunch), client)
}

// ParseNATPunch returns the client address of a NATPunch payload.
func ParseNATPunch(b []byte) (netip.AddrPort, error) {
	return parseAddr4(b)
}

// NATPingMsg returns a NATPing message.
func NATPingMsg() []byte {
	return AppendNAT(nil, NATPing)
}

// NATRelayOpenMsg returns a NATRelayOpen message.
func NATRelayOpenMsg(port uint16, session uint32) []byte {
	b := binary.LittleEndian.AppendUint16(AppendNAT(nil, NATRelayOpen), port)
	return binary.LittleEndian.AppendUint32(b, session)
}

// ParseNATRelayOpen returns the relay port and the session of a NATRelayOpen payload.
func ParseNATRelayOpen(b []byte) (uint16, uint32, error) {
	if len(b) < 6 {
		return 0, 0, errNATShort
	}
	return binary.LittleEndian.Uint16(b), binary.LittleEndian.Uint32(b[2:]), nil
}

// NATRelayHelloMsg returns a NATRelayHello message.
func NATRelayHelloMsg(session uint32) []byte {
	return binary.LittleEndian.AppendUint32(AppendNAT(nil, NATRelayHello), session)
}

// ParseNATRelayHello returns the session of a NATRelayHello payload.
func ParseNATRelayHello(b []byte) (uint32, error) {
	if len(b) < 4 {
		return 0, errNATShort
	}
	return binary.LittleEndian.Uint32(b), nil
}

// NATJoinMsg returns a NATJoin message for a host address, as the lobby lists it.
func NATJoinMsg(host netip.AddrPort) []byte {
	return appendAddr4(AppendNAT(nil, NATJoin, NATVersion), host)
}

// ParseNATJoin returns the host address of a NATJoin payload.
func ParseNATJoin(b []byte) (netip.AddrPort, error) {
	if len(b) < 1 {
		return netip.AddrPort{}, errNATShort
	}
	return parseAddr4(b[1:])
}
