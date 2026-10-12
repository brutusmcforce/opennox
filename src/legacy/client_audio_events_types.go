package legacy

import (
	"github.com/brutusmcforce/opennox/legacy/timer"
	"unsafe"
)

type audioEventMetadata struct {
	Enabled, Flags, Field8, Field12 uint32
	Volume                          timer.Timer
	Priority                        int32
	Active, Limit, Loops            uint32
	Distance, DelayMin, DelayMax    uint32
	PitchMin, PitchMax              int32
	Name                            *byte
	Ranking                         legacyListNode
	Frame, BucketGeneration         uint32
	BucketPriority                  int32
	BucketNode                      legacyListNode
	BucketVolume                    uint32
	Samples                         [32]int16
	SampleCount, Field196           uint32
}
type audioEvent struct {
	Node, Rank              legacyListNode
	Flags, State, NextState uint32
	Metadata                *audioEventMetadata
	Samples                 [32]*audioStreamCacheEntry
	Loaded, LastSample      int32
	Voice                   *audioStreamVoice
	Field180                uint32
	Timers                  timer.TimerGroup
	Serial, Delay           uint32
	Deadline                uint64
	Pending                 *audioStreamBuffer
	PriorityOffset          int32
	Order                   [32]int32
	Remaining, Iteration    int32
	ReloadOrder             [32]int32
	ReloadRemaining         int32
	Field572                uint32
}
type audioEventHandle struct {
	Event    *audioEvent
	Serial   uint32
	Metadata *audioEventMetadata
}

var (
	_ = [1]struct{}{}[200-unsafe.Sizeof(audioEventMetadata{})]
	_ = [1]struct{}{}[16-unsafe.Offsetof(audioEventMetadata{}.Volume)]
	_ = [1]struct{}{}[88-unsafe.Offsetof(audioEventMetadata{}.Ranking)]
	_ = [1]struct{}{}[112-unsafe.Offsetof(audioEventMetadata{}.BucketNode)]
	_ = [1]struct{}{}[128-unsafe.Offsetof(audioEventMetadata{}.Samples)]
	_ = [1]struct{}{}[192-unsafe.Offsetof(audioEventMetadata{}.SampleCount)]
	_ = [1]struct{}{}[576-unsafe.Sizeof(audioEvent{})]
	_ = [1]struct{}{}[168-unsafe.Offsetof(audioEvent{}.Loaded)]
	_ = [1]struct{}{}[176-unsafe.Offsetof(audioEvent{}.Voice)]
	_ = [1]struct{}{}[184-unsafe.Offsetof(audioEvent{}.Timers)]
	_ = [1]struct{}{}[280-unsafe.Offsetof(audioEvent{}.Serial)]
	_ = [1]struct{}{}[288-unsafe.Offsetof(audioEvent{}.Deadline)]
	_ = [1]struct{}{}[296-unsafe.Offsetof(audioEvent{}.Pending)]
	_ = [1]struct{}{}[304-unsafe.Offsetof(audioEvent{}.Order)]
	_ = [1]struct{}{}[440-unsafe.Offsetof(audioEvent{}.ReloadOrder)]
	_ = [1]struct{}{}[568-unsafe.Offsetof(audioEvent{}.ReloadRemaining)]
	_ = [1]struct{}{}[12-unsafe.Sizeof(audioEventHandle{})]
)
