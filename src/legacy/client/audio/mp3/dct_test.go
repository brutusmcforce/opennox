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

func TestDCTFrozenVectors(t *testing.T) {
	f, err := os.Open("testdata/dct_c.bin.gz")
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
	if string(magic[:]) != "NMP3DCT1" {
		t.Fatal("bad DCT fixture magic")
	}
	count := 0
	for record := 0; ; record++ {
		op, err := r.ReadByte()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		if op != 1 {
			t.Fatal("invalid opcode")
		}
		u := func() uint32 { return readU32(t, r, record) }
		pat, seed, mark, off, n := u(), u(), u(), u(), u()
		if pat > 5 || mark >= 1216 || off > 576 || n > 18 {
			t.Fatal("invalid DCT input")
		}
		var data [1216]float32
		imdctPattern(data[:], pat, seed, mark)
		before := data
		dctII(data[off:], int(n))
		// Independent extent check: only the first n columns of this channel may change.
		for i := range data {
			j := i - int(off)
			if j < 0 || j >= 576 || j%18 >= int(n) {
				if math.Float32bits(data[i]) != math.Float32bits(before[i]) {
					t.Fatalf("record%d changed inactive word%d", record, i)
				}
			}
		}
		got := spectrumFloats(nil, data[:])
		got = binary.LittleEndian.AppendUint32(got, 1)
		var want [4868]byte
		readFixture(t, r, want[:], record)
		if !bytes.Equal(got, want[:]) {
			at := 0
			for at < len(got) && got[at] == want[at] {
				at++
			}
			word := at &^ 3
			t.Fatalf("record%d byte%d got%08x want%08x", record, at, binary.LittleEndian.Uint32(got[word:]), binary.LittleEndian.Uint32(want[word:]))
		}
		count++
	}
	if count != 2172 {
		t.Fatalf("DCT count%d", count)
	}
	const wantSHA = "94e0c6fc6283b89b0846dee8a44a40eca979374861b877b8e3241a354af482d6"
	if hex.EncodeToString(hash.Sum(nil)) != wantSHA {
		t.Fatal("DCT fixture hash mismatch")
	}
	// A constant input column has only a DC transform coefficient.
	var dc [576]float32
	for i := 0; i < 32; i++ {
		dc[i*18] = 1
	}
	dctII(dc[:], 1)
	for i, v := range dc {
		want := float32(0)
		if i == 0 {
			want = 32
		}
		if v != want {
			t.Fatalf("DC response index%d: %g", i, v)
		}
	}
	t.Logf("validated%d synthesis DCT vectors; SHA256%s", count, wantSHA)
}
