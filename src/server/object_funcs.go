package server

import (
	"reflect"
	"runtime"
	"unsafe"

	"github.com/brutusmcforce/opennox/libs/types"

	"github.com/brutusmcforce/opennox/common/sound"
	"github.com/brutusmcforce/opennox/legacy/common/ccall"
)

var (
	createFuncs    = make(map[string]unsafe.Pointer)
	initFuncs      = make(map[string]objectDefFunc)
	initParseFuncs = map[string]ObjectParseFunc{
		"SkullInit":     objectDirectionInitParse,
		"DirectionInit": objectDirectionInitParse,
		"BreakInit": func(objt *ObjectType, _ []string) error {
			objt.Field9 = 2
			return nil
		},
	}
	updateFuncs       = make(map[string]objectDefFunc)
	updateParseFuncs  = make(map[string]ObjectParseFunc)
	collideFuncs      = make(map[string]objectDefFunc)
	collideParseFuncs = make(map[string]ObjectParseFunc)
	useFuncs          = make(map[string]objectDefFunc)
	useParseFuncs     = make(map[string]ObjectParseFunc)
	damageFuncs       = make(map[string]unsafe.Pointer)
	damageSoundFuncs  = make(map[string]unsafe.Pointer)
	deathFuncs        = make(map[string]objectDefFunc)
	deathParseFuncs   = make(map[string]ObjectParseFunc)
	dropFuncs         = make(map[string]unsafe.Pointer)
	dropParseFuncs    = map[string]ObjectParseFunc{
		"AudEventDrop": func(t *ObjectType, args []string) error {
			if len(args) != 0 {
				if snd := sound.ByName(args[0]); snd != 0 {
					t.s.dropSoundTable[t.ind] = snd
				}
			}
			return nil
		},
	}
	pickupFuncs      = make(map[string]unsafe.Pointer)
	pickupParseFuncs = map[string]ObjectParseFunc{
		"AudEventPickup": func(t *ObjectType, args []string) error {
			if len(args) != 0 {
				if snd := sound.ByName(args[0]); snd != 0 {
					t.s.pickupSoundTable[t.ind] = snd
				}
			}
			return nil
		},
	}
	xferFuncs = make(map[string]unsafe.Pointer)
)

var (
	DefaultDamage      unsafe.Pointer
	DefaultDamageSound unsafe.Pointer
	DefaultXfer        unsafe.Pointer
)

func init() {
	RegisterObjectCreate("NoCreate", nil)
	RegisterObjectCreate("PlayerCreate", nil)

	RegisterObjectInit("NoInit", nil, NoData)
	RegisterObjectInit("ProjectileInit", nil, NoData)
	RegisterObjectInit("GruntInit", nil, NoData)
	RegisterObjectInit("SkeletonInit", nil, NoData)
	RegisterObjectInit("TowerInit", nil, NoData)
	RegisterObjectInit("RewardMarkerInit", nil, DataOf[RewardMarkerInitData]())
	RegisterObjectInit("AnkhInit", nil, DataOf[AnkhInitData]())

	RegisterObjectUpdate("NoUpdate", nil, NoData)
	RegisterObjectUpdate("HomingProjectileUpdate", nil, DataOf[HomingProjectileUpdateData]())
	RegisterObjectUpdate("SpikeBlockUpdate", nil, RawData(2200))
	RegisterObjectUpdate("TowerUpdate", nil, DataOf[TowerUpdateData]())
	RegisterObjectUpdate("WeaponArmorUpdate", nil, DataOf[WeaponArmorUpdateData]())
	RegisterObjectUpdate("DamageRoundoffUpdate", nil, DataOf[WeaponArmorUpdateData]()) // used in demo instead of WeaponArmorUpdate

	RegisterObjectCollide("NoCollide", nil, NoData)

	RegisterObjectUseC("AmmoUse", nil, DataOf[AmmoUseData]())
	RegisterObjectUseC("BowUse", nil, DataOf[BowUseData]())
}

type ObjectParseFunc func(objt *ObjectType, args []string) error

type objectDefFunc struct {
	Func     unsafe.Pointer
	DataSize uintptr
	DataType reflect.Type // nil for raw data
}

// DataDef describes the data block of an object function.
type DataDef struct {
	Size uintptr
	Type reflect.Type // nil for raw data, which keeps the 32-bit layout of the original game on every platform
}

// NoData is used for object functions without data.
var NoData DataDef

// DataOf describes a data block that holds a T.
func DataOf[T any]() DataDef {
	return DataDef{Size: reflect.TypeFor[T]().Size(), Type: reflect.TypeFor[T]()}
}

// RawData describes a data block of size bytes accessed by offset.
func RawData(size uintptr) DataDef {
	return DataDef{Size: size}
}

type ObjectCreateFunc func(obj *Object)

var objectCreateFuncs = ccall.NewFuncs(func(cfnc unsafe.Pointer) ObjectCreateFunc {
	return func(obj *Object) { panic("unregistered object create callback") }
})

func RegisterObjectCreateGo(name string, cfnc unsafe.Pointer, fnc ObjectCreateFunc) {
	RegisterObjectCreate(name, cfnc)
	objectCreateFuncs.Register(cfnc, fnc)
}

func RegisterObjectCreate(name string, fnc unsafe.Pointer) {
	if _, ok := createFuncs[name]; ok {
		panic("already registered")
	}
	createFuncs[name] = fnc
}

type ObjectInitFunc func(obj *Object)

var objectInitGoFuncs = make(map[unsafe.Pointer]ObjectInitFunc)

// Argument-aware callbacks preserve the second initializer argument.
var objectInitArgGoFuncs = make(map[unsafe.Pointer]func(*Object, unsafe.Pointer))

func RegisterObjectInitGo(name string, cfnc unsafe.Pointer, fnc ObjectInitFunc, sz DataDef) {
	RegisterObjectInit(name, cfnc, sz)
	objectInitGoFuncs[cfnc] = fnc
}

func RegisterObjectInit(name string, fnc unsafe.Pointer, sz DataDef) {
	if _, ok := initFuncs[name]; ok {
		panic("already registered")
	}
	initFuncs[name] = objectDefFunc{Func: fnc, DataSize: sz.Size, DataType: sz.Type}
}

func RegisterObjectUpdate(name string, fnc unsafe.Pointer, sz DataDef) {
	if _, ok := updateFuncs[name]; ok {
		panic("already registered")
	}
	updateFuncs[name] = objectDefFunc{Func: fnc, DataSize: sz.Size, DataType: sz.Type}
}

type UpdateFunc func(obj *Object)

// HasUpdate checks if the object uses the update function registered with a given name.
func (obj *Object) HasUpdate(name string) bool {
	def, ok := updateFuncs[name]
	return ok && obj.Update == def.Func
}

var objUpdate = ccall.NewFuncs(func(cfnc unsafe.Pointer) UpdateFunc {
	return func(obj *Object) { panic("unregistered object update callback") }
})

func RegisterObjectUpdateGo(name string, cfnc unsafe.Pointer, fnc UpdateFunc, sz DataDef) {
	RegisterObjectUpdate(name, cfnc, sz)
	objUpdate.Register(cfnc, fnc)
}

// RegisterObjectUpdateCallbackGo registers an update callback under a native Go identity.
func RegisterObjectUpdateCallbackGo(key unsafe.Pointer, fnc UpdateFunc) {
	objUpdate.Register(key, fnc)
}

func RegisterObjectUpdateParse(name string, fnc ObjectParseFunc) {
	if _, ok := updateParseFuncs[name]; ok {
		panic("already registered")
	}
	updateParseFuncs[name] = fnc
}

type CollideFunc func(obj *Object, a2, a3 uintptr)

// CollideResultFunc preserves the declared result bits of a native collision
// owner. Owners whose original callback was void return zero.
type CollideResultFunc func(obj *Object, a2, a3 uintptr) uint32

type objectCollideGoEntry struct {
	call   CollideFunc
	result CollideResultFunc
}

var objectCollideGoFuncs = make(map[unsafe.Pointer]objectCollideGoEntry)

func RegisterObjectCollideGo(name string, cfnc unsafe.Pointer, fnc CollideFunc, sz DataDef) {
	RegisterObjectCollide(name, cfnc, sz)
	objectCollideGoFuncs[cfnc] = objectCollideGoEntry{call: fnc}
}

// RegisterObjectCollideNative binds a stable, non-executable identity to its Go
// owner. The caller owns the identity's lifetime, including all stored copies.
func RegisterObjectCollideNative(name string, key unsafe.Pointer, fnc CollideResultFunc, sz DataDef) {
	if key == nil || fnc == nil {
		panic("nil native collision identity or function")
	}
	RegisterObjectCollide(name, key, sz)
	objectCollideGoFuncs[key] = objectCollideGoEntry{
		call:   func(obj *Object, a2, a3 uintptr) { fnc(obj, a2, a3) },
		result: fnc,
	}
}

func RegisterObjectCollide(name string, fnc unsafe.Pointer, sz DataDef) {
	if _, ok := collideFuncs[name]; ok {
		panic("already registered")
	}
	collideFuncs[name] = objectDefFunc{Func: fnc, DataSize: sz.Size, DataType: sz.Type}
}

func RegisterObjectCollideParse(name string, fnc ObjectParseFunc) {
	if _, ok := collideParseFuncs[name]; ok {
		panic("already registered")
	}
	collideParseFuncs[name] = fnc
}

type UseFuncPtr struct {
	Ptr unsafe.Pointer
}

func (p UseFuncPtr) Get() UseFunc {
	if p.Ptr == nil {
		return nil
	}
	return objUse.Get(p.Ptr)
}

type UseFunc func(obj, obj2 *Object) bool

// UseResultFunc preserves the full signed C-int result for result-aware callers.
type UseResultFunc func(obj, obj2 *Object) int32

var objectUseNativeFuncs = make(map[unsafe.Pointer]UseResultFunc)

// CallResult requires a registered callback and preserves its signed result.
func (p UseFuncPtr) CallResult(obj, obj2 *Object) int32 {
	var result int32
	if fn := objectUseNativeFuncs[p.Ptr]; fn != nil {
		result = fn(obj, obj2)
	} else {
		panic("unregistered object use callback")
	}
	runtime.KeepAlive(obj)
	runtime.KeepAlive(obj2)
	return result
}

// CallDiscard requires a registered callback and ignores its result.
func (p UseFuncPtr) CallDiscard(obj, obj2 *Object) {
	if fn := objectUseNativeFuncs[p.Ptr]; fn != nil {
		fn(obj, obj2)
	} else {
		panic("unregistered object use callback")
	}
	runtime.KeepAlive(obj)
	runtime.KeepAlive(obj2)
}

// RegisterObjectUseNative binds a stable identity to both exact-result and
// nonzero/boolean dispatch. The caller owns the identity and its stored copies.
func RegisterObjectUseNative(name string, key unsafe.Pointer, fn UseResultFunc, sz DataDef) {
	if key == nil || fn == nil {
		panic("nil native use identity or function")
	}
	RegisterObjectUse(name, key, func(obj, obj2 *Object) bool { return fn(obj, obj2) != 0 }, sz)
	objectUseNativeFuncs[key] = fn
}

var objUse = ccall.NewFuncs(func(cfnc unsafe.Pointer) UseFunc {
	return func(obj, obj2 *Object) bool {
		panic("unregistered object use callback")
	}
})

func RegisterObjectUseC(name string, cfnc unsafe.Pointer, sz DataDef) {
	if _, ok := useFuncs[name]; ok {
		panic("already registered")
	}
	useFuncs[name] = objectDefFunc{Func: cfnc, DataSize: sz.Size, DataType: sz.Type}
}

func RegisterObjectUse(name string, cfnc unsafe.Pointer, fnc UseFunc, sz DataDef) {
	if _, ok := useFuncs[name]; ok {
		panic("already registered")
	}
	useFuncs[name] = objectDefFunc{Func: cfnc, DataSize: sz.Size, DataType: sz.Type}
	objUse.Register(cfnc, fnc)
}

func RegisterObjectUseParse(name string, fnc ObjectParseFunc) {
	if _, ok := useParseFuncs[name]; ok {
		panic("already registered")
	}
	useParseFuncs[name] = fnc
}

func RegisterObjectDamage(name string, fnc unsafe.Pointer) {
	if _, ok := damageFuncs[name]; ok {
		panic("already registered")
	}
	damageFuncs[name] = fnc
}

type DamageFunc func(obj, source, weapon *Object, amount, kind int32) bool

type DamageValueFunc func(obj, source, weapon *Object, amount, kind int32) int32

var objDamage = ccall.NewFuncs(func(cfnc unsafe.Pointer) DamageFunc {
	return func(obj, source, weapon *Object, amount, kind int32) bool {
		panic("unregistered object damage callback")
	}
})

var objDamageValue = ccall.NewFuncs(func(cfnc unsafe.Pointer) DamageValueFunc {
	return func(obj, source, weapon *Object, amount, kind int32) int32 {
		panic("unregistered object damage callback")
	}
})

func RegisterObjectDamageGo(name string, cfnc unsafe.Pointer, fnc DamageFunc) {
	RegisterObjectDamage(name, cfnc)
	objDamage.Register(cfnc, fnc)
}

func RegisterObjectDamageValueGo(name string, cfnc unsafe.Pointer, fnc DamageValueFunc) {
	RegisterObjectDamage(name, cfnc)
	objDamageValue.Register(cfnc, fnc)
	objDamage.Register(cfnc, func(obj, source, weapon *Object, amount, kind int32) bool {
		return fnc(obj, source, weapon, amount, kind) != 0
	})
}

type DamageSoundFunc func(obj, other *Object)

var objectDamageSoundGoFuncs = make(map[unsafe.Pointer]DamageSoundFunc)

func RegisterObjectDamageSoundGo(name string, cfnc unsafe.Pointer, fnc DamageSoundFunc) {
	RegisterObjectDamageSound(name, cfnc)
	objectDamageSoundGoFuncs[cfnc] = fnc
}

func RegisterObjectDamageSound(name string, fnc unsafe.Pointer) {
	if _, ok := damageSoundFuncs[name]; ok {
		panic("already registered")
	}
	damageSoundFuncs[name] = fnc
}

func RegisterObjectDeath(name string, fnc unsafe.Pointer, sz DataDef) {
	if _, ok := deathFuncs[name]; ok {
		panic("already registered")
	}
	deathFuncs[name] = objectDefFunc{Func: fnc, DataSize: sz.Size, DataType: sz.Type}
}

type DeathFunc func(obj *Object)

var objDeath = ccall.NewFuncs(func(cfnc unsafe.Pointer) DeathFunc {
	return func(obj *Object) {
		panic("unregistered object death callback")
	}
})

func RegisterObjectDeathGo(name string, cfnc unsafe.Pointer, fnc DeathFunc, sz DataDef) {
	RegisterObjectDeath(name, cfnc, sz)
	objDeath.Register(cfnc, fnc)
}

func RegisterObjectDeathParse(name string, fnc ObjectParseFunc) {
	if _, ok := deathParseFuncs[name]; ok {
		panic("already registered")
	}
	deathParseFuncs[name] = fnc
}

type DropFuncPtr struct {
	Ptr unsafe.Pointer
}

func (p DropFuncPtr) Get() DropFunc {
	if p.Ptr == nil {
		return nil
	}
	return objDrop.Get(p.Ptr)
}

type DropFunc func(obj, obj2 *Object, pos types.Pointf) bool

var objDrop = ccall.NewFuncs(func(cfnc unsafe.Pointer) DropFunc {
	return func(obj, obj2 *Object, pos types.Pointf) bool {
		panic("unregistered object drop callback")
	}
})

func RegisterObjectDropC(name string, cfnc unsafe.Pointer) {
	if _, ok := dropFuncs[name]; ok {
		panic("already registered")
	}
	dropFuncs[name] = cfnc
}

func RegisterObjectDrop(name string, cfnc unsafe.Pointer, fnc DropFunc) {
	if _, ok := dropFuncs[name]; ok {
		panic("already registered")
	}
	dropFuncs[name] = cfnc
	objDrop.Register(cfnc, fnc)
}

type PickupFuncPtr struct {
	Ptr unsafe.Pointer
}

func (p PickupFuncPtr) Get() PickupFunc {
	if p.Ptr == nil {
		return nil
	}
	return objPickup.Get(p.Ptr)
}

type PickupFunc func(who, it *Object, a3, a4 int) bool

var objPickup = ccall.NewFuncs(func(cfnc unsafe.Pointer) PickupFunc {
	return func(who, it *Object, a3, a4 int) bool {
		panic("unregistered object pickup callback")
	}
})

func RegisterObjectPickupC(name string, cfnc unsafe.Pointer) {
	if _, ok := pickupFuncs[name]; ok {
		panic("already registered")
	}
	pickupFuncs[name] = cfnc
}

func RegisterObjectPickup(name string, cfnc unsafe.Pointer, fnc PickupFunc) {
	if _, ok := pickupFuncs[name]; ok {
		panic("already registered")
	}
	pickupFuncs[name] = cfnc
	objPickup.Register(cfnc, fnc)
}

type XferFunc func(obj *Object, arg unsafe.Pointer) int

var objectXferGoFuncs = make(map[unsafe.Pointer]XferFunc)

func RegisterObjectXferGo(name string, cfnc unsafe.Pointer, fnc XferFunc) {
	RegisterObjectXfer(name, cfnc)
	objectXferGoFuncs[cfnc] = fnc
}

func RegisterObjectXfer(name string, fnc unsafe.Pointer) {
	if _, ok := xferFuncs[name]; ok {
		panic("already registered")
	}
	xferFuncs[name] = fnc
}
