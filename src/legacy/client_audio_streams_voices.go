package legacy

import (
	"unsafe"

	"github.com/brutusmcforce/opennox/legacy/common/alloc"
)

func audioStreamVoiceNew(ctx *AudioStreamContext) *audioStreamVoice {
	p, _ := alloc.New(audioStreamVoice{})
	listInit(&p.Node)
	p.Flags = 0
	p.Loops = 0
	p.Effective.Init()
	audioStreamVoiceInit(p)
	p.Context = ctx
	p.API = ctx.VoiceAPI
	if AudioStreamCallbackInt(p.API.Init, unsafe.Pointer(p)) == 0 {
		return p
	}
	audioStreamVoiceFree(p)
	return nil
}
func audioStreamVoiceFree(p *audioStreamVoice) {
	AudioStreamCallbackVoid(p.API.Free, unsafe.Pointer(p))
	alloc.Free(p)
}
func audioStreamVoiceInit(p *audioStreamVoice) unsafe.Pointer {
	p.DataCallback = unsafe.Pointer(&audioStreamCallbackKeys[1])
	p.LoopCallback = unsafe.Pointer(&audioStreamCallbackKeys[2])
	p.EndCallback = unsafe.Pointer(&audioStreamCallbackKeys[3])
	p.OnData = nil
	p.OnLoop = nil
	p.OnEnd = nil
	p.Fields152[0] = 0
	p.Kind = 1
	p.Flags = 0
	p.Loops = 0
	p.Priority = 0
	p.GlobalTimers = AudioGlobalTimers
	p.ExtraTimers = nil
	p.Timers.Init()
	p.Buffer = nil
	return nil // Existing TimerGroup initializer returns null through its C adapter.
}
func audioStreamVoiceMix(p *audioStreamVoice) {
	p.Effective.Init()
	p.Effective.Mix(&p.Timers)
	p.Timers.ClearUpdated()
	if p.ExtraTimers != nil {
		p.Effective.Mix(p.ExtraTimers)
		p.ExtraTimers.ClearUpdated()
	}
	if p.GlobalTimers != nil {
		p.Effective.Mix(p.GlobalTimers)
	}
	p.Effective.Mix(&p.Context.Timers)
	if p.Context.ParentTimers != nil {
		p.Effective.Mix(p.Context.ParentTimers)
	}
}
func audioStreamVoiceData(p *audioStreamVoice) int32 {
	if p.OnData != nil {
		result := int32(AudioStreamCallbackInt(p.OnData, unsafe.Pointer(p)))
		if result != 0 {
			p.Remaining = 0
			p.Length = 0
			p.Data = 0
			return result
		}
	} else {
		if p.Chunk != nil {
			p.Chunk = (*audioStreamChunk)(unsafe.Pointer(listNext(&p.Chunk.Node)))
			if p.Chunk != nil {
				p.Data = p.Chunk.Data
				p.Remaining = p.Chunk.Length
				p.Length = p.Chunk.Length
				return 0
			}
		}
		p.Remaining = 0
	}
	return 0
}
func audioStreamVoiceLoop(p *audioStreamVoice) int32 {
	if p.Loops != 0 {
		if p.Loops != -1 {
			p.Loops--
		}
		audioStreamVoiceBind(p, p.Buffer)
	} else {
		audioStreamVoiceBind(p, nil)
	}
	if p.OnLoop != nil {
		AudioStreamCallbackVoid(p.OnLoop, unsafe.Pointer(p))
	}
	if p.Buffer != nil {
		AudioStreamCallbackVoid(p.API.Restart, unsafe.Pointer(p))
	}
	return 0
}
func audioStreamVoiceEnd(p *audioStreamVoice) int32 {
	p.Buffer = nil
	p.Flags &^= 5
	p.Loops = 0
	p.Timers.Init()
	if p.OnEnd != nil {
		return int32(AudioStreamCallbackInt(p.OnEnd, unsafe.Pointer(p)))
	}
	return 0
}
func audioStreamVoiceDestroy(p *audioStreamVoice) {
	audioStreamVoiceStop(p)
	audioStreamVoiceUnlink(p)
	audioStreamVoiceFree(p)
}
func audioStreamVoiceStop(p *audioStreamVoice) int32 {
	if p.Flags&5 != 0 {
		AudioStreamCallbackVoid(p.API.Stop, unsafe.Pointer(p))
	}
	result := int32(0)
	if p.OnStop != nil {
		result = int32(AudioStreamCallbackInt(p.OnStop, unsafe.Pointer(p)))
	}
	p.Buffer = nil
	return result
}
func audioStreamVoiceReserve(p *audioStreamVoice) *audioStreamVoice   { p.Flags |= 0x10; return p }
func audioStreamVoiceUnreserve(p *audioStreamVoice) *audioStreamVoice { p.Flags &^= 0x10; return p }
func audioStreamVoiceStart(p *audioStreamVoice) int32 {
	if p.Flags&5 != 0 {
		return -2146500608
	}
	if p.Buffer == nil {
		return -2147024896
	}
	p.Timers.Update()
	audioStreamVoiceMix(p)
	result := int32(AudioStreamCallbackInt(p.API.Start, unsafe.Pointer(p)))
	if result == 0 {
		p.Flags |= 1
	}
	return result
}
func audioStreamVoiceBind(p *audioStreamVoice, b *audioStreamBuffer) {
	p.Buffer = b
	if b == nil {
		return
	}
	p.Chunk = audioStreamBufferFirst(b)
	if p.Chunk != nil {
		p.Data = p.Chunk.Data
		p.Remaining = p.Chunk.Length
		p.Length = p.Chunk.Length
		b.Data = 0
	} else {
		p.Data = b.Data
		p.Remaining = b.Length
		p.Length = b.Length
	}
}
func audioStreamVoiceLink(ctx *AudioStreamContext, p *audioStreamVoice) int32 {
	p.Context = ctx
	ctx.VoiceCount++
	ctx.Mutating++
	listAppend(&ctx.Voices, &p.Node)
	result := ctx.Mutating - 1
	ctx.Mutating = result
	if result < 0 {
		ctx.Mutating = 0
	}
	return result
}
func audioStreamVoiceUnlink(p *audioStreamVoice) int32 {
	ctx := p.Context
	listRemove(&p.Node)
	ctx.VoiceCount--
	ctx.Mutating++
	listRemove(&p.Node)
	result := ctx.Mutating - 1
	ctx.Mutating = result
	if result < 0 {
		ctx.Mutating = 0
	}
	return result
}
func audioStreamVoiceCreate(ctx *AudioStreamContext) *audioStreamVoice {
	if ctx.VoiceCount >= ctx.VoiceCapacity {
		return nil
	}
	p := audioStreamVoiceNew(ctx)
	if p == nil {
		return nil
	}
	audioStreamVoiceLink(ctx, p)
	return p
}
func audioStreamVoiceCreateMany(ctx *AudioStreamContext, count int32) int32 {
	result := int32(0)
	for audioStreamVoiceCreate(ctx) != nil {
		result++
		count--
		if count == 0 {
			break
		}
	}
	return result
}
func audioStreamVoiceFirst(ctx *AudioStreamContext, it **audioStreamVoice) *audioStreamVoice {
	*it = (*audioStreamVoice)(unsafe.Pointer(listNext(&ctx.Voices)))
	return *it
}
func audioStreamVoiceNext(it **audioStreamVoice) *audioStreamVoice {
	if *it != nil {
		*it = (*audioStreamVoice)(unsafe.Pointer(listNext(&(*it).Node)))
	}
	return *it
}
func audioStreamVoiceSelect(ctx *AudioStreamContext, kind int32) *audioStreamVoice {
	if kind == -1 {
		kind = 1
	}
	activePriority, busyPriority := int32(127), int32(127)
	level := uint32(0xffffffff)
	var active, busy *audioStreamVoice
	for n := listNext(&ctx.Voices); n != nil; n = listNext(n) {
		p := (*audioStreamVoice)(unsafe.Pointer(n))
		if p.Kind != kind {
			continue
		}
		if p.Flags&0x15 == 0 {
			return p
		}
		if p.Flags&1 != 0 {
			v := p.Effective.Timers[0].Current
			if p.Priority < activePriority || p.Priority == activePriority && v < level && level-v >= 0x666 {
				active, activePriority, level = p, p.Priority, v
			}
		} else if p.Priority < busyPriority {
			busy, busyPriority = p, p.Priority
		}
	}
	if busy != nil && busyPriority <= activePriority {
		return busy
	}
	return active
}
func audioStreamVoiceDeleteKind(ctx *AudioStreamContext, kind int32) int32 {
	for n := listNext(&ctx.Voices); n != nil; {
		next := listNext(n)
		p := (*audioStreamVoice)(unsafe.Pointer(n))
		if kind == -1 || p.Kind == kind {
			audioStreamVoiceDestroy(p)
		}
		n = next
	}
	return 0
}
func audioStreamVoiceStopKind(ctx *AudioStreamContext, kind int32) uintptr {
	var result uintptr
	for n := listNext(&ctx.Voices); n != nil; {
		next := listNext(n)
		result = uintptr(unsafe.Pointer(next))
		p := (*audioStreamVoice)(unsafe.Pointer(n))
		if kind == -1 || p.Kind == kind {
			result = uintptr(uint32(audioStreamVoiceStop(p)))
		}
		n = next
	}
	return result
}
