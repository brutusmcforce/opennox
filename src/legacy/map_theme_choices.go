package legacy

func mapThemeAppend(head *uint32, p uint32, next int) {
	if *head == 0 {
		*head = p
		return
	}
	q := *head
	for *populationWord(q, next) != 0 {
		q = *populationWord(q, next)
	}
	*populationWord(q, next) = p
}
func mapThemeChoices(f uint32) uint32 {
	var head uint32
	remaining, stars := int32(100), int32(0)
	for {
		p := mapThemeAlloc(1, 2060)
		if p == 0 || !mapThemeNext(f) {
			break
		}
		if mapThemeText() == "*" {
			*populationWord(p, 0) = ^uint32(0)
			stars++
		} else {
			n := mapThemeInt(mapThemeText())
			*populationWord(p, 0) = uint32(n)
			remaining -= n
		}
		linked := false
		for mapThemeNext(f) {
			key := mapThemeLower(mapThemeText())
			if key == "end" || key == "or" {
				mapThemeAppend(&head, p, 2056)
				linked = true
				if key == "end" {
					if stars != 0 {
						for q := head; q != 0; q = *populationWord(q, 2056) {
							if int32(*populationWord(q, 0)) < 0 {
								*populationWord(q, 0) = uint32(remaining/stars + 1)
							}
						}
					}
					return head
				}
				break
			}
			kind := mapThemeTable(253216, mapThemeText())
			if kind < 0 {
				mapThemeFree(p)
				break
			}
			index := *populationWord(p, 2052)
			*populationWord(p, 4+64*int(index)) = uint32(kind)
			if mapThemeRead(f, mapThemeByte(p, 8+64*int(index))) == 0 {
				break
			}
			index++
			if index >= 32 {
				index = 31
			}
			*populationWord(p, 2052) = index
		}
		if !linked {
			break
		}
	}
	for p := head; p != 0; {
		q := *populationWord(p, 2056)
		mapThemeFree(p)
		p = q
	}
	return 0
}
func mapThemeForeach(f uint32) uint32 {
	if !mapThemeNext(f) {
		return 0
	}
	p := mapThemeAlloc(1, 12)
	if p == 0 {
		return 0
	}
	*populationWord(p, 0) = uint32(GetServer().S().Types.IndByID(mapThemeText()))
	if mapThemeNext(f) && mapThemeLower(mapThemeText()) == "contains" {
		q := mapThemeChoices(f)
		*populationWord(p, 4) = q
		if q != 0 {
			return p
		}
	}
	mapThemeFree(p)
	return 0
}
func mapThemeDecorSet(p, f uint32) uint32 {
	set := mapThemeAlloc(1, 24)
	if set == 0 {
		return 0
	}
	if !mapThemeNext(f) {
		return 0
	}
	*populationWord(set, 0) = uint32(mapThemeInt(mapThemeText()))
	if !mapThemeNext(f) {
		return 0
	}
	*populationWord(set, 4) = uint32(mapThemeInt(mapThemeText()))
	var item uint32
	for mapThemeNext(f) {
		key := mapThemeLower(mapThemeText())
		switch key {
		case "end":
			mapThemeAppend(populationWord(p, 92), set, 20)
			*populationWord(set, 20) = 0
			*populationWord(p, 96)++
			return 1
		case "foreach":
			if item == 0 || *populationWord(item, 0) != 1 {
				return 0
			}
			q := mapThemeForeach(f)
			if q == 0 {
				return 0
			}
			*populationWord(q, 8) = *populationWord(item, 84)
			*populationWord(item, 84) = q
			continue
		case "contains":
			if item == 0 || *populationWord(item, 0) != 0 {
				return 0
			}
			q := mapThemeChoices(f)
			*populationWord(item, 80) = q
			if q == 0 {
				return 0
			}
			continue
		}
		item = mapThemeAlloc(1, 100)
		if item == 0 {
			return 0
		}
		kind := mapThemeTable(253216, mapThemeText())
		if kind < 0 {
			kind = 6
		}
		*populationWord(item, 0) = uint32(kind)
		if kind == 6 {
			mapThemeFree(item)
			return 0
		}
		read := func() bool {
			if mapThemeNext(f) {
				return true
			}
			mapThemeFree(item)
			return false
		}
		switch kind {
		case 0, 1:
			if !read() {
				return 0
			}
			mapThemeCopy(item+4, mapThemeText())
			if !read() {
				return 0
			}
			if mapThemeLower(mapThemeText()) == "density" {
				*populationWord(item, 64) = 1
				if !read() {
					return 0
				}
				*populationFloat(item, 76) = float32(mapThemeFloat(mapThemeText()))
				if !read() {
					return 0
				}
				*populationWord(item, 68) = mapThemeBound(mapThemeText(), 0)
				if !read() {
					return 0
				}
				*populationWord(item, 72) = mapThemeBound(mapThemeText(), 999999)
			} else {
				*populationWord(item, 64) = 0
				*populationWord(item, 68) = uint32(mapThemeInt(mapThemeText()))
				if !read() {
					return 0
				}
				*populationWord(item, 72) = uint32(mapThemeInt(mapThemeText()))
			}
		case 3, 4, 5:
			if !read() {
				return 0
			}
			mapThemeCopy(item+4, mapThemeText())
		}
		*populationWord(item, 88) = *populationWord(set, 8)
		*populationWord(set, 8) = item
		*populationWord(set, 12)++
	}
	return 0
}
