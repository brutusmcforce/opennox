package legacy

import (
	"github.com/brutusmcforce/opennox/legacy/timer"
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
