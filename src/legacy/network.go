package legacy

import (
	"net"
	"net/netip"
	"unsafe"

	"github.com/brutusmcforce/opennox/libs/noxnet/netmsg"

	"github.com/brutusmcforce/opennox/client"
	"github.com/brutusmcforce/opennox/common/ntype"
	"github.com/brutusmcforce/opennox/internal/netstr"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"github.com/brutusmcforce/opennox/server"
)

var (
	NetworkLogPrint                   func(str string)
	ClientSetServerHost               func(addr string)
	Nox_client_joinGame_438A90        func() int
	SendXXX_5550D0                    func(addr netip.AddrPort, data []byte) (int, error)
	Nox_xxx_netStatsMultiplier_4D9C20 func(a1p *server.Object) int
	Sub_554240                        func(a1 ntype.PlayerInd) int
	Nox_xxx_netOnPacketRecvCli_48EA70 func(ind ntype.PlayerInd, buf *byte, sz int) int
	Sub_43C6E0                        func() int
	Sub_43CF40                        func()
	Sub_43CF70                        func()
)

func int2ip(v uint32) netip.Addr {
	b := (*[4]byte)(unsafe.Pointer(&v))[:]
	ip := net.IPv4(b[0], b[1], b[2], b[3])
	addr, _ := netip.AddrFromSlice(ip.To4())
	return addr
}

func ip2int(ip netip.Addr) uint32 {
	if !ip.IsValid() {
		return 0
	}
	b := ip.As4()
	v := *(*uint32)(unsafe.Pointer(&b[0]))
	return v
}

func nox_xxx_net_getIP_554200(a1 int) uint32 {
	if a1 < 0 || a1 >= 31 {
		panic("unexpected index")
	}
	var conn *netstr.Conn
	if a1 == 0 {
		conn = GetServer().S().NetStr.Host()
	} else {
		conn = GetServer().S().NetStr.ConnByPlayerInd(ntype.PlayerInd(a1) + 1)
	}
	return ip2int(GetServer().S().GetExtIP(conn))
}

func ClientGetServerPort() int {
	return int(nox_client_getServerPort_43B320())
}

func Sub_43AF90(v int) {
	browserUI.connectionState = uint32(v)
}

func Nox_xxx_netSendPacket_4E5030(a1 int, buf []byte, a4, a5, a6 int) int {
	return reliableEnqueue(a1, buf, (*server.Object)(unsafe.Pointer(uintptr(uint32(a4)))), a5, byte(a6))
}

func Nox_client_getServerAddr_43B300() netip.Addr {
	return int2ip(uint32(nox_client_getServerAddr_43B300()))
}

func Nox_xxx_netSendLineMessage_4D9EB0(u *server.Object, s string) bool {
	_ = netmsg.MSG_TEXT_MESSAGE
	cstr, free := CWString(s)
	defer free()
	return textFormatLine(u, (*uint16)(unsafe.Pointer(cstr))) != 0
}

func Nox_server_makeServerInfoPacket_554040(src, dst []byte) int {
	return serverTextListing(src, dst)
}

func Sub_40A740() int {
	return int(int32(serverConfigSpecialMode()))
}

func Sub_417DE0() int {
	return int(teamRuntimeGroupCount())
}

func Nox_xxx_countObserverPlayers_425BF0() int {
	return runtimeObserverCount()
}

func Sub_43C650() {
	clientFrameSample()
}

func Sub_41D6C0() {
	// The former service list has no population path.
}

func Sub_49C7A0() {
	sub_49C7A0()
}
func Sub_467CA0() {
	uiInventoryResetClosedScroll()
}
func Sub_48D660() {
	clientSequencePoll()
}
func Sub_40A220() int {
	return int(int32(serverConfigTimerGet()))
}
func Sub_40A230() uint32 {
	return uint32(serverConfigTimerLeft())
}

func convSendToServerErr(n int, err error) int {
	if err == client.ErrLobbyNoSocket {
		return -17
	} else if err != nil {
		return -1
	}
	return n
}
func Nox_xxx_netClientSend2_4E53C0(a1 int, a2 unsafe.Pointer, a3 int, a4 int, a5 int) {
	reliableClientSend(a1, unsafe.Slice((*byte)(a2), a3), (*server.Object)(unsafe.Pointer(uintptr(uint32(a4)))), a5)
}
func Sub_57B920(a1 unsafe.Pointer) {
	resetNetworkAliases((*[255]server.PlayerNetData)(a1))
}
func Nox_xxx_cliSetSettingsAcquired_4169D0(a1 int) {
	serverConfigAcquiredSet(int32(a1))
}
func Sub_457140(a1 int, a2 *uint16) {
	teamUIPlayerAdd(a1, alloc.GoString16(a2))
}
func Sub_455920(a1 *uint16) {
	serverPanelsPlayerAdd(a1)
}
func Sub_456DF0(a1 int) {
	teamUIPlayerRemove(a1)
}
func Sub_455950(a1 *uint16) {
	serverPanelsPlayerRemove(a1)
}
func Nox_xxx_netChangeTeamMb_419570(a1 *server.ObjectTeam, a2 uint32) {
	teamRuntimeLeave(a1, int(a2))
}
func Sub_49BB80(a1 byte) {
	presentationChantStart(a1)
}
func Nox_xxx_netOnPacketRecvCli_48EA70_switch(a1 ntype.PlayerInd, a2 netmsg.Op, data []byte) int {
	if n, handled := clientGameState(a2, data); handled {
		return n
	}
	if n, handled := clientGameProgress(a2, data); handled {
		return n
	}
	if n, handled := clientGameEffects(int(a1), a2, data); handled {
		return n
	}
	return clientGameSession(int(a1), a2, data)
}
func Sub_4DDE10(a1 int, a2 *server.Player) {
	matchRosterInventory(a1, a2)
}
func Nox_xxx_netPlayerObjSend_518C30(a1 *server.Object, a2 *server.Object, a3 int, a4 int) int {
	return objectReportPlayer(a1, a2, a3, a4)
}
func Nox_xxx_gameServerReadyMB_4DD180(a1 int) {
	sessionPlayerReady(int32(a1))
}
func Nox_xxx_teamCompare2_419180(t *server.ObjectTeam, id server.TeamID) int {
	return int(teamRuntimeBool(teamRuntimeContains(t, id)))
}
func Sub_4D12A0(a1 int) int {
	return int(sessionRosterContains(int32(a1)))
}
func Sub_4D1210(a1 int) {
	sessionRosterAdd(int32(a1))
}
func Nox_net_importantACK_4E55A0(a1 int, a2 int) {
	reliableACK(a1, uint32(a2))
}
func Sub_4196D0(a1 unsafe.Pointer, a2 unsafe.Pointer, a3 int, a4 int) {
	teamRuntimeSwitch(teamRuntimeMember(a1), asTeamP(a2), a3, a4)
}

func Nox_xxx_netReportAcquireCreature_4D91A0(pli int, obj *server.Object) {
	gameplayReportAcquireCreature(int(int32(pli)), obj)
}
func Nox_xxx_netSendSimpleObject2_4DF360(pli int, obj *server.Object) {
	matchRosterSimpleObject(int(pli), obj)
}
func Nox_xxx_netCode2ChatBubble_48D850(a1 int) int {
	return int(uintptr(unsafe.Pointer(chatBubbleLookup(uint32(a1)))))
}
