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
	sideInfoFixturePath              = "testdata/sideinfo_c.bin.gz"
	sideInfoFixtureMagic             = "NMP3SID1"
	wantSideInfoFixtureSHA256        = "8bea472e931894a2649b33e1753f8cb768ee829493ce6d2c8870ad25479d5f1e"
	wantSideInfoFixtureCases  uint64 = 36864
)

func TestSideInfoFrozenVectors(t *testing.T) {
	checkSideInfoInvariants(t)
	f, err := os.Open(sideInfoFixturePath)
	if err != nil {
		t.Fatalf("open frozen actual-C fixture %q: %v", sideInfoFixturePath, err)
	}
	defer f.Close()
	gz, err := gzip.NewReader(f)
	if err != nil {
		t.Fatalf("open gzip fixture: %v", err)
	}
	defer gz.Close()
	h := sha256.New()
	r := bufio.NewReader(io.TeeReader(gz, h))
	var magic [8]byte
	readSideInfoFixture(t, r, magic[:], -1)
	if string(magic[:]) != sideInfoFixtureMagic {
		t.Fatalf("fixture magic %q, want %q", magic, sideInfoFixtureMagic)
	}
	cases := uint64(0)
	for record := 0; ; record++ {
		var hdr [4]byte
		first, err := r.ReadByte()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("record %d: read header: %v", record, err)
		}
		hdr[0] = first
		readSideInfoFixture(t, r, hdr[1:], record)
		length := readSideInfoU32(t, r, record)
		start := readSideInfoU32(t, r, record)
		seed := readSideInfoU32(t, r, record)
		var inputBytes [64]byte
		readSideInfoFixture(t, r, inputBytes[:], record)
		if !headerValid(hdr) || (hdr[1]>>1)&3 != 1 || length > 4096 || start > length*8+32 || start > 512 || seed > 255 {
			t.Fatalf("record %d: unsupported input header=% x length=%d start=%d seed=%d", record, hdr, length, start, seed)
		}

		var data [4096 + 64]byte
		copy(data[:64], inputBytes[:])
		var bs bitReader
		bsInit(&bs, data[:], int(length))
		bs.pos = int(start)
		sentinel := [40]uint8{}
		for i := range sentinel {
			sentinel[i] = uint8(seed)
		}
		var gr [4]grInfo
		for i := range gr {
			gr[i] = seededGrInfo(uint8(seed), sentinel[:])
		}
		ret := readSideInfo(&bs, &gr, hdr)

		wantRet := readSideInfoU32(t, r, record)
		wantPos := readSideInfoU32(t, r, record)
		wantLimit := readSideInfoU32(t, r, record)
		if uint32(ret) != wantRet || uint32(bs.pos) != wantPos || uint32(bs.limit) != wantLimit {
			t.Fatalf("record %d: return/reader got (%d,%d,%d), want bits (%d,%d,%d)", record,
				ret, bs.pos, bs.limit, int32(wantRet), wantPos, wantLimit)
		}
		for i := range gr {
			wantScalars := [21]uint32{}
			for j := range wantScalars {
				wantScalars[j] = readSideInfoU32(t, r, record)
			}
			wantUnchanged := readSideInfoU32(t, r, record)
			var wantTable [40]uint8
			readSideInfoFixture(t, r, wantTable[:], record)
			gotScalars := sideInfoScalars(gr[i])
			if gotScalars != wantScalars {
				t.Fatalf("record %d granule %d: scalar fields got %v, want %v", record, i, gotScalars, wantScalars)
			}
			unchanged := sameSideInfoTable(gr[i].sfbTable, sentinel[:])
			gotUnchanged := uint32(0)
			if unchanged {
				gotUnchanged = 1
			}
			if wantUnchanged > 1 || gotUnchanged != wantUnchanged {
				t.Fatalf("record %d granule %d: sentinel-table alias=%d, want %d", record, i, gotUnchanged, wantUnchanged)
			}
			gotTable := normalizeSideInfoTable(gr[i].sfbTable, sentinel[:], unchanged)
			if gotTable != wantTable {
				t.Fatalf("record %d granule %d: normalized sfb table got %v, want %v", record, i, gotTable, wantTable)
			}
		}
		cases++
	}
	if got := hex.EncodeToString(h.Sum(nil)); got != wantSideInfoFixtureSHA256 {
		t.Fatalf("uncompressed fixture SHA256 %s, want %s", got, wantSideInfoFixtureSHA256)
	}
	if cases != wantSideInfoFixtureCases {
		t.Fatalf("fixture cases %d, want %d", cases, wantSideInfoFixtureCases)
	}
	t.Logf("validated %d frozen side-info vectors; uncompressed SHA256 %s", cases, wantSideInfoFixtureSHA256)
}

func checkSideInfoInvariants(t *testing.T) {
	t.Helper()
	for i, row := range scfLong {
		if bandSum(row[:]) != 576 || !zeroAfterTerminator(row[:]) {
			t.Fatalf("long scalefactor row %d is malformed", i)
		}
	}
	for name, rows := range map[string][8][40]uint8{"short": scfShort, "mixed": scfMixed} {
		for i, row := range rows {
			if bandSum(row[:]) != 576 || !zeroAfterTerminator(row[:]) {
				t.Fatalf("%s scalefactor row %d is malformed", name, i)
			}
		}
	}

	// All-zero MPEG-1 mono side information begins after a 16-bit CRC field.
	hMonoMPEG1 := [4]byte{0xff, 0xfb, 0x90, 0xc0}
	data := make([]byte, 64)
	var bs bitReader
	bsInit(&bs, data, len(data))
	bs.pos = 16
	var sentinel [40]uint8
	for i := range sentinel {
		sentinel[i] = 0xa5
	}
	var gr [4]grInfo
	for i := range gr {
		gr[i] = seededGrInfo(0xa5, sentinel[:])
	}
	if ret := readSideInfo(&bs, &gr, hMonoMPEG1); ret != 0 || bs.pos != 152 || bs.limit != 512 {
		t.Fatalf("zero side-info parse after CRC: ret=%d pos=%d limit=%d, want 0/152/512", ret, bs.pos, bs.limit)
	}
	if !sameSideInfoTable(gr[0].sfbTable, scfLong[5][:]) || !sameSideInfoTable(gr[1].sfbTable, scfLong[5][:]) ||
		!sameSideInfoTable(gr[2].sfbTable, sentinel[:]) || !sameSideInfoTable(gr[3].sfbTable, sentinel[:]) {
		t.Fatalf("zero side-info parse assigned an unexpected table row or touched an unused granule")
	}
	for i := 0; i < 2; i++ {
		if gr[i].subblockGain != [3]uint8{0xa5, 0xa5, 0xa5} {
			t.Fatalf("non-switched long block overwrote seeded subblock gains: %v", gr[i].subblockGain)
		}
		if gr[i].regionCount[2] != 255 {
			t.Fatalf("non-switched long block did not set region_count[2]: %d", gr[i].regionCount[2])
		}
	}

	// A valid short-block parse selects a specific row; a later big_values
	// error occurs before C assigns sfbtab and must retain that exact alias.
	hRate2 := [4]byte{0xff, 0xfb, 0x98, 0xc0}
	shortBits := mpeg1MonoSideBits(
		granuleSideBits{globalGain: 77, switched: true, blockType: 2},
		granuleSideBits{},
	)
	shortData := packSideInfoBits(shortBits, 64)
	gr[0].regionCount[2] = 0x6b
	bsInit(&bs, shortData, len(shortData))
	bs.pos = 16
	if ret := readSideInfo(&bs, &gr, hRate2); ret != 0 || !sameSideInfoTable(gr[0].sfbTable, scfShort[7][:]) || gr[0].regionCount[2] != 0x6b {
		t.Fatalf("short-block setup failed: ret=%d row-is-short-rate2=%v", ret, sameSideInfoTable(gr[0].sfbTable, scfShort[7][:]))
	}
	gr[0].globalGain = 77
	badBigValues := mpeg1MonoPrefixAndFirstGranule(7, 289)
	badData := packSideInfoBits(badBigValues, 64)
	bsInit(&bs, badData, len(badData))
	bs.pos = 16
	if ret := readSideInfo(&bs, &gr, hMonoMPEG1); ret != -1 {
		t.Fatalf("big_values=289 returned %d, want -1", ret)
	}
	if !sameSideInfoTable(gr[0].sfbTable, scfShort[7][:]) || gr[0].part23Length != 7 || gr[0].bigValues != 289 || gr[0].globalGain != 77 {
		t.Fatalf("big_values early return lost prior table alias or had incorrect partial writes: %+v", gr[0])
	}

	// A switched block with block_type=0 writes fields up to blockType, after
	// selecting the long table, then returns before mixed flags and regions.
	gr = [4]grInfo{}
	var sentinel5a [40]uint8
	for i := range sentinel5a {
		sentinel5a[i] = 0x5a
	}
	for i := range gr {
		gr[i] = seededGrInfo(0x5a, sentinel5a[:])
	}
	zeroBlock := mpeg1MonoPrefixAndBlockTypeZero()
	zeroBlockData := packSideInfoBits(zeroBlock, 64)
	bsInit(&bs, zeroBlockData, len(zeroBlockData))
	bs.pos = 16
	if ret := readSideInfo(&bs, &gr, hMonoMPEG1); ret != -1 {
		t.Fatalf("block_type=0 returned %d, want -1", ret)
	}
	if gr[0].part23Length != 3 || gr[0].globalGain != 77 || gr[0].blockType != 0 ||
		gr[0].mixedBlockFlag != 0x5a || gr[0].regionCount != [3]uint8{0x5a, 0x5a, 0x5a} ||
		!sameSideInfoTable(gr[0].sfbTable, scfLong[5][:]) {
		t.Fatalf("block_type=0 partial-write contract mismatch: %+v", gr[0])
	}
}

func sideInfoScalars(g grInfo) [21]uint32 {
	return [21]uint32{
		uint32(g.part23Length), uint32(g.bigValues), uint32(g.scalefacCompress),
		uint32(g.globalGain), uint32(g.blockType), uint32(g.mixedBlockFlag), uint32(g.nLongSFB), uint32(g.nShortSFB),
		uint32(g.tableSelect[0]), uint32(g.tableSelect[1]), uint32(g.tableSelect[2]),
		uint32(g.regionCount[0]), uint32(g.regionCount[1]), uint32(g.regionCount[2]),
		uint32(g.subblockGain[0]), uint32(g.subblockGain[1]), uint32(g.subblockGain[2]),
		uint32(g.preflag), uint32(g.scalefacScale), uint32(g.count1Table), uint32(g.scfsi),
	}
}

func seededGrInfo(seed uint8, sentinel []uint8) grInfo {
	word := uint16(seed) | uint16(seed)<<8
	return grInfo{
		sfbTable:     sentinel,
		part23Length: word, bigValues: word, scalefacCompress: word,
		globalGain: seed, blockType: seed, mixedBlockFlag: seed, nLongSFB: seed, nShortSFB: seed,
		tableSelect: [3]uint8{seed, seed, seed}, regionCount: [3]uint8{seed, seed, seed},
		subblockGain: [3]uint8{seed, seed, seed}, preflag: seed, scalefacScale: seed, count1Table: seed, scfsi: seed,
	}
}

func sameSideInfoTable(a, b []uint8) bool {
	return len(a) != 0 && len(b) != 0 && &a[0] == &b[0]
}

func normalizeSideInfoTable(table, sentinel []uint8, unchanged bool) [40]uint8 {
	var out [40]uint8
	if unchanged {
		copy(out[:], sentinel)
		return out
	}
	for i, value := range table {
		if i == len(out) {
			break
		}
		out[i] = value
		if value == 0 {
			break
		}
	}
	return out
}

func bandSum(values []uint8) uint32 {
	var sum uint32
	for _, value := range values {
		sum += uint32(value)
	}
	return sum
}

func zeroAfterTerminator(values []uint8) bool {
	seenZero := false
	for _, value := range values {
		if seenZero && value != 0 {
			return false
		}
		if value == 0 {
			seenZero = true
		}
	}
	return seenZero
}

type granuleSideBits struct {
	part23Length     uint32
	bigValues        uint32
	globalGain       uint32
	scalefacCompress uint32
	switched         bool
	blockType        uint32
	mixed            uint32
	tables           uint32
	region0          uint32
	region1          uint32
	subblock0        uint32
	subblock1        uint32
	subblock2        uint32
	preflag          uint32
	scalefacScale    uint32
	count1Table      uint32
}

type sideInfoBit struct {
	value uint32
	width int
}

func mpeg1MonoSideBits(a, b granuleSideBits) []sideInfoBit {
	bits := []sideInfoBit{{0, 9}, {0, 9}} // main_data_begin and SCFSI
	for _, g := range [...]granuleSideBits{a, b} {
		bits = append(bits,
			sideInfoBit{g.part23Length, 12}, sideInfoBit{g.bigValues, 9},
			sideInfoBit{g.globalGain, 8}, sideInfoBit{g.scalefacCompress, 4},
		)
		if g.switched {
			bits = append(bits, sideInfoBit{1, 1}, sideInfoBit{g.blockType, 2}, sideInfoBit{g.mixed, 1},
				sideInfoBit{g.tables, 10}, sideInfoBit{g.subblock0, 3}, sideInfoBit{g.subblock1, 3}, sideInfoBit{g.subblock2, 3})
		} else {
			bits = append(bits, sideInfoBit{0, 1}, sideInfoBit{g.tables, 15}, sideInfoBit{g.region0, 4}, sideInfoBit{g.region1, 3})
		}
		bits = append(bits, sideInfoBit{g.preflag, 1}, sideInfoBit{g.scalefacScale, 1}, sideInfoBit{g.count1Table, 1})
	}
	return bits
}

func mpeg1MonoPrefixAndFirstGranule(part23, bigValues uint32) []sideInfoBit {
	return []sideInfoBit{
		{0, 9}, {0, 9},
		{part23, 12}, {bigValues, 9},
	}
}

func mpeg1MonoPrefixAndBlockTypeZero() []sideInfoBit {
	return []sideInfoBit{
		{0, 9}, {0, 9},
		{3, 12}, {0, 9}, {77, 8}, {0, 4},
		{1, 1}, {0, 2},
	}
}

func packSideInfoBits(bits []sideInfoBit, bytes int) []byte {
	out := make([]byte, bytes)
	pos := 16 // model the CRC field already consumed by the production caller
	for _, field := range bits {
		for shift := field.width - 1; shift >= 0; shift-- {
			if (field.value>>uint(shift))&1 != 0 {
				out[pos>>3] |= 1 << uint(7-(pos&7))
			}
			pos++
		}
	}
	return out
}

func readSideInfoFixture(t *testing.T, r io.Reader, dst []byte, record int) {
	t.Helper()
	if _, err := io.ReadFull(r, dst); err != nil {
		t.Fatalf("record %d: truncated payload: %v", record, err)
	}
}

func readSideInfoU32(t *testing.T, r io.Reader, record int) uint32 {
	t.Helper()
	var b [4]byte
	readSideInfoFixture(t, r, b[:], record)
	return binary.LittleEndian.Uint32(b[:])
}
