package cnxz

import (
	"encoding/binary"
	"errors"
	"io"

	"github.com/brutusmcforce/opennox/libs/ifs"
)

func DecompressFile(src, dst string) error {
	if src == "" {
		return errors.New("empty source path")
	}
	if dst == "" {
		return errors.New("empty destination path")
	}
	r, err := ifs.Open(src)
	if err != nil {
		return err
	}
	defer r.Close()
	fi, err := r.Stat()
	if err != nil {
		return err
	}
	var buf [4]byte
	if _, err = io.ReadFull(r, buf[:]); err != nil {
		return err
	}
	dstSize := binary.LittleEndian.Uint32(buf[:])
	if uint64(dstSize) > uint64(^uint(0)>>1) || fi.Size()-4 > int64(^uint(0)>>1) {
		return errors.New("nxz: file exceeds addressable size")
	}
	sbuf := make([]byte, int(fi.Size()-4))
	if _, err = io.ReadFull(r, sbuf); err != nil {
		return err
	}
	dbuf := make([]byte, int(dstSize))
	if err := newMapDecoder(sbuf).decode(dbuf); err != nil {
		return err
	}

	w, err := ifs.Create(dst)
	if err != nil {
		return err
	}
	defer w.Close()
	if _, err = w.Write(dbuf); err != nil {
		return err
	}
	return w.Close()
}

func CompressFile(src, dst string) error {
	if src == "" {
		return errors.New("empty source path")
	}
	if dst == "" {
		return errors.New("empty destination path")
	}
	r, err := ifs.Open(src)
	if err != nil {
		return err
	}
	defer r.Close()
	fi, err := r.Stat()
	if err != nil {
		return err
	}

	if fi.Size() < 0 || fi.Size() > int64(^uint(0)>>1)-5 || fi.Size() > int64(^uint32(0)) {
		return errors.New("nxz: file exceeds addressable size")
	}
	srcSz := int(fi.Size())
	// Block-end rolling hashes inspect five following bytes. Preserve actual
	// following file bytes and bound the final lookahead with zero padding.
	sbuf := make([]byte, srcSz+5)
	if _, err = io.ReadFull(r, sbuf[:srcSz]); err != nil {
		return err
	}
	encoder := newMapEncoder()
	var dbuf []byte
	for i := 0; i < srcSz; {
		n := min(srcSz-i, 500000)
		dbuf = append(dbuf, encoder.block(sbuf[i:], n)...)
		// Advancing by the final partial chunk cannot overflow the source size.
		i += n
	}

	w, err := ifs.Create(dst)
	if err != nil {
		return err
	}
	defer w.Close()

	var buf [4]byte
	binary.LittleEndian.PutUint32(buf[:4], uint32(srcSz))
	if _, err := w.Write(buf[:4]); err != nil {
		return err
	}
	if _, err := w.Write(dbuf); err != nil {
		return err
	}
	return w.Close()
}
