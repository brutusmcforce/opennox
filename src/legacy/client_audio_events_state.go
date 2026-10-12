package legacy

import (
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/common/sound"
	"github.com/brutusmcforce/opennox/legacy/timer"
	"unsafe"
)

var (
	audioEventEnabled      = (*uint32)(unsafe.Pointer(&dword_5d4594_1045432))
	audioEventPlayback     = (*uint32)(unsafe.Pointer(&dword_587000_126996))
	audioEventCatalogWord  = (*uint32)(unsafe.Pointer(&dword_5d4594_1045420))
	audioEventCacheStorage uint32
	audioEventCacheWord    = &audioEventCacheStorage
	audioEventContextWord  = (*uint32)(unsafe.Pointer(&dword_5d4594_1045428))
	audioEventPoolStorage  uint32
	audioEventPoolWord     = &audioEventPoolStorage
	audioEventMusicCount   = (*uint32)(unsafe.Pointer(&dword_5d4594_816368))
	audioEventMusicLevel   = (*uint32)(unsafe.Pointer(&dword_5d4594_816372))
)

func audioEventCache() *audioStreamCache {
	return (*audioStreamCache)(unsafe.Pointer(uintptr(*audioEventCacheWord)))
}
func audioEventPool() *unsafe.Pointer {
	return (*unsafe.Pointer)(unsafe.Pointer(uintptr(*audioEventPoolWord)))
}
func audioEventContext() *audioStreamContext {
	return (*audioStreamContext)(unsafe.Pointer(uintptr(*audioEventContextWord)))
}
func audioEventGlobalTimers() *timer.TimerGroup {
	return (*timer.TimerGroup)(legacyGlobals.dword_587000_127004)
}
func audioEventRoot() *legacyListNode { return (*legacyListNode)(memmap.PtrOff(0x5D4594, 840612)) }
func audioEventContextTimers() *timer.TimerGroup {
	return (*timer.TimerGroup)(memmap.PtrOff(0x5D4594, 1045228))
}
func audioEventDefaults(m *audioEventMetadata) int32 {
	m.Enabled = 0
	m.Flags = 0
	m.Field8 = 0
	m.Limit = 0
	m.Loops = 0
	m.PitchMin = 0
	m.PitchMax = 0
	m.Priority = 1
	m.SampleCount = 0
	m.DelayMax = 0
	m.DelayMin = 0
	m.Frame = 0
	m.BucketGeneration = 0
	m.Distance = 600
	if m.Volume.Init(0x4000) {
		return 1
	}
	return 0
}
func audioEventInit(ctx *audioStreamContext, cat *audioStreamCatalog) int32 {
	for i := 0; i < 1023; i++ {
		m := (*audioEventMetadata)(memmap.PtrOff(0x5D4594, 840628+uintptr(i)*200))
		audioEventDefaults(m)
		m.Name = (*byte)(unsafe.Pointer(internCStr(sound.ID(i).String())))
	}
	*audioEventCatalogWord = uint32(uintptr(unsafe.Pointer(cat)))
	*audioEventContextWord = uint32(uintptr(unsafe.Pointer(ctx)))
	if cat != nil {
		*audioEventCacheWord = uint32(uintptr(unsafe.Pointer(audioStreamCacheNew(cat, 0x100000, 200, 0x2000))))
		*audioEventPoolWord = uint32(uintptr(unsafe.Pointer(audioStreamPoolNew(200, 576))))
	}
	if audioEventCache() == nil || cat == nil || ctx == nil || audioEventPool() == nil {
		return 0
	}
	listClear(audioEventRoot())
	audioEventContextTimers().Init()
	ctx.ParentTimers = audioEventContextTimers()
	*audioEventEnabled = 1
	return 1
}
func audioEventFree() {
	audioEventStopAll()
	audioEventReclaim()
	if p := audioEventCache(); p != nil {
		audioStreamCacheFree(p)
		*audioEventCacheWord = 0
	}
	if p := audioEventPool(); p != nil {
		audioStreamPoolFree(unsafe.Pointer(p))
		*audioEventPoolWord = 0
	}
	*audioEventEnabled = 0
}
func audioEventCallbacks(v *audioStreamVoice) {
	v.OnLoop = audioBridgeCallbackKey(13)
	v.OnEnd = audioBridgeCallbackKey(14)
	v.OnStop = audioBridgeCallbackKey(15)
}
func audioEventMusicDisable() { dword_587000_93156 = 0 }
func audioEventMusicEnable() int32 {
	v := dword_5d4594_816376
	if v != 0 {
		dword_587000_93156 = 1
	}
	return int32(v)
}
func audioEventMusicEnabled() int32 { return int32(dword_587000_93156) }
func audioEventDialogDisable()      { dword_587000_122848 = 0 }
func audioEventDialogEnable() int32 {
	v := dword_5d4594_831092
	if v != 0 {
		dword_587000_122848 = 1
	}
	return int32(v)
}
func audioEventDialogEnabled() int32 { return int32(dword_587000_122848) }
