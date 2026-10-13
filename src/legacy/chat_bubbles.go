package legacy

import (
	"encoding/binary"
	"github.com/brutusmcforce/opennox/client"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"unsafe"
)

type chatBubble struct {
	Text                     [318]uint16
	Length, Expiry, Position uint32
	X, Y                     int32
	Code, Visible, Arrow     uint32
	Drawable                 *client.Drawable
	Width, Height            int32
	Ordinal                  uint32
	Next, Previous           *chatBubble
}

// ChatBubbleSize is the size of a chat bubble in the "Chat" pool.
const ChatBubbleSize = unsafe.Sizeof(chatBubble{})

var chatBubblePool unsafe.Pointer
var chatBubbleTail *chatBubble

func chatBubbleHead() **chatBubble { return &chatBubbleList }
func chatBubbleLookup(code uint32) *chatBubble {
	for b := *chatBubbleHead(); b != nil; b = b.Next {
		if b.Code == code {
			return b
		}
	}
	return nil
}
func chatBubbleDestroy() {
	if chatBubblePool != nil {
		alloc.AsClass(chatBubblePool).Free()
	}
	chatBubblePool = nil
	*chatBubbleHead() = nil
}
func chatBubbleClear() {
	alloc.AsClass(chatBubblePool).FreeAllObjects()
	*chatBubbleHead() = nil
}
func chatBubbleUnlink(b *chatBubble) {
	if b.Previous != nil {
		b.Previous.Next = b.Next
	} else {
		*chatBubbleHead() = b.Next
	}
	if b.Next != nil {
		b.Next.Previous = b.Previous
	} else {
		chatBubbleTail = b.Previous
	}
	alloc.AsClass(chatBubblePool).FreeObjectFirst(unsafe.Pointer(b))
}
func chatBubbleRemove(code uint32) {
	if b := chatBubbleLookup(code); b != nil {
		chatBubbleUnlink(b)
	}
}
func chatBubbleCreate(header unsafe.Pointer, text *uint16) {
	data := unsafe.Slice((*byte)(header), 11)
	code := uint32(binary.LittleEndian.Uint16(data[1:]))
	b := chatBubbleLookup(code)
	fresh := b == nil
	if fresh {
		b = (*chatBubble)(alloc.AsClass(chatBubblePool).NewObject())
		if b == nil {
			return
		}
	}
	// The message decoder supplies the existing terminated UTF-16 text buffer.
	n := alloc.StrLen(text) + 1
	copy(unsafe.Slice(&b.Text[0], n), unsafe.Slice(text, n))
	b.Length = uint32(data[8])
	b.Position = binary.LittleEndian.Uint32(data[4:])
	b.Code = code
	ttl := uint32(binary.LittleEndian.Uint16(data[9:]))
	if ttl == 0 {
		ttl = GetServer().S().TickRate() * (min(b.Length/8, 8) + 2)
	}
	b.Expiry = GetServer().S().Frame() + ttl
	if fresh {
		if *chatBubbleHead() != nil {
			chatBubbleTail.Next = b
			b.Previous = chatBubbleTail
		} else {
			*chatBubbleHead() = b
			b.Previous = nil
		}
		b.Next = nil
		chatBubbleTail = b
	}
}
