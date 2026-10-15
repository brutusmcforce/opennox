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
	data, err := io.ReadAll(r)
	_ = r.Close()
	if err != nil {
		return err
	}
	out, err := Compress(data)
	if err != nil {
		return err
	}
	w, err := ifs.Create(dst)
	if err != nil {
		return err
	}
	defer func() { _ = w.Close() }()
	if _, err := w.Write(out); err != nil {
		return err
	}
	return w.Close()
}

// Compress returns the NXZ file for the data of a map file.
func Compress(data []byte) ([]byte, error) {
	if len(data) > int(^uint(0)>>1)-5 || uint64(len(data)) > uint64(^uint32(0)) {
		return nil, errors.New("nxz: file exceeds addressable size")
	}
	srcSz := len(data)
	// Block-end rolling hashes inspect five following bytes. Preserve actual
	// following file bytes and bound the final lookahead with zero padding.
	sbuf := make([]byte, srcSz+5)
	copy(sbuf, data)
	encoder := newMapEncoder()
	out := binary.LittleEndian.AppendUint32(nil, uint32(srcSz))
	for i := 0; i < srcSz; {
		n := min(srcSz-i, 500000)
		out = append(out, encoder.block(sbuf[i:], n)...)
		// Advancing by the final partial chunk cannot overflow the source size.
		i += n
	}
	return out, nil
}
