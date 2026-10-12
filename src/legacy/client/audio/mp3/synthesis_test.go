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

func TestSynthesisFrozenVectors(t *testing.T) {
	checkSynthesisInvariants(t)
	f, err := os.Open("testdata/synthesis_c.bin.gz")
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
	if string(magic[:]) != "NMP3SYN1" {
		t.Fatal("bad synthesis magic")
	}
	var counts [4]int
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
		var got []byte
		flag := func(v bool) { got = binary.LittleEndian.AppendUint32(got, streamBool(v)) }
		if op == 1 {
			b := u()
			if b&0x7fffffff > 0x7f800000 {
				t.Fatal("NaN outside C cast contract")
			}
			got = binary.LittleEndian.AppendUint32(got, uint32(int32(scalePCM(math.Float32frombits(b)))))
		} else {
			pat, seed, mark, nch, shift := u(), u(), u(), u(), u()
			if pat > 5 || mark >= 2176 || nch < 1 || nch > 2 || shift > 24 {
				t.Fatal("invalid synthesis request")
			}
			var data [1216]float32
			var lins [2176]float32
			var qmf [1024]float32
			var pcm [1216]int16
			synthesisPattern(data[:], pat, seed, mark, shift)
			synthesisPattern(lins[:], pat, seed^0x4c494e, mark, shift)
			synthesisPattern(qmf[:], pat, seed^0x514d46, mark, shift)
			for i := range pcm {
				pcm[i] = -13108
			}
			switch op {
			case 2:
				off := u()
				if off > 125 {
					t.Fatal("invalid pair offset")
				}
				before := lins
				synthPair(pcm[:], int(nch), lins[off:])
				got = synthesisPCM(got, pcm[:96])
				flag(bytes.Equal(spectrumFloats(nil, lins[:]), spectrumFloats(nil, before[:])))
				flag(true)
			case 3:
				off := u()
				if off > 16 || off&1 != 0 {
					t.Fatal("invalid synth offset")
				}
				before := data
				synth(data[off:], pcm[:], int(nch), lins[:])
				if !bytes.Equal(spectrumFloats(nil, data[:]), spectrumFloats(nil, before[:])) {
					t.Fatal("synth changed spectral input")
				}
				got = spectrumFloats(got, data[:])
				got = spectrumFloats(got, lins[:])
				got = synthesisPCM(got, pcm[:160])
				flag(true)
			case 4:
				steps := u()
				if steps < 1 || steps > 4 {
					t.Fatal("invalid steps")
				}
				var seq [4][2]uint32
				for i := uint32(0); i < steps; i++ {
					seq[i] = [2]uint32{u(), u()}
					if seq[i][0] > 18 || seq[i][0]&1 != 0 || seq[i][1] < 1 || seq[i][1] > 2 {
						t.Fatal("invalid granule request")
					}
				}
				for step := uint32(0); step < steps; step++ {
					synthesisPattern(data[:], pat, seed+step*37, (mark+step*17)%1152, shift)
					for i := range pcm {
						pcm[i] = -13108
					}
					before := qmf
					bands, channels := int(seq[step][0]), int(seq[step][1])
					synthGranule(qmf[:], data[:], bands, channels, pcm[:], lins[:])
					for i := range qmf {
						if i >= 960 || (channels == 1 && i&1 != 0) {
							if math.Float32bits(qmf[i]) != math.Float32bits(before[i]) {
								t.Fatalf("record%d step%d retained qmf%d changed", record, step, i)
							}
						}
					}
					for i := 32 * bands * channels; i < len(pcm); i++ {
						if pcm[i] != -13108 {
							t.Fatal("PCM tail changed")
						}
					}
					got = spectrumFloats(got, data[:])
					got = spectrumFloats(got, lins[:])
					got = spectrumFloats(got, qmf[:])
					got = synthesisPCM(got, pcm[:])
					flag(true)
				}
				stepsTotal += int(steps)
			default:
				t.Fatal("invalid opcode")
			}
		}
		want := make([]byte, len(got))
		readFixture(t, r, want, record)
		if !bytes.Equal(got, want) {
			at := 0
			for at < len(got) && got[at] == want[at] {
				at++
			}
			t.Fatalf("record%d op%d byte%d got%02x want%02x", record, op, at, got[at], want[at])
		}
		counts[op-1]++
	}
	if counts != [4]int{196616, 390, 408, 360} || stepsTotal != 576 {
		t.Fatalf("counts%v steps%d", counts, stepsTotal)
	}
	const wantSHA = "303913b591fb7498ab55e24fa945db58d0cf63dd5cb16936bf7123e5528bb764"
	if hex.EncodeToString(hash.Sum(nil)) != wantSHA {
		t.Fatal("synthesis hash mismatch")
	}
	t.Logf("validated%v synthesis records and%d granule steps; SHA256%s", counts, stepsTotal, wantSHA)
}
func synthesisPCM(dst []byte, v []int16) []byte {
	for _, s := range v {
		dst = binary.LittleEndian.AppendUint16(dst, uint16(s))
	}
	return dst
}
func synthesisPattern(v []float32, pat, seed, mark, shift uint32) {
	imdctPattern(v, pat, seed, mark)
	scale := float32(1)
	for i := uint32(0); i < shift; i++ {
		scale *= 0.5
	}
	for i := range v {
		v[i] = float32(v[i] * scale)
	}
}
func checkSynthesisInvariants(t *testing.T) {
	t.Helper()
	// Preserve the actual scalar negative rounding quirk rather than replacing it
	// with a generic round-to-nearest operation.
	for _, row := range []struct {
		x float32
		y int16
	}{{0, 0}, {0.5, 1}, {-0.5, 0}, {-1, 0}, {-1.5, -2}, {32766.5, 32767}, {-32767.5, -32768}} {
		if got := scalePCM(row.x); got != row.y {
			t.Fatalf("scalePCM(%g)=%d want%d", row.x, got, row.y)
		}
	}
	var qmf [960]float32
	var data [1152]float32
	var lins [2112]float32
	var pcm [1152]int16
	for _, channels := range []int{1, 2} {
		for i := range pcm {
			pcm[i] = -13108
		}
		synthGranule(qmf[:], data[:], 18, channels, pcm[:], lins[:])
		for i := 0; i < 576*channels; i++ {
			if pcm[i] != 0 {
				t.Fatal("silent synthesis")
			}
		}
	}
}
