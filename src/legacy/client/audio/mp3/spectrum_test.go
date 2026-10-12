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

func TestSpectrumFrozenVectors(t *testing.T) {
	checkSpectrumInvariants(t)
	f, err := os.Open("testdata/spectrum_c.bin.gz")
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
	if string(magic[:]) != "NMP3SPC1" {
		t.Fatal("bad spectrum fixture magic")
	}
	var counts [6]int
	var gotStore, wantStore [7076]byte
	for record := 0; ; record++ {
		op, err := r.ReadByte()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		u := func() uint32 { return readU32(t, r, record) }
		got := gotStore[:0]
		expectedBytes := 0
		var data [1192]float32
		var scratch [576]float32
		for i := range scratch {
			scratch[i] = math.Float32frombits(0x4b123456 + uint32(i))
		}
		switch op {
		case 1, 2:
			n, off, pat := u(), u(), u()
			if n > 576 || off > 576-n || pat > 5 {
				t.Fatal("invalid stereo band fixture")
			}
			spectrumPattern(data[:], pat, 0)
			if op == 1 {
				midsideStereo(data[off:], int(n))
			} else {
				kl, kr := u(), u()
				intensityStereoBand(data[off:], int(n), math.Float32frombits(kl), math.Float32frombits(kr))
			}
			got = spectrumFloats(got, data[:])
			got = binary.LittleEndian.AppendUint32(got, 1)
			expectedBytes = 4772
		case 3, 4, 5:
			var hdr [4]byte
			var side [64]byte
			readFixture(t, r, hdr[:], record)
			readFixture(t, r, side[:], record)
			var gr [4]grInfo
			var bs bitReader
			bsInit(&bs, side[:], 64)
			if !headerValid(hdr) || (hdr[1]>>1)&3 != 1 || readSideInfo(&bs, &gr, hdr) < 0 {
				t.Fatal("invalid sideinfo setup")
			}
			pat, mark := u(), u()
			if pat > 5 || mark > 575 {
				t.Fatal("invalid spectrum pattern")
			}
			spectrumPattern(data[:], pat, mark)
			switch op {
			case 3:
				before := spectrumFloats(nil, data[:])
				max := stereoTopBand(data[576:], gr[0].sfbTable, int(gr[0].nLongSFB)+int(gr[0].nShortSFB))
				for _, v := range max {
					got = binary.LittleEndian.AppendUint32(got, uint32(int32(v)))
				}
				got = binary.LittleEndian.AppendUint32(got, streamBool(bytes.Equal(before, spectrumFloats(nil, data[:]))))
				got = binary.LittleEndian.AppendUint32(got, 1)
				expectedBytes = 20
			case 4:
				shift, ip := u(), u()
				if shift > 1 || ip > 7 {
					t.Fatal("invalid intensity fixture")
				}
				gr[1].scalefacCompress = uint16(shift)
				before := gr
				var ist [39]byte
				for i := range ist {
					ist[i] = spectrumIPos(ip, i)
				}
				intensityStereo(data[:], ist[:], (*[2]grInfo)(gr[:2]), hdr)
				unchanged := true
				for i := range gr {
					if sideInfoScalars(gr[i]) != sideInfoScalars(before[i]) || len(gr[i].sfbTable) != len(before[i].sfbTable) || (len(gr[i].sfbTable) > 0 && !sameSideInfoTable(gr[i].sfbTable, before[i].sfbTable)) {
						unchanged = false
					}
				}
				got = spectrumFloats(got, data[:])
				got = append(got, ist[:]...)
				got = binary.LittleEndian.AppendUint32(got, 1)
				got = binary.LittleEndian.AppendUint32(got, streamBool(unchanged))
				expectedBytes = 4815
			case 5:
				ch, off := u(), u()
				if ch > 1 || off > 72 {
					t.Fatal("invalid reorder fixture")
				}
				reorder(data[ch*576+off:], scratch[:], gr[0].sfbTable[gr[0].nLongSFB:])
				got = spectrumFloats(got, data[:])
				got = spectrumFloats(got, scratch[:])
				got = binary.LittleEndian.AppendUint32(got, 1)
				expectedBytes = 7076
			}
		case 6:
			bands, ch, pat := int32(u()), u(), u()
			if bands < -1 || bands > 31 || ch > 1 || pat > 5 {
				t.Fatal("invalid antialias fixture")
			}
			spectrumPattern(data[:], pat, 0)
			antialias(data[ch*576:], int(bands))
			got = spectrumFloats(got, data[:])
			got = binary.LittleEndian.AppendUint32(got, 1)
			expectedBytes = 4772
		default:
			t.Fatalf("record%d unknown opcode%d", record, op)
		}
		want := wantStore[:expectedBytes]
		readFixture(t, r, want, record)
		if !bytes.Equal(got, want) {
			at := 0
			for at < len(got) && at < len(want) && got[at] == want[at] {
				at++
			}
			t.Fatalf("record%d op%d differs at byte%d (got length%d want%d)", record, op, at, len(got), len(want))
		}
		counts[op-1]++
	}
	if counts != [6]int{60, 126, 378, 3600, 144, 56} {
		t.Fatalf("spectrum counts%v", counts)
	}
	const wantSHA = "29aff6bae8ad530905c01fc96adfafdcccd42cbcfb31ca8d25ae6c0e6cfff3f8"
	if hex.EncodeToString(hash.Sum(nil)) != wantSHA {
		t.Fatal("spectrum fixture hash mismatch")
	}
	t.Logf("validated4364 spectrum vectors with exactfloat bits; SHA256%s", wantSHA)
}

func spectrumFloats(dst []byte, v []float32) []byte {
	for _, f := range v {
		dst = binary.LittleEndian.AppendUint32(dst, math.Float32bits(f))
	}
	return dst
}
func spectrumPattern(v []float32, id, mark uint32) {
	for i := range v {
		b := uint32(0)
		if id == 1 || id == 4 || id == 5 {
			b = 0x3f000000 + uint32(i%257)*0x800
		}
		if id == 2 {
			b = 0x3f800000 + uint32(i%257)*0x1000
			if i&1 != 0 {
				b |= 0x80000000
			}
		}
		if id == 3 {
			b = [4]uint32{1, 0x007fffff, 0x00800000, 0x80000001}[i%4]
		}
		if (id == 4 || id == 5) && i >= 576 && i < 1152 {
			b = 0
			if id == 5 && i == 576+int(mark) {
				b = 0x3f800000
			}
		}
		v[i] = math.Float32frombits(b)
	}
}
func spectrumIPos(id uint32, i int) byte {
	switch id {
	case 0:
		return 0
	case 1:
		return byte(i % 7)
	case 2:
		return 6
	case 3:
		return 7
	case 4:
		return 63
	case 5:
		return 64
	case 6:
		return 255
	default:
		return byte(i*7 + 3)
	}
}

func checkSpectrumInvariants(t *testing.T) {
	t.Helper()
	var data [1192]float32
	data[0] = 2
	data[576] = 1
	data[1] = 99
	midsideStereo(data[:], 1)
	if data[0] != 3 || data[576] != 1 || data[1] != 99 {
		t.Fatal("midside must use both old channel values and retain tails")
	}
	data[0] = 4
	data[576] = 99
	intensityStereoBand(data[:], 1, 0.25, 0.75)
	if data[0] != 1 || data[576] != 3 {
		t.Fatal("intensity gains must both use the old left sample")
	}
	right := []float32{0, math.Float32frombits(0x80000000), 0, 1, 0, math.Float32frombits(1)}
	if got := stereoTopBand(right, []uint8{2, 2, 2, 0}, 3); got != [3]int{-1, 1, 2} {
		t.Fatalf("top bands%v", got)
	}
	small := []float32{1, 2, 3, 4, 5, 6}
	tmp := make([]float32, 6)
	reorder(small, tmp, []uint8{2, 2, 2, 0})
	for i, want := range []float32{1, 3, 5, 2, 4, 6} {
		if small[i] != want || tmp[i] != want {
			t.Fatal("short-window reorder permutation")
		}
	}
	data = [1192]float32{}
	data[18] = 1
	antialias(data[:], 1)
	if data[18] != antialiasCoefficients[0][0] || data[17] != antialiasCoefficients[1][0] {
		t.Fatal("antialias coefficient orientation")
	}
	for i, v := range data {
		if i != 17 && i != 18 && v != 0 {
			t.Fatalf("antialias changed unrelated index%d", i)
		}
	}
	data = [1192]float32{}
	data[0] = 2
	stereoProcess(data[:], []uint8{3}, []uint8{2, 0}, [4]byte{255, 251, 144, 0x20}, [3]int{-1, -1, -1}, 0)
	if data[0] != float32(1.41421356) || data[576] != float32(1.41421356) {
		t.Fatal("stereo processing must use raw MS extension bit with symmetric pan")
	}
	// Preserve the source8kHz mixed-block extent explicitly within a sized workspace.
	var scratch [576]float32
	for ch := 0; ch < 2; ch++ {
		for i := range data {
			data[i] = float32(i + 1)
		}
		reorder(data[ch*576+72:], scratch[:], scfMixed[1][6:])
		if data[ch*576+576] != float32(ch*576+541) || data[ch*576+600] != float32(ch*576+601) {
			t.Fatal("legacy mixed-block reorder extent changed")
		}
	}
}
