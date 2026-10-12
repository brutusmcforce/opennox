package legacy

import (
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"unsafe"
)

var audioStreamsRoot *audioStreamSystem

func audioStreamRoot() *audioStreamSystem {
	if audioStreamsRoot == nil {
		audioStreamsRoot = (*audioStreamSystem)(memmap.PtrOff(0x5D4594, 1193204))
	}
	return audioStreamsRoot
}
func audioStreamDeviceNew(desc *audioStreamDescriptor) *audioStreamDevice {
	p, _ := alloc.New(audioStreamDevice{})
	listInit(&p.Node)
	p.References = 0
	p.Descriptor = desc
	if AudioStreamCallbackInt(desc.Init, unsafe.Pointer(p)) == 0 {
		return p
	}
	audioStreamDeviceFree(p)
	return nil
}
func audioStreamDeviceRegister(desc *audioStreamDescriptor) *audioStreamDevice {
	p := audioStreamDeviceNew(desc)
	if p == nil {
		return nil
	}
	desc.InUse |= 1
	audioStreamDeviceLink(p)
	if desc.Flags&2 != 0 {
		*memmap.PtrUint32(0x5D4594, 1193332) = 1
	}
	return p
}
func audioStreamDeviceFree(p *audioStreamDevice) {
	AudioStreamCallbackVoid(p.Descriptor.Free, unsafe.Pointer(p))
	p.Descriptor.InUse &^= 1
	alloc.Free(p)
}
func audioStreamDeviceLink(p *audioStreamDevice)   { listAppend(&audioStreamRoot().Devices, &p.Node) }
func audioStreamDeviceUnlink(p *audioStreamDevice) { listRemove(&p.Node) }
func audioStreamDeviceDestroy(p *audioStreamDevice) {
	audioStreamDeviceUnlink(p)
	audioStreamDeviceFree(p)
	*memmap.PtrUint32(0x5D4594, 1193332) = 0
}
func audioStreamDeviceFirst(it **audioStreamDevice) *audioStreamDevice {
	*it = (*audioStreamDevice)(unsafe.Pointer(listNext(&audioStreamRoot().Devices)))
	return *it
}
func audioStreamDeviceNext(it **audioStreamDevice) *audioStreamDevice {
	if *it != nil {
		*it = (*audioStreamDevice)(unsafe.Pointer(listNext(&(*it).Node)))
	}
	return *it
}
func audioStreamDeviceDestroyAll() {
	for n := listNext(&audioStreamRoot().Devices); n != nil; {
		next := listNext(n)
		audioStreamDeviceDestroy((*audioStreamDevice)(unsafe.Pointer(n)))
		n = next
	}
}
func audioStreamDeviceSlot(index int32, device **audioStreamDevice, slot *int32) unsafe.Pointer {
	var it *audioStreamDevice
	for p := audioStreamDeviceFirst(&it); p != nil; p = audioStreamDeviceNext(&it) {
		if index < p.Slots {
			*device = p
			*slot = index
			return unsafe.Pointer(slot)
		}
		index -= p.Slots
	}
	*device = nil
	*slot = -1
	return nil
}
func audioStreamContextAcquire(index int32, format *audioStreamFormat) *audioStreamContext {
	if index == -1 {
		index = 0
	}
	var device *audioStreamDevice
	var slot int32
	audioStreamDeviceSlot(index, &device, &slot)
	if device == nil {
		return nil
	}
	p := device.Contexts[slot]
	if p == nil {
		p = audioStreamContextNew(device, slot, format)
		if p == nil {
			return nil
		}
		p.Index = index
		audioStreamContextLink(p)
	}
	p.References++
	return p
}
func audioStreamContextNew(device *audioStreamDevice, slot int32, format *audioStreamFormat) *audioStreamContext {
	p, _ := alloc.New(audioStreamContext{})
	listInit(&p.Node)
	p.Slot = slot
	p.Device = device
	p.References = 0
	device.References++
	device.Contexts[slot] = p
	p.VoiceAPI = device.Descriptor.VoiceAPI
	listClear(&p.Voices)
	p.Timers.Init()
	p.Mutating = 0
	p.Period = 33
	p.MaxElapsed = 0
	p.Elapsed = 0
	p.LastTick = 0
	p.Tick = unsafe.Pointer(&audioStreamCallbackKeys[0])
	if format != nil {
		audioStreamContextFormat(p, format)
	}
	if AudioStreamCallbackInt(device.Descriptor.ContextInit, unsafe.Pointer(p)) == 0 {
		return p
	}
	audioStreamContextFree(p)
	return nil
}
func audioStreamContextFree(p *audioStreamContext) {
	audioStreamVoiceDeleteKind(p, -1)
	AudioStreamCallbackVoid(p.Device.Descriptor.ContextFree, unsafe.Pointer(p))
	p.Device.Contexts[p.Slot] = nil
	p.Device.References--
	if p.Device.References < 0 {
		p.Device.References = 0
	}
	alloc.Free(p)
}
func audioStreamContextLink(p *audioStreamContext) int32 {
	root := audioStreamRoot()
	root.Mutating++
	listAppend(&root.Contexts, &p.Node)
	result := root.Mutating - 1
	root.Mutating = result
	if result < 0 {
		root.Mutating = 0
	}
	return result
}
func audioStreamContextUnlink(p *audioStreamContext) uintptr {
	root := audioStreamRoot()
	root.Mutating++
	listRemove(&p.Node)
	result := root.Mutating - 1
	root.Mutating = result
	if result < 0 {
		root.Mutating = 0
		return uintptr(unsafe.Pointer(root))
	}
	return uintptr(uint32(result))
}
func audioStreamContextDestroy(p *audioStreamContext) {
	audioStreamContextUnlink(p)
	audioStreamContextFree(p)
}
func audioStreamContextFirst(it **audioStreamContext) *audioStreamContext {
	*it = (*audioStreamContext)(unsafe.Pointer(listNext(&audioStreamRoot().Contexts)))
	return *it
}
func audioStreamContextNext(it **audioStreamContext) *audioStreamContext {
	if *it != nil {
		*it = (*audioStreamContext)(unsafe.Pointer(listNext(&(*it).Node)))
	}
	return *it
}
func audioStreamContextDestroyAll() int32 {
	root := audioStreamRoot()
	root.Mutating++
	for n := listNext(&root.Contexts); n != nil; {
		next := listNext(n)
		audioStreamContextDestroy((*audioStreamContext)(unsafe.Pointer(n)))
		n = next
	}
	result := root.Mutating - 1
	root.Mutating = result
	if result < 0 {
		root.Mutating = 0
	}
	return result
}
func audioStreamContextFormat(p *audioStreamContext, format *audioStreamFormat) *audioStreamContext {
	p.RequestedFormat = *format
	return p
}
func audioStreamContextTick(p *audioStreamContext) int32 {
	if p.Mutating != 0 {
		return -2146304000
	}
	now := uint64(uint32(PlatformTicks()))
	delta := now - p.LastTick
	if delta < p.Period {
		return 0
	}
	p.Elapsed = delta
	if p.MaxElapsed > 10*p.Period {
		p.MaxElapsed = 0
	}
	if delta > p.MaxElapsed {
		p.MaxElapsed = delta
	}
	p.Timers.Update()
	if p.ParentTimers != nil {
		p.ParentTimers.Update()
	}
	changed := p.ParentTimers != nil && p.ParentTimers.IsUpdated() || p.Timers.IsUpdated()
	p.LastTick = now
	for node := p.Voices.next; node != &p.Voices; {
		next := node.next
		voice := (*audioStreamVoice)(unsafe.Pointer(node))
		if voice.Flags&1 != 0 && voice.Buffer != nil {
			voice.Timers.Update()
			if changed || voice.Timers.IsUpdated() || voice.GlobalTimers != nil && voice.GlobalTimers.IsUpdated() || voice.ExtraTimers != nil && voice.ExtraTimers.IsUpdated() {
				audioStreamVoiceMix(voice)
				AudioStreamCallbackVoid(voice.API.Update, unsafe.Pointer(voice))
			}
		}
		node = next
	}
	if p.ParentTimers != nil {
		p.ParentTimers.ClearUpdated()
	}
	p.Timers.ClearUpdated()
	return 0
}
