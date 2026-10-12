package legacy

import (
	"github.com/brutusmcforce/opennox/legacy/timer"
	"unsafe"
)

type audioStreamFormat struct {
	Kind, Encoding, Rate, Channels, Width, ByteRate, Extra uint32
}
type audioStreamCatalog struct {
	Entries       *audioStreamEntry
	Count         uint32
	Directory     [260]byte
	Bag, Override *FILE
	UseOverride   uint32
	Active        *FILE
	Remaining     uint32
}
type audioStreamEntry struct {
	Name                               [16]byte
	Offset, Length, Rate, Flags, Extra uint32
}

// Data addresses are raw ABI words: copying or resetting them must not ask the
// Go write barrier to interpret an opaque value as a managed pointer.
type audioStreamBuffer struct {
	Data    uint32
	Length  uint32
	Chunks  legacyListNode
	Format  *audioStreamFormat
	Field24 uint32
}
type audioStreamChunk struct {
	Node   legacyListNode
	Data   uint32
	Length uint32
	Owner  *audioStreamBuffer
}
type audioStreamCache struct {
	Catalog         *audioStreamCatalog
	Blocks, Entries *unsafe.Pointer
	List            legacyListNode
	Payload         int32
}
type audioStreamCacheEntry struct {
	Node       legacyListNode
	References uint32
	Index      int32
	Valid      uint32
	Buffer     audioStreamBuffer
	Cache      *audioStreamCache
	Format     audioStreamFormat
}
type audioStreamDescriptor struct {
	Name                                 *byte
	Field4, Flags, InUse, Field16        uint32
	Init, Free, ContextInit, ContextFree unsafe.Pointer
	VoiceAPI                             *audioStreamVoiceAPI
}
type audioStreamDevice struct {
	Node       legacyListNode
	Descriptor *audioStreamDescriptor
	References int32
	Slots      int32
	Contexts   [16]*audioStreamContext
}
type audioStreamContext struct {
	Node                                  legacyListNode
	Flags, References                     uint32
	Device                                *audioStreamDevice
	Slot                                  int32
	Field28                               uint32
	Format, RequestedFormat               audioStreamFormat
	Timers                                timer.TimerGroup
	ParentTimers                          *timer.TimerGroup
	Index                                 int32
	VoiceCount, VoiceCapacity             int32
	Voices                                legacyListNode
	Mutating                              int32
	Tick                                  unsafe.Pointer
	Field220                              uint32
	Period, Elapsed, MaxElapsed, LastTick uint64
	VoiceAPI                              *audioStreamVoiceAPI
	Field260                              uint32
}
type audioStreamVoiceAPI struct {
	Field0                  uint32
	Init, Free, Start, Stop unsafe.Pointer
	Fields20                [3]unsafe.Pointer
	Update, Restart         unsafe.Pointer
}
type audioStreamVoice struct {
	Node                                    legacyListNode
	Kind                                    int32
	Timers                                  timer.TimerGroup
	ExtraTimers, GlobalTimers               *timer.TimerGroup
	Priority                                int32
	Flags                                   uint32
	Loops                                   int32
	Context                                 *audioStreamContext
	OnData, OnLoop, OnEnd, OnStop           unsafe.Pointer
	Fields152                               [5]uint32
	API                                     *audioStreamVoiceAPI
	Effective                               timer.TimerGroup
	Sample                                  unsafe.Pointer
	DataCallback, LoopCallback, EndCallback unsafe.Pointer
	Buffer                                  *audioStreamBuffer
	Chunk                                   *audioStreamChunk
	Data                                    uint32
	Remaining, Length                       uint32
	Field308                                uint32
}
type audioStreamSystem struct {
	Devices, Contexts legacyListNode
	Mutating          int32
	Field28           uint32
	Timers            timer.TimerGroup
}
