package legacy

import "unsafe"

func audioEventRandom(min, max int32) int32 {
	return int32(GetServer().S().Rand.Other.Int(int(min), int(max)))
}
func audioEventSelectFirst(p *audioEvent) *audioStreamBuffer {
	p.Remaining = p.Loaded
	if p.Loaded == 0 {
		return nil
	}
	for i := int32(0); i < p.Remaining; i++ {
		p.Order[i] = i
	}
	p.LastSample = -1
	return audioEventSelectNext(p)
}
func audioEventSelectNext(p *audioEvent) *audioStreamBuffer {
	m := p.Metadata
	if p.Remaining != 0 {
		if m.Flags&2 != 0 {
			i := audioEventRandom(0, p.Remaining-1)
			p.LastSample = p.Order[i]
			for ; i < p.Remaining-1; i++ {
				p.Order[i] = p.Order[i+1]
			}
		} else {
			p.LastSample++
		}
		p.Remaining--
		return audioStreamCacheBuffer(p.Samples[p.LastSample])
	}
	if m.Flags&1 != 0 {
		if m.Loops != 0 {
			p.Iteration++
			if p.Iteration >= int32(m.Loops) {
				return nil
			}
		}
		return audioEventSelectFirst(p)
	}
	return nil
}
func audioEventLoadSample(p *audioEvent, index int32) int32 {
	e := audioStreamCacheLoad(audioEventCache(), int32(p.Metadata.Samples[index]))
	p.Samples[p.Loaded] = e
	if e == nil {
		return 0
	}
	audioStreamCacheRef(e)
	p.Loaded++
	return p.Loaded
}
func audioEventReleaseSamples(p *audioEvent) int32 {
	n := p.Loaded
	for i := int32(0); i < n; i++ {
		audioStreamCacheUnref(p.Samples[i])
		p.Samples[i] = nil
	}
	p.Loaded = 0
	return n
}
func audioEventLoad(p *audioEvent) int32 {
	m := p.Metadata
	if p.Loaded != 0 {
		if m.DelayMin < 33 {
			return p.Loaded
		}
		audioEventReleaseSamples(p)
	}
	if m.Flags&4 != 0 {
		if m.DelayMin >= 33 {
			return audioEventLoadSample(p, audioEventReloadIndex(p))
		}
		for i := uint32(0); i < m.SampleCount; i++ {
			audioEventLoadSample(p, int32(i))
		}
		return int32(m.SampleCount)
	}
	if m.Flags&2 != 0 {
		return audioEventLoadSample(p, audioEventRandom(0, int32(m.SampleCount)-1))
	}
	return audioEventLoadSample(p, 0)
}
func audioEventReloadIndex(p *audioEvent) int32 {
	if p.ReloadRemaining <= 0 {
		p.ReloadRemaining = int32(p.Metadata.SampleCount)
		for i := int32(0); i < p.ReloadRemaining; i++ {
			p.ReloadOrder[i] = p.ReloadRemaining - i - 1
		}
	}
	p.ReloadRemaining--
	if p.Metadata.Flags&2 == 0 {
		return p.ReloadOrder[p.ReloadRemaining]
	}
	i := audioEventRandom(0, p.ReloadRemaining)
	out := p.ReloadOrder[i]
	for ; i < p.ReloadRemaining; i++ {
		p.ReloadOrder[i] = p.ReloadOrder[i+1]
	}
	return out
}
func audioEventVolume(p *audioEvent, percent int32) uint32 {
	if percent < 0 {
		percent = 0
	}
	if percent > 100 {
		percent = 100
	}
	return (163 * uint32(percent) * (p.Metadata.Volume.Current >> 16)) >> 14
}
func audioEventSetVolume(p *audioEvent, percent int32) int32 {
	p.Timers.Timers[0].SetRaw(audioEventVolume(p, percent))
	if p.Timers.Timers[0].Update() {
		return 1
	}
	return 0
}
func audioEventFadeVolume(p *audioEvent, percent int32) uintptr {
	p.Timers.Timers[0].SetInterp(audioEventVolume(p, percent))
	return 0
}
func audioEventPan(pan int32) int32 {
	if pan < -50 {
		pan = -50
	}
	if pan > 50 {
		pan = 50
	}
	return pan*8192/50 + 8192
}
func audioEventSetPan(p *audioEvent, pan int32) unsafe.Pointer {
	t := &p.Timers.Timers[2]
	t.SetRaw(uint32(audioEventPan(pan)))
	return nil
}
func audioEventFadePan(p *audioEvent, pan int32) uintptr {
	t := &p.Timers.Timers[2]
	t.SetInterp(uint32(audioEventPan(pan)))
	return 0
}
func audioEventHandleSet(h *audioEventHandle, p *audioEvent) uintptr {
	h.Event = p
	if p == nil {
		return 0
	}
	h.Serial = p.Serial
	h.Metadata = p.Metadata
	return uintptr(unsafe.Pointer(p.Metadata))
}
func audioEventHandleGet(h *audioEventHandle) *audioEvent {
	p := h.Event
	if p != nil && (h.Serial != p.Serial || h.Metadata != p.Metadata) {
		h.Event = nil
		return nil
	}
	return p
}
