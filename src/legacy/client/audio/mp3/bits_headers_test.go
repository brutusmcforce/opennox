package mp3

import (
	"bufio"
	"compress/gzip"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"io"
	"os"
	"testing"
)

const (
	integerFixturePath       = "testdata/integer_c.bin.gz"
	integerFixtureMagic      = "NMP3INT1"
	integerPatternBytes      = 2815 + 32
	wantIntegerFixtureSHA256 = "d247a1f9f026f1d9b39d57efa5bba5a3be014ee41a536b40e0252cc035285f37"
)

// Frozen after three identical actual-C captures and a matching UBSan capture.
var wantIntegerFixtureCounts = [4]uint64{196864, 524288, 30240, 67815}

func TestBitHeaderFrozenVectors(t *testing.T) {
	checkBitHeaderInvariants(t)
	f, err := os.Open(integerFixturePath)
	if err != nil {
		t.Fatalf("open frozen actual-C fixture %q: %v", integerFixturePath, err)
	}
	defer f.Close()
	gz, err := gzip.NewReader(f)
	if err != nil {
		t.Fatalf("open gzip fixture: %v", err)
	}
	defer gz.Close()

	fixtureHash := sha256.New()
	r := bufio.NewReader(io.TeeReader(gz, fixtureHash))
	var magic [8]byte
	if _, err := io.ReadFull(r, magic[:]); err != nil {
		t.Fatalf("read fixture magic: %v", err)
	}
	if string(magic[:]) != integerFixtureMagic {
		t.Fatalf("fixture magic %q, want %q", magic, integerFixtureMagic)
	}
	var patterns [5][]byte
	for i := range patterns {
		patterns[i] = integerPattern(i, integerPatternBytes)
	}
	var counts [4]uint64
	for record := 0; ; record++ {
		op, err := r.ReadByte()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("read operation for record %d: %v", record, err)
		}
		switch op {
		case 1:
			var h [4]byte
			readFixture(t, r, h[:], record)
			want := readU32(t, r, record)
			if want > 1 {
				t.Fatalf("record %d: op1 expected-valid must be 0 or 1, got %d", record, want)
			}
			got := uint32(0)
			if headerValid(h) {
				got = 1
			}
			if got != want {
				t.Fatalf("record %d op1 header=% x: valid=%d, want %d", record, h, got, want)
			}
			counts[0]++
		case 2:
			var h1, h2 [4]byte
			readFixture(t, r, h1[:], record)
			readFixture(t, r, h2[:], record)
			want := readU32(t, r, record)
			if want > 1 {
				t.Fatalf("record %d: op2 expected-compare must be 0 or 1, got %d", record, want)
			}
			got := uint32(0)
			if headerCompare(h1, h2) {
				got = 1
			}
			if got != want {
				t.Fatalf("record %d op2 h1=% x h2=% x: compare=%d, want %d", record, h1, h2, got, want)
			}
			counts[1]++
		case 3:
			var h [4]byte
			readFixture(t, r, h[:], record)
			freeFormat := int32(readU32(t, r, record))
			want := [5]int32{
				int32(readU32(t, r, record)),
				int32(readU32(t, r, record)),
				int32(readU32(t, r, record)),
				int32(readU32(t, r, record)),
				int32(readU32(t, r, record)),
			}
			got := [5]int32{
				int32(headerBitrateKbps(h)),
				int32(headerSampleRateHz(h)),
				int32(headerFrameSamples(h)),
				headerFrameBytes(h, freeFormat),
				headerPadding(h),
			}
			if got != want {
				t.Fatalf("record %d op3 header=% x free=%d: got %v, want %v", record, h, freeFormat, got, want)
			}
			counts[2]++
		case 4:
			patternID := readU32(t, r, record)
			length := readU32(t, r, record)
			startBit := readU32(t, r, record)
			width := readU32(t, r, record)
			want := [4]uint32{
				readU32(t, r, record),
				readU32(t, r, record),
				readU32(t, r, record),
				readU32(t, r, record),
			}
			if patternID > 4 || length > 2815 || width > 32 || startBit > length*8+32 {
				t.Fatalf("record %d op4 unsupported parameters pattern=%d length=%d start=%d width=%d", record, patternID, length, startBit, width)
			}
			data := patterns[patternID]
			bs := bitReader{buf: []byte{0xcc}, pos: 17, limit: 8}
			bsInit(&bs, data, int(length))
			if bs.pos != int(want[3]) || bs.limit != int(want[2]) {
				t.Fatalf("record %d op4 bsInit state pos=%d limit=%d, want pos=%d limit=%d", record, bs.pos, bs.limit, want[3], want[2])
			}
			bs.pos = int(startBit)
			gotValue := getBits(&bs, int(width))
			got := [4]uint32{gotValue, uint32(bs.pos), uint32(bs.limit), uint32(0)}
			if got != want {
				t.Fatalf("record %d op4 pattern=%d length=%d start=%d width=%d: got %v, want %v", record, patternID, length, startBit, width, got, want)
			}
			counts[3]++
		default:
			t.Fatalf("record %d: unknown op %d", record, op)
		}
	}
	if got := hex.EncodeToString(fixtureHash.Sum(nil)); got != wantIntegerFixtureSHA256 {
		t.Fatalf("uncompressed fixture SHA256 %s, want %s", got, wantIntegerFixtureSHA256)
	}
	if counts != wantIntegerFixtureCounts {
		t.Fatalf("fixture op counts %v, want %v", counts, wantIntegerFixtureCounts)
	}
	t.Logf("validated %d frozen integer-helper vectors; uncompressed SHA256 %s", counts[0]+counts[1]+counts[2]+counts[3], wantIntegerFixtureSHA256)
}

func checkBitHeaderInvariants(t *testing.T) {
	t.Helper()
	var msb bitReader
	bsInit(&msb, []byte{0xb2}, 1)
	if got := getBits(&msb, 4); got != 0xb || msb.pos != 4 {
		t.Fatalf("MSB-first read got value=%#x pos=%d, want value=0xb pos=4", got, msb.pos)
	}
	if got := getBits(&msb, 0); got != 0 || msb.pos != 4 {
		t.Fatalf("zero-width read got value=%d pos=%d, want 0 and unchanged pos 4", got, msb.pos)
	}
	var overrun bitReader
	bsInit(&overrun, []byte{0xff}, 1)
	if got := getBits(&overrun, 9); got != 0 || overrun.pos != 9 {
		t.Fatalf("overrun got value=%d pos=%d, want 0 and advanced pos 9", got, overrun.pos)
	}
	if got := getBits(&overrun, 1); got != 0 || overrun.pos != 10 {
		t.Fatalf("repeated overrun got value=%d pos=%d, want 0 and advanced pos 10", got, overrun.pos)
	}
	borrowed := []byte{0x80}
	var alias bitReader
	bsInit(&alias, borrowed, 1)
	borrowed[0] = 0x40
	if got := getBits(&alias, 2); got != 1 {
		t.Fatalf("bit reader failed to retain borrowed buffer alias: got %#x, want 1", got)
	}

	mpeg1 := [4]byte{0xff, 0xfb, 0x90, 0x00}
	if !headerValid(mpeg1) || headerBitrateKbps(mpeg1) != 128 || headerSampleRateHz(mpeg1) != 44100 ||
		headerFrameSamples(mpeg1) != 1152 || headerFrameBytes(mpeg1, -123) != 417 || headerPadding(mpeg1) != 0 {
		t.Fatalf("MPEG-1 reference header arithmetic mismatch")
	}
	mpeg2 := [4]byte{0xff, 0xf3, 0x80, 0x00}
	if !headerValid(mpeg2) || headerBitrateKbps(mpeg2) != 64 || headerSampleRateHz(mpeg2) != 22050 ||
		headerFrameSamples(mpeg2) != 576 || headerFrameBytes(mpeg2, -123) != 208 || headerPadding(mpeg2) != 0 {
		t.Fatalf("MPEG-2 reference header arithmetic mismatch")
	}
	if headerFrameBytes([4]byte{0xff, 0xfb, 0x00, 0x00}, -123) != -123 {
		t.Fatalf("free-format fallback did not preserve signed int32 value")
	}
	if headerPadding([4]byte{0xff, 0xff, 0x92, 0x00}) != 4 ||
		headerPadding([4]byte{0xff, 0xfb, 0x92, 0x00}) != 1 {
		t.Fatalf("padding invariants do not match Layer I/other-layer values")
	}
	base := mpeg1
	ignored := mpeg1
	ignored[1] ^= 1   // CRC bit is not compared.
	ignored[2] ^= 2   // Padding bit is not compared.
	ignored[3] = 0xc0 // Stereo mode is not compared.
	if !headerCompare(base, ignored) {
		t.Fatalf("header compare stopped ignoring CRC/padding/stereo bits")
	}
	rateChanged := mpeg1
	rateChanged[2] ^= 4
	if headerCompare(base, rateChanged) {
		t.Fatalf("header compare accepted differing sample-rate bits")
	}
	freeChanged := mpeg1
	freeChanged[2] &= 0x0f
	if headerCompare(base, freeChanged) {
		t.Fatalf("header compare accepted differing free-format status")
	}
	badSync := mpeg1
	badSync[0] = 0
	if !headerCompare(badSync, mpeg1) || headerCompare(mpeg1, badSync) {
		t.Fatalf("header compare lost asymmetric validation of only its second argument")
	}
}

func integerPattern(id, n int) []byte {
	b := make([]byte, n)
	for i := range b {
		switch id {
		case 0:
			b[i] = 0
		case 1:
			b[i] = 0xff
		case 2:
			if i&1 == 0 {
				b[i] = 0x55
			} else {
				b[i] = 0xaa
			}
		case 3:
			b[i] = byte(i)
		case 4:
			b[i] = byte((i*73 + 19) ^ (i >> 2))
		default:
			panic("invalid integer fixture pattern")
		}
	}
	return b
}

func readFixture(t *testing.T, r io.Reader, dst []byte, record int) {
	t.Helper()
	if _, err := io.ReadFull(r, dst); err != nil {
		t.Fatalf("record %d: truncated payload: %v", record, err)
	}
}

func readU32(t *testing.T, r io.Reader, record int) uint32 {
	t.Helper()
	var b [4]byte
	readFixture(t, r, b[:], record)
	return binary.LittleEndian.Uint32(b[:])
}
