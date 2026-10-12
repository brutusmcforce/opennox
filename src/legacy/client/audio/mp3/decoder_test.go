package mp3

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"io"
	"os"
	"testing"
)

func TestFrameFrozenVectors(t *testing.T) {
	f, err := os.Open("testdata/frame_c.bin.gz")
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
	if string(magic[:]) != "NMP3FRM1" {
		t.Fatal("bad frame magic")
	}
	sequences, calls := 0, 0
	for record := 0; ; record++ {
		op, err := r.ReadByte()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		if op != 1 {
			t.Fatal("invalid frame opcode")
		}
		u := func() uint32 { return readU32(t, r, record) }
		steps := u()
		if steps < 1 || steps > 8 {
			t.Fatal("invalid frame sequence length")
		}
		type request struct {
			flags  uint32
			packet []byte
		}
		seq := make([]request, steps)
		for i := range seq {
			flags, n := u(), u()
			if flags > 3 || n > 49152 {
				t.Fatal("invalid frame request")
			}
			seq[i] = request{flags, make([]byte, n)}
			readFixture(t, r, seq[i].packet, record)
		}
		var dec, poison Decoder
		for step, req := range seq {
			before := append([]byte(nil), req.packet...)
			run := func(d *Decoder, poisonState bool) []byte {
				if req.flags&1 != 0 {
					d.Init()
				}
				if poisonState {
					for i := 0; i < 15; i++ {
						d.state.qmfState[896+4*i+2] = float32(12345 + i)
						d.state.qmfState[896+4*i+3] = float32(-12345 - i)
					}
				}
				var backing [2384]int16
				for i := range backing {
					backing[i] = -13108
				}
				pcm := backing[8:2376:2376]
				info := FrameInfo{0x5a5a5a5a, 0x5a5a5a5a, 0x5a5a5a5a, 0x5a5a5a5a, 0x5a5a5a5a}
				out := pcm
				if req.flags&2 != 0 {
					out = nil
				}
				samples := d.DecodeFrame(req.packet, out, &info)
				var got []byte
				for _, v := range []int{samples, info.FrameBytes, info.Channels, info.Hz, info.Layer, info.BitrateKbps} {
					got = binary.LittleEndian.AppendUint32(got, uint32(v))
				}
				got = frameStateBytes(got, &d.state)
				got = synthesisPCM(got, pcm)
				got = binary.LittleEndian.AppendUint32(got, streamBool(bytes.Equal(before, req.packet)))
				intact := true
				for _, v := range backing[:8] {
					intact = intact && v == -13108
				}
				for _, v := range backing[2376:] {
					intact = intact && v == -13108
				}
				got = binary.LittleEndian.AppendUint32(got, streamBool(intact))
				return got
			}
			got := run(&dec, false)
			other := run(&poison, true)
			if !bytes.Equal(got, other) {
				t.Fatalf("sequence%d step%d unused QMF poison affected observable output", record, step)
			}
			var want [11435]byte
			readFixture(t, r, want[:], record)
			if !bytes.Equal(got, want[:]) {
				at := 0
				for at < len(got) && at < len(want) && got[at] == want[at] {
					at++
				}
				t.Fatalf("sequence%d step%d differs byte%d (gotlen%d)", record, step, at, len(got))
			}
			calls++
		}
		sequences++
	}
	if sequences != 544 || calls != 2682 {
		t.Fatalf("frame counts%d/%d", sequences, calls)
	}
	const wantSHA = "bde0b75144e2caf02f62f2e56fbbfa9dfbdb8481a1c2789d8592b92f587bdbc2"
	if hex.EncodeToString(hash.Sum(nil)) != wantSHA {
		t.Fatal("frame fixture hash mismatch")
	}
	t.Logf("validated%d complete frame calls in%d sequences; SHA256%s", calls, sequences, wantSHA)
}
func frameStateBytes(dst []byte, d *decoderState) []byte {
	dst = spectrumFloats(dst, d.mdctOverlap[0][:])
	dst = spectrumFloats(dst, d.mdctOverlap[1][:])
	qmf := d.qmfState
	// Only these30 original uninitialized words are overwritten before any read.
	for i := 0; i < 15; i++ {
		qmf[896+4*i+2] = 0
		qmf[896+4*i+3] = 0
	}
	dst = spectrumFloats(dst, qmf[:])
	dst = binary.LittleEndian.AppendUint32(dst, uint32(d.reserv))
	dst = binary.LittleEndian.AppendUint32(dst, uint32(d.freeFormatBytes))
	dst = append(dst, d.header[:]...)
	return append(dst, d.reservBuf[:]...)
}
