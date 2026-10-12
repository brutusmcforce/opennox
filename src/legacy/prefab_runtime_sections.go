package legacy

import (
	"bytes"
	"encoding/binary"
	"io"
	"math"
	"path/filepath"
	"strings"
	"unsafe"

	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/internal/cryptfile"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"github.com/brutusmcforce/opennox/libs/ifs"
	"github.com/brutusmcforce/opennox/libs/types"
	"github.com/brutusmcforce/opennox/server"
)

type prefabWire struct{ bad bool }

func (w *prefabWire) data(b []byte) {
	n, e := cryptfile.Global().ReadWrite(b)
	if e != nil || n != len(b) {
		w.bad = true
	}
}
func (w *prefabWire) byte(v byte) byte { b := [1]byte{v}; w.data(b[:]); return b[0] }
func (w *prefabWire) short(v uint16) uint16 {
	var b [2]byte
	binary.LittleEndian.PutUint16(b[:], v)
	w.data(b[:])
	return binary.LittleEndian.Uint16(b[:])
}
func (w *prefabWire) word(v uint32) uint32 {
	var b [4]byte
	binary.LittleEndian.PutUint32(b[:], v)
	w.data(b[:])
	return binary.LittleEndian.Uint32(b[:])
}
func (w *prefabWire) result() uint32 { return uint32(bool2int(!w.bad)) }
func prefabByteString(b []byte) string {
	if i := bytes.IndexByte(b, 0); i >= 0 {
		b = b[:i]
	}
	return string(b)
}
func prefabGroupSection() uint32 {
	w := prefabWire{}
	version := int16(w.short(3))
	if w.bad || version > 3 {
		return 0
	}
	groups := &GetServer().S().MapGroups
	blocked := memmap.Uint8(0x5D4594, 739992)&4 != 0
	if cryptfile.Global().ReadOnly() {
		count := int32(w.word(0))
		if w.bad {
			return 0
		}
		for i := int32(0); i < count; i++ {
			length := int(w.byte(0))
			if length > 75 || w.bad {
				return 0
			}
			b := make([]byte, length)
			w.data(b)
			name := prefabByteString(b)
			current := alloc.GoString((*byte)(memmap.PtrOff(0x5D4594, 2598188)))
			if version < 2 {
				if len(current)+1+len(name) >= 53 {
					return 0
				}
				name = current + ".map:" + name
			} else if version == 2 {
				fields := strings.FieldsFunc(name, func(c rune) bool { return c == ':' })
				if len(fields) < 2 || len(current)+1+len(fields[1]) >= 53 {
					return 0
				}
				name = current + ":" + fields[1]
			}
			kind, index := w.byte(0), w.word(0)
			if w.bad {
				return 0
			}
			if !blocked {
				if noxflags.HasGame(0x400000) {
					groups.Sub504600(name, index, kind)
				} else if noxflags.HasGame(0x200001) {
					groups.MapLoadAddGroup57C0C0(name, index, kind)
				}
			}
			countItems := int32(w.word(0))
			for j := int32(0); j < countItems; j++ {
				if kind > 3 {
					return 0
				}
				item := [2]uint32{w.word(0), 0}
				if kind == 2 {
					item[1] = w.word(0)
				}
				if w.bad {
					return 0
				}
				if !blocked {
					if noxflags.HasGame(0x400000) {
						groups.Sub5046A0(item[:], index)
					} else {
						groups.Sub57C130(item[:], index)
					}
				}
			}
		}
		return w.result()
	}
	count := uint32(0)
	for g := groups.GetFirstMapGroup(); g != nil; g = g.Next() {
		count++
	}
	w.word(count)
	for g := groups.GetFirstMapGroup(); g != nil; g = g.Next() {
		name := append([]byte(g.ID()), 0)
		w.byte(byte(len(name)))
		w.data(name)
		kind := byte(g.GroupType())
		w.byte(kind)
		w.word(g.Index())
		count = 0
		for p := g.List; p != nil; p = p.Next8 {
			count++
		}
		w.word(count)
		for p := g.List; p != nil; p = p.Next8 {
			if kind > 3 {
				return 0
			}
			w.word(p.Raw0)
			if kind == 2 {
				w.word(p.Raw4)
			}
		}
	}
	return w.result()
}
func prefabWaypointSection(bounds *[8]uint32) uint32 {
	w := prefabWire{}
	version := int16(w.short(4))
	if w.bad || version > 4 {
		return 0
	}
	s := GetServer().S()
	if cryptfile.Global().ReadOnly() {
		count := int32(w.word(0))
		if w.bad {
			return 0
		}
		for i := int32(0); i < count; i++ {
			index := w.word(0)
			xb, yb := w.word(0), w.word(0)
			x, y := math.Float32frombits(xb), math.Float32frombits(yb)
			if version < 4 {
				x, y = float32(xb), float32(yb)
			}
			var name []byte
			if version >= 3 {
				n := int(w.byte(0))
				if n > 75 {
					return 0
				}
				name = make([]byte, n)
				w.data(name)
			}
			if bounds != nil {
				var rect [4]int32
				geometryWallBounds(bounds, &rect)
				x = float32(float64(x) - float64(int32(23*memmap.Uint32(0x5D4594, 739980))) + float64(rect[0]) - 11)
				y = float32(float64(y) - float64(int32(23*memmap.Uint32(0x5D4594, 739984))) + float64(rect[1]) - 11)
			}
			if w.bad {
				return 0
			}
			var wp *server.Waypoint
			if noxflags.HasGame(0x400000) {
				n := prefabWaypointNew(index, x, y)
				if n == 0 {
					return 0
				}
				wp = prefabWaypoint(*prefabWord(n, 0))
			} else {
				wp = s.WPs.Nox_xxx_waypointNewNotMap_579970(int(index), types.Ptf(x, y))
			}
			if wp == nil {
				return 0
			}
			if version >= 3 {
				copy(wp.NameBuf[:], []byte(prefabByteString(name)))
			}
			wp.Flags = w.word(0)
			if version < 4 {
				wp.PointsCnt = byte(w.word(0))
			} else {
				wp.PointsCnt = w.byte(0)
			}
			if wp.PointsCnt > 32 {
				return 0
			}
			for j := byte(0); j < wp.PointsCnt; j++ {
				wp.Field348[j] = w.word(0)
				wp.Points[j].Ind = 2
				if version >= 2 {
					wp.Points[j].Ind = w.byte(0)
				}
			}
			if w.bad {
				return 0
			}
		}
		return w.result()
	}
	includes := func(wp *server.Waypoint) bool {
		if bounds == nil {
			return true
		}
		point := [2]int32{int32(int64(wp.PosVec.X)), int32(int64(wp.PosVec.Y))}
		return geometryWallPoint(&point, (*[8]int32)(unsafe.Pointer(bounds))) != 0
	}
	count := uint32(0)
	for p := s.WPs.List; p != nil; p = p.WpNext {
		if includes(p) {
			count++
		}
	}
	w.word(count)
	for p := s.WPs.List; p != nil; p = p.WpNext {
		if !includes(p) {
			continue
		}
		w.word(p.Index)
		w.word(math.Float32bits(p.PosVec.X))
		w.word(math.Float32bits(p.PosVec.Y))
		name := p.ID()
		w.byte(byte(len(name)))
		w.data([]byte(name))
		w.word(p.Flags & 1)
		w.byte(p.PointsCnt)
		if p.PointsCnt > 32 {
			return 0
		}
		for i := byte(0); i < p.PointsCnt; i++ {
			if p.Points[i].Waypoint == nil {
				return 0
			}
			w.word(p.Points[i].Waypoint.Index)
			w.byte(p.Points[i].Ind)
		}
	}
	return w.result()
}
func prefabIntroFree() uint32 {
	p := (*prefabGlobal(prefabIntro))
	if p != 0 {
		mapRoomRelease(mapRoomPointer(p))
		(*prefabGlobal(prefabIntro)) = 0
	}
	return p
}
func prefabIntroSection() uint32 {
	prefabIntroFree()
	w := prefabWire{}
	if int16(w.short(1)) < 1 || w.bad {
		return 0
	}
	name := Nox_xxx_mapGetMapName_409B40()
	path := filepath.Join(saveRoot(), "maps", name, name+".txt")
	editor := noxflags.HasGame(0x200000)
	if cryptfile.Global().ReadOnly() {
		length := int32(w.word(0))
		if w.bad {
			return 0
		}
		if length <= 0 {
			return 1
		}
		if noxflags.HasGame(0x400000) {
			return uint32(bool2int(cryptfile.Global().Seek(int64(length), io.SeekCurrent) == nil))
		}
		if editor {
			f, err := ifs.Create(path)
			if err != nil {
				return 0
			}
			defer f.Close()
			for i := int32(0); i < length; i++ {
				b := w.byte(0)
				if _, err := f.Write([]byte{b}); err != nil {
					return 0
				}
			}
		} else {
			p := mapRoomCalloc(1, uintptr(length))
			if p == nil {
				return 0
			}
			(*prefabGlobal(prefabIntro)) = mapRoomRaw(p)
			w.data(unsafe.Slice((*byte)(p), int(length)))
		}
		return w.result()
	}
	if editor {
		f, err := ifs.Open(path)
		if err == nil {
			defer f.Close()
			stat, err := f.Stat()
			if err != nil {
				return 0
			}
			size := uint32(stat.Size())
			w.word(size)
			var b [1]byte
			for i := int32(0); i < int32(size); i++ {
				if _, err := io.ReadFull(f, b[:]); err != nil {
					return 0
				}
				w.data(b[:])
			}
			return w.result()
		}
	}
	w.word(0)
	return w.result()
}
