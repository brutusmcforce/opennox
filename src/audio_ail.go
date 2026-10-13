package opennox

import (
	"encoding/binary"
	"fmt"
	"math"
	"sort"
	"strings"
	"unsafe"

	"github.com/brutusmcforce/opennox/common/layout"

	"github.com/brutusmcforce/opennox/libs/ifs"

	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/internal/binfile"
	"github.com/brutusmcforce/opennox/legacy"
	"github.com/brutusmcforce/opennox/legacy/client/audio/ail"
	"github.com/brutusmcforce/opennox/legacy/common/alloc"
	"github.com/brutusmcforce/opennox/legacy/timer"
)

var (
	nox_enable_audio     = 1
	dword_5d4594_805980  *legacy.AudioStreamCatalog
	dword_5d4594_1193336 uint32
)

var (
	audioTimer93944 ail.Timer = math.MaxUint32
	audioDev        ail.Driver
)

func sub_43EFD0(a1 unsafe.Pointer) int {
	s := legacy.AudioVoiceSample(a1)
	s.Smp.End()
	if s.Flag7 == 0 {
		voice := s.Field1
		legacy.AudioStreamCallbackVoid(legacy.AudioVoiceEndCallback(voice), voice)
		s.Flag7 = 1
	}
	return 0
}

func sub_43EC30(a1 unsafe.Pointer) int {
	smp := audioDev.AllocateSample()
	s, _ := alloc.New(legacy.AudioSample{})
	legacy.AudioVoiceSetSample(a1, s)
	s.Dev = audioDev
	s.Field1 = a1
	s.Smp = smp
	b1, _ := alloc.Make([]byte{}, 16384)
	b2, _ := alloc.Make([]byte{}, 16384)
	s.Data1 = &b1[0]
	s.Data2 = &b2[0]
	if smp == 0 {
		return -2147221504 // 0x80040000
	}
	smp.SetUserData(s)
	return 0
}

func sub_43ECB0(a1 unsafe.Pointer) int {
	p := legacy.AudioVoiceSample(a1)
	if p.Smp != 0 {
		p.Smp.Release()
	}
	if p2 := p.Data1; p2 != nil {
		alloc.Free(p2)
	}
	if p2 := p.Data2; p2 != nil {
		alloc.Free(p2)
	}
	alloc.Free(p)
	return 0
}

func sub_43E940(a1 unsafe.Pointer) int {
	ail.Startup()
	audioTimer93944 = ail.RegisterTimer(func(u uint32) {
		sub_486EF0()
		legacy.MusicModule.Sub_43D2D0()
		(*timer.TimerGroup)(legacy.Get_dword_587000_127004()).ClearUpdated()
	})
	if audioTimer93944 == math.MaxUint32 {
		return -2147221504 // 0x80040000
	}
	audioTimer93944.SetFrequency(30)
	audioTimer93944.Start()

	*(*uint32)(unsafe.Pointer(&layout.As[legacy.AudioStreamDevice](a1).Slots)) = 1
	return 0
}

func sub_43E9F0() {
	if audioTimer93944 != math.MaxUint32 {
		audioTimer93944.Stop()
		audioTimer93944.Release()
		audioTimer93944 = math.MaxUint32
	}
	ail.Shutdown()
}

func noxAudioServeT(a1 int) {
	noxAudioServe()
}

func noxAudioServe() { // sub_43F1A0
	if audioDev != 0 {
		ail.Serve()
	}
}

func sub_43EA20(a1 unsafe.Pointer) int {
	//char v2[16]; // [esp+4h] [ebp-10h]

	v2 := sub_43EA90(unsafe.Pointer(&layout.As[legacy.AudioStreamContext](a1).RequestedFormat))
	audioDev = sub_43EAD0(v2)
	if audioDev == 0 {
		return -2147221504 // 0x80040000
	}
	*memmap.PtrInt32(0x587000, 93948) = 0
	layout.As[legacy.AudioStreamContext](a1).VoiceCapacity = 24
	return 0
}

type audioStruct1 struct {
	Field0  uint16
	Field2  uint16
	Field4  uint32
	Field8  uint32
	Field12 uint16
	Field14 uint16
}

func sub_43EA90(a2 unsafe.Pointer) *audioStruct1 {
	return &audioStruct1{
		Field0:  1,
		Field2:  *(*uint16)(unsafe.Pointer(&layout.As[legacy.AudioStreamFormat](a2).Channels)),
		Field4:  layout.As[legacy.AudioStreamFormat](a2).Rate,
		Field8:  layout.As[legacy.AudioStreamFormat](a2).ByteRate,
		Field12: *(*uint16)(unsafe.Pointer(&layout.As[legacy.AudioStreamFormat](a2).Width)) * *(*uint16)(unsafe.Pointer(&layout.As[legacy.AudioStreamFormat](a2).Channels)),
		Field14: 8 * *(*uint16)(unsafe.Pointer(&layout.As[legacy.AudioStreamFormat](a2).Width)),
	}
}

func sub_43EAD0(a1 *audioStruct1) ail.Driver {
	if dr := ail.WaveOutOpen(); dr != 0 {
		return dr
	}
	_ = ail.LastError()
	return 0
}

func sub_43EC10() int {
	if audioDev != 0 {
		audioDev.Close()
		audioDev = 0
	}
	return 0
}

func sub_43F130() ail.Driver {
	return audioDev
}

func sub_43ED00(voice unsafe.Pointer) int {
	s := legacy.AudioVoiceSample(voice)
	smp := s.Smp
	smp.Init()
	s.Field4 = legacy.AudioVoiceFormat(voice)
	v3 := int32(legacy.Sub_43F0E0(s.Field4))
	smp.SetType(v3, 0)
	if format := (*legacy.AudioStreamFormat)(s.Field4); format.Encoding == 2 {
		smp.SetADPCMBlockSize(format.Extra)
	}
	sub_43F060(voice)
	s.Flag7 = 0
	s.Field3 = 0
	smp.RegisterEOBCallback(func() {
		v := smp.UserData().(*legacy.AudioSample)
		legacy.Sub_43EE00(v)
	})
	smp.RegisterEOSCallback(func() {
		legacy.Sub_43EDB0(smp)
	})
	legacy.Sub_43EE00(s)
	return 0
}

func sub_43F060(voice unsafe.Pointer) int {
	s := legacy.AudioVoiceSample(voice)
	eff := legacy.AudioVoiceEffective(voice)
	smp := s.Smp
	smp.SetVolume(int((127 * (eff.Timers[0].Current >> 16)) >> 14))
	smp.SetPan(int((127 * (eff.Timers[2].Current >> 16)) >> 14))
	p2 := s.Field4
	v2 := legacy.Sub_486640(eff, int((*legacy.AudioStreamFormat)(p2).Rate))
	smp.SetPlaybackRate(v2)
	return 0
}

func nox_xxx_parseSoundSetBin_424170(path string) error {
	if legacy.Nox_xxx_parseSoundSetBin_424170(path) == 0 {
		return fmt.Errorf("failed to load sounds")
	}
	return nil
}

func nox_audio_initall(a3 int) int {
	if nox_enable_audio == 0 {
		return 1
	}
	if a3 != 0 {
		sub_486F30()
		if sub_4311F0() != 0 {
			legacy.Set_dword_587000_81128(unsafe.Pointer(&layout.As[legacy.AudioStreamContext](legacy.Get_dword_5d4594_805984()).Timers))
			dword_5d4594_805980 = sub_4866F0("audio", "audio")
		}
	}
	(*timer.TimerGroup)(memmap.PtrOff(0x5D4594, 805884)).Init()
	(*timer.TimerGroup)(legacy.Get_dword_587000_93164()).Init()
	(*timer.TimerGroup)(legacy.Get_dword_587000_122852()).Init()
	(*timer.TimerGroup)(legacy.Get_dword_587000_127004()).Init()
	legacy.Dialogs.Nox_xxx_WorkerHurt_44D810()
	legacy.MusicModule.Init()
	legacy.Sub_451850(legacy.Get_dword_5d4594_805984(), unsafe.Pointer(dword_5d4594_805980))
	v1 := configGetVolume(VolumeMusic)
	if v1 == 0 {
		legacy.Sub_43DC00()
	}
	(*timer.Timer)(legacy.Get_dword_587000_93164()).SetRaw(uint32(v1))
	v2 := configGetVolume(VolumeDialog)
	if v2 == 0 {
		legacy.Sub_44D960()
	}
	(*timer.Timer)(legacy.Get_dword_587000_122852()).SetRaw(uint32(v2))
	v3 := configGetVolume(VolumeFX)
	if v3 == 0 {
		legacy.Sub_453050()
	}
	(*timer.Timer)(legacy.Get_dword_587000_127004()).SetRaw(uint32(v3))
	return 1
}

func sub_4311F0() int {
	legacy.RegisterMilesAudioDevice()
	v2, free := alloc.Make([]int32{}, 7)
	defer free()
	v2[2] = 22050
	v2[1] = 0
	v2[3] = 2
	v2[4] = 2
	v2[0] = 4
	legacy.Sub_487D00(unsafe.Pointer(&v2[0]))
	v0 := legacy.Sub_487150(-1, unsafe.Pointer(&v2[0]))
	legacy.Set_dword_5d4594_805984(v0)
	return bool2int(v0 != nil && legacy.Sub_487790(v0, 16) == 16)
}

func sub_486F30() int {
	legacy.AudioStreamSystemReset()
	dword_5d4594_1193336 = 1
	return 0
}

func sub_486EF0() {
	if dword_5d4594_1193336 != 0 {
		legacy.AudioStreamSystemTick()
	}
}

// freeAudioCatalog closes the sound bag files and frees the catalog.
func freeAudioCatalog(p *legacy.AudioStreamCatalog) {
	if p == nil {
		return
	}
	if p.Bag != nil {
		legacy.Nox_fs_close(p.Bag)
		p.Bag = nil
	}
	if p.Override != nil {
		legacy.Nox_fs_close(p.Override)
		p.Override = nil
	}
	if p.Entries != nil {
		alloc.Free(p.Entries)
		p.Entries = nil
	}
	alloc.Free(p)
}

func sub_4866F0(path1 string, path2 string) *legacy.AudioStreamCatalog {
	idxPath := path1
	if i := strings.LastIndexByte(idxPath, '.'); i >= 0 {
		idxPath = idxPath[:i]
	}
	idxPath += ".idx"
	fi, err := ifs.Open(idxPath)
	if err != nil {
		return nil
	}
	bf1 := binfile.NewFile(fi)
	cf1 := legacy.NewFileHandle(bf1)
	defer legacy.Nox_fs_close(cf1)

	bagPath := path1
	if i := strings.LastIndexByte(bagPath, '.'); i >= 0 {
		bagPath = bagPath[:i]
	}
	bagPath += ".bag"
	fb, err := ifs.Open(bagPath)
	if err != nil {
		return nil
	}

	p, _ := alloc.New(legacy.AudioStreamCatalog{})
	p.Bag = legacy.NewFileHandle(binfile.NewFile(fb))
	if p.Bag == nil {
		freeAudioCatalog(p)
		return nil
	}
	var hdr [12]byte
	if _, err := bf1.Read(hdr[:12]); err != nil {
		freeAudioCatalog(p)
		return nil
	}
	_ = binary.LittleEndian.Uint32(hdr[0:])
	vers := binary.LittleEndian.Uint32(hdr[4:])
	p.Count = binary.LittleEndian.Uint32(hdr[8:])
	var arr []legacy.AudioStreamEntry
	// key returns the first 32 bytes of an entry (the name and the fields up to Flags), which the
	// catalog is sorted by, as in the original.
	key := func(e *legacy.AudioStreamEntry) []byte { return (*[32]byte)(unsafe.Pointer(e))[:] }
	if p.Count > 0 {
		arr, _ = alloc.Make([]legacy.AudioStreamEntry{}, p.Count)
		p.Entries = &arr[0]
		var buf [36]byte
		n := 32
		if vers != 1 {
			n = 36
		}
		for i := range arr {
			if _, err := bf1.Read(buf[:n]); err != nil {
				freeAudioCatalog(p)
				return nil
			}
			copy(key(&arr[i]), buf[:32])
			arr[i].Extra = 0
			if n == 36 {
				arr[i].Extra = binary.LittleEndian.Uint32(buf[32:])
			}
		}
	}
	sort.Slice(arr, func(i, j int) bool {
		s1, s2 := alloc.GoStringS(key(&arr[i])), alloc.GoStringS(key(&arr[j]))
		s1, s2 = strings.ToLower(s1), strings.ToLower(s2)
		return s1 < s2
	})
	p.UseOverride = 0
	if path2 != "" {
		alloc.StrCopyZero(p.Directory[:], path2)
		if i := alloc.StrLenS(p.Directory[:]); p.Directory[i] == '\\' {
			p.Directory[i] = 0
		}
		if info, err := ifs.Stat(alloc.GoStringS(p.Directory[:])); err == nil && info.IsDir() {
			p.UseOverride = 1
		}
		// TODO: strlen()-1 ?
		if i := alloc.StrLenS(p.Directory[:]); p.Directory[i] != '\\' {
			p.Directory[i] = '\\'
			p.Directory[i+1] = 0
		}
	}
	return p
}
