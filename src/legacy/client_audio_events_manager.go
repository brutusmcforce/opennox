package legacy

import (
	"github.com/brutusmcforce/opennox/common/memmap"
	"unsafe"
)

func audioEventFromRank(n *legacyListNode) *audioEvent {
	return (*audioEvent)(unsafe.Add(unsafe.Pointer(n), -12))
}
func audioEventRank(p *audioEvent) int32 {
	m := p.Metadata
	volume := p.Timers.Timers[0].Current >> 16
	at := m.Ranking.next
	for at != &m.Ranking {
		other := audioEventFromRank(at)
		ov := other.Timers.Timers[0].Current >> 16
		diff := int32(ov - volume)
		if diff < 0 {
			diff = int32(volume - ov)
		}
		if uint32(diff) >= (m.Volume.Current>>16)/10 {
			if ov < volume {
				break
			}
		} else if m.Flags&16 != 0 {
			if other.State != 0 {
				break
			}
		} else if other.State == 0 {
			break
		}
		at = at.next
	}
	listInit(&p.Rank)
	listAppend(at, &p.Rank)
	m.Active++
	result := m.Limit
	if result != 0 && int32(m.Active) > int32(result) {
		last := m.Ranking.prev
		listRemove(last)
		audioEventDelete(audioEventFromRank(last))
		m.Active--
		result = m.Active
	}
	return int32(result)
}
func audioEventBucket(priority int32, volume uint32) *legacyListNode {
	return &audioEventBuckets[priority][volume]
}
func audioEventBucketsClear() uint32 {
	for i := int32(0); i < 6; i++ {
		for j := uint32(0); j < 10; j++ {
			listClear(audioEventBucket(i, j))
		}
	}
	g := memmap.PtrUint32(0x5D4594, 1045444)
	*g++
	return *g
}
func audioEventBucketAdd(p *audioEvent) {
	m := p.Metadata
	priority := m.Priority + p.PriorityOffset
	volume := (p.Timers.Timers[0].Current >> 16) / 0x666
	generation := *memmap.PtrUint32(0x5D4594, 1045444)
	if m.BucketGeneration == generation {
		if priority < m.BucketPriority || priority == m.BucketPriority && volume <= m.BucketVolume {
			return
		}
		m.BucketPriority = priority
		m.BucketVolume = volume
		listRemove(&m.BucketNode)
	} else {
		m.BucketGeneration = generation
		m.BucketPriority = priority
		m.BucketVolume = volume
		listInit(&m.BucketNode)
	}
	listAppend(audioEventBucket(priority, volume), &m.BucketNode)
}
func audioEventBucketRemove(m *audioEventMetadata) { listRemove(&m.BucketNode) }
func audioEventBucketFirst(limit int32) *audioEventMetadata {
	for i := int32(0); i < limit; i++ {
		for j := uint32(0); j < 10; j++ {
			if n := listNext(audioEventBucket(i, j)); n != nil {
				return (*audioEventMetadata)(unsafe.Add(unsafe.Pointer(n), -112))
			}
		}
	}
	return nil
}
func audioEventEvict(p *audioEvent) uintptr {
	m := audioEventBucketFirst(p.Metadata.Priority + p.PriorityOffset)
	if m == nil {
		return 0
	}
	audioEventBucketRemove(m)
	result := uintptr(0)
	root := audioEventRoot()
	for n := root.next; n != root; {
		next := n.next
		e := (*audioEvent)(unsafe.Pointer(n))
		if e.Metadata == m {
			audioEventDelete(e)
			audioEventUnlink(e)
			result = 1
		}
		n = next
	}
	return result
}
func audioEventUpdate() {
	busy := memmap.PtrUint32(0x5D4594, 1045448)
	if *audioEventEnabled == 0 || *busy != 0 {
		return
	}
	*busy = 1
	audioEventGlobalTimers().Update()
	frame := memmap.PtrUint32(0x5D4594, 1045440)
	*frame++
	root := audioEventRoot()
	for n := root.next; n != root; n = n.next {
		p := (*audioEvent)(unsafe.Pointer(n))
		m := p.Metadata
		if m.Frame != *frame {
			listClear(&m.Ranking)
			m.Active = 0
			m.Frame = *frame
		}
		p.Timers.Update()
		if p.State != 4 {
			audioEventRank(p)
		}
	}
	for n := root.next; n != root; n = n.next {
		audioEventStep((*audioEvent)(unsafe.Pointer(n)))
	}
	audioEventBucketsClear()
	total := int32(0)
	for n := root.next; n != root; {
		next := n.next
		p := (*audioEvent)(unsafe.Pointer(n))
		v := p.Voice
		if v == nil || audioEventVoiceOwner(v) != p {
			audioEventDelete(p)
		}
		if p.Flags&1 != 0 {
			audioEventUnlink(p)
		} else {
			total += int32((33 * (p.Metadata.Volume.Current >> 16)) >> 14)
			audioEventBucketAdd(p)
		}
		n = next
	}
	gain := uint32(0x4000)
	if total > 100 {
		gain = 0x190000 / uint32(total)
	}
	audioEventContextTimers().Timers[0].SetInterp(gain)
	audioEventContextTimers().Update()
	for n := root.next; n != root; {
		next := n.next
		p := (*audioEvent)(unsafe.Pointer(n))
		if p.State == 1 {
			audioEventLoad(p)
			p.Pending = audioEventSelectFirst(p)
			if p.Pending == nil {
				for {
					if audioEventEvict(p) == 0 {
						break
					}
					next = n.next
					audioEventLoad(p)
					p.Pending = audioEventSelectFirst(p)
					if p.Pending != nil {
						break
					}
				}
			}
			// The second selection is observable through RNG consumption.
			p.Pending = audioEventSelectFirst(p)
			if p.Pending == nil || audioEventStart(p) == 0 {
				audioEventDelete(p)
				audioEventUnlink(p)
			}
		}
		n = next
	}
	*busy = 0
}
