package legacy

import (
	"encoding/binary"
	"github.com/brutusmcforce/opennox/libs/types"
	noxflags "github.com/brutusmcforce/opennox/common/flags"
	"github.com/brutusmcforce/opennox/common/ntype"
	"github.com/brutusmcforce/opennox/internal/netlist"
	"github.com/brutusmcforce/opennox/server"
	"math"
)

func motionSentryUpdate(u *server.Object) uint32 {
	result := uint32(u.ObjFlags)
	data := (*[3]float32)(u.UpdateData)
	if int32(result) >= 0 {
		u.Field125 = nil
		u.InvNextItem = motionObject(motionSentryHead)
		if u.InvNextItem != nil {
			u.InvNextItem.Field125 = u
		}
		motionSentryHead = motionAddress(u)
		u.ObjFlags |= 0x80000000
		result = uint32(u.ObjFlags)
	}
	if u.ObjFlags&0x20 != 0 {
		result = motionSentryRemove(u)
	}
	if u.ObjFlags&0x1000000 != 0 {
		start := u.PosVec
		end := types.Pointf{X: float32(math.Cos(float64(data[0]))*600 + float64(start.X)), Y: float32(math.Sin(float64(data[0]))*600 + float64(start.Y))}
		var hit types.Pointf
		if !GetServer().S().MapTraceRayAt(start, end, &hit, nil, 5) {
			end = hit
		}
		u.Pos39 = end
		data[0] = float32(float64(data[2]) + float64(data[0]))
		// Keep comparison order for unordered coordinates, as in the C rectangle.
		min, max := start, end
		if start.X >= end.X {
			min.X, max.X = end.X, start.X
		}
		if start.Y >= end.Y {
			min.Y, max.Y = end.Y, start.Y
		}
		ray := [4]float32{start.X, start.Y, end.X, end.Y}
		GetServer().S().Map.EachObjInRect(types.Rectf{Min: min, Max: max}, func(t *server.Object) bool { motionSentryContact(t, u, &ray); return true })
	} else {
		data[0] = data[1]
	}
	return result
}
func motionSentryContact(target, sentry *server.Object, ray *[4]float32) {
	f := target.ObjFlags
	if f&0x41 != 0 || f&0x10 != 0 && !(noxflags.HasGame(noxflags.GameModeQuest) && target.ObjClass&2 != 0 && f&0x8000 == 0) || target.HealthData == nil {
		return
	}
	var p types.Pointf
	if !projectLine(ray, &target.PosVec, &p) {
		return
	}
	dx, dy := float64(target.PosVec.X)-float64(p.X), float64(target.PosVec.Y)-float64(p.Y)
	radius := float64(target.Shape.Circle.R)
	if radius*radius > dy*dy+dx*dx {
		_ = target.CallDamageValue(sentry.FindOwnerChainPlayer(), sentry, 500, 16)
		GetServer().S().Audio.EventObj(298, target, 0, 0)
	}
}
func motionSentryReport(player int32) {
	pl := GetServer().S().Players.ByInd(ntype.PlayerInd(player))
	x, y := pl.Pos3632Vec.X, pl.Pos3632Vec.Y
	left, right := float32(float64(x)-float64(pl.Field10)), float32(float64(x)+float64(pl.Field10))
	top, bottom := float32(float64(y)-float64(pl.Field12)), float32(float64(y)+float64(pl.Field12))
	for p := motionSentryHead; p != 0; {
		u := motionObject(p)
		if u.ObjFlags&0x1000000 != 0 {
			loX, hiX := u.PosVec.X, u.Pos39.X
			if u.PosVec.X >= u.Pos39.X {
				loX, hiX = u.Pos39.X, u.PosVec.X
			}
			loY, hiY := u.PosVec.Y, u.Pos39.Y
			if u.PosVec.Y >= u.Pos39.Y {
				loY, hiY = u.Pos39.Y, u.PosVec.Y
			}
			if left < hiX && right > loX && top < hiY && bottom > loY {
				motionSentryPacket(player, u)
			}
		} else {
			d := (*[3]uint32)(u.UpdateData)
			d[0] = d[1]
		}
		p = motionAddress(u.InvNextItem)
	}
}
func motionSentryPacket(player int32, u *server.Object) int32 {
	msg := [9]byte{0x95}
	for i, v := range [4]float32{u.PosVec.X, u.PosVec.Y, u.Pos39.X, u.Pos39.Y} {
		binary.LittleEndian.PutUint16(msg[1+i*2:], uint16(worldRoundScratch(float32(v))))
	}
	return int32(bool2int(GetServer().S().NetList.AddToMsgListCli(ntype.PlayerInd(player), netlist.Kind(1), msg[:])))
}
