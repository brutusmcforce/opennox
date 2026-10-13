package legacy

import (
	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/common/unit/ai"
	"github.com/brutusmcforce/opennox/server"
)

func monsterControlAnimation(u *server.Object) {
	ud := u.UpdateDataMonster()
	if u.ObjSubClass&0x10 != 0 {
		action := monsterControlHead(u)
		if action == 16 || action == 17 {
			ud.Field120_3 = 0
			return
		}
	}
	if ud.Field120_3 != 0 {
		return
	}
	p := monsterNPCAnim(u)
	if p == nil {
		return
	}
	frames := p.Frames
	delay := p.Delay
	ud.Field120_0 = frames
	if frames == 0 {
		ud.Field120_3 = 1
		return
	}
	ud.Field120_2++
	if int(ud.Field120_2) < int(delay)+1 {
		return
	}
	ud.Field120_2 = 0
	ud.Field120_1++
	if ud.Field120_1 >= frames {
		if p.Loop != 0 {
			ud.Field120_1 = 0
		} else {
			ud.Field120_1 = frames - 1
			ud.Field120_3 = 1
		}
	}
}
func monsterControlRefresh(u *server.Object) int32 {
	ud := u.UpdateDataMonster()
	if target := ud.Field304; target != nil && target.ObjFlags&0x8020 != 0 {
		ud.Field304 = nil
	}
	if ud.AIStackInd < 0 {
		return int32(ud.AIStackInd)
	}
	for i := int(ud.AIStackInd); i >= 0; i-- {
		a := &ud.AIStack[i]
		count := memmap.Uint32(0x587000, 230388+16*uintptr(a.Action))
		for j := uint32(0); j < count; j++ {
			if memmap.Uint32(0x587000, 230392+4*uintptr(j+4*a.Action)) == 1 && a.Args[j*2] != 0 {
				if t := a.ArgObj(int(j * 2)); t.ObjFlags&0x20 != 0 {
					a.Args[j*2] = 0
				}
			}
		}
		var target *server.Object
		switch a.Type() {
		case ai.ACTION_ESCORT:
			target = a.ArgObj(2)
		case ai.ACTION_MOVE_TO, ai.ACTION_FAR_MOVE_TO:
			t := a.ArgObj(2)
			if t != nil {
				if GetServer().S().CanInteract(u, t, 0) || ud.HasAction(ai.ACTION_ESCORT) {
					target = t
				} else {
					a.Args[2] = 0
				}
			}
		case ai.ACTION_FIGHT:
			target = ud.CurrentEnemy
		case ai.ACTION_MISSILE_ATTACK:
			t := a.ArgObj(2)
			if t != nil && GetServer().S().CanInteract(u, t, 0) {
				target = t
			}
		}
		if target != nil {
			a.SetArgs(target.PosVec)
		}
	}
	return 0
}
func monsterControlRevive(u *server.Object) uint32 {
	if u == nil {
		return 0
	}
	if u.ObjClass&2 != 0 && u.ObjFlags&0x8000 != 0 {
		u.UpdateDataMonster().StatusFlags &^= 0x100000
		return lifecycleRaiseZombie(u)
	}
	return motionAddress(u)
}
