package seat

import "image"

// Cursor is a system mouse cursor created by CursorScreen.
type Cursor interface {
	// Free releases the cursor. It must not be the current cursor.
	Free()
}

// CursorScreen is implemented by screens that can show a system (hardware) mouse cursor.
// The system draws such cursor, so it follows the mouse without waiting for the next frame.
type CursorScreen interface {
	// NewCursor creates a cursor from an image. Hot is the position of the pointer in the image.
	NewCursor(img image.Image, hot image.Point) (Cursor, error)
	// SetCursor shows a given cursor, or hides the cursor if it's nil.
	SetCursor(c Cursor)
}
