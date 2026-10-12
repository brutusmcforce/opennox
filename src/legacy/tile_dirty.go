package legacy

import "unsafe"

// tileBufMark marks the tile buffer cells (see nox_video_tileBuf_dirty) covered by n bytes written at the given
// offset in the tile buffer. Every writer of the tile buffer must call it, so the GPU renderer only uploads the
// cells that changed.
func tileBufMark(off, n int) {
	dirty := nox_video_tileBuf_dirty
	if dirty == nil || n <= 0 {
		return
	}
	const cell = 46
	pitch := int(int32(dword_5d4594_3798804))
	cols := int(int32(dword_5d4594_3798812))
	size := int(uintptr(legacyGlobals.nox_video_tileBuf_end_3798844) - uintptr(legacyGlobals.nox_video_tileBuf_ptr_3798796))
	if pitch <= 0 || off < 0 || off >= size {
		return
	}
	n = min(n, size-off)
	for n > 0 {
		x := off % pitch
		l := min(pitch-x, n)
		row := unsafe.Add(dirty, (off/pitch/cell)*cols)
		for c := x / (2 * cell); c <= (x+l-1)/(2*cell); c++ {
			*(*byte)(unsafe.Add(row, c)) = 1
		}
		off += l
		n -= l
	}
}

// tileBufMarkTile marks a whole tile (46x46 pixels) drawn at the given offset by the packed tile writers.
func tileBufMarkTile(off int) {
	pitch := int(int32(dword_5d4594_3798804))
	for i := 0; i < 46; i++ {
		tileBufMark(off+i*pitch, 2*46)
	}
}
