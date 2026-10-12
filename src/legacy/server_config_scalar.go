package legacy

import (
	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"unsafe"
)

func serverConfigFlagsSet(v int32) int32 {
	if uint32(dword_5d4594_3484) != uint32(v) {
		dword_5d4594_3484 = uint32(v)
		serverConfigUpdatedSet()
	}
	return v
}
func serverConfigFlagsGet() int32 { return int32(dword_5d4594_3484) }
func serverConfigFlagsAdd(v int32) int32 {
	result := v & serverConfigFlagsGet()
	if result != v {
		dword_5d4594_3484 |= uint32(v)
		result = int32(bool2int(noxflags.HasGame(1)))
		if result != 0 && v&0x2000 != 0 {
			result = int32(gameplayReportResetAll())
		}
		serverConfigUpdatedSet()
	}
	return result
}
func serverConfigFlagsRemove(v int32) int32 {
	result := v
	if serverConfigFlagsGet()&v != 0 {
		result = ^v
		dword_5d4594_3484 &= uint32(^v)
		serverConfigUpdatedSet()
	}
	return result
}
func serverConfigFlagsToggle(v int32) int32 {
	dword_5d4594_3484 ^= uint32(v)
	result := int32(bool2int(noxflags.HasGame(1)))
	if result != 0 && v&0x2000 != 0 {
		result = serverConfigFlagsGet()
		if result&0x2000 != 0 {
			result = int32(gameplayReportResetAll())
		}
	}
	serverConfigUpdatedSet()
	return result
}
func serverConfigFlagsQuery(v int32) int32 {
	return int32(bool2int(v == 0x2000 && noxflags.HasGame(1056) || serverConfigFlagsGet()&v != 0))
}
func serverConfigLimitSet(v int32) int32 {
	p := memmap.PtrInt32(0x5D4594, 3464)
	if *p != v {
		*p = v
		serverConfigUpdatedSet()
	}
	return v
}
func serverConfigLimitGet() int32 { return memmap.Int32(0x5D4594, 3464) }
func serverConfigScore(v int16) int16 {
	return memmap.Int16(0x5D4594, 3488+2*uintptr(Sub_409A70(int(v))))
}
func serverConfigMinutes(v int16) byte {
	return memmap.Uint8(0x5D4594, 3500+uintptr(Sub_409A70(int(v))))
}
func serverConfigTimerSet(v int32) int32 {
	*memmap.PtrInt32(0x587000, 4660) = v
	if noxflags.HasGame(1) {
		return int32(gameplayReportTimer(159, uint32(v)))
	}
	return 0
}
func serverConfigTimerGet() int32 { return memmap.Int32(0x587000, 4660) }
func serverConfigTimerLeft() int32 {
	return int32(memmap.Uint32(0x5D4594, 3468) - uint32(PlatformTicks()))
}
func serverConfigTimerInit() int64 {
	index := Sub_409A70(int(int16(noxflags.GetGame())))
	ticks := uint64(uint32(PlatformTicks()))
	duration := int64(memmap.Uint8(0x5D4594, 3500+uintptr(index))) * 60000
	*memmap.PtrUint64(0x5D4594, 3468) = uint64(duration) + ticks
	return duration
}
func serverConfigTimerReset(v int32) int64 {
	*memmap.PtrUint64(0x5D4594, 3468) = uint64(int64(v) + int64(uint32(PlatformTicks())))
	return int64(v)
}
func serverConfig3512Set(v int32) { *memmap.PtrInt32(0x5D4594, 3512) = v }
func serverConfig3512Get() int32  { return memmap.Int32(0x5D4594, 3512) }
func serverConfigModeStore(v uint32) uint32 {
	if v > 0x800 && v < 0x8000 {
		*memmap.PtrUint32(0x587000, 4652) = v
	}
	return v
}
func serverConfigNameSet(p *byte) *byte {
	dst := memmap.PtrUint8(0x5D4594, 1324)
	if p == nil {
		*dst = 0
		serverConfigUpdatedSet()
		return nil
	}
	var value [16]byte
	src := alloc.GoString(p)
	if len(src) > 15 {
		src = src[:15]
	}
	copy(value[:], src)
	// Match the game's byte-wise case folding, not Unicode case equivalence.
	if serverConfigEqualFoldBytes(alloc.GoString(dst), alloc.GoString(&value[0])) {
		return nil
	}
	copy(unsafe.Slice(dst, 16), value[:])
	serverConfigUpdatedSet()
	return dst
}
func serverConfigNameGet() *byte           { return memmap.PtrUint8(0x5D4594, 1324) }
func serverConfigRespawnSet(v int32) int32 { *memmap.PtrInt32(0x5D4594, 3584) = v; return v }
func serverConfigRespawnGet() int32 {
	if noxflags.HasGame(4096) {
		return 0
	}
	return memmap.Int32(0x5D4594, 3584)
}
func serverConfigPasswordSet(p *uint16) *uint16 {
	dst := memmap.PtrUint16(0x5D4594, 3540)
	for i := uintptr(0); ; i += 2 {
		c := *(*uint16)(unsafe.Add(unsafe.Pointer(p), i))
		*(*uint16)(unsafe.Add(unsafe.Pointer(dst), i)) = c
		if c == 0 {
			break
		}
	}
	return dst
}
func serverConfigPasswordGet() *uint16       { return memmap.PtrUint16(0x5D4594, 3540) }
func serverConfigUpdatedSet()                { legacyGlobals.nox_server_gameSettingsUpdated = 1 }
func serverConfigUpdatedGet() int32          { return int32(legacyGlobals.nox_server_gameSettingsUpdated) }
func serverConfigUpdatedClear()              { legacyGlobals.nox_server_gameSettingsUpdated = 0 }
func serverConfigRateDirtySet(v int32) int32 { *memmap.PtrInt32(0x5D4594, 3588) = v; return v }
func serverConfigRateDirtyGet() int32        { return memmap.Int32(0x5D4594, 3588) }
func serverConfigRateGet() int32             { return memmap.Int32(0x587000, 4728) }
func serverConfigRateSet(v int32) int32 {
	result := serverConfigRateGet()
	if v != result {
		result = int32(bool2int(noxflags.HasGame(0x20000)))
		if result == 1 {
			result = int32(gameplayReportRate(159))
		}
	}
	*memmap.PtrInt32(0x587000, 4728) = v
	return result
}
func serverConfigConnectionRate(v int32) int32 {
	for i := uintptr(0); i < 5; i++ {
		if memmap.Int32(0x587000, 4664+8*i) == v {
			return memmap.Int32(0x587000, 4668+8*i)
		}
	}
	return 1
}
func serverConfigSpecialMode() int32 {
	return int32(bool2int(noxflags.HasGame(128) && *(*int8)(unsafe.Add(unsafe.Pointer(serverConfigSlotCurrent()), 53)) < 0))
}
func serverConfigSlotIndex() int32 { return memmap.Int32(0x5D4594, 371688) }
func serverConfigSlot(index int32) *byte {
	return (*byte)(memmap.PtrOff(0x5D4594, 371380+uintptr(uint32(index)*58)))
}
func serverConfigSlotCurrent() *byte { return serverConfigSlot(serverConfigSlotIndex()) }
func serverConfigSlotSelect(index int32) *byte {
	*memmap.PtrInt32(0x5D4594, 371688) = index
	return serverConfigSlot(index)
}
func serverConfigSlotCopy(from, to int32) int32 {
	copy(unsafe.Slice(serverConfigSlot(to), 58), unsafe.Slice(serverConfigSlot(from), 58))
	return to
}
func serverConfigAdmissionData() *byte      { return (*byte)(memmap.PtrOff(0x5D4594, 371616)) }
func serverConfigSettings() unsafe.Pointer  { return memmap.PtrOff(0x5D4594, 371516) }
func serverConfigRecordState() int32        { return memmap.Int32(0x5D4594, 371700) }
func serverConfigAcquiredGet() int32        { return memmap.Int32(0x5D4594, 371704) }
func serverConfigAcquiredSet(v int32) int32 { *memmap.PtrInt32(0x5D4594, 371704) = v; return v }
func serverConfigOpenAdmission() *byte {
	p := serverConfigSettings()
	*(*byte)(unsafe.Add(p, 100)) &= 0xef
	return (*byte)(p)
}

func serverConfigEqualFoldBytes(a, b string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		x, y := a[i], b[i]
		if x >= 'A' && x <= 'Z' {
			x += 'a' - 'A'
		}
		if y >= 'A' && y <= 'Z' {
			y += 'a' - 'A'
		}
		if x != y {
			return false
		}
	}
	return true
}
