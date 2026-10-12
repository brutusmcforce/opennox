package legacy

import (
	"encoding/binary"
	"github.com/brutusmcforce/opennox/common/ntype"
	"github.com/brutusmcforce/opennox/internal/netlist"
	"github.com/brutusmcforce/opennox/server"
)

var sessionOperators []*server.Player
var sessionOperatorsInitialized bool
var sessionLatencyPlayer *server.Player

func sessionRosterInit() {
	if !sessionOperatorsInitialized {
		sessionOperators = nil
		sessionOperatorsInitialized = true
	}
}
func sessionRosterClear() { sessionOperators = nil }
func sessionRosterContains(index int32) int {
	for _, pl := range sessionOperators {
		if int32(pl.PlayerInd) == index {
			return 1
		}
	}
	return 0
}
func sessionRosterAdd(index int32) {
	if sessionRosterContains(index) != 0 {
		return
	}
	if pl := GetServer().S().Players.ByInd(ntype.PlayerInd(index)); pl != nil {
		sessionOperators = append(sessionOperators, pl)
	}
}
func sessionRosterRemove(index int32) {
	for i, pl := range sessionOperators {
		if int32(pl.PlayerInd) == index {
			copy(sessionOperators[i:], sessionOperators[i+1:])
			sessionOperators[len(sessionOperators)-1] = nil
			sessionOperators = sessionOperators[:len(sessionOperators)-1]
			return
		}
	}
}
func sessionReportLatency() {
	s := GetServer().S()
	if sessionLatencyPlayer != nil {
		sessionLatencyPlayer = s.Players.Next(sessionLatencyPlayer)
	}
	if sessionLatencyPlayer == nil {
		sessionLatencyPlayer = s.Players.First()
	}
	if sessionLatencyPlayer == nil {
		return
	}
	for i := 0; sessionLatencyPlayer.PlayerInd != 31 && i < 32; i++ {
		if Sub_554240(ntype.PlayerInd(sessionLatencyPlayer.PlayerInd)) != 0 {
			break
		}
		sessionLatencyPlayer = s.Players.Next(sessionLatencyPlayer)
		if sessionLatencyPlayer == nil {
			sessionLatencyPlayer = s.Players.First()
		}
	}
	if sessionLatencyPlayer == nil {
		return
	}
	var data [5]byte
	data[0] = 215
	binary.LittleEndian.PutUint16(data[1:], uint16(sessionLatencyPlayer.NetCodeVal))
	binary.LittleEndian.PutUint16(data[3:], uint16(Sub_554240(ntype.PlayerInd(sessionLatencyPlayer.PlayerInd))))
	for pl := s.Players.First(); pl != nil; pl = s.Players.Next(pl) {
		s.NetList.AddToMsgListCli(ntype.PlayerInd(pl.PlayerInd), netlist.Kind1, data[:])
	}
}
