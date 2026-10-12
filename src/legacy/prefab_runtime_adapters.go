package legacy

import (
	"math"
)

func sub_502DF0() uint32 { return prefabClose() }

func sub_5044B0(a0 int32, a1 float32, a2 float32) *uint32 {
	return (*uint32)(mapRoomPointer(prefabWaypointNew(uint32(a0), math.Float32frombits(math.Float32bits(float32(a1))), math.Float32frombits(math.Float32bits(float32(a2))))))
}

func nox_xxx_unitAddToList_5048A0(a0 int32) *uint32 {
	return (*uint32)(mapRoomPointer(prefabObjectNew(uint32(a0))))
}

func sub_51D0E0() { prefabResetWaypoint() }
