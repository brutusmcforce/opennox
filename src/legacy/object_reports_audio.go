package legacy

import (
	"image"
	"unsafe"

	"github.com/brutusmcforce/opennox/libs/types"
	"github.com/brutusmcforce/opennox/server"
)

func objectReportPolygon(p types.Pointf) int8 {
	level, _ := minimapPolygonLevel(image.Pt(int(floatToInt32(p.X)), int(floatToInt32(p.Y))))
	return int8(level)
}
func objectReportSoundLevel(p types.Pointf, u *server.Object) int {
	if u != nil {
		if uint32(u.ObjClass)&4 != 0 {
			if level := *(*int8)(unsafe.Pointer(&u.UpdateDataPlayer().Player.Field3668)); level != 0 {
				return int(level)
			}
		} else if uint32(u.ObjClass)&2 != 0 {
			poly := mapPolygonGet(*(*uint32)(u.UpdateData))
			if poly != nil {
				if level := *(*int8)(unsafe.Pointer(&poly.Level)); level != 0 {
					return int(level)
				}
			}
		}
	}
	return int(objectReportPolygon(p))
}
func objectReportRemoteAudio(viewer *server.Object) {
	data := viewer.UpdateData
	pl := viewer.UpdateDataPlayer().Player
	var level int8
	if follow := pl.CameraFollowObj; pl.Field3680&3 != 0 && follow != nil {
		if uint32(follow.ObjClass)&4 != 0 {
			level = *(*int8)(unsafe.Pointer(&follow.UpdateDataPlayer().Player.Field3668))
		} else {
			level = objectReportPolygon(follow.PosVec)
		}
	} else {
		if pl.Field3664 == 0xdeadface {
			Nox_xxx_questCheckSecretArea_421C70(viewer)
		}
		level = *(*int8)(unsafe.Pointer(&viewer.UpdateDataPlayer().Player.Field3668))
	}
	GetServer().NetUpdateRemotePlrAudioEvents(viewer, data, level)
}
