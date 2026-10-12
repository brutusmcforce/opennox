package mp3

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"io"
	"math"
	"os"
	"testing"
)

func TestIMDCTFrozenVectors(t *testing.T) {
	checkIMDCTInvariants(t)
	f, err := os.Open("testdata/imdct_c.bin.gz")
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
	if string(magic[:]) != "NMP3IMD1" {
		t.Fatal("bad IMDCT magic")
	}
	var counts [7]int
	stepsTotal := 0
	for record := 0; ; record++ {
		op, err := r.ReadByte()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		u := func() uint32 { return readU32(t, r, record) }
		pat, seed, mark := u(), u(), u()
		if pat > 5 || mark > 639 {
			t.Fatal("invalid pattern")
		}
		var data [640]float32
		var overlap [320]float32
		imdctPattern(data[:], pat, seed, mark)
		ovpat := pat
		if pat >= 2 {
			ovpat = 4
		}
		imdctPattern(overlap[:], ovpat, seed^0x50434d, mark)
		var got []byte
		flag := func(v bool) { got = binary.LittleEndian.AppendUint32(got, streamBool(v)) }
		state := func() { got = spectrumFloats(got, data[:]); got = spectrumFloats(got, overlap[:]); flag(true) }
		switch op {
		case 1:
			dct3Nine((*[9]float32)(data[:9]))
			got = spectrumFloats(got, data[:18])
			flag(true)
		case 2:
			before := data
			idct3(data[0], data[1], data[2], (*[3]float32)(overlap[:3]))
			got = spectrumFloats(got, overlap[:9])
			flag(bytes.Equal(spectrumFloats(nil, before[:3]), spectrumFloats(nil, data[:3])))
			flag(true)
		case 3:
			off := u()
			if off > 2 {
				t.Fatal("invalid offset")
			}
			before := data
			var dst [12]float32
			imdctPattern(dst[:], 4, seed^0x445354, 0)
			imdct12(data[off:], dst[:], overlap[:])
			got = spectrumFloats(got, dst[:])
			got = spectrumFloats(got, overlap[:9])
			flag(bytes.Equal(spectrumFloats(nil, before[:20]), spectrumFloats(nil, data[:20])))
			flag(true)
		case 4, 5:
			bands := u()
			if bands > 32 {
				t.Fatal("invalid bands")
			}
			if op == 4 {
				win := u()
				if win > 2 {
					t.Fatal("invalid window")
				}
				var window [18]float32
				for i := range window {
					switch win {
					case 0:
						window[i] = 1
					case 1:
						window[i] = float32(i & 1)
					case 2:
						window[i] = float32(i+1) * (1.0 / 32)
					}
				}
				imdct36(data[:], overlap[:], window[:], int(bands))
			} else {
				imdctShort(data[:], overlap[:], int(bands))
			}
			state()
		case 6:
			changeSign(data[:])
			got = spectrumFloats(got, data[:])
			flag(true)
		case 7:
			steps := u()
			if steps < 1 || steps > 6 {
				t.Fatal("invalid steps")
			}
			var seq [6][2]uint32
			for i := uint32(0); i < steps; i++ {
				seq[i] = [2]uint32{u(), u()}
				if seq[i][0] > 3 || seq[i][1] > 32 {
					t.Fatal("invalid sequence")
				}
			}
			for i := uint32(0); i < steps; i++ {
				imdctPattern(data[:], pat, seed+i*37, (mark+i*17)%576)
				imdctGranule(data[:], overlap[:], int(seq[i][0]), int(seq[i][1]))
				state()
			}
			stepsTotal += int(steps)
		default:
			t.Fatalf("unknown opcode%d", op)
		}
		want := make([]byte, len(got))
		readFixture(t, r, want, record)
		if !bytes.Equal(got, want) {
			at := 0
			for at < len(got) && got[at] == want[at] {
				at++
			}
			word := at &^ 3
			t.Fatalf("record%d op%d byte%d got%08x want%08x", record, op, at, binary.LittleEndian.Uint32(got[word:]), binary.LittleEndian.Uint32(want[word:]))
		}
		counts[op-1]++
	}
	if counts != [7]int{98, 50, 558, 720, 240, 34, 680} || stepsTotal != 1258 {
		t.Fatalf("counts%v steps%d", counts, stepsTotal)
	}
	const wantSHA = "5f5456faf5c7002062876ae24d0969de7383452e9e7bc61e82d7735e65f2a3b1"
	if hex.EncodeToString(hash.Sum(nil)) != wantSHA {
		t.Fatal("fixture hash mismatch")
	}
	t.Logf("validated2380 inverse-transform records and1258 wrapper steps; SHA256%s", wantSHA)
}
func imdctPattern(v []float32, id, seed, mark uint32) {
	state := seed
	if state == 0 {
		state = 1
	}
	for i := range v {
		var b uint32
		switch id {
		case 1:
			b = 0x80000000
		case 2:
			x := float32(i+1+int(seed%17)) * (1.0 / 16)
			if i&1 != 0 {
				x = -x
			}
			b = math.Float32bits(x)
		case 3:
			if uint32(i) == mark {
				b = 0x3f800000
			}
		case 4:
			state ^= state << 13
			state ^= state >> 17
			state ^= state << 5
			b = (state & 0x807fffff) | ((100 + (state>>24)%40) << 23)
		case 5:
			b = [4]uint32{1, 0x007fffff, 0x00800000, 0x80000001}[(uint32(i)+seed)%4]
		}
		v[i] = math.Float32frombits(b)
	}
}
func checkIMDCTInvariants(t *testing.T) {
	t.Helper()
	var data [640]float32
	imdctPattern(data[:], 4, 42, 0)
	before := data
	changeSign(data[:])
	for i := range data {
		want := math.Float32bits(before[i])
		if i < 576 && (i/18)&1 == 1 && i&1 == 1 {
			want ^= 0x80000000
		}
		if math.Float32bits(data[i]) != want {
			t.Fatalf("sign flip index%d", i)
		}
	}
	changeSign(data[:])
	if data != before {
		t.Fatal("double sign flip must restore all words")
	}
	var overlap [320]float32
	imdctPattern(overlap[:], 4, 17, 0)
	old := overlap
	imdct36(data[:], overlap[:], mdctWindow[0][:], 0)
	imdctShort(data[:], overlap[:], 0)
	if data != before || overlap != old {
		t.Fatal("zero bands must retain state")
	}
	imdctShort(data[:], overlap[:], 1)
	for i := 0; i < 6; i++ {
		if data[i] != old[i] {
			t.Fatal("short block must emit previous overlap prefix")
		}
	}
	if !bytes.Equal(spectrumFloats(nil, data[18:]), spectrumFloats(nil, before[18:])) || !bytes.Equal(spectrumFloats(nil, overlap[9:]), spectrumFloats(nil, old[9:])) {
		t.Fatal("one short band changed later bands")
	}
	var y [3]float32
	idct3(1, 0, 0, &y)
	if y != [3]float32{1, 1, 1} {
		t.Fatal("DC butterfly")
	}
}
