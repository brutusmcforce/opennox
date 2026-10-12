package legacy

import (
	"github.com/brutusmcforce/opennox/libs/types"
	"github.com/brutusmcforce/opennox/server"
	"math"
)

func geometryWallHorizontal(u *server.Object, center, previous *types.Pointf, bounds *[4]float32, edge *types.Pointf, length float32) int32 {
	lo := bounds[0]
	if bounds[0] <= edge.X {
		lo = edge.X
	}
	end := float64(length) + float64(edge.X)
	hi := float64(bounds[2])
	if float64(bounds[2]) >= end {
		hi = end
	}
	if float64(lo) > end || hi < float64(edge.X) {
		return 0
	}
	delta, prior := (float64(center.Y) - float64(edge.Y)), (float64(previous.Y) - float64(edge.Y))
	if prior*delta < 0 {
		if float32(prior) >= 0 {
			center.Y = float32(float64(edge.Y) + 2)
		} else {
			center.Y = float32(float64(edge.Y) - 2)
		}
		u.NewPos.X = float32((float64(center.X) - float64(center.Y)) * 0.70710677)
		u.NewPos.Y = float32((float64(center.X) + float64(center.Y)) * 0.70710677)
		bounds[1] = float32(float64(center.Y) - float64(u.Shape.Box.H)*0.5)
		bounds[3] = float32(float64(u.Shape.Box.H)*0.5 + float64(center.Y))
		delta = float64(center.Y - edge.Y)
	}
	if bounds[1] > edge.Y || bounds[3] < edge.Y {
		return 0
	}
	ratio := float32((hi - float64(lo)) / (float64(bounds[2]) - float64(bounds[0])))
	velocity := float32((float64(u.VelVec.Y) - float64(u.VelVec.X)) * 0.70710677)
	var depth float32
	var spring float64
	if delta >= 0 {
		wide := float64(edge.Y) - float64(bounds[1])
		depth = float32(wide)
		spring = float64(geometryWallForce()) * wide
	} else {
		wide := float64(bounds[3]) - float64(edge.Y)
		depth = float32(wide)
		spring = -float64(geometryWallForce()) * wide
	}
	spring32 := float32(spring)
	force := (math.Sqrt(float64(u.Mass)*float64(geometryWallForce())*4)*float64(-velocity)*0.5 + float64(spring32)) * float64(ratio)
	response := types.Pointf{float32(force * -0.70710677), float32(force * 0.70710677)}
	var nx float32
	var ny float64
	if depth >= 0 {
		nx = -0.70710677
		ny = 0.70710677
	} else {
		nx = 0.70710677
		ny = -0.70710677
	}
	inward := float64(nx)*float64(u.VelVec.X) + ny*float64(u.VelVec.Y)
	if inward < 0 {
		u.VelVec.X = float32(float64(u.VelVec.X) - inward*float64(nx))
		u.VelVec.Y = float32(float64(u.VelVec.Y) - inward*ny)
	}
	tangentX := -ny
	tx := float32(tangentX)
	speed := float32(tangentX*float64(u.VelVec.X) + float64(nx)*float64(u.VelVec.Y))
	response.X = float32(float64(response.X) - float64(u.Mass)*float64(speed)*float64(tx)*0.69999999)
	response.Y = float32(float64(response.Y) - float64(u.Mass)*float64(speed)*float64(nx)*0.69999999)
	u.Sub548600(response)
	geometryHit(u, nil, &response)
	return 1
}
func geometryWallVertical(u *server.Object, center, previous *types.Pointf, bounds *[4]float32, edge *types.Pointf, length float32) int32 {
	lo := bounds[1]
	if edge.Y >= bounds[1] {
		lo = edge.Y
	}
	end := float64(length) + float64(edge.Y)
	hi := float64(bounds[3])
	if float64(bounds[3]) >= end {
		hi = end
	}
	if float64(lo) > end || hi < float64(edge.Y) {
		return 0
	}
	delta, prior := (float64(center.X) - float64(edge.X)), (float64(previous.X) - float64(edge.X))
	if prior*delta < 0 {
		if float32(prior) >= 0 {
			center.X = float32(float64(edge.X) + 2)
		} else {
			center.X = float32(float64(edge.X) - 2)
		}
		u.NewPos.X = float32((float64(center.X) - float64(center.Y)) * 0.70710677)
		u.NewPos.Y = float32((float64(center.Y) + float64(center.X)) * 0.70710677)
		bounds[0] = float32(float64(center.X) - float64(u.Shape.Box.W)*0.5)
		bounds[2] = float32(float64(u.Shape.Box.W)*0.5 + float64(center.X))
		delta = float64(center.X - edge.X)
	}
	if bounds[0] > edge.X || bounds[2] < edge.X {
		return 0
	}
	ratio := float32((hi - float64(lo)) / (float64(bounds[3]) - float64(bounds[1])))
	velocity := float32((float64(u.VelVec.X) + float64(u.VelVec.Y)) * 0.70710677)
	var depth float32
	var spring float64
	if delta >= 0 {
		wide := float64(edge.X) - float64(bounds[0])
		depth = float32(wide)
		spring = float64(geometryWallForce()) * wide
	} else {
		wide := float64(bounds[2]) - float64(edge.X)
		depth = float32(wide)
		spring = -float64(geometryWallForce()) * wide
	}
	spring32 := float32(spring)
	force := (math.Sqrt(float64(u.Mass)*float64(geometryWallForce())*4)*float64(-velocity)*0.5 + float64(spring32)) * float64(ratio) * 0.70710677
	response := types.Pointf{float32(force), float32(force)}
	var nx float32
	var ny float64
	if depth >= 0 {
		nx = 0.70710677
		ny = 0.70710677
	} else {
		nx = -0.70710677
		ny = -0.70710677
	}
	inward := float64(nx)*float64(u.VelVec.X) + ny*float64(u.VelVec.Y)
	if inward < 0 {
		u.VelVec.X = float32(float64(u.VelVec.X) - inward*float64(nx))
		u.VelVec.Y = float32(float64(u.VelVec.Y) - inward*ny)
	}
	tangentX := -ny
	tx := float32(tangentX)
	speed := float32(tangentX*float64(u.VelVec.X) + float64(nx)*float64(u.VelVec.Y))
	response.X = float32(float64(response.X) - float64(u.Mass)*float64(speed)*float64(tx)*0.69999999)
	response.Y = float32(float64(response.Y) - float64(u.Mass)*float64(speed)*float64(nx)*0.69999999)
	u.Sub548600(response)
	geometryHit(u, nil, &response)
	return 1
}
