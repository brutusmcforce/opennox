package legacy

import (
	"github.com/brutusmcforce/opennox/client/gui"
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/server"
	"unsafe"
)

func clientFrameSample() {
	// Preserve the old platform C.uint boundary before widening the timestamp.
	now := uint64(uint32(PlatformTicks()))
	last := memmap.PtrUint64(0x5D4594, 815740)
	cursor := memmap.PtrUint64(0x5D4594, 815756)
	*memmap.PtrUint64(0x5D4594, 815220+8*uintptr(uint32(*cursor))) = now - *last
	*cursor = (*cursor + 1) % 60
	*last = now
	dword_5d4594_815748++
}
func clientFrameAverage() {
	n := min(int32(dword_5d4594_815748), 60)
	average := uint64(33)
	if n > 10 {
		var sum uint64
		for i := int32(0); i < n; i++ {
			sum += memmap.Uint64(0x5D4594, 815220+8*uintptr(i))
		}
		average = sum / uint64(n)
	}
	*memmap.PtrUint64(0x587000, 91880) = average
}
func clientModalWindow(command int) {
	if command != 0 {
		win := GetClient().Cli().GUI.NewWindowRaw(nil, 552, 0, 0, int(nox_win_width), int(nox_win_height), nil)
		dword_5d4594_816412 = uint32(uintptr(unsafe.Pointer(win)))
		win.DrawData().BgColorVal = uint32(nox_color_black_2650656)
	} else if dword_5d4594_816412 != 0 {
		win := (*gui.Window)(unsafe.Pointer(uintptr(dword_5d4594_816412)))
		win.Destroy()
		dword_5d4594_816412 = 0
	}
}
func clientPlayerColors(pl *server.Player) {
	pl.InitColors()
	pl.Colors.UnkColor = uint32(nox_color_white_2523948)
}
func clientAllPlayerColors() {
	players := &GetServer().S().Players
	for pl := players.First(); pl != nil; pl = players.Next(pl) {
		clientPlayerColors(pl)
	}
}
func clientAnimationCachesLoad() {
	ConfusedBirdiesAnim = Nox_xxx_gLoadAnim("ConfusedBirdies")
	SphericalShieldAnim = Nox_xxx_gLoadAnim("SphericalShieldAnim")
}
func clientAnimationCachesClear() {
	ConfusedBirdiesAnim = nil
	SphericalShieldAnim = nil
}
