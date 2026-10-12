package legacy

import (
	"fmt"
	"image"
	"unsafe"

	noxcolor "github.com/brutusmcforce/opennox/libs/color"
	"github.com/brutusmcforce/opennox/libs/datapath"
	"github.com/brutusmcforce/opennox/libs/ifs"
	"github.com/brutusmcforce/opennox/libs/strman"
	"github.com/brutusmcforce/opennox/client/gui"
	flags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
)

var characterUI = struct {
	defaults                      uint32
	classHost, colorHost          unsafe.Pointer
	classRoot, colorRoot, palette *gui.Window
	classAnim, colorAnim          *gui.Anim
	target                        uint32
	controls                      [15]*gui.Window
	modifiers                     [3]unsafe.Pointer // pants, shirt, shoes
}{defaults: 1}

func characterWord(w *gui.Window) *uint32 { return (*uint32)(unsafe.Add(w.C(), 32)) }
func characterColor(w *gui.Window) []byte {
	v := *characterWord(w)
	return unsafe.Slice(memmap.PtrUint8(0x5D4594, 1307796+3*uintptr((v>>16)+32*uint32(uint16(v)))), 3)
}
func characterColor555(w *gui.Window) noxcolor.RGBA5551 {
	c := characterColor(w)
	return noxcolor.RGB5551Color(c[0], c[1], c[2])
}
func characterName(dst *uint16) int {
	n := alloc.StrLen(dst)
	s := unsafe.Slice(dst, n+1)
	space := func(v uint16) bool { return v == 32 || v >= 9 && v <= 13 }
	start := 0
	for start < n && space(s[start]) {
		start++
	}
	if start == n {
		return 0
	}
	switch s[start] {
	case '*', '?', '<', '>', '\\', '/', ':', '"', '|':
		s[start] = '-'
	}
	// Copy the full suffix first: bytes beyond the new terminator are observable.
	copy(s, s[start:n+1])
	end := n - start
	for end > 0 && space(s[end-1]) {
		end--
	}
	s[end] = 0
	return 1
}
func characterDefaultName() *uint16 {
	return alloc.InternCString16(GetServer().S().Strings().GetStringInFile("DefaultName", "SelColor.c"))
}
func characterCopyName(dst *uint16) {
	src := (*uint16)(unsafe.Pointer(uintptr(optionsSend(characterUI.controls[14], 16413, 0, 0))))
	if *src == 0 {
		characterCopyUTF16(src, characterDefaultName())
	}
	characterCopyUTF16(dst, src)
	if characterName(dst) == 0 {
		characterCopyUTF16(dst, characterDefaultName())
	}
}
func characterCopyUTF16(dst, src *uint16) {
	n := alloc.StrLen(src) + 1
	copy(unsafe.Slice(dst, n), unsafe.Slice(src, n))
}
func characterPaletteMatch(w *gui.Window, bank int, rgb *byte) uintptr {
	c := unsafe.Slice(rgb, 3)
	idx := 0
	var ret uintptr
	for ; idx < 32; idx++ {
		p := memmap.PtrUint8(0x5D4594, 1307796+uintptr(96*bank+3*idx))
		v := unsafe.Slice(p, 3)
		if v[0] == c[0] && v[1] == c[1] && v[2] == c[2] {
			break
		}
	}
	ret = uintptr(memmap.PtrOff(0x5D4594, 1307797+uintptr(96*bank+3*idx)))
	if bank == 1 {
		if idx == 32 {
			ret = uintptr(uiWindowEnable(characterUI.colorRoot.ChildByID(w.ID()-10), 0))
			idx = 9
		} else {
			uiWindowEnable(characterUI.colorRoot.ChildByID(w.ID()-10), 1)
			other := characterUI.colorRoot.ChildByID(w.ID() + 10)
			other.DrawData().Field0 |= 6
			ret = uintptr(other.C())
		}
	}
	*characterWord(w) = uint32(uint16(bank)) | uint32(uint16(idx))<<16
	return ret
}
func characterPaletteFill(bank uint16) *gui.Window {
	*memmap.PtrUint32(0x5D4594, 1307788) = uint32(bank)
	var w *gui.Window
	for id := uint(761); id <= 792; id++ {
		w = characterUI.palette.ChildByID(id)
		if w != nil {
			*characterWord(w) = uint32(bank) | uint32(id-761)<<16
		}
	}
	return w
}
func characterPaletteClose(index uint16) uintptr {
	characterUI.palette.StackPop()
	characterUI.palette.Hide()
	if index >= 32 {
		return 0
	}
	w := characterUI.palette.ChildByID(uint(characterUI.target))
	if w != nil {
		*characterWord(w) = uint32(memmap.Uint16(0x5D4594, 1307788)) | uint32(index)<<16
	}
	return uintptr(w.C())
}
func characterPaletteDraw(w *gui.Window) int {
	p := w.GlobalPos()
	r := GetClient().R2()
	cl := characterColor555(w)
	r.Data().SetColor2(cl)
	r.DrawRectFilledOpaque(p.X, p.Y, w.EndPos.X-w.Off.X, w.EndPos.Y-w.Off.Y, cl)
	return 1
}
func characterAppearance() unsafe.Pointer {
	characterCopyName((*uint16)(characterUI.colorHost))
	out := unsafe.Slice((*byte)(characterUI.colorHost), 128)
	skin := characterColor(characterUI.controls[0])
	copy(out[71:74], skin)
	for i, off := range []int{68, 74, 77, 80} {
		col := skin
		if characterUI.controls[10+i].Flags&8 != 0 {
			col = characterColor(characterUI.controls[1+i])
		}
		copy(out[off:off+3], col)
	}
	for i := 0; i < 5; i++ {
		out[83+i] = byte(*characterWord(characterUI.controls[5+i]) >> 16)
	}
	return characterUI.colorHost
}
func characterSetup() uint32 {
	Sub_4A5E90_A()
	for i := 0; i < 10; i++ {
		characterUI.controls[i] = characterUI.colorRoot.ChildByID(uint(720 + i))
	}
	for i := 0; i < 4; i++ {
		characterUI.controls[10+i] = characterUI.colorRoot.ChildByID(uint(711 + i))
	}
	characterUI.controls[14] = characterUI.colorRoot.ChildByID(751)
	*characterWord(characterUI.controls[0]) = 131074
	for i := 1; i < 5; i++ {
		*characterWord(characterUI.controls[i]) = 589825
	}
	for i := 0; i < 5; i++ {
		*characterWord(characterUI.controls[5+i]) = uint32(memmap.Uint16(0x587000, 171372+uintptr(2*i))) << 16
	}
	if characterUI.defaults != 0 {
		return characterUI.defaults
	}
	optionsSend(characterUI.controls[14], 16414, uintptr(characterUI.colorHost), 0)
	host := unsafe.Slice((*byte)(characterUI.colorHost), 128)
	characterPaletteMatch(characterUI.controls[0], 2, &host[71])
	for i, off := range []int{68, 74, 77, 80} {
		characterPaletteMatch(characterUI.controls[i+1], 1, &host[off])
	}
	var result uint32
	for i := 0; i < 5; i++ {
		result = uint32(host[83+i]) << 16
		*characterWord(characterUI.controls[5+i]) = result
	}
	return result
}
func characterQuickbar(mode int) uintptr {
	class := *(*byte)(unsafe.Add(characterUI.classHost, 66))
	table := *(*unsafe.Pointer)(memmap.PtrOff(0x587000, 170156+4*uintptr(class)))
	var count byte
	result := uintptr(class)
	for *(*byte)(unsafe.Add(table, 5*uintptr(count))) != 0 {
		count++
		result = uintptr(count)
		if count == 0 {
			break
		}
	}
	if count == 0 {
		return result
	}
	row := byte(GetServer().S().Rand.Other.Int(0, int(count)-1))
	rows, kind := 1, uintptr(3)
	if class != 0 {
		rows, kind = 5, 2
	}
	for i := 0; i < rows; i++ {
		quickbarSelectRow(i)
		for j := 0; j < 5; j++ {
			var id byte
			if mode != 1 {
				id = *(*byte)(unsafe.Add(table, 5*uintptr(row)+uintptr(j)))
			}
			result = quickbarBookSlot(kind, uint32(id), j)
		}
	}
	if class != 0 {
		result = uintptr(quickbarSelectRow(0))
	}
	return result
}
func characterClassEvent(_ *gui.Window, ev gui.WindowEvent) gui.WindowEventResp {
	// Construction notifications carry a numeric child ID, not a window pointer.
	if ev.EventCode() != 16389 && ev.EventCode() != 16391 {
		return gui.RawEventResp(0)
	}
	a, _ := ev.EventArgsC()
	child := (*gui.Window)(unsafe.Pointer(a))
	id := int(child.ID())
	switch ev.EventCode() {
	case 16389:
		if id >= 601 && id <= 603 {
			uiWindowEnable(bookWindow(memmap.Uint32(0x5D4594, 1307728)), 1)
			cl := byte(id - 89)
			*(*byte)(unsafe.Add(characterUI.classHost, 66)) = cl
			key := alloc.GoString(*(**byte)(memmap.PtrOff(0x587000, 170208+4*uintptr(cl))))
			label := GetServer().S().Strings().GetStringInFile(strman.ID(key), "SelClass.c")
			optionsSend(characterUI.classRoot.ChildByID(605), 16385, uintptr(unsafe.Pointer(alloc.InternCString16(label))), 0)
			*memmap.PtrUint32(0x5D4594, 1307740) = uint32(id)
		}
		Nox_xxx_clientPlaySoundSpecial_452D80(920, 100)
	case 16391:
		if id >= 601 && id <= 603 {
			return gui.RawEventResp(1)
		}
		if id == 610 {
			if flags.HasGame(0x2000) && !flags.HasGame(4096) {
				mode := 0
				if questRuntimeWord(1556160) != 0 || questRuntimeWord(1556164) != 0 {
					mode = 1
				}
				characterQuickbar(mode)
				characterClassStart()
			} else {
				characterClassStart()
			}
			characterClassNextColor()
		}
		Nox_xxx_clientPlaySoundSpecial_452D80(921, 100)
	default:
		return gui.RawEventResp(0)
	}
	return gui.RawEventResp(1)
}
func characterClassDraw(w *gui.Window) int {
	if memmap.Uint32(0x5D4594, 1307740) != uint32(w.ID()) {
		p := w.GlobalPos()
		GetClient().R2().DrawRectFilledAlpha(p.X, p.Y, w.EndPos.X-w.Off.X, w.EndPos.Y-w.Off.Y)
	}
	return 1
}
func characterClassStart() int {
	characterUI.classAnim.SetState(gui.AnimOut)
	gui.SetAnimGlobalState(gui.AnimOut)
	Nox_xxx_clientPlaySoundSpecial_452D80(923, 100)
	return 1
}
func characterColorStart() int {
	characterUI.colorAnim.SetState(gui.AnimOut)
	gui.SetAnimGlobalState(gui.AnimOut)
	Nox_xxx_clientPlaySoundSpecial_452D80(923, 100)
	characterAppearance()
	return 1
}
func characterClassDone() int {
	fn := characterUI.classAnim.Func13Ptr
	characterUI.classAnim.Free()
	characterUI.classRoot.Destroy()
	gui.CallAnimationCallback(fn)
	return 1
}
func characterColorDone() int {
	fn := characterUI.colorAnim.Func13Ptr
	characterUI.colorAnim.Free()
	characterUI.colorRoot.Destroy()
	characterUI.palette.Destroy()
	if fn != nil {
		gui.CallAnimationCallback(fn)
	} else {
		Nox_client_resetScreenParticles_431510()
		GetClient().Cli().GUI.Draw()
		if !flags.HasGame(0x2000) {
			Nox_client_guiXxxDestroy_4A24A0()
		} else {
			Nox_xxx_serverHost_43B4D0()
		}
	}
	return 1
}
func characterPaletteOutside(w *gui.Window, ev gui.WindowEvent) gui.WindowEventResp {
	if ev.EventCode() != 5 {
		return gui.RawEventResp(0)
	}
	a, _ := ev.EventArgsC()
	p := image.Pt(int(uint16(a)), int(uint32(a)>>16))
	if p.X >= w.Off.X && p.X <= w.EndPos.X && p.Y >= w.Off.Y && p.Y <= w.EndPos.Y {
		return gui.RawEventResp(0)
	}
	characterPaletteClose(0xdead)
	return gui.RawEventResp(1)
}
func characterPaletteShow(x, y int) int {
	w := characterUI.palette
	w.ShowModal()
	w.StackPush()
	w.SetPos(image.Pt(x-w.SizeVal.X, y-w.SizeVal.Y/2))
	return 0
}
func characterSetMap(class byte) {
	if int(class) < 3 {
		sessionSetMapPath(alloc.InternCString([...]string{"war01a.map", "wiz01a.map", "con01a.map"}[class]))
	}
}
func characterColorEvent(_ *gui.Window, ev gui.WindowEvent) gui.WindowEventResp {
	a, b := ev.EventArgsC()
	switch ev.EventCode() {
	case 16389:
		Nox_xxx_clientPlaySoundSpecial_452D80(920, 100)
	case 16391:
		id := int((*gui.Window)(unsafe.Pointer(a)).ID())
		switch {
		case id >= 720 && id <= 729:
			bank := uint16(0)
			if id == 720 {
				bank = 2
			} else if id <= 724 {
				bank = 1
			}
			characterUI.target = uint32(id)
			characterPaletteFill(bank)
			characterPaletteShow(int(uint16(b)), int(uint32(b)>>16))
		case id >= 731 && id <= 734:
			w := characterUI.colorRoot.ChildByID(uint(id - 20))
			if w != nil {
				uiWindowEnable(w, int((^uint32(w.Flags)>>3)&1))
			}
		case id >= 761 && id <= 792:
			characterPaletteClose(uint16(id - 761))
		case id == 799:
			if memmap.Uint32(0x5D4594, 1308168) == 1 {
				nox_game_decStateInd_43BDC0()
			}
			nox_game_decStateInd_43BDC0()
			nox_game_decStateInd_43BDC0()
			characterUI.defaults = 1
			if characterCreateFile() != 0 {
				characterSetMap(*(*byte)(unsafe.Add(characterUI.colorHost, 66)))
				Sub_4A24C0(0)
				characterColorStart()
				characterUI.colorAnim.Func13Ptr = nil
			}
		}
		Nox_xxx_clientPlaySoundSpecial_452D80(921, 100)
	default:
		return gui.RawEventResp(0)
	}
	return gui.RawEventResp(1)
}
func characterPreview(w *gui.Window) int {
	r := GetClient().R2()
	d := r.Data()
	p := w.GlobalPos()
	skin := characterColor555(characterUI.controls[0])
	d.SetMaterial(6, skin)
	d.SetMaterial(1, characterColor555(characterUI.controls[1]))
	for _, v := range [][3]int{{12, 3, 2}, {13, 4, 4}, {11, 2, 5}} {
		c := skin
		if characterUI.controls[v[0]].Flags&8 != 0 {
			c = characterColor555(characterUI.controls[v[1]])
		}
		d.SetMaterial(v[2], c)
	}
	d.SetMaterial(3, skin)
	sex := uintptr(*(*byte)(unsafe.Add(characterUI.colorHost, 67)))
	off := uintptr(24)
	if characterUI.controls[10].Flags&8 != 0 {
		off = 16
	}
	bookDrawImage(memmap.Uint32(0x973A20, off+4*sex), p)
	// Only the shoes, pants and shirt layers are part of this preview.
	for _, layer := range []struct {
		bit, mod      int
		colors, slots []int
	}{
		{0, 2, []int{8, 9}, []int{40, 36}},
		{2, 0, []int{5}, []int{40}},
		{10, 1, []int{6, 7}, []int{40, 44}},
	} {
		mod := characterUI.modifiers[layer.mod]
		for i := 1; i <= 6; i++ {
			c := unsafe.Slice((*byte)(unsafe.Add(mod, 12+3*i)), 3)
			d.SetMaterialRGB(i, int(c[0]), int(c[1]), int(c[2]))
		}
		for i, control := range layer.colors {
			slot := *(*int32)(unsafe.Add(mod, layer.slots[i]))
			d.SetMaterial(int(slot), characterColor555(characterUI.controls[control]))
		}
		bookDrawImage(memmap.Uint32(0x973A20, 32+104*sex+4*uintptr(layer.bit)), p)
	}
	return 1
}
func characterCreateFile() int {
	campaign := flags.HasGame(2048)
	if campaign {
		_ = Nox_savegame_rm("WORKING", false)
	}
	// All defined fields are initialized; zero the two unused tail padding bytes too.
	rec, free := alloc.Make([]byte{}, 1280)
	defer free()
	characterCopyName((*uint16)(unsafe.Pointer(&rec[1224])))
	rec[1274] = *(*byte)(unsafe.Add(characterUI.colorHost, 66))
	rec[1276] = 1
	skin := characterColor(characterUI.controls[0])
	copy(rec[1204:1207], skin)
	for i := 0; i < 4; i++ {
		c := skin
		if characterUI.controls[10+i].Flags&8 != 0 {
			c = characterColor(characterUI.controls[1+i])
		}
		copy(rec[1207+3*i:1210+3*i], c)
	}
	for i := 0; i < 5; i++ {
		rec[1219+i] = byte(*characterWord(characterUI.controls[5+i]) >> 16)
	}
	path := saveRoot() + alloc.GoString(memmap.PtrUint8(0x587000, 171764))
	if campaign {
		path += alloc.GoString(memmap.PtrUint8(0x587000, 171772)) + alloc.GoString(memmap.PtrUint8(0x587000, 171780))
	}
	_ = ifs.Mkdir(path)
	_ = ifs.Chdir(path)
	name := "Player.plr"
	slot := 0
	if !campaign {
		// The original formats UTF-16 storage as a byte string, capped at six bytes.
		prefix := alloc.GoStringS(rec[1224:1274])
		if len(prefix) > 6 {
			prefix = prefix[:6]
		}
		for ; slot < 100; slot++ {
			name = fmt.Sprintf("%s%02d.plr", prefix, slot)
			f, err := ifs.Open(name)
			if err != nil {
				break
			}
			f.Close()
		}
	}
	_ = ifs.Chdir(datapath.Data())
	if slot == 100 {
		return 0
	}
	alloc.StrCopy(rec[4:1028], path+name)
	if campaign {
		characterSetMap(rec[1274])
	}
	return playerFileClientWrite(unsafe.Pointer(&rec[0]), 1)
}
