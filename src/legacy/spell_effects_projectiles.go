package legacy

import (
	"unsafe"

	"github.com/brutusmcforce/opennox/common/memmap"
	"github.com/brutusmcforce/opennox/libs/types"
	"github.com/brutusmcforce/opennox/server"
)

func spellEffectDirection(dir int16) types.Pointf {
	return types.Ptf(*memmap.PtrFloat32(0x587000, 194136+uintptr(int32(dir)*8)), *memmap.PtrFloat32(0x587000, 194140+uintptr(int32(dir)*8)))
}
func spellEffectPlacementFlags() int {
	if controlFlags(2048) {
		return 9
	}
	return 73
}
func spellEffectBurn(id int32, a, b, c *server.Object, record unsafe.Pointer, level int32) int32 {
	typ := spellEffectGlyphType()
	if record == nil || c == nil {
		return 0
	}
	pos := c.PosVec
	if uint32(c.TypeInd) != typ {
		pos = (*server.SpellAcceptArg)(record).Pos
		if !spellEffectPlacement(c.PosVec, pos) {
			spellEffectInform(c)
			return 0
		}
	}
	if u := spellEffectNew(stateType(2487732, "MediumFlame")); u != nil {
		spellEffectCreate(u, c, pos)
		motionDecaySet(u, floatToInt32(float32(spellEffectScalar("BurnDuration"))))
		visibilityFXSpark(u.PosVec, 64)
	}
	spellEffectPosAudio(id, 0, (*server.SpellAcceptArg)(record).Pos)
	return 1
}
func spellEffectFireball(id int32, a, b, c *server.Object, record unsafe.Pointer, level int32) int32 {
	name := (*int8)(memmap.LoadPtr(0x587000, 258864+uintptr(level)*4))
	u := GetServer().S().NewObjectByTypeID(GoStringP(unsafe.Pointer(name)))
	if u == nil {
		return 1
	}
	dir := spellEffectDirection(int16(c.Direction1))
	radius := float64(c.Shape.Circle.R) * 2
	// C stores the X position before adding inherited velocity; Y stays in x87.
	x := float32(radius*float64(dir.X) + float64(c.PosVec.X))
	y := radius*float64(dir.Y) + float64(c.PosVec.Y)
	pos := types.Ptf(float32(float64(x)+float64(c.VelVec.X)), float32(y+float64(c.VelVec.Y)))
	if !spellEffectTrace(c.PosVec, pos, 5) {
		pos = c.PosVec
	}
	spellEffectCreate(u, c, pos)
	speed := spellEffectTable("FireballSpeedCoeff", level-1) * float64(u.SpeedCur)
	u.SpeedCur = float32(speed)
	vx := float32(speed * float64(dir.X))
	vy := speed * float64(dir.Y)
	u.VelVec.X = float32(float64(vx) + float64(c.VelVec.X))
	u.VelVec.Y = float32(vy + float64(c.VelVec.Y))
	u.Direction1 = c.Direction1
	u.Direction2 = c.Direction1
	spellEffectAudio(id, 0, c)
	return 1
}
func spellEffectOwnedTypes(u *server.Object, start uintptr, count int) bool {
	if u == nil {
		return false
	}
	for it := u.Field129; it != nil; it = it.Field128 {
		for i := 0; i < count; i++ {
			if uint32(it.TypeInd) == *memmap.PtrUint32(0x5d4594, start+uintptr(4*i)) {
				return true
			}
		}
	}
	return false
}
func spellEffectFist(id int32, a, b, c *server.Object, record unsafe.Pointer, level int32) int32 {
	if *memmap.PtrUint32(0x5d4594, 2487736) == 0 {
		for i, name := range []string{"SmallFist", "MediumFist", "LargeFist", "LargeFist", "LargeFist"} {
			*memmap.PtrUint32(0x5d4594, 2487740+uintptr(i*4)) = uint32(GetServer().S().Types.IndByID(name))
		}
		*memmap.PtrUint32(0x5d4594, 2487736) = 1
	}
	if spellEffectOwnedTypes(b, 2487740, 5) {
		resourcePriority(b, "ExecSpel.c:TooManyFists")
		return 0
	}
	pos := (*server.SpellAcceptArg)(record).Pos
	if !spellEffectTrace(c.PosVec, pos, spellEffectPlacementFlags()) {
		spellEffectInform(c)
		return 0
	}
	if u := spellEffectNew(*memmap.PtrUint32(0x5d4594, 2487736+uintptr(level)*4)); u != nil {
		u.UpdateDataDamage().Damage = uint32(floatToInt32(float32(spellEffectTable("FistOfVengeanceDamage", level-1))))
		spellEffectCreate(u, c, pos)
		u.Field5 |= 0x20
		stateRaise(u, 255)
		u.Field27 = float32(-spellEffectScalar("FistSpeed"))
		u.ObjFlags |= 0x800000
		u.Field29 = 1091567616
		spellEffectAudio(id, 0, c)
	}
	return 1
}
func spellEffectCleansingFlame(id int32, a, b, c *server.Object, record unsafe.Pointer, level int32) int32 {
	s := GetServer().S()
	if *memmap.PtrUint32(0x5d4594, 2487760) == 0 {
		for i, name := range []string{"SmallFlameCleanse", "SmallFlameCleanse", "MediumFlameCleanse", "FlameCleanse", "LargeFlameCleanse", "SmallBlueFlameCleanse", "SmallBlueFlameCleanse", "MediumBlueFlameCleanse", "BlueFlameCleanse", "LargeBlueFlameCleanse"} {
			*memmap.PtrUint32(0x5d4594, 2487760+uintptr(i*4)) = uint32(s.Types.IndByID(name))
		}
	}
	if spellEffectOwnedTypes(a, 2487760, 10) {
		resourcePriority(b, "plyrspel.c:TooManySpells")
		return 0
	}
	if !controlFlags(2048) {
		level = 4
	}
	for i := 0; i < 48; i++ {
		n := level - int32(s.Rand.Logic.IntClamp(0, 1))
		if n < 1 {
			continue
		}
		bank := int32(5)
		if id == 10 {
			bank = 0
		}
		u := spellEffectNew(*memmap.PtrUint32(0x5d4594, 2487756+uintptr(bank+n)*4))
		if u == nil {
			continue
		}
		u.Direction1 = server.Dir16(s.Rand.Logic.IntClamp(0, 255))
		dir := spellEffectDirection(int16(u.Direction1))
		radius := float64(u.Shape.Circle.R) + float64(c.Shape.Circle.R) + 4
		pos := types.Ptf(float32(radius*float64(dir.X)+float64(c.PosVec.X)), float32(radius*float64(dir.Y)+float64(c.PosVec.Y)))
		if spellEffectTrace(c.PosVec, pos, 65) {
			spellEffectCreate(u, c, pos)
			u.Direction2 = u.Direction1
			u.VelVec = types.Ptf(dir.X*4, dir.Y*4)
			u.Field34 = s.Frame() + uint32(s.Rand.Logic.IntClamp(int(3*s.TickRate()), int(6*s.TickRate())))
			u.Pos39 = c.PosVec
			u.Update = flameCleanseCallbackKey()
			s.Objs.AddToUpdatable(u)
			u.ObjClass |= 0x40000000
			u.Float28 = 0
			visibilityFXPrediction(u)
		} else {
			GetServer().DelayedDelete(u)
		}
	}
	spellEffectAudio(id, 0, c)
	return 1
}
func spellEffectMeteorShower(id int32, a, b, c *server.Object, record unsafe.Pointer, level int32) int32 {
	typ := stateType(2487800, "MeteorShower")
	pos := (*server.SpellAcceptArg)(record).Pos
	if !spellEffectTrace(c.PosVec, pos, spellEffectPlacementFlags()) {
		spellEffectInform(c)
		return 0
	}
	if u := spellEffectNew(typ); u != nil {
		u.UpdateDataDamage().Damage = uint32(floatToInt32(float32(spellEffectTable("MeteorDamage", level-1))))
		spellEffectCreate(u, b, pos)
		spellEffectAudio(id, 0, b)
	}
	return 1
}
func spellEffectMeteor(id int32, a, b, c *server.Object, record unsafe.Pointer, level int32) int32 {
	if dword_5d4594_2487804 == 0 {
		dword_5d4594_2487804 = uint32(GetServer().S().Types.IndByID("Meteor"))
	}
	typ := uint32(dword_5d4594_2487804)
	for it := b.Field129; it != nil; it = it.Field128 {
		if uint32(it.TypeInd) == typ {
			resourcePriority(b, "ExecSpel.c:TooManyMeteors")
			return 0
		}
	}
	pos := (*server.SpellAcceptArg)(record).Pos
	if !spellEffectTrace(c.PosVec, pos, spellEffectPlacementFlags()) {
		spellEffectInform(c)
		return 0
	}
	if u := spellEffectNew(typ); u != nil {
		u.UpdateDataDamage().Damage = uint32(floatToInt32(float32(spellEffectTable("MeteorDamage", level-1))))
		spellEffectCreate(u, b, pos)
		u.Field5 |= 0x20
		stateRaise(u, 255)
		u.Field27 = float32(-spellEffectScalar("MeteorSpeed"))
		spellEffectPosAudio(id, 0, pos)
	}
	return 1
}
func spellEffectToxicCloud(id int32, a, b, c *server.Object, record unsafe.Pointer, level int32) int32 {
	typ := stateType(2487808, "ToxicCloud")
	pos := (*server.SpellAcceptArg)(record).Pos
	if !spellEffectPlacement(c.PosVec, pos) {
		spellEffectInform(c)
		return 0
	}
	u := spellEffectNew(typ)
	if u != nil {
		spellEffectCreate(u, b, pos)
		u.UpdateDataToxicCloud().Life = uint32(floatToInt32(float32(spellEffectScalar("ToxicCloudLifetime") * float64(int32(GetServer().S().TickRate())))))
	}
	spellEffectAudio(id, 0, u)
	return 1
}
func spellEffectArachna(id int32, a, b, c *server.Object, record unsafe.Pointer, level int32) int32 {
	typ := stateType(2487812, "ArachnaphobiaFocus")
	pos := (*server.SpellAcceptArg)(record).Pos
	if !spellEffectPlacement(c.PosVec, pos) {
		spellEffectInform(c)
		return 0
	}
	if u := spellEffectNew(typ); u != nil {
		spellEffectCreate(u, b, pos)
	}
	return 1
}
