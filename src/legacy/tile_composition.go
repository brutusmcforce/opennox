package legacy

import (
	"github.com/brutusmcforce/opennox/client/noxrender"
	"github.com/brutusmcforce/opennox/common/memmap"
	"image"
	"unsafe"
)

var tileDrawCallback = func(pos image.Point, img noxrender.ImageHandle, tile uint16) { tileRasterTexture(pos, img) }
var tileEdgeCallback = tileCompositionOverlay

func SetTileDrawCallbacks(textured bool) {
	if textured {
		tileDrawCallback = func(pos image.Point, img noxrender.ImageHandle, tile uint16) { tileRasterTexture(pos, img) }
		tileEdgeCallback = tileCompositionOverlay
	} else {
		tileDrawCallback = func(pos image.Point, img noxrender.ImageHandle, tile uint16) { tileRasterFill(pos, tile) }
		tileEdgeCallback = func(image.Point, *[5]uint32) {}
	}
}
func tileCompositionReset() { nox_xxx_waypointCounterMB_587000_154948 = 0xffffffff }
func tileCompositionEdges(pos image.Point, edge *[5]uint32) {
	for edge != nil {
		tileEdgeCallback(pos, edge)
		edge = (*[5]uint32)(memmap.WordPtr(edge[4]))
	}
}
func tileCompositionCell(x, y int32, pos image.Point) {
	cell := mapPaintCell(x, y)
	for _, half := range []int{1, 0} {
		if cell[0]&(1<<half) == 0 {
			continue
		}
		off := 1 + 5*half
		tile := uint16(cell[off])
		def := &tileDefinitionsAll()[tile]
		img := *(*noxrender.ImageHandle)(unsafe.Add(def.Data32, 4*uintptr(cell[off+1]+uint32(def.Field46))))
		p := pos
		if half == 1 {
			p.Y += 23
		} else {
			p.X += 23
		}
		tileDrawCallback(p, img, tile)
		*memmap.PtrUint32(0x85B3FC, 228+4*uintptr(tile)) = 1
		tileCompositionEdges(p, (*[5]uint32)(memmap.WordPtr(cell[off+4])))
	}
}
func tileCompositionBounds(x, y int32) (lx, ly, hx, hy int32) {
	lx = max(0, x)
	ly = max(0, y)
	cols, rows := int32(dword_5d4594_3798812), int32(dword_5d4594_3798816)
	hx = cols + lx - 1
	if hx >= 128 {
		hx = 127
		lx = 127 - cols
	}
	hy = rows + ly
	if hy >= 128 {
		hy = 127
		ly = 127 - rows
	}
	return
}

// tileCompositionFullBounds is tileCompositionBounds for a full redraw.
func tileCompositionFullBounds(x, y int32) (lx, ly, hx, hy int32) {
	lx = max(0, x)
	ly = max(0, y)
	cols, rows := int32(dword_5d4594_3798812), int32(dword_5d4594_3798816)
	hx = cols + lx - 1
	if hx >= 128 {
		// OpenNox: the original clamped it to 127, but tileCompositionHorizontal only stops scrolling at 128,
		// causing a full redraw on each frame near the map edge
		hx = 128
		lx = 129 - cols
	}
	hy = rows + ly
	if hy >= 128 {
		// OpenNox: see above, tileCompositionVertical stops scrolling at 128
		hy = 128
		ly = 128 - rows
	}
	return
}
func tileCompositionRedraw(vp *noxrender.Viewport) int {
	dp := vp.ToWorldPos(image.Point{})
	// The predicate's X expression is unsigned in the original owner.
	lx, ly, hx, hy := tileCompositionBounds(int32((uint32(dp.X)-11)/46), int32(dp.Y-11)/46-1)
	for y := ly; y < hy; y++ {
		for x := lx; x < hx; x++ {
			c := mapPaintCell(x, y)
			if c[0]&1 != 0 && tileDefinitionsAll()[c[1]].Field58&1 != 0 {
				return 1
			}
			if c[0]&2 != 0 && tileDefinitionsAll()[c[6]].Field58&1 != 0 {
				return 1
			}
		}
	}
	return 0
}
func tileCompositionFull(vp *noxrender.Viewport) {
	dword_5d4594_3798836 = 0
	dword_5d4594_3798840 = 0
	tileCompositionReset()
	clear(unsafe.Slice(memmap.PtrUint32(0x85B3FC, 228), 176))
	clear(unsafe.Slice(memmap.PtrUint32(0x5D4594, 2523980), 64))
	dp := vp.ToWorldPos(image.Point{})
	lx, ly, hx, hy := tileCompositionFullBounds(int32(dp.X-11)/46, int32(dp.Y-11)/46-1)
	dword_5d4594_3798828 = uint32(lx)
	dword_5d4594_3798832 = uint32(ly)
	dword_5d4594_3798820 = uint32(46*lx - 11)
	dword_5d4594_3798824 = uint32(46*ly - 11)
	for y := ly; y < hy; y++ {
		for x := lx; x < hx; x++ {
			tileCompositionCell(x, y, image.Pt(int(46*x-11), int(46*y-11)))
		}
	}
}
func tileCompositionHorizontal(vp *noxrender.Viewport, view int) {
	ox := int32(dword_5d4594_3798820)
	width := int32(dword_5d4594_3798800)
	height := int32(dword_5d4594_3798808)
	tx := int32(dword_5d4594_3798828)
	cols := int32(dword_5d4594_3798812)
	sx, sy := int32(dword_5d4594_3798836), int32(dword_5d4594_3798840)
	var column, px int32
	if int32(view) >= ox+23 {
		right := int32(Nox_getBackbufWidth()) + int32(view)
		if right <= width+ox-46 || cols+tx-1 >= 128 {
			return
		}
		if right > width+ox {
			tileCompositionFull(vp)
			return
		}
		tx++
		ox += 46
		column = cols + tx - 2
		sx += 46
		if sx >= width {
			sx -= width
			sy++
			if sy >= height {
				sy -= height
			}
		}
		px = width + ox - 92
	} else {
		if tx <= 0 {
			return
		}
		if int32(view) < ox-23 {
			tileCompositionFull(vp)
			return
		}
		ox -= 46
		tx--
		column = tx
		px = ox
		sx -= 46
		if sx < 0 {
			sx += width
			sy--
			if sy < 0 {
				sy += height
			}
		}
	}
	dword_5d4594_3798820 = uint32(ox)
	dword_5d4594_3798828 = uint32(tx)
	dword_5d4594_3798836 = uint32(sx)
	dword_5d4594_3798840 = uint32(sy)
	tileCompositionReset()
	py := int32(dword_5d4594_3798824)
	start := int32(dword_5d4594_3798832)
	end := start + int32(dword_5d4594_3798816)
	for y := start; y < end; y++ {
		tileCompositionCell(column, y, image.Pt(int(px), int(py)))
		py += 46
	}
}
func tileCompositionVertical(vp *noxrender.Viewport, view int) {
	oy := int32(dword_5d4594_3798824)
	height := int32(dword_5d4594_3798808)
	rows := int32(dword_5d4594_3798816)
	ty := int32(dword_5d4594_3798832)
	sy := int32(dword_5d4594_3798840)
	var row, py int32
	if int32(view) >= oy+23 {
		bottom := int32(view) + int32(Nox_getBackbufHeight())
		if bottom <= oy+height {
			return
		}
		if ty+rows >= 128 {
			return
		}
		if bottom > oy+height+46 {
			tileCompositionFull(vp)
			return
		}
		row = rows + ty
		ty++
		oy += 46
		sy += 46
		if sy >= height {
			sy -= height
		}
		py = oy + height - 46
	} else {
		if ty <= 0 {
			return
		}
		if int32(view) < oy-23 {
			tileCompositionFull(vp)
			return
		}
		ty--
		row = ty
		oy -= 46
		sy -= 46
		if sy < 0 {
			sy += height
		}
		py = oy
	}
	dword_5d4594_3798824 = uint32(oy)
	dword_5d4594_3798832 = uint32(ty)
	dword_5d4594_3798840 = uint32(sy)
	tileCompositionReset()
	px := int32(dword_5d4594_3798820)
	start := int32(dword_5d4594_3798828)
	end := start + int32(dword_5d4594_3798812) - 1
	for x := start; x < end; x++ {
		tileCompositionCell(x, row, image.Pt(int(px), int(py)))
		px += 46
	}
}
func tileCompositionOverlay(pos image.Point, edge *[5]uint32) {
	def := &tileDefinitionsAll()[edge[0]]
	tileImage := *(*noxrender.ImageHandle)(unsafe.Add(def.Data32, 4*uintptr(edge[1]+uint32(def.Field46))))
	table := edgeAssets[edge[2]].Data32
	index := edge[3] + uint32(edgeAssets[edge[2]].Field46)
	edgeImage := *(*noxrender.ImageHandle)(unsafe.Add(table, 4*uintptr(index)))
	*memmap.PtrUint32(0x5D4594, 2523980+4*uintptr(edge[2])) = 1
	data := GetClient().R2().GetBag().AsImage(edgeImage).Pixdata()
	if len(data) == 0 {
		return
	}
	source := GetClient().R2().GetBag().AsImage(tileImage).Pixdata()
	if len(source) == 0 {
		return
	}
	first, last := int(data[0]), int(data[1])
	data = data[2:]
	base, offset, size, stride := tileRasterOffset(pos)
	shift := memmap.Uint8(0x973F18, 7696)
	src := int(memmap.Uint32(0x973CE0, 4*uintptr(first)) << shift)
	for y := first; y <= last; y++ {
		count := int(memmap.Uint32(0x973CE0, 384+4*uintptr(y)))
		left := int(memmap.Uint32(0x973CE0, 192+4*uintptr(y)) << shift)
		at := (offset + y*stride + left) % size
		for count > 0 {
			op, n := data[0], int(data[1])
			data = data[2:]
			bytes := n << shift
			var pixels []byte
			switch op {
			case 1:
			case 2:
				pixels = data[:bytes]
				data = data[bytes:]
			case 3:
				pixels = source[src : src+bytes]
			default:
				return
			}
			if pixels != nil {
				split := min(bytes, size-at)
				tileBufMark(at, split)
				copy(unsafe.Slice((*byte)(unsafe.Add(base, at)), split), pixels[:split])
				if split < bytes {
					tileBufMark(0, bytes-split)
					copy(unsafe.Slice((*byte)(base), bytes-split), pixels[split:])
				}
			}
			count -= n
			at = (at + bytes) % size
			src += bytes
		}
	}
}
