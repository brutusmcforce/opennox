package legacy

func chatBubblePlace(rect *[4]int32, arrow *uint32) {
	w, h := int32(nox_win_width), int32(nox_win_height)
	if rect[0] < 0 || rect[1] < 0 || rect[2] > w || rect[3] > h {
		*arrow = 0
	}
	inside := func(zone [4]int32) bool {
		a, b := [2]int32{rect[0], rect[1]}, [2]int32{rect[2], rect[3]}
		return geometryRectInt(&a, &zone) != 0 || geometryRectInt(&b, &zone) != 0
	}
	top := int32(55)
	if uiInventoryWindowOpenState() {
		top = 279
	}
	if inside([4]int32{0, 0, 563, top}) {
		rect[1] = top
		*arrow = 0
	}
	if Nox_client_getRenderGUI() != 0 {
		for _, zone := range [][4]int32{{0, h - 127, 111, h}, {w/2 - 160, h - 74, w/2 + 160, h}, {w - 91, h - 201, w, h}} {
			if inside(zone) {
				rect[1] += zone[1] - rect[3]
				*arrow = 0
			}
		}
		if summonFirst() != nil && inside([4]int32{w - 87, 0, w, 145}) {
			rect[1] = 145
			*arrow = 0
		}
	}
}
func chatBubbleRegion(x, y int32) int32 {
	w, h := int32(nox_win_width), int32(nox_win_height)
	var region int32
	if x < w/3 {
		region = 16
	} else if x < 2*w/3 {
		region = 32
	} else {
		region = 64
	}
	if y < h/3 {
		return region | 1
	}
	if y >= 2*h/3 {
		return region | 4
	}
	return region | 2
}
func chatBubbleOverlap(a, b *chatBubble) bool {
	pad := uint32(GetClient().R2().FontHeight(nil))
	ax, ay, bx, by := uint32(a.X), uint32(a.Y), uint32(b.X), uint32(b.Y)
	// C's mixed signed/unsigned comparisons use unsigned 32-bit bounds.
	return ax-pad < pad+bx+uint32(b.Width) && pad+ax+uint32(a.Width) > bx-pad && ay-pad < pad+by+uint32(b.Height) && int32(pad+ay+uint32(a.Height)) > b.Y-int32(pad)
}
func chatBubbleShift(mask byte, a, b *chatBubble) (x, y int32) {
	gap := 2 * int32(GetClient().R2().FontHeight(nil))
	left, right := b.X-a.Width-gap, b.X+b.Width+gap
	top, bottom := b.Y-a.Height-gap, b.Y+b.Height+gap
	switch mask {
	case 1:
		return left, top
	case 2:
		return right, top
	case 4:
		return left, bottom
	case 8:
		return right, bottom
	case 16:
		return a.X, top
	case 32:
		return a.X, bottom
	case 64:
		return left, a.Y
	case 128:
		return right, a.Y
	}
	return
}
func chatBubbleCandidate(b *chatBubble, x, y int32) bool {
	rect := [4]int32{x, y, x + b.Width, y + b.Height}
	arrow := uint32(1)
	chatBubblePlace(&rect, &arrow)
	if arrow == 0 {
		return false
	}
	for p := *chatBubbleHead(); p != nil; p = p.Next {
		// Preserve the original moving-bubble visibility gate.
		if b.Visible != 0 {
			if p.Ordinal >= b.Ordinal {
				return true
			}
			oldX, oldY := b.X, b.Y
			b.X, b.Y = x, y
			overlaps := chatBubbleOverlap(b, p)
			b.X, b.Y = oldX, oldY
			if overlaps {
				return false
			}
		}
	}
	return true
}
func chatBubbleArrange(b *chatBubble) {
	if b == *chatBubbleHead() {
		return
	}
	var obstacle *chatBubble
	for p := *chatBubbleHead(); p != nil; p = p.Next {
		if p.Visible != 0 {
			if p.Ordinal >= b.Ordinal {
				return
			}
			if chatBubbleOverlap(b, p) {
				obstacle = p
				break
			}
		}
	}
	if obstacle == nil {
		return
	}
	var first, second byte
	switch chatBubbleRegion(b.X+b.Width/2, b.Y+b.Height/2) {
	case 17:
		first, second = 160, 8
	case 18:
		first, second = 48, 138
	case 20:
		first, second = 144, 2
	case 33:
		first, second = 192, 44
	case 34:
		first, second = 240, 15
	case 36:
		first, second = 192, 19
	case 65:
		first, second = 96, 4
	case 66:
		first, second = 48, 69
	case 68:
		first, second = 80, 1
	}
	for _, priority := range []byte{first, second} {
		for i := 0; i < 8; i++ {
			mask := byte(1 << i)
			if priority&mask == 0 {
				continue
			}
			x, y := chatBubbleShift(mask, b, obstacle)
			if chatBubbleCandidate(b, x, y) {
				b.X, b.Y = x, y
				return
			}
		}
	}
}
