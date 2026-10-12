package legacy

import (
	"encoding/binary"
	"unsafe"

	"github.com/brutusmcforce/opennox/libs/types"
	"github.com/brutusmcforce/opennox/internal/netlist"
	"github.com/brutusmcforce/opennox/server"
)

func visibilityFXSend(pos types.Pointf, packet []byte) int {
	GetServer().S().Nox_xxx_netSendFxAllCli_523030(pos, packet)
	return 0
}
func visibilityFXPoint(code byte, pos types.Pointf) int {
	b := []byte{code, 0, 0, 0, 0}
	binary.LittleEndian.PutUint16(b[1:], uint16(floatToInt32(pos.X)))
	binary.LittleEndian.PutUint16(b[3:], uint16(floatToInt32(pos.Y)))
	return visibilityFXSend(pos, b)
}
func visibilityFXPointExtra(code, extra byte, pos types.Pointf) int {
	b := []byte{code, extra, 0, 0, 0, 0}
	binary.LittleEndian.PutUint16(b[2:], uint16(floatToInt32(pos.X)))
	binary.LittleEndian.PutUint16(b[4:], uint16(floatToInt32(pos.Y)))
	return visibilityFXSend(pos, b)
}
func visibilityFXSpark(pos types.Pointf, extra byte) int {
	b := []byte{147, 0, 0, 0, 0, extra}
	binary.LittleEndian.PutUint16(b[1:], uint16(floatToInt32(pos.X)))
	binary.LittleEndian.PutUint16(b[3:], uint16(floatToInt32(pos.Y)))
	return visibilityFXSend(pos, b)
}
func visibilityFXGeneratorBreak(pos types.Pointf, extra byte) {
	b := []byte{240, 25, 0, 0, 0, 0, extra}
	binary.LittleEndian.PutUint16(b[2:], uint16(floatToInt32(pos.X)))
	binary.LittleEndian.PutUint16(b[4:], uint16(floatToInt32(pos.Y)))
	visibilityFXSend(pos, b)
}
func visibilityFXVampire(code byte, words [4]int32, id uint16) int {
	b := []byte{code}
	for _, v := range words {
		b = binary.LittleEndian.AppendUint16(b, uint16(v))
	}
	b = binary.LittleEndian.AppendUint16(b, id)
	return visibilityFXSend(types.Pointf{X: float32(uint16(words[2])), Y: float32(uint16(words[3]))}, b)
}
func visibilityFXBroadcast(b []byte) int {
	s := GetServer().S()
	for u := s.Players.FirstUnit(); u != nil; u = s.Players.NextUnit(u) {
		s.NetList.AddToMsgListCli(u.UpdateDataPlayer().Player.PlayerIndex(), netlist.Kind1, b)
	}
	return 0
}
func visibilityFXPrediction(u *server.Object) int {
	b := []byte{181}
	b = binary.LittleEndian.AppendUint16(b, gameplayReportCode(u))
	b = binary.LittleEndian.AppendUint16(b, u.TypeInd)
	b = binary.LittleEndian.AppendUint16(b, uint16(int64(u.PosVec.X)))
	b = binary.LittleEndian.AppendUint16(b, uint16(int64(u.PosVec.Y)))
	b = binary.LittleEndian.AppendUint16(b, uint16(u.Direction1))
	b = append(b, byte(int64(float64(u.Float28)*16)), byte(int64(float64(u.VelVec.X)*16)), byte(int64(float64(u.VelVec.Y)*16)))
	return visibilityFXBroadcast(b)
}
func visibilityFXShield(u *server.Object, pos *types.Pointf) int {
	direction := int32(int16(u.Direction1))
	if pos != nil {
		delta := types.Pointf{X: float32(float64(u.PosVec.X) - float64(pos.X)), Y: float32(float64(u.PosVec.Y) - float64(pos.Y))}
		direction = int32(int32(geometryVectorAngle((*types.Pointf)(unsafe.Pointer(unsafe.Pointer(&delta))))))
	}
	b := []byte{128}
	b = binary.LittleEndian.AppendUint16(b, gameplayReportCode(u))
	b = append(b, byte(int32(geometryDirection4Index(int32(direction)))))
	return visibilityFXSend(u.PosVec, b)
}
func visibilityFXSummonStart(id uint16, pos types.Pointf, extra byte, a, b uint16) int {
	packet := []byte{126}
	packet = binary.LittleEndian.AppendUint16(packet, uint16(int64(pos.X)))
	packet = binary.LittleEndian.AppendUint16(packet, uint16(int64(pos.Y)))
	packet = binary.LittleEndian.AppendUint16(packet, id)
	packet = binary.LittleEndian.AppendUint16(packet, a)
	packet = append(packet, extra)
	packet = binary.LittleEndian.AppendUint16(packet, b)
	return gameplayReportSend(255, packet, false, 1)
}
func visibilityFXSummonCancel(id uint16) int {
	return gameplayReportSend(255, []byte{127, byte(id), byte(id >> 8)}, false, 1)
}
func visibilityFXGeneratorSpawn(words [4]int32, id uint16) {
	b := []byte{240, 16}
	for _, v := range words {
		b = binary.LittleEndian.AppendUint16(b, uint16(v))
	}
	b = binary.LittleEndian.AppendUint16(b, id)
	visibilityFXSend(types.Pointf{X: float32(words[2]), Y: float32(words[3])}, b)
}
func visibilityFXArrowTrap(pos types.Pointf, extra byte) {
	b := []byte{161, 0, 0, 0, 0, extra}
	binary.LittleEndian.PutUint16(b[1:], uint16(floatToInt32(pos.X)))
	binary.LittleEndian.PutUint16(b[3:], uint16(floatToInt32(pos.Y)))
	visibilityFXSend(pos, b)
}
