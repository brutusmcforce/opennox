package mp3

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"math"
	"os"
	"testing"
)

func TestStreamFrozenVectors(t *testing.T) {
	checkStreamInvariants(t)
	f, err := os.Open("testdata/stream_c.bin.gz")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	gz, err := gzip.NewReader(f)
	if err != nil {
		t.Fatal(err)
	}
	defer gz.Close()
	hash := sha256.New()
	r := bufio.NewReader(io.TeeReader(gz, hash))
	var magic [8]byte
	readFixture(t, r, magic[:], -1)
	if string(magic[:]) != "NMP3STR1" {
		t.Fatal("invalid stream fixture magic")
	}
	var counts [5]int
	for record := 0; ; record++ {
		op, err := r.ReadByte()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		u := func() uint32 { return readU32(t, r, record) }
		switch op {
		case 1, 2:
			n, hint := u(), u()
			if n > 49152 || hint > 2304 {
				t.Fatalf("record%d: unsupported scan bounds", record)
			}
			packet := make([]byte, n)
			readFixture(t, r, packet, record)
			if op == 1 {
				free, frame := int32(hint), int32(0x12345678)
				offset := findFrame(packet, int32(n), &free, &frame)
				want := [3]int32{int32(u()), int32(u()), int32(u())}
				got := [3]int32{offset, free, frame}
				if got != want {
					t.Fatalf("record%d finder: got%v want%v n%d hint%d", record, got, want, n, hint)
				}
			} else {
				if n < 4 || !headerValid(headerAt(packet, 0)) {
					t.Fatalf("record%d: unsupported match header", record)
				}
				got := streamBool(matchFrame(packet, int32(n), int32(hint)))
				want := u()
				if got != want {
					t.Fatalf("record%d matcher: got%d want%d", record, got, want)
				}
			}
		case 3:
			seed := u()
			if seed > 255 {
				t.Fatal("invalid init seed")
			}
			dec := seedStreamDecoder(byte(seed))
			before := dec
			decoderInit(&dec)
			before.header[0] = 0
			wantHeader, wantUnchanged := u(), u()
			if uint32(dec.header[0]) != wantHeader || streamBool(sameStreamDecoder(&dec, &before)) != wantUnchanged {
				t.Fatalf("record%d init changed extra state", record)
			}
		case 4:
			have, begin, length, pos, seed := u(), u(), u(), u(), u()
			if have > 511 || begin > 511 || length > 2304 || pos > length*8 || seed > 255 {
				t.Fatal("invalid restore fixture")
			}
			dec := seedStreamDecoder(byte(seed))
			dec.reserv = int32(have)
			streamPattern(dec.reservBuf[:], seed)
			before := dec
			var s scratchState
			for i := range s.mainData {
				s.mainData[i] = byte(seed ^ 255)
			}
			s.bs = bitReader{buf: []byte{1}, pos: 17, limit: 8}
			data := make([]byte, length)
			streamPattern(data, seed+19)
			dataBefore := append([]byte(nil), data...)
			var bs bitReader
			bsInit(&bs, data, len(data))
			bs.pos = int(pos)
			oldPos, oldLimit := bs.pos, bs.limit
			ret := restoreReservoir(&dec, &bs, &s, int32(begin))
			alias := cap(s.bs.buf) > 0 && &s.bs.buf[:1][0] == &s.mainData[0]
			outerSame := bs.pos == oldPos && bs.limit == oldLimit && len(bs.buf) == len(data)
			if len(data) > 0 {
				outerSame = outerSame && &bs.buf[0] == &data[0]
			}
			got := [6]uint32{streamBool(ret), uint32(s.bs.pos), uint32(s.bs.limit), streamBool(alias), streamBool(sameStreamDecoder(&dec, &before)), streamBool(outerSame)}
			want := [6]uint32{u(), u(), u(), u(), u(), u()}
			var wantData [2815]byte
			readFixture(t, r, wantData[:], record)
			if got != want || s.mainData != wantData || !bytes.Equal(data, dataBefore) {
				t.Fatalf("record%d restore: got%v want%v; scratch match%v", record, got, want, s.mainData == wantData)
			}
		case 5:
			length, pos, seed := u(), u(), u()
			if length > 2815 || pos > length*8+32 || seed > 255 {
				t.Fatal("invalid save fixture")
			}
			dec := seedStreamDecoder(byte(seed))
			streamPattern(dec.reservBuf[:], seed)
			before := dec
			var s scratchState
			streamPattern(s.mainData[:], seed+19)
			bsInit(&s.bs, s.mainData[:], int(length))
			s.bs.pos = int(pos)
			beforeData := s.mainData
			saveReservoir(&dec, &s)
			wantReserv := int32(u())
			var wantBuf [511]byte
			readFixture(t, r, wantBuf[:], record)
			wantUnchanged := u()
			unchanged := s.mainData == beforeData && s.bs.pos == int(pos) && s.bs.limit == int(length)*8 && cap(s.bs.buf) > 0 && &s.bs.buf[:1][0] == &s.mainData[0]
			before.reserv = dec.reserv
			before.reservBuf = dec.reservBuf
			if dec.reserv != wantReserv || dec.reservBuf != wantBuf || streamBool(unchanged) != wantUnchanged || !sameStreamDecoder(&dec, &before) {
				t.Fatalf("record%d save mismatch: reservoir%d want%d", record, dec.reserv, wantReserv)
			}
		default:
			t.Fatalf("record%d unknown op%d", record, op)
		}
		counts[op-1]++
	}
	if counts != [5]int{1650, 190, 5, 425, 504} {
		t.Fatalf("stream fixture counts%v", counts)
	}
	const wantSHA = "893c6b12f0b8d2f5ea8cfc04eef864d92e7976cc3b92e216f6472f4590cfe515"
	if hex.EncodeToString(hash.Sum(nil)) != wantSHA {
		t.Fatal("stream fixture hash mismatch")
	}
	t.Logf("validated2774 frozen stream-state vectors; SHA256%s", wantSHA)
}

func streamBool(v bool) uint32 {
	if v {
		return 1
	}
	return 0
}
func streamPattern(b []byte, seed uint32) {
	for i := range b {
		b[i] = byte((i*73 + int(seed)) ^ (i >> 2))
	}
}
func seedStreamDecoder(seed byte) decoderState {
	var d decoderState
	word := uint32(seed) * 0x01010101
	v := math.Float32frombits(word)
	for i := range d.mdctOverlap {
		for j := range d.mdctOverlap[i] {
			d.mdctOverlap[i][j] = v
		}
	}
	for i := range d.qmfState {
		d.qmfState[i] = v
	}
	d.reserv = int32(word)
	d.freeFormatBytes = int32(word)
	for i := range d.header {
		d.header[i] = seed
	}
	for i := range d.reservBuf {
		d.reservBuf[i] = seed
	}
	return d
}
func sameStreamDecoder(a, b *decoderState) bool {
	if a.reserv != b.reserv || a.freeFormatBytes != b.freeFormatBytes || a.header != b.header || a.reservBuf != b.reservBuf {
		return false
	}
	for i := range a.mdctOverlap {
		for j := range a.mdctOverlap[i] {
			if math.Float32bits(a.mdctOverlap[i][j]) != math.Float32bits(b.mdctOverlap[i][j]) {
				return false
			}
		}
	}
	for i := range a.qmfState {
		if math.Float32bits(a.qmfState[i]) != math.Float32bits(b.qmfState[i]) {
			return false
		}
	}
	return true
}

func checkStreamInvariants(t *testing.T) {
	t.Helper()
	for n := 0; n < 20; n++ {
		free, frame := int32(17), int32(-1)
		offset := findFrame(make([]byte, n), int32(n), &free, &frame)
		want := n - 4
		if want < 0 {
			want = 0
		}
		if offset != int32(want) || frame != 0 || free != 17 {
			t.Fatalf("no-header scan n%d got%d/%d/%d", n, offset, free, frame)
		}
	}
	packet := make([]byte, 417)
	copy(packet, []byte{255, 251, 144, 0})
	free, frame := int32(0), int32(-1)
	if off := findFrame(packet, 417, &free, &frame); off != 0 || frame != 417 {
		t.Fatalf("single exact frame not accepted: %d/%d", off, frame)
	}
	if matchFrame(packet, 417, 0) {
		t.Fatal("single frame cannot establish following-frame sync")
	}
	// Init must preserve NaN payloads and all persistent state except header byte0.
	dec := seedStreamDecoder(255)
	before := dec
	before.header[0] = 0
	decoderInit(&dec)
	if !sameStreamDecoder(&dec, &before) {
		t.Fatal("init changed persistent state beyond first header byte")
	}
	dec = decoderState{reserv: 3}
	copy(dec.reservBuf[:], []byte{10, 11, 12})
	var s scratchState
	for i := range s.mainData {
		s.mainData[i] = 0xa5
	}
	var bs bitReader
	bsInit(&bs, []byte{20, 21, 22}, 3)
	bs.pos = 8
	if !restoreReservoir(&dec, &bs, &s, 2) || !bytes.Equal(s.mainData[:4], []byte{11, 12, 21, 22}) || s.bs.pos != 0 || s.bs.limit != 32 || s.mainData[4] != 0xa5 || bs.pos != 8 {
		t.Fatal("reservoir concatenation/order/tail contract")
	}
	if restoreReservoir(&dec, &bs, &s, 4) || !bytes.Equal(s.mainData[:5], []byte{10, 11, 12, 21, 22}) || s.bs.limit != 40 {
		t.Fatal("insufficient history must still initialize copied data before false")
	}
	// Ceil bit rounding discards a partly consumed byte; destination tail stays put.
	s = scratchState{}
	copy(s.mainData[:], []byte{30, 31, 32})
	bsInit(&s.bs, s.mainData[:], 3)
	s.bs.pos = 1
	dec = seedStreamDecoder(0xa5)
	saveReservoir(&dec, &s)
	if dec.reserv != 2 || dec.reservBuf[0] != 31 || dec.reservBuf[1] != 32 || dec.reservBuf[2] != 0xa5 {
		t.Fatal("save byte rounding/tail contract")
	}
	for i := range s.mainData {
		s.mainData[i] = byte(i)
	}
	bsInit(&s.bs, s.mainData[:], 513)
	saveReservoir(&dec, &s)
	if dec.reserv != 511 || !bytes.Equal(dec.reservBuf[:], s.mainData[2:513]) {
		t.Fatal("save must retain last511 bytes")
	}
	s.bs.pos = s.bs.limit + 1
	prior := dec.reservBuf
	saveReservoir(&dec, &s)
	if dec.reserv != -1 || dec.reservBuf != prior {
		t.Fatal("negative save remainder changed buffer or stored length")
	}
}
