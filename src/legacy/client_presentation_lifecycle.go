package legacy

import (
	"github.com/brutusmcforce/opennox/common/memmap"
	"strings"
	"unsafe"
)

func clientServerAddressCopy(address string) {
	dst := unsafe.Slice((*byte)(memmap.PtrOff(0x5D4594, 806060)), 23)
	if end := strings.IndexByte(address, 0); end >= 0 {
		address = address[:end]
	}
	clear(dst)
	copy(dst, address)
}
func clientAudioContextDestroy() {
	if p := Get_dword_5d4594_805984(); p != nil {
		audioStreamContextDestroy((*audioStreamContext)(p))
		Set_dword_5d4594_805984(nil)
	}
}
func clientAudioContextStop() {
	if p := Get_dword_5d4594_805984(); p != nil {
		audioStreamVoiceStopKind((*audioStreamContext)(p), -1)
	}
}

var clientRandomNames = [...]string{"Dweezle", "Glork", "Floogle", "Goombah", "Kraun", "Kloog", "Zurg", "Darg", "Arfingle", "Buurl", "Gurgin", "Grok", "Hurlong", "Luric", "Lupis", "Mallik", "Thrall", "Norwood", "Nulik", "Orin", "Olaf", "Orguk", "Pervis", "Paavik", "Qix", "Xevin", "Xurcon", "Markoan", "Yuric", "Yoovis", "Yalek", "Zug", "Zivik"}

func clientRandomName() string {
	count := memmap.PtrUint32(0x5D4594, 814516)
	if *count == 0 {
		*count = uint32(len(clientRandomNames))
	}
	return clientRandomNames[GetServer().S().Rand.Other.Int(0, int(*count)-1)]
}
