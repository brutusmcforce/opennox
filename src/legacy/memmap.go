package legacy

import "github.com/brutusmcforce/opennox/common/memmap"

func init() {
	memmap.RegisterBlobData(0x581450, "byte_581450", byte_581450)
	memmap.RegisterBlobData(0x587000, "byte_587000", byte_587000)
	memmap.RegisterBlobData(0x5D4594, "byte_5D4594", byte_5D4594)
	memmap.RegisterBlobData(0x973CE0, "byte_973CE0", byte_973CE0)
	memmap.RegisterBlobData(0x973F18, "byte_973F18", byte_973F18)
	memmap.RegisterBlobData(0x85B3FC, "byte_85B3FC", byte_85B3FC)
	memmap.RegisterBlobData(0x852978, "byte_852978", byte_852978)
	memmap.RegisterBlobData(0x973A20, "byte_973A20", byte_973A20)
}
