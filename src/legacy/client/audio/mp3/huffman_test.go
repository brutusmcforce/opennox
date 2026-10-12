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

func TestHuffmanFrozenVectors(t *testing.T) {
	f, err := os.Open("testdata/huffman_c.bin.gz")
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
	if string(magic[:]) != "NMP3HUF1" {
		t.Fatal("bad Huffman magic")
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
		var got []byte
		flag := func(v bool) { got = binary.LittleEndian.AppendUint32(got, streamBool(v)) }
		switch op {
		case 1:
			x := int(int32(u()))
			if x < -16 || x > 8206 {
				t.Fatal("invalid power input")
			}
			got = binary.LittleEndian.AppendUint32(got, math.Float32bits(pow43(x)))
		case 2, 3:
			var hdr [4]byte
			var side [64]byte
			readFixture(t, r, hdr[:], record)
			readFixture(t, r, side[:], record)
			var gr [4]grInfo
			var bs bitReader
			bsInit(&bs, side[:], 64)
			if !headerValid(hdr) || readSideInfo(&bs, &gr, hdr) < 0 {
				t.Fatal("invalid side info")
			}
			tab, big, count1, pos, length, pat, seed, scpat := u(), u(), u(), u(), u(), u(), u(), u()
			if tab > 31 || big > 288 || count1 > 1 || pos > 22480 || length > 4095 || pat > 3 || scpat > 2 || pos+length > 22520 {
				t.Fatal("invalid Huffman request")
			}
			gr[0].bigValues = uint16(big)
			gr[0].count1Table = uint8(count1)
			gr[0].tableSelect = [3]uint8{uint8(tab), uint8((tab + 7) % 32), uint8((tab + 13) % 32)}
			gr[0].regionCount = [3]uint8{0, 0, 255}
			before := gr[0]
			var input [2815]byte
			huffmanBytePattern(input[:], pat, seed)
			if op == 3 {
				n := u()
				if n > 64 {
					t.Fatal("invalid directed byte count")
				}
				readFixture(t, r, input[:n], record)
			}
			oldInput := input
			var dst [640]float32
			for i := range dst {
				dst[i] = math.Float32frombits(0x4b123456 + uint32(i))
			}
			var scf [40]float32
			for i := range scf {
				b := uint32(0x3f800000)
				if scpat == 1 {
					b = uint32(110+i%24) << 23
				} else if scpat == 2 {
					b = 1
					if i&1 != 0 {
						b = 0x80000000
					}
				}
				scf[i] = math.Float32frombits(b)
			}
			oldSCF := scf
			// Deliberately expose only logical bytes in len, leaving physical lookahead in cap.
			bsInit(&bs, input[:], int((pos+length+7)/8))
			bs.pos = int(pos)
			bs.limit = int(pos + length)
			huffman(dst[:], &bs, &gr[0], scf[:], int(pos+length))
			got = spectrumFloats(got, dst[:])
			got = binary.LittleEndian.AppendUint32(got, uint32(bs.pos))
			got = binary.LittleEndian.AppendUint32(got, uint32(bs.limit))
			flag(input == oldInput)
			flag(bytes.Equal(spectrumFloats(nil, scf[:]), spectrumFloats(nil, oldSCF[:])))
			flag(sideInfoScalars(gr[0]) == sideInfoScalars(before) && sameSideInfoTable(gr[0].sfbTable, before.sfbTable))
			flag(true)
			for i := 576; i < len(dst); i++ {
				if math.Float32bits(dst[i]) != 0x4b123456+uint32(i) {
					t.Fatalf("record%d changed spectral tail%d", record, i)
				}
			}
		default:
			t.Fatal("invalid opcode")
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
	if counts != [3]int{8223, 13224, 12906} {
		t.Fatalf("Huffman counts%v", counts)
	}
	const wantSHA = "307f6e6d42de8f0bf0f7a94a95603dc1d4b3772681bc96c13195c4456653bc28"
	if hex.EncodeToString(hash.Sum(nil)) != wantSHA {
		t.Fatal("Huffman hash mismatch")
	}
	for _, x := range []int{0, 1, 8, 27, 64, 125} {
		if pow43(x) < 0 {
			t.Fatal("nonnegative power")
		}
	}
	if pow43(0) != 0 || pow43(1) != 1 || pow43(8) != 16 || pow43(27) != 81 || pow43(64) != 256 || pow43(125) != 625 {
		t.Fatal("exact cube fourth powers")
	}
	t.Logf("validated%v power/Huffman records; SHA256%s", counts, wantSHA)
}
func huffmanBytePattern(v []byte, pat, seed uint32) {
	state := seed
	if state == 0 {
		state = 1
	}
	for i := range v {
		state ^= state << 13
		state ^= state >> 17
		state ^= state << 5
		switch pat {
		case 0:
			v[i] = 0
		case 1:
			v[i] = 255
		case 2:
			v[i] = 0xaa
		default:
			v[i] = byte(state >> 24)
		}
	}
}
