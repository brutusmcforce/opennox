package mp3

import (
	"bufio"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"math"
	"os"
	"testing"
)

func TestScalefactorFrozenVectors(t *testing.T) {
	checkScalefactorInvariants(t)
	f, err := os.Open("testdata/scalefactors_c.bin.gz")
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
	if string(magic[:]) != "NMP3SCF1" {
		t.Fatal("bad scalefactor fixture magic")
	}
	var counts [3]int
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
		case 1:
			var sizes, partitions [4]uint8
			readFixture(t, r, sizes[:], record)
			readFixture(t, r, partitions[:], record)
			scfsi, length, start, pattern := int32(u()), u(), u(), u()
			sum := 0
			for i, n := range partitions {
				sum += int(n)
				if sizes[i] > 5 {
					t.Fatal("invalid read width")
				}
			}
			if sum > 36 || length > 64 || start > 32 || pattern > 4 || (scfsi < 0 && scfsi != -16) || scfsi > 15 {
				t.Fatal("invalid scalefactor read fixture")
			}
			var scf, ist [40]uint8
			for i := range scf {
				scf[i] = byte(0xa5 ^ (i * 13))
				ist[i] = byte(0x80 + i*7)
			}
			data := integerPattern(int(pattern), 128)
			var bs bitReader
			bsInit(&bs, data, int(length))
			bs.pos = int(start)
			readScalefactors(&scf, (*[39]uint8)(ist[:39]), sizes, partitions, &bs, scfsi)
			wantPos, wantLimit := u(), u()
			var wantSCF, wantIST [40]byte
			readFixture(t, r, wantSCF[:], record)
			readFixture(t, r, wantIST[:], record)
			canaries := u()
			if bs.pos != int(wantPos) || bs.limit != int(wantLimit) || scf != wantSCF || ist != wantIST || canaries != 1 {
				t.Fatalf("record%d scalefactor byte state mismatch: pos%d want%d", record, bs.pos, wantPos)
			}
		case 2:
			bits, exp := u(), u()
			if bits&0x7f800000 == 0x7f800000 || exp > 1024 {
				t.Fatal("invalid quarter-power fixture")
			}
			want := u()
			got := math.Float32bits(ldexpQ2(math.Float32frombits(bits), int32(exp)))
			if got != want {
				t.Fatalf("record%d quarter-power y%08x exp%d got%08x want%08x", record, bits, exp, got, want)
			}
		case 3:
			var hdr [4]byte
			readFixture(t, r, hdr[:], record)
			length, start, pattern, compress, gain, layout, sub, flags, ch := u(), u(), u(), u(), u(), u(), u(), u(), u()
			maxCompress := uint32(511)
			if hdr[1]&8 != 0 {
				maxCompress = 15
			}
			if !headerValid(hdr) || (hdr[1]>>1)&3 != 1 || length > 64 || start > 32 || pattern > 4 || compress > maxCompress || gain > 255 || layout > 2 || sub > 511 || flags > 0xf03 || ch > 1 {
				t.Fatal("invalid scalefactor decode fixture")
			}
			gr := grInfo{scalefacCompress: uint16(compress), globalGain: uint8(gain), preflag: uint8(flags & 1), scalefacScale: uint8(flags >> 1 & 1), scfsi: uint8(flags >> 8 & 15), subblockGain: [3]uint8{uint8(sub & 7), uint8(sub >> 3 & 7), uint8(sub >> 6 & 7)}}
			switch layout {
			case 0:
				gr.nLongSFB = 22
			case 1:
				gr.nLongSFB = 6
				if hdr[1]&8 != 0 {
					gr.nLongSFB = 8
				}
				gr.nShortSFB = 30
			case 2:
				gr.nShortSFB = 39
			}
			before := gr
			var scf [40]float32
			var ist [40]byte
			for i := range scf {
				scf[i] = math.Float32frombits(0x4b123456 + uint32(i))
				ist[i] = byte(0x80 + i*7)
			}
			data := integerPattern(int(pattern), 128)
			var bs bitReader
			bsInit(&bs, data, int(length))
			bs.pos = int(start)
			decodeScalefactors(hdr, (*[39]uint8)(ist[:39]), &bs, &gr, &scf, int(ch))
			wantPos, wantLimit := u(), u()
			if bs.pos != int(wantPos) || bs.limit != int(wantLimit) {
				t.Fatalf("record%d decoded reader state mismatch", record)
			}
			for i, v := range scf {
				want := u()
				if got := math.Float32bits(v); got != want {
					t.Fatalf("record%d float%d got%08x want%08x hdr%x compress%d layout%d gain%d flags%x", record, i, got, want, hdr, compress, layout, gain, flags)
				}
			}
			var wantIST [40]byte
			readFixture(t, r, wantIST[:], record)
			canaries, unchanged := u(), u()
			if ist != wantIST || canaries != 1 || unchanged != 1 || sideInfoScalars(gr) != sideInfoScalars(before) || gr.sfbTable != nil {
				t.Fatalf("record%d intensity bytes/retained granule mismatch", record)
			}
		default:
			t.Fatalf("record%d unknown opcode%d", record, op)
		}
		counts[op-1]++
	}
	if counts != [3]int{34272, 10250, 12480} {
		t.Fatalf("scalefactor fixture counts%v", counts)
	}
	const wantSHA = "44d7df50676a00eb2379474ab61bdd2dcffd0ec7fac2905b7f8f335d27063707"
	if hex.EncodeToString(hash.Sum(nil)) != wantSHA {
		t.Fatal("scalefactor fixture hash mismatch")
	}
	t.Logf("validated57002 scalefactor vectors with exactfloat32 bits; SHA256%s", wantSHA)
}

func checkScalefactorInvariants(t *testing.T) {
	t.Helper()
	if ldexpQ2(8, 4) != 4 || ldexpQ2(2048, 44) != 1 {
		t.Fatal("quarter-power whole exponent identity")
	}
	if math.Float32bits(ldexpQ2(math.Float32frombits(1), 0)) != 1 || math.Float32bits(ldexpQ2(math.Float32frombits(0x80000000), 4)) != 0x80000000 {
		t.Fatal("subnormal identity/signed zero")
	}
	var scf [40]uint8
	var ist [39]uint8
	for i := range scf {
		scf[i] = 0xa5
	}
	for i := range ist {
		ist[i] = byte(i + 2)
	}
	oldIST := ist
	var bs bitReader
	bsInit(&bs, []byte{0xff}, 1)
	readScalefactors(&scf, &ist, [4]uint8{3}, [4]uint8{2}, &bs, 8)
	if scf[0] != 2 || scf[1] != 3 || scf[2] != 0 || scf[3] != 0 || scf[4] != 0 || scf[5] != 0xa5 || ist != oldIST || bs.pos != 0 {
		t.Fatal("reuse must copy prior values without reading, zero only three tail bytes")
	}
	bsInit(&bs, []byte{0xff}, 1)
	readScalefactors(&scf, &ist, [4]uint8{2}, [4]uint8{2}, &bs, -16)
	if scf[0] != 3 || scf[1] != 3 || ist[0] != 255 || ist[1] != 255 || bs.pos != 4 {
		t.Fatal("intensity sentinel must preserve decoded value in scalefactors")
	}
	bsInit(&bs, []byte{0xff}, 1)
	readScalefactors(&scf, &ist, [4]uint8{}, [4]uint8{2}, &bs, 0)
	if scf[0] != 0 || scf[1] != 0 || ist[0] != 0 || ist[1] != 0 || bs.pos != 0 {
		t.Fatal("zero-width partitions must clear outputs without reading")
	}
	for _, hdr := range [][4]byte{{255, 251, 144, 0}, {255, 243, 128, 0}} {
		var out [40]float32
		for i := range out {
			out[i] = 123
		}
		for i := range ist {
			ist[i] = 0xa5
		}
		bsInit(&bs, []byte{0xff}, 1)
		gr := grInfo{globalGain: 214, nLongSFB: 22}
		decodeScalefactors(hdr, &ist, &bs, &gr, &out, 0)
		for i, v := range out {
			want := float32(123)
			if i < 22 {
				want = 1
			}
			if v != want {
				t.Fatalf("zero-compressor unit gain/tail at%d got%g want%g", i, v, want)
			}
		}
		if bs.pos != 0 || ist[20] != 0 || ist[21] != 0xa5 {
			t.Fatal("zero-compressor read/tail contract")
		}
	}
}
