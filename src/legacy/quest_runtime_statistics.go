package legacy

import (
	"encoding/binary"
	"math"
	"unsafe"

	"github.com/brutusmcforce/opennox/common/layout"

	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/common/ntype"
	"github.com/brutusmcforce/opennox/server"
)

// questRuntimePlayerWord returns the quest statistics word of pl at 386 offset off (4652-4684).
func questRuntimePlayerWord(pl *server.Player, off int) *uint32 {
	return (*uint32)(unsafe.Add(unsafe.Pointer(pl), layoutServerPlayer.Off(uintptr(off))))
}
func questRuntimeReset(u *server.Object) uint32 {
	if u == nil {
		return 0
	}
	pl := controlPlayer(u)
	for off := 4652; off <= 4684; off += 4 {
		*questRuntimePlayerWord(pl, off) = 0
	}
	pl.Field4688 = memmap.Uint32(0x587000, 202028)
	pl.Field4692 = 63
	return memmap.PtrWord(unsafe.Pointer(pl))
}
func questRuntimeResetAll() uint32 {
	players := &GetServer().S().Players
	for u := players.FirstUnit(); u != nil; u = players.NextUnit(u) {
		questRuntimeReset(u)
	}
	return 0
}
func questRuntimeStageComplete(u *server.Object) uint32 {
	if u.ObjFlags&0x20 != 0 {
		return controlRaw(u)
	}
	pl := controlPlayer(u)
	if pl.Field4792 == 1 {
		*questRuntimePlayerWord(pl, 4652)++
		pl.Field4692 |= 1
	}
	return memmap.PtrWord(unsafe.Pointer(pl))
}
func questRuntimeIncrement(u *server.Object, off int, mask uint32) uint32 {
	if u == nil || u.ObjFlags&0x20 != 0 {
		return controlRaw(u)
	}
	pl := controlPlayer(u)
	*questRuntimePlayerWord(pl, off)++
	if off == 4672 {
		pl.Field4676++
	}
	pl.Field4692 |= mask
	return memmap.PtrWord(unsafe.Pointer(pl))
}
func questRuntimeCount() int {
	count := 0
	players := &GetServer().S().Players
	for u := players.FirstUnit(); u != nil; u = players.NextUnit(u) {
		if !worldQuestIgnoreHost(u) && controlPlayer(u).Field4792 == 1 {
			count++
		}
	}
	return count
}
func questRuntimeRoom() bool {
	count := 0
	players := &GetServer().S().Players
	for u := players.FirstUnit(); u != nil; u = players.NextUnit(u) {
		if !worldQuestIgnoreHost(u) && controlPlayer(u).Field4792 != 0 {
			count++
		}
	}
	return count < 6
}
func questRuntimeScore(a, b, c, stage uint32) uint32 {
	factor := float32(math.Pow(float64(stage), memmap.Float64(0x581450, 10088)))
	weighted := float64(a)*10 + float64(b)*35 + float64(c)*0.1
	return uint32(floatToInt32(float32(float64(factor) * weighted)))
}
func questRuntimePlayerScore(to int) uint32 {
	count := questRuntimeCount()
	pl := GetServer().S().Players.ByInd(ntype.PlayerInd(to))
	if pl == nil || pl.Field4792 == 0 {
		return 0
	}
	p := pl.C()
	score := func() uint32 {
		return questRuntimeScore(layout.As[server.Player](p).Field4668, layout.As[server.Player](p).Field4672, layout.As[server.Player](p).Field4664, layout.As[server.Player](p).Field4688)
	}
	var result uint32
	if count == 1 {
		result = score()
	} else {
		var maximum, a, b uint32
		stage := uint32(1)
		players := &GetServer().S().Players
		// C gates on the requesting player's participation here, already checked above;
		// it includes every unit in the aggregation, not only participating units.
		for u := players.FirstUnit(); u != nil; u = players.NextUnit(u) {
			q := controlPlayer(u)
			v := questRuntimeScore(q.Field4668, q.Field4672, 0, q.Field4688)
			if v > maximum {
				maximum = v
			}
			a += q.Field4668
			b += q.Field4672
			stage = q.Field4688
		}
		total := questRuntimeScore(a, b, 0, stage)
		ratio := float32(1)
		if maximum > 0 {
			ratio = float32(float64(total) / float64(maximum))
		}
		result = uint32(floatToInt32(float32(float64(score()) * float64(ratio))))
	}
	if result > 999999999 {
		result = 999999999
	}
	return result
}
func questRuntimeScoreboard(to int) int {
	pl := GetServer().S().Players.ByInd(ntype.PlayerInd(to))
	data := make([]byte, 90)
	data[0] = 240
	data[1] = 12
	binary.LittleEndian.PutUint16(data[2:], uint16(memmap.Uint32(0x5D4594, 1556132)))
	count := 0
	players := &GetServer().S().Players
	for u := players.FirstUnit(); u != nil; u = players.NextUnit(u) {
		p := controlPlayer(u)
		if p.Field4792 == 1 && count < 6 {
			at := 6 + 14*count
			binary.LittleEndian.PutUint16(data[at:], uint16(u.NetCode))
			for i, off := range []int{4668, 4672, 4680, 4664} {
				binary.LittleEndian.PutUint16(data[at+2+2*i:], uint16(*equipmentWord(unsafe.Pointer(p), off)))
			}
			binary.LittleEndian.PutUint32(data[at+10:], questRuntimePlayerScore(int(p.PlayerInd)))
			count++
			binary.LittleEndian.PutUint16(data[4:], uint16(pl.Field4688))
		}
	}
	return gameplayReportSend(to, data, true, 1)
}
func questRuntimeClearRelations(u *server.Object, slot int) {
	d := u.UpdateDataPlayer()
	d.QuestLevels[slot] = 0
	d.Field81[slot] = 0
	d.SilverKeys[slot] = 0
	d.GoldKeys[slot] = 0
}
func questRuntimeReconnect(u *server.Object) uint32 {
	d := u.UpdateData
	gameplayReportQuestObject(255, u)
	questRuntimeReset(u)
	players := &GetServer().S().Players
	for it := players.FirstUnit(); it != nil; it = players.NextUnit(it) {
		questRuntimeClearRelations(it, int(layout.As[server.PlayerUpdateData](d).Player.PlayerInd))
	}
	return 0
}
func questRuntimeDepartureStamp(slot int) int {
	*memmap.PtrUint32(0x5D4594, 1556172+4*uintptr(slot)) = GetServer().S().Frame()
	return slot
}
func questRuntimeDepartureTick() int {
	core := GetServer().S()
	for i := 0; i < 32; i++ {
		p := memmap.PtrUint32(0x5D4594, 1556172+4*uintptr(i))
		pl := core.Players.ByInd(ntype.PlayerInd(i))
		if pl != nil && pl.PlayerUnit != nil && pl.Active != 0 && pl.Field4792 == 1 {
			*p = 0
		} else if *p != 0 && core.Frame()-*p > 30*uint32(core.TickRate()) {
			for u := core.Players.FirstUnit(); u != nil; u = core.Players.NextUnit(u) {
				questRuntimeClearRelations(u, i)
			}
			*p = 0
		}
	}
	return 32
}
func questRuntimeDepartureReset() uint32 {
	clear(unsafe.Slice(memmap.PtrUint32(0x5D4594, 1556172), 32))
	return 0
}
