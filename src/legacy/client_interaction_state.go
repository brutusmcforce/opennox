package legacy

import "unsafe"

// Private client-interaction state. Window/data handles retain their legacy
// 32-bit representation while ownership resides in Go.
var (
	interactionChatRootWord         uint32
	interactionChatEditWord         uint32
	interactionChatDataWord         uint32
	interactionChatActive           uint32
	interactionHoverDepth           uint32
	interactionHoverDrawable        unsafe.Pointer
	interactionConversationActive   uint32
	interactionConversationRoot     unsafe.Pointer
	interactionObserverIcon         uint32
	interactionGameOverRootWord     uint32
	interactionHelpRootWord         uint32
	interactionKeyState             uint32
	interactionKeyRootWord          uint32
	interactionVoteIcon             uint32
	interactionPrimaryKey           uint32
	interactionDrawToggle           uint32
	interactionMessageHead          uint32
	interactionChatIcon             uint32
	interactionUsableCursorDrawable unsafe.Pointer
	interactionCursorMode           uint32 = 5
)
