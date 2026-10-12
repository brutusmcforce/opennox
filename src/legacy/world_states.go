package legacy

import (
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/server"
	"unsafe"
)

func worldAngleQueue(ud unsafe.Pointer) {
	if *equipmentWord(ud, 28) == 0 {
		*equipmentWord(ud, 36) = uint32(collisionAngleHead)
		collisionAngleHead = memmap.PtrWord(ud)
		*equipmentWord(ud, 28) = 1
	}
}
func worldAngle(u *server.Object, delta int16) {
	p := (*uint16)(unsafe.Add(u.UpdateData, 40))
	// The original unsigned-short loops reduce the wrapped sum to its low byte.
	*p = uint16(uint16(*p)+uint16(delta)) & 255
	worldAngleQueue(u.UpdateData)
}
func worldDoor(u *server.Object) byte {
	ud := u.UpdateData
	desired, prior, current := int32(*equipmentWord(ud, 4)), int32(*equipmentWord(ud, 8)), int32(*equipmentWord(ud, 12))
	sound := func(closed bool) {
		id := 237
		if !closed {
			id = 239
		}
		if u.ObjSubClass&4 != 0 {
			if u.Material&8 != 0 {
				if closed {
					id = 245
				} else {
					id = 246
				}
			} else {
				if closed {
					id = 241
				} else {
					id = 243
				}
			}
		} else if u.ObjSubClass&1 != 0 {
			if closed {
				id = 247
			} else {
				id = 248
			}
		} else if u.ObjSubClass&0x1000 != 0 {
			if closed {
				id = 1014
			} else {
				id = 1015
			}
		}
		inventorySound(id, u, 0, 0)
	}
	if prior == desired {
		if current != desired {
			sound(true)
		}
	} else if current == desired {
		sound(false)
		GetServer().S().Objs.RemoveFromUpdatable(u)
	}
	current = int32(*equipmentWord(ud, 12))
	if *equipmentWord(ud, 8) != uint32(current) {
		u.NeedSync()
		current = int32(*equipmentWord(ud, 12))
		*equipmentWord(ud, 8) = uint32(current)
	}
	core := GetServer().S()
	if u.ObjFlags&0x1000000 != 0 && core.Frame()-*equipmentWord(ud, 44) > uint32(int32(core.TickRate())>>1) {
		desired = int32(*equipmentWord(ud, 4))
		if current != desired {
			diff := current - desired
			if diff < 0 {
				diff += 32
			}
			if diff >= 16 {
				worldAngle(u, 2)
			} else {
				worldAngle(u, -2)
			}
			return byte(collisionActivate(u))
		}
	}
	return byte(current)
}
func worldAnimate(u *server.Object, frame int) {
	stateAnimation(u, uint32(int32(frame)))
}
func worldScript(u *server.Object, offset int, caller *server.Object, event int) {
	GetServer().NoxScriptC().ScriptCallback((*server.ScriptCallback)(unsafe.Add(u.UpdateData, offset)), caller, u, server.ScriptEventType(event))
}
func worldToggle(u *server.Object) byte {
	ud := u.UpdateData
	bits := equipmentWord(ud, 0)
	state := (*byte)(unsafe.Add(ud, 8))
	core := GetServer().S()
	if u.ObjFlags&0x1000000 == 0 {
		*bits &^= 9
		return byte(*bits)
	}
	if *bits&8 == 0 {
		*state = 0
		u.Field34 = core.Frame()
		worldAnimate(u, 0)
	}
	switch *state {
	case 0:
		if core.Frame() > u.Field34 && *bits&1 != 0 {
			inventorySound(int(*equipmentWord(ud, 36)), u, 0, 0)
			worldAnimate(u, 1)
			worldScript(u, 20, *temporaryRefWord(ud, 4), 10)
			*state = 3
			u.Field34 = core.Frame() + uint32(core.TickRate())
		}
	case 1:
		if core.Frame() > u.Field34 && *bits&1 != 0 {
			inventorySound(int(*equipmentWord(ud, 40)), u, 0, 0)
			worldAnimate(u, 0)
			worldScript(u, 28, nil, 11)
			*state = 0
			if *bits&2 != 0 {
				*state = 5
			}
			u.Field34 = core.Frame() + uint32(core.TickRate())
		}
	case 3:
		if core.Frame() > u.Field34 && *bits&1 == 0 {
			*state = 1
		}
	}
	result := *bits
	if result&1 != 0 {
		result |= 4
	} else {
		result &^= 4
	}
	*bits = result&^1 | 8
	return byte(result)
}
func worldTrigger(u *server.Object) byte {
	core := GetServer().S()
	ud := u.UpdateData
	bits := equipmentWord(ud, 0)
	state := (*byte)(unsafe.Add(ud, 8))
	trigger, plate := memmap.PtrUint32(0x5d4594, 2488680), memmap.PtrUint32(0x5d4594, 2488676)
	if *trigger == 0 {
		*trigger = uint32(core.Types.IndByID("Trigger"))
		*plate = uint32(core.Types.IndByID("PressurePlate"))
	}
	delay := uint32(core.TickRate())
	if uint32(u.TypeInd) == *trigger || uint32(u.TypeInd) == *plate {
		delay = 0
	}
	if u.ObjFlags&0x1000000 == 0 {
		*bits &^= 9
		return byte(*bits)
	}
	if *bits&8 == 0 {
		*state = 0
		u.Field34 = core.Frame()
		worldAnimate(u, 0)
	}
	if *state == 0 {
		if *bits&1 != 0 {
			inventorySound(int(*equipmentWord(ud, 36)), u, 0, 0)
			worldAnimate(u, 1)
			worldScript(u, 20, *temporaryRefWord(ud, 4), 8)
			*state = 1
			u.Field34 = delay + core.Frame()
		}
	} else if *state == 1 && core.Frame() > u.Field34 && *bits&1 == 0 {
		inventorySound(int(*equipmentWord(ud, 40)), u, 0, 0)
		worldAnimate(u, 0)
		worldScript(u, 28, nil, 9)
		*state = 0
		if *bits&2 != 0 {
			*state = 5
		}
	}
	result := *bits
	if result&1 != 0 {
		result |= 4
	} else {
		result &^= 4
	}
	result = result&^1 | 8
	*bits = result
	return byte(result)
}
func worldEnabledCollision(u *server.Object) byte {
	if u.ObjFlags&0x1000000 != 0 {
		u.ObjFlags &^= 0x40
	} else {
		u.ObjFlags |= 0x40
	}
	return byte(u.ObjFlags)
}
func worldSwitch(u *server.Object) byte {
	flags := u.ObjFlags
	if flags&0x1000000 != 0 {
		u.ObjFlags = flags &^ 0x40
		if u.Collide != nil && flags&0x40 != 0 {
			return byte(collisionActivate(u))
		}
	} else {
		if u.Field33 == 0 {
			u.NeedSync()
			u.Field33 = 1
		}
		u.ObjFlags |= 0x40
	}
	return byte(u.ObjFlags)
}
func worldTrapDoor(u *server.Object) uint32 {
	data := u.CollideData
	stamp := equipmentWord(data, 16)
	core := GetServer().S()
	if u.ObjFlags&0x1000000 != 0 {
		result := u.Field5
		if result&2 != 0 {
			u.UnsetXStatus(2)
			u.SetXStatus(8)
		} else if result&4 != 0 {
			result = *stamp
			if core.Frame() >= result {
				u.UnsetXStatus(4)
				u.SetXStatus(8)
			}
		} else {
			*equipmentWord(data, 24) = 0
		}
		return result
	}
	result := *stamp
	if result != 0 && core.Frame() >= result {
		stateOnOff(u, true)
		u.UnsetXStatus(2)
		u.SetXStatus(4)
		*stamp += 5 * uint32(core.TickRate())
		inventorySound(874, u, 0, 0)
	}
	return result
}
func worldPhantom(u *server.Object) {
	ud := u.UpdateData
	owner := *temporaryRefWord(ud, 0)
	x := float64(owner.PosVec.X) - float64(*temporaryFloat(ud, 4))
	yd := float64(owner.PosVec.Y) - float64(*temporaryFloat(ud, 8))
	y := float32(yd)
	if yd*float64(y)+x*x <= 160000 {
		u.NewPos.X = float32(float64(owner.PosVec.X) - x)
		u.NewPos.Y = float32(float64(owner.PosVec.Y) - float64(y))
		u.Direction2 = server.Dir16((uint16(owner.Direction1) + 128) & 255)
	} else {
		GetServer().DelayedDelete(u)
	}
}
