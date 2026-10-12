package legacy

import (
	"unsafe"

	"github.com/brutusmcforce/opennox/common/ntype"
	"github.com/brutusmcforce/opennox/internal/netlist"
)

var (
	GetNetPlayerBufSize         func() int
	Nox_netlist_addToMsgListSrv func(ind ntype.PlayerInd, buf []byte) bool
)

func nox_netlist_addToMsgListCli_40EBC0(ind1, ind2 int, buf *uint8, sz int) int {
	return bool2int(GetServer().S().NetList.AddToMsgListCli(ntype.PlayerInd(ind1), netlist.Kind(ind2), unsafe.Slice((*byte)(unsafe.Pointer(buf)), sz)))
}

func Nox_xxx_netImportant_4E5770(a1 byte, a2 int) {
	reliableDeliver(a1, a2)
}
