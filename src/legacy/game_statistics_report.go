package legacy

import (
	"github.com/brutusmcforce/opennox/common/memmap"
	"unsafe"
)

func statisticsReport(report unsafe.Pointer, mode int, eventCount int16) []byte {
	records := statisticsRecords{}
	for _, f := range []struct {
		tag  string
		kind uint16
		off  int
	}{
		{"MXPL", 6, 8}, {"IDNO", 6, 12}, {"GSKU", 6, 16}, {"GSTY", 6, 20}, {"CLGM", 2, 24},
		{"LIMT", 6, 32}, {"TLMT", 6, 36}, {"RSTC", 2, 40}, {"MINE", 6, 44}, {"MAXE", 6, 48},
		{"MINP", 6, 52}, {"MAXP", 6, 56}, {"VIDM", 2, 60}, {"SVRS", 2, 61}, {"NTMS", 2, 25},
		{"SCEN", 7, 96}, {"GNAM", 7, 352}, {"SPL1", 6, 64}, {"SPL2", 6, 68}, {"SPL3", 6, 72},
		{"ARMR", 6, 88}, {"WPN1", 2, 84}, {"WPN2", 2, 85}, {"WPN3", 2, 86}, {"STAF", 6, 92},
		{"DURA", 6, 28}, {"FINI", 2, -1}, {"TRNY", 2, 26},
	} {
		if f.kind == 7 {
			records.add(7, f.tag, 0, []byte(statisticsString(unsafe.Add(report, f.off))))
			continue
		}
		value := int32(1)
		if f.off >= 0 {
			value = int32(*controlByte(report, f.off))
		}
		records.add(f.kind, f.tag, value, nil)
	}
	sequence := memmap.PtrUint32(0x5D4594, 741668)
	switch mode {
	case 0:
		*sequence = 0
		records.add(6, "SEQU", 0, nil)
		records.add(2, "ENDF", 0, nil)
		records.add(3, statisticsTag(71872), -1, nil)
	case 1, 2:
		*sequence++
		records.add(6, "SEQU", int32(*sequence), nil)
		records.add(2, "ENDF", int32(mode-1), nil)
		count := int16(*controlHalf(report, 6))
		tagOff, eventOff := uintptr(71928), uintptr(71936)
		if mode == 2 {
			tagOff = 71896
			eventOff = 71904
		}
		records.add(3, statisticsTag(tagOff), int32(count), nil)
		*memmap.PtrUint32(0x5D4594, 741660) = 0
		for i := 0; i < int(count); i++ {
			*memmap.PtrUint8(0x587000, 71491) = byte(i + 48)
			tag := "LGL?"
			if mode == 1 {
				tag = statisticsTag(71488)
			}
			records.add(7, tag, 0, []byte(statisticsString(*controlPtr(*controlPtr(report, 608), 4*i))))
			for _, f := range []struct {
				tag    string
				kind   uint16
				array  int
				suffix uintptr
			}{
				{"IPL?", 6, 612, 71499}, {"CNL?", 6, 616, 71515}, {"CLL?", 2, 620, 71507},
				{"CMP?", 2, 624, 71523}, {"DUR?", 6, 628, 71531}, {"PAR?", 2, 632, 71539},
			} {
				*memmap.PtrUint8(0x587000, f.suffix) = byte(i + 48)
				index := i
				if f.kind == 6 {
					index *= 4
				}
				value := int32(*controlByte(*controlPtr(report, f.array), index))
				records.add(f.kind, f.tag, value, nil)
			}
			*memmap.PtrUint32(0x5D4594, 741660) = uint32(i + 1)
		}
		records.add(20, statisticsTag(eventOff), 0, unsafe.Slice((*byte)(*controlPtr(report, 636)), int(uint16(2*int32(eventCount)))))
	}
	return records.serialize()
}
