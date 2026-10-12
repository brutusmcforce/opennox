package legacy

import (
	"unsafe"

	"github.com/brutusmcforce/opennox/internal/protection"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"github.com/brutusmcforce/opennox/server"
)

func objectProtectionChecksum(obj *server.Object) uint32 {
	if obj == nil {
		return 0
	}
	var sum uint32
	sum ^= obj.NetCode
	sum ^= uint32(obj.TypeInd)
	if obj.HealthData != nil {
		sum ^= uint32(obj.HealthData.Cur)
	}
	if obj.InitData != nil {
		typ := GetServer().S().Types.ByInd(int(obj.TypeInd))
		if typ != nil {
			if size := int32(typ.InitDataSize); size > 0 {
				sum ^= protection.Checksum(unsafe.Slice((*byte)(obj.InitData), int(size)))
			}
		}
	}
	if obj.IDPtr != nil {
		n := alloc.StrLen((*byte)(obj.IDPtr))
		if n != 0 {
			sum ^= protection.Checksum(unsafe.Slice((*byte)(obj.IDPtr), n))
		}
	}
	return sum
}

func toggleProtectionObject(id int32, obj *server.Object) uint32 {
	if id < 657757279 {
		return uint32(id)
	}
	key := uint32(dword_5d4594_2516348)
	r := protection.Find(protectionHead(), key, uint32(id))
	if r == nil {
		return uint32(id)
	}
	digest := objectProtectionChecksum(obj)
	r.Value ^= digest
	dword_5d4594_2516328 ^= uint32(digest)
	return uint32(dword_5d4594_2516328)
}
