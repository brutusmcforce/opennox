package legacy

import (
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/server"
	"unsafe"
)

func nox_xxx_netClientSend2_4E53C0(to int32, data unsafe.Pointer, size, related, priority int32) int32 {
	return int32(reliableClientSend(int(to), unsafe.Slice((*byte)(data), int(size)), (*server.Object)(memmap.WordPtr(uint32(related))), int(priority)))
}

func nox_xxx_netSendPacket0_4E5420(to int32, data unsafe.Pointer, size, related, priority int32) int32 {
	return int32(reliableEnqueue(int(to), unsafe.Slice((*byte)(data), int(size)), (*server.Object)(memmap.WordPtr(uint32(related))), int(priority), 0))
}

func nox_net_importantACK_4E55A0(to, frame int32) int32 {
	return int32(reliableACK(int(to), uint32(frame)))
}
