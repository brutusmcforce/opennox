package legacy

import (
	"github.com/brutusmcforce/opennox/common/memmap"
	"unsafe"
)

func audioEventNew(m *audioEventMetadata) *audioEvent {
	if *audioEventEnabled == 0 || *audioEventPlayback == 0 || m.Enabled == 0 {
		return nil
	}
	p := (*audioEvent)(audioStreamPoolPop(audioEventPool()))
	if p == nil {
		audioEventReclaim()
		p = (*audioEvent)(audioStreamPoolPop(audioEventPool()))
		if p == nil {
			return nil
		}
	}
	*p = audioEvent{Metadata: m}
	listInit(&p.Node)
	p.Timers.Init()
	listAppend(audioEventRoot(), &p.Node)
	serial := memmap.PtrUint32(0x587000, 127000)
	p.Serial = *serial
	*serial++
	return p
}
func audioEventDelete(p *audioEvent) uint32 {
	if p.Flags&1 != 0 {
		return 0
	}
	audioEventDetach(p)
	audioEventReleaseSamples(p)
	p.State = 4
	p.Serial = 0
	p.Flags |= 1
	return p.Flags
}
func audioEventUnlink(p *audioEvent) uintptr {
	listRemove(&p.Node)
	p.Serial = 0
	return uintptr(audioStreamPoolPush(audioEventPool(), unsafe.Pointer(p)))
}
func audioEventStopAll() uintptr {
	result := uintptr(*audioEventEnabled)
	if *audioEventEnabled != 0 {
		root := audioEventRoot()
		for n := root.next; n != root; {
			next := n.next
			p := (*audioEvent)(unsafe.Pointer(n))
			audioEventDelete(p)
			result = audioEventUnlink(p)
			n = next
		}
	}
	return result
}
func audioEventReclaim() uintptr {
	result := uintptr(*audioEventEnabled)
	if *audioEventEnabled != 0 {
		root := audioEventRoot()
		result = uintptr(unsafe.Pointer(root.next))
		for n := root.next; n != root; {
			next := n.next
			p := (*audioEvent)(unsafe.Pointer(n))
			if p.Flags&1 != 0 {
				audioEventUnlink(p)
			}
			n = next
			result = uintptr(unsafe.Pointer(n))
		}
	}
	return result
}
func audioEventVoiceOwner(v *audioStreamVoice) *audioEvent {
	return (*audioEvent)(memmap.WordPtr(v.Fields152[0]))
}
func audioEventDetach(p *audioEvent) uintptr {
	v := p.Voice
	if v != nil && audioEventVoiceOwner(v) == p {
		if p.Flags&2 != 0 {
			audioStreamVoiceStop(v)
		}
		audioStreamVoiceUnreserve(v)
		v.Fields152[0] = 0
		v.OnStop = nil
		v.OnLoop = nil
		v.OnEnd = nil
		v.ExtraTimers = nil
		p.Voice = nil
	}
	return uintptr(unsafe.Pointer(v))
}
func audioEventStart(p *audioEvent) int32 {
	v := p.Voice
	if audioEventVoiceOwner(v) != p {
		return 0
	}
	pending := p.Pending
	audioStreamVoiceBind(v, pending)
	p.State = 3
	p.Flags |= 2
	p.Pending = nil
	if audioStreamVoiceStart(v) == 0 {
		return 1
	}
	p.State = 1
	p.Pending = pending
	p.Flags &^= 2
	return 0
}
func audioEventStep(p *audioEvent) {
	if *audioEventPlayback == 0 {
		p.State = 4
	}
	for {
		switch p.State {
		case 0:
			if audioEventReserve(p) == 0 {
				audioEventDelete(p)
			}
			return
		case 2:
			if uint64(uint32(PlatformTicks())) <= p.Deadline {
				return
			}
			p.State = p.NextState
		case 4:
			audioEventDelete(p)
			return
		default:
			return
		}
	}
}
func audioEventSchedule(p *audioEvent, delay uint64, next uint32) uint64 {
	p.NextState = next
	p.Deadline = delay + uint64(uint32(PlatformTicks()))
	p.State = 2
	return p.Deadline
}
func audioEventVoiceStopped(v *audioStreamVoice) int32 { audioEventVoiceOwner(v).State = 4; return 0 }
func audioEventVoiceEnded(v *audioStreamVoice) int32 {
	p := audioEventVoiceOwner(v)
	p.Flags &^= 2
	if p.State != 4 {
		next := uint32(4)
		if p.Pending != nil || p.ReloadRemaining != 0 {
			next = 1
		} else {
			p.Delay = 0
		}
		if p.Delay != 0 {
			audioEventSchedule(p, uint64(p.Delay), next)
			p.Delay = 0
			return 0
		}
		p.State = next
	}
	return 0
}
func audioEventChooseVoice(priority int32, flags int8) *audioStreamVoice {
	ctx := audioEventContext()
	if ctx == nil {
		return nil
	}
	v := audioStreamVoiceSelect(ctx, 1)
	if v == nil {
		return nil
	}
	if v.Flags&0x15 != 0 && v.Priority > priority {
		return nil
	}
	audioStreamVoiceStop(v)
	v.GlobalTimers = audioEventGlobalTimers()
	v.Priority = priority
	v.Loops = 0
	if flags&1 != 0 {
		v.Loops = -1
	}
	v.Timers.Timers[0].SetRaw(0x4000)
	return v
}
func audioEventReserve(p *audioEvent) int32 {
	m := p.Metadata
	if m.SampleCount == 0 {
		return 0
	}
	p.Iteration = 0
	p.Voice = audioEventChooseVoice(m.Priority+p.PriorityOffset, 0)
	v := p.Voice
	if v == nil {
		return 0
	}
	v.Timers.Timers[1].SetRaw(uint32(audioEventRandom(m.PitchMin, m.PitchMax) + 100))
	audioStreamVoiceReserve(v)
	v.Fields152[0] = memmap.PtrWord(unsafe.Pointer(p))
	audioEventCallbacks(v)
	p.State = 1
	v.ExtraTimers = &p.Timers
	if m.Flags&8 != 0 {
		delay := audioEventRandom(int32(m.DelayMin), int32(m.DelayMax))
		if delay > 33 {
			audioEventSchedule(p, uint64(int64(delay)), 1)
		}
	}
	return 1
}
func audioEventVoiceLoop(v *audioStreamVoice) int32 {
	p := audioEventVoiceOwner(v)
	next := audioEventSelectNext(p)
	m := p.Metadata
	if m.DelayMax < 33 {
		audioStreamVoiceBind(v, next)
		return 0
	}
	audioStreamVoiceBind(v, nil)
	if m.Flags&8 == 0 || next != nil || p.ReloadRemaining != 0 {
		delay := uint32(audioEventRandom(int32(m.DelayMin), int32(m.DelayMax)))
		if delay < 33 {
			audioStreamVoiceBind(v, next)
			return 0
		}
		p.Delay = delay
		p.Pending = next
	}
	return 0
}

var audioEventPlayObserver func(int, int)

func audioEventPlay(id, volume, pan int32, mode int) {
	if mode == 0 && audioEventPlayObserver != nil {
		audioEventPlayObserver(int(id), int(volume))
	}
	m := audioAssetSlot(id)
	if m == nil {
		return
	}
	p := audioEventNew(m)
	if p == nil {
		return
	}
	audioEventSetVolume(p, volume)
	if mode != 0 {
		audioEventSetPan(p, pan)
	}
	if mode == 2 {
		p.PriorityOffset = 2
	}
	audioEventStep(p)
}
