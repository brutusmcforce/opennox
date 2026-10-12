package legacy

import (
	"github.com/brutusmcforce/opennox/internal/cryptfile"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"unsafe"
)

func mapPolygonSection(bypass int) int {
	if bypass != 0 {
		return 1
	}
	r := objectXferStream{cryptfile.Global()}
	version := int16(r.short(4))
	if version > 4 {
		return 0
	}
	if r.read() {
		count := r.word(0)
		for i := uint32(1); i <= count; i++ {
			old := r.word(0)
			var xy [2]float32
			r.raw(unsafe.Pointer(&xy), 8)
			if mapPolygonVertexSet(xy[0], xy[1], i, old) == nil {
				return 0
			}
			if version < 3 {
				r.word(0)
			}
		}
		count = r.word(0)
		var secrets uint32
		for i := uint32(0); i < count; i++ {
			p := mapPolygonNew()
			if p == nil {
				return 0
			}
			n := r.byte(0)
			r.raw(unsafe.Pointer(&p.Name[0]), int(n))
			*(*byte)(unsafe.Add(unsafe.Pointer(&p.Name[0]), uintptr(n))) = 0
			if version < 3 {
				r.word(0)
				for j := range p.Color {
					p.Color[j] = byte(r.word(0))
				}
			} else {
				for j := range p.Color {
					p.Color[j] = r.byte(0)
				}
			}
			p.Level = r.byte(p.Level)
			p.Count = r.short(p.Count)
			p.Vertices = (*uint32)(legacyCalloc(uintptr(p.Count), 4))
			if p.Vertices == nil {
				return 0
			}
			r.raw(unsafe.Pointer(p.Vertices), 4*int(p.Count))
			mapPolygonRemapIDs(p)
			mapPolygonBounds(p)
			if version >= 2 {
				mapPolygonScriptTransfer(p)
			}
			if version >= 4 {
				p.Flags = r.word(p.Flags)
				if p.Flags&1 != 0 {
					secrets++
				}
			}
		}
		questRuntimePreviousStage(secrets)
		return 1
	}
	var count uint32
	for p := mapPolygonVertexFirst(); p != nil; p = mapPolygonVertexAfter(p.ID) {
		count++
	}
	r.word(count)
	for p := mapPolygonVertexFirst(); p != nil; p = mapPolygonVertexAfter(p.ID) {
		r.raw(unsafe.Pointer(p), 12)
	}
	count = 0
	for p := mapPolygonFirst(); p != nil; p = mapPolygonAfter(p.ID) {
		count++
	}
	r.word(count)
	for p := mapPolygonFirst(); p != nil; p = mapPolygonAfter(p.ID) {
		n := r.byte(byte(len(alloc.GoString(&p.Name[0]))))
		r.raw(unsafe.Pointer(&p.Name[0]), int(n))
		for _, v := range p.Color {
			r.byte(v)
		}
		r.byte(p.Level)
		r.short(p.Count)
		r.raw(unsafe.Pointer(p.Vertices), 4*int(p.Count))
		mapPolygonScriptTransfer(p)
		r.word(p.Flags)
	}
	return 1
}
func mapPolygonScriptTransfer(p *mapPolygon) {
	objectXferScript(unsafe.Pointer(&p.Enter), p.Metadata)
	var leaveName unsafe.Pointer
	if p.Metadata != nil {
		leaveName = unsafe.Add(p.Metadata, 128)
	}
	objectXferScript(unsafe.Pointer(&p.Leave), leaveName)
}
