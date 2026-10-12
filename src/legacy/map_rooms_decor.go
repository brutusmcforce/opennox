package legacy

import "unsafe"

func mapRoomAssignDecoration(config unsafe.Pointer, r *mapRoom) unsafe.Pointer {
	r.Decoration = mapRoomSelectDecoration(r, unsafe.Add(config, 184))
	return r.Decoration
}
func mapRoomSelectDecoration(r *mapRoom, list unsafe.Pointer) unsafe.Pointer {
	weight := int32(mapRoomRaw(unsafe.Pointer(r)))
	for attempt := 0; ; attempt++ {
		switch r.ThemeFlags {
		case 0:
			// Backdrops are created without a theme phase. Only unrestricted
			// decorations match them; their weights must not depend on r's address.
			weight = 0
			for d := *mapRoomRef(list, 0); d != nil; d = *mapRoomRef(d, 220) {
				if mapRoomDecorMatchesFlags(d, r) != 0 {
					weight += *mapRoomWord(d, 72)
				}
			}
		case 1:
			weight = *mapRoomWord(list, 8)
		case 2:
			weight = *mapRoomWord(list, 12)
		case 4:
			weight = *mapRoomWord(list, 16)
		case 8:
			weight = *mapRoomWord(list, 20)
		case 16:
			weight = *mapRoomWord(list, 24)
		case 32:
			weight = *mapRoomWord(list, 28)
		}
		if weight <= 0 {
			return nil
		}
		draw := mapRoomRandomInt(0, weight)
		d := *mapRoomRef(list, 0)
		if d == nil {
			continue
		}
		for d != nil {
			if mapRoomDecorMatchesFlags(d, r) != 0 {
				draw -= *mapRoomWord(d, 72)
				if draw <= 0 {
					break
				}
			}
			d = *mapRoomRef(d, 220)
		}
		if d == nil {
			continue
		}
		if mapRoomDecorFitsRoom(d, r) != 0 && *mapRoomByte(d, 66) == 0 {
			mapRoomConsumeDecoration(list, d)
			return d
		}
		if attempt >= 100 {
			return nil
		}
	}
}
func mapRoomDecorMatchesFlags(d unsafe.Pointer, r *mapRoom) uint32 {
	flags := *mapRoomByte(d, 64)
	return uint32(bool2int(byte(r.ThemeFlags)&flags != 0 || flags == 0))
}
func mapRoomDecorFitsRoom(d unsafe.Pointer, r *mapRoom) uint32 {
	size := r.Width
	if size <= r.Height {
		size = r.Height
	}
	return uint32(bool2int(size >= *mapRoomWord(d, 76) && size <= *mapRoomWord(d, 80)))
}
func mapRoomConsumeDecoration(list, d unsafe.Pointer) int8 {
	count := mapRoomByte(d, 65)
	if *count == 0 {
		return 0
	}
	*count--
	if *count != 0 {
		return int8(*count)
	}
	for i := 1; i <= 4; i++ {
		if *mapRoomByte(d, 64)&(1<<i) != 0 {
			*mapRoomWord(list, 8+4*i) -= *mapRoomWord(d, 72)
		}
	}
	var prev unsafe.Pointer
	cur := *mapRoomRef(list, 0)
	for cur != nil {
		if cur == d {
			break
		}
		prev = cur
		cur = *mapRoomRef(cur, 220)
	}
	if prev != nil {
		cur = *mapRoomRef(d, 220)
		*mapRoomRef(prev, 220) = cur
	} else {
		*mapRoomRef(list, 0) = *mapRoomRef(d, 220)
	}
	tail := *mapRoomRef(list, 0)
	if tail != nil {
		for cur = *mapRoomRef(tail, 220); cur != nil; cur = *mapRoomRef(cur, 220) {
			tail = cur
		}
		*mapRoomRef(tail, 220) = d
	} else {
		*mapRoomRef(list, 0) = d
	}
	*mapRoomRef(d, 220) = nil
	*mapRoomByte(d, 66) = 1
	return int8(mapRoomRaw(cur))
}
func mapRoomAssignRequiredDecorations(config unsafe.Pointer) uint32 {
	lists := [2]unsafe.Pointer{unsafe.Add(config, 88), unsafe.Add(config, 120)}
	for kind, list := range lists {
		for d := *mapRoomRef(list, 0); d != nil; {
			next := *mapRoomRef(d, 220)
			if *mapRoomByte(d, 67) != 0 && *mapRoomByte(d, 68) == 0 {
				r := mapRoomHead()
				for r != nil {
					if r.Decoration == nil && mapRoomIsHall(r) == uint32(kind) && mapRoomDecorMatchesFlags(d, r) != 0 && mapRoomDecorFitsRoom(d, r) != 0 && *mapRoomByte(d, 66) == 0 {
						break
					}
					r = r.Next
				}
				if r == nil {
					return 0
				}
				mapRoomConsumeDecoration(list, d)
				r.Decoration = d
				*mapRoomByte(d, 68) = 1
			}
			d = next
		}
	}
	for r := mapRoomHead(); r != nil; r = r.Next {
		if r.Decoration == nil {
			r.Decoration = mapRoomSelectDecoration(r, lists[mapRoomIsHall(r)])
			if r.Decoration == nil {
				return 0
			}
		}
	}
	return 1
}
