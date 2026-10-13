package server

import "unsafe"

// Secret wall states (SecretWallData.State).
const (
	SecretWallClosed  = 1 // closed
	SecretWallClosing = 2 // closing: Progress counts down to 0
	SecretWallOpen    = 3 // open
	SecretWallOpening = 4 // opening: Progress counts up to SecretWallProgressMax
)

// SecretWallProgressMax is the Progress value of a fully open secret wall.
const SecretWallProgressMax = 23

// SecretWallData is the record of a secret wall (Wall.Data when the wall has wall.FlagSecret).
// It is a pointer-free block with the 32-bit layout of the map file; references are memmap pointer words.
// Records are linked into the global secret wall list through Next.
type SecretWallData struct {
	Next      uint32 // 0, 0: next record in the secret wall list (memmap word)
	GridX     uint32 // 1, 4: wall grid X
	GridY     uint32 // 2, 8: wall grid Y
	Wall      uint32 // 3, 12: owning *Wall (memmap word)
	OpenDelay uint32 // 4, 16: seconds before an automatic open/close
	Flags     byte   // 5, 20: 0x2 opens on touch, 0x4 closes automatically, 0x8 opens automatically (initially open in the map)
	State     byte   // 5, 21: SecretWallClosed..SecretWallOpening
	Progress  byte   // 5, 22: open progress, 0 (closed) to SecretWallProgressMax (open)
	Field23   byte   // 5, 23
	Timer     uint32 // 6, 24: frames left until the automatic open/close
	OpenSent  uint32 // 7, 28: bitmask of players that were told the wall is open
}

var _ = [1]struct{}{}[32-unsafe.Sizeof(SecretWallData{})]

// AsSecretWallData converts a secret wall record pointer (Wall.Data or a list node).
func AsSecretWallData(p unsafe.Pointer) *SecretWallData {
	return (*SecretWallData)(p)
}

// SecretData returns the secret wall record of the wall. The wall must have wall.FlagSecret.
func (w *Wall) SecretData() *SecretWallData {
	return (*SecretWallData)(w.Data)
}
