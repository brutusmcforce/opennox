package server

import (
	"unsafe"

	"github.com/brutusmcforce/opennox/legacy/common/alloc"
)

// ReadUseData is the use data of readable items (ReadUse, WarpReadUse).
type ReadUseData struct {
	TextBuf  [256]byte // 0, 0: string ID of the text shown to the reader (ReadUse)
	LastRead uint32    // 64, 256: frame of the last read, used to throttle repeated reads
}

var _ = [1]struct{}{}[260-unsafe.Sizeof(ReadUseData{})]

func (d *ReadUseData) UseDataPtr() unsafe.Pointer {
	return unsafe.Pointer(d)
}

func (d *ReadUseData) Text() string {
	return alloc.GoStringS(d.TextBuf[:])
}

func (p UseDataPtr) AsRead() *ReadUseData {
	return useDataPtrAs[ReadUseData](p)
}

// UseDataRead returns the use data of readable items (ReadUse, WarpReadUse).
func (obj *Object) UseDataRead() *ReadUseData {
	// TODO: verify this conversion by checking ObjectType
	return obj.UseData.AsRead()
}

// AmmoUseData is the use data of quivers and shuriken (AmmoUse).
type AmmoUseData struct {
	MaxCharges byte // 0, 0
	Charges    byte // 0, 1
	Infinite   byte // 0, 2: nonzero: shots do not consume charges (and 1 allows shooting with no charges)
}

var _ = [1]struct{}{}[3-unsafe.Sizeof(AmmoUseData{})]

func (d *AmmoUseData) UseDataPtr() unsafe.Pointer {
	return unsafe.Pointer(d)
}

func (p UseDataPtr) AsAmmo() *AmmoUseData {
	return useDataPtrAs[AmmoUseData](p)
}

// UseDataAmmo returns the use data of quivers and shuriken (AmmoUse).
func (obj *Object) UseDataAmmo() *AmmoUseData {
	// TODO: verify this conversion by checking ObjectType
	return obj.UseData.AsAmmo()
}

// BowUseData is the use data of bows and crossbows (BowUse).
type BowUseData struct {
	Reload byte // 0, 0: reload countdown: while nonzero, an attack only decrements it
}

var _ = [1]struct{}{}[1-unsafe.Sizeof(BowUseData{})]

func (d *BowUseData) UseDataPtr() unsafe.Pointer {
	return unsafe.Pointer(d)
}

func (p UseDataPtr) AsBow() *BowUseData {
	return useDataPtrAs[BowUseData](p)
}

// UseDataBow returns the use data of bows and crossbows (BowUse).
func (obj *Object) UseDataBow() *BowUseData {
	// TODO: verify this conversion by checking ObjectType
	return obj.UseData.AsBow()
}
