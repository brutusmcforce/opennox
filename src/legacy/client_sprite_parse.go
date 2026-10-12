package legacy

import (
	"github.com/brutusmcforce/opennox/client"
	"github.com/brutusmcforce/opennox/client/noxrender"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/internal/binfile"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"unsafe"
)

// Draw-data buffers share the legacy allocator domain with spriteDataFree.
// Preserve the existing zero-count allocation behavior.
func spriteDataAlloc(count int, size uintptr) unsafe.Pointer {
	return legacyCalloc(uintptr(count), uintptr(size))
}
func spriteReadKind(f *binfile.MemFile, scratch []byte) client.AnimKind {
	n := int(f.ReadU8())
	f.Read(scratch[:n])
	scratch[n] = 0
	return client.ParseAnimKind(alloc.GoString(&scratch[0]))
}
func spriteReadImage(f *binfile.MemFile, scratch []byte, empty uintptr) noxrender.ImageHandle {
	id := int(f.ReadI32())
	scratch[0] = *memmap.PtrUint8(0x5D4594, empty)
	var typ byte
	var name string
	if id == -1 {
		typ = f.ReadU8()
		n := int(f.ReadU8())
		f.Read(scratch[:n])
		scratch[n] = 0
		name = alloc.GoString(&scratch[0])
	}
	return GetClient().R2().GetBag().ImageRef(id, typ, name).C()
}
func spriteReadFrames(f *binfile.MemFile, scratch []byte, empty uintptr, count int) *noxrender.ImageHandle {
	frames := (*noxrender.ImageHandle)(spriteDataAlloc(count, 4))
	if frames == nil {
		return nil
	}
	for i := 0; i < count; i++ {
		*(*noxrender.ImageHandle)(unsafe.Add(unsafe.Pointer(frames), 4*i)) = spriteReadImage(f, scratch, empty)
	}
	return frames
}
func spriteParseAnimate(obj *client.ObjectType, f *binfile.MemFile, scratch []byte) bool {
	data := (*spriteAnimationData)(spriteDataAlloc(1, 16))
	data.Size = 16
	data.Count = f.ReadU8()
	data.Delay = f.ReadU8()
	data.Kind = spriteReadKind(f, scratch)
	data.Frames = spriteReadFrames(f, scratch, 830832, int(data.Count))
	if data.Frames == nil {
		return false
	}
	obj.DrawData = unsafe.Pointer(data)
	obj.DrawFunc = drawableDrawKey(drawKey_nox_thing_animate_draw)
	return true
}
func spriteParseConditional(obj *client.ObjectType, f *binfile.MemFile, scratch []byte) bool {
	data := (*spriteConditionalData)(spriteDataAlloc(1, 56))
	data.Size = 16 // Historical header value is not the allocation size.
	states := int(f.ReadU8())
	for i := 0; i < states; i++ {
		data.Count[i] = f.ReadU8()
		data.Delay[i] = f.ReadU8()
		data.Kind[i] = spriteReadKind(f, scratch)
		data.Frames[i] = spriteReadFrames(f, scratch, 830836, int(data.Count[i]))
		if data.Frames[i] == nil {
			return false
		}
	}
	obj.DrawData = unsafe.Pointer(data)
	obj.DrawFunc = drawableDrawKey(drawKey_nox_thing_cond_animate_draw)
	obj.Field_60 = 0
	return true
}
func spriteStaticRandomData(f *binfile.MemFile, scratch []byte) unsafe.Pointer {
	data := spriteDataAlloc(1, 12)
	*(*uint32)(data) = 12
	count := f.ReadU8()
	*(*byte)(unsafe.Add(data, 8)) = count
	frames := spriteReadFrames(f, scratch, 830852, int(count))
	*(**noxrender.ImageHandle)(unsafe.Add(data, 4)) = frames
	if frames == nil {
		return nil
	}
	return data
}
func spriteParseStatic(obj *client.ObjectType, f *binfile.MemFile, scratch []byte) bool {
	data := spriteDataAlloc(1, 8)
	if data == nil {
		return false
	}
	*(*uint32)(data) = 8
	*(*noxrender.ImageHandle)(unsafe.Add(data, 4)) = spriteReadImage(f, scratch, 830856)
	obj.DrawData = data
	obj.DrawFunc = drawableDrawKey(drawKey_nox_thing_static_draw)
	return true
}
func spriteParseRandom(obj *client.ObjectType, f *binfile.MemFile, scratch []byte, slave bool) bool {
	if slave {
		obj.DrawFunc = drawableDrawKey(drawKey_nox_thing_slave_draw)
	} else {
		obj.DrawFunc = drawableDrawKey(drawKey_nox_thing_static_random_draw)
	}
	obj.DrawData = spriteStaticRandomData(f, scratch)
	if slave {
		obj.Field_60 = 0
	}
	return obj.DrawData != nil
}
func spriteVectorHeader(vector *client.AnimationVector, f *binfile.MemFile) int {
	vector.Cnt40 = uint16(f.ReadU8())
	vector.Val42 = uint16(f.ReadU8())
	var scratch [256]byte
	vector.Kind = spriteReadKind(f, scratch[:])
	return 1
}
func spriteVectorFrames(vector *client.AnimationVector, f *binfile.MemFile) int {
	var scratch [256]byte
	vector.Frames[0] = spriteReadFrames(f, scratch[:], 830848, int(int16(vector.Cnt40)))
	if vector.Frames[0] == nil {
		return 0
	}
	return 1
}
func spriteParseState(obj *client.ObjectType, f *binfile.MemFile) bool {
	data := spriteDataAlloc(1, 148)
	*(*uint32)(data) = 148
	for {
		cmd := f.ReadU32()
		if cmd == 0x454e4420 {
			break
		}
		params := f.ReadU32()
		if params&14 == 0 {
			return false
		}
		f.Skip(int(f.ReadU8()))
		f.Skip(int(f.ReadU8()))
		state := 0
		if params&2 != 0 {
			state = 0
		} else if params&4 != 0 {
			state = 1
		} else {
			state = 2
		}
		vector := (*client.AnimationVector)(unsafe.Add(data, 4+state*48))
		spriteVectorHeader(vector, f)
		if spriteVectorFrames(vector, f) == 0 {
			return false
		}
	}
	obj.Field_54 = 2
	obj.DrawFunc = drawableDrawKey(drawKey_nox_thing_animate_state_draw)
	obj.DrawData = data
	return true
}
