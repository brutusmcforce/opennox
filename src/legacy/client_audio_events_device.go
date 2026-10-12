package legacy

import (
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/legacy/client/audio/ail"
	"unsafe"
)

var audioEventDeviceUser = func(s ail.Sample) *AudioSample {
	p := s.UserData()
	if p == nil {
		return nil
	}
	return p.(*AudioSample)
}
var audioEventDeviceReady = func(s ail.Sample) int { return s.BufferReady() }
var audioEventDeviceLoad = func(s ail.Sample, index uint32, data []byte) { s.LoadBuffer(index, data) }

func audioEventSampleEnded(sample ail.Sample) {
	p := audioEventDeviceUser(sample)
	if p.Flag7 == 0 {
		v := (*audioStreamVoice)(p.Field1)
		AudioStreamCallbackInt(v.EndCallback, unsafe.Pointer(v))
		p.Flag7 = 1
	}
}
func audioEventSampleRefill(p *AudioSample) int32 {
	for index := audioEventDeviceReady(p.Smp); index != -1; index = audioEventDeviceReady(p.Smp) {
		v := (*audioStreamVoice)(p.Field1)
		scratch := p.Data1
		if index != 0 {
			scratch = p.Data2
		}
		var data, first uint32
		var total, firstLength int32
		copied := false
		if p.Field3 == 0 {
			for total < 16384 {
				remaining := int32(v.Remaining)
				if remaining == 0 {
					AudioStreamCallbackVoid(v.DataCallback, unsafe.Pointer(v))
					remaining = int32(v.Remaining)
					if remaining == 0 {
						AudioStreamCallbackVoid(v.LoopCallback, unsafe.Pointer(v))
						remaining = int32(v.Remaining)
						if remaining == 0 {
							p.Field3 = 1
							break
						}
					}
				}
				if total == 0 {
					data = v.Data
				}
				n := remaining
				if total == 0 && remaining >= 16384 && first == 0 {
					// Large contiguous input is passed directly, even when larger than scratch.
				} else {
					if n+total > 16384 {
						n = 16384 - total
					}
					if n != 0 {
						if first != 0 {
							if firstLength != 0 {
								copy(unsafe.Slice(scratch, int(firstLength)), unsafe.Slice((*byte)(memmap.WordPtr(first)), int(firstLength)))
								data = memmap.PtrWord(unsafe.Pointer(scratch))
								first = 0
								copied = true
							}
							copy(unsafe.Slice((*byte)(unsafe.Add(unsafe.Pointer(scratch), uintptr(uint32(total)))), int(n)), unsafe.Slice((*byte)(memmap.WordPtr(v.Data)), int(n)))
						} else if !copied {
							first = v.Data
							firstLength = n
						} else {
							copy(unsafe.Slice((*byte)(unsafe.Add(unsafe.Pointer(scratch), uintptr(uint32(total)))), int(n)), unsafe.Slice((*byte)(memmap.WordPtr(v.Data)), int(n)))
						}
					}
				}
				v.Length -= uint32(n)
				v.Remaining -= uint32(n)
				v.Data += uint32(n)
				total += n
				if p.Field3 != 0 {
					break
				}
			}
		}
		audioEventDeviceLoad(p.Smp, uint32(index), unsafe.Slice((*byte)(memmap.WordPtr(data)), int(total)))
	}
	return -1
}
func audioEventSampleFormat(p *audioStreamFormat) int32 {
	if p.Encoding == 2 {
		if p.Channels == 2 {
			return 7
		}
		return 5
	}
	if p.Encoding != 0 {
		return 0
	}
	result := int32(0)
	if p.Width == 2 {
		result = 1
	}
	if p.Channels == 2 {
		result += 2
	}
	return result
}
