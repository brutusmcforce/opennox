package server

import "unsafe"

// Script name slots of ObjectEditorData.ScriptNames.
const (
	ScriptNamePickup            = 0  // 0: Object.ScriptPickup (event 14)
	ScriptNameHole              = 1  // 128: hole collide script (event 12)
	ScriptNameTriggerActivate   = 2  // 256: trigger ScriptActivate (event 1)
	ScriptNameTriggerDeactivate = 3  // 384: trigger ScriptDeactivate (event 2)
	ScriptNameTriggerCollide    = 4  // 512: trigger ScriptCollide (event 0)
	ScriptNameEnemySighted      = 5  // 640: monster ScriptEnemySighted (event 3)
	ScriptNameLookingForEnemy   = 6  // 768: monster ScriptLookingForEnemy (event 4)
	ScriptNameDeath             = 7  // 896: monster ScriptDeath (event 5)
	ScriptNameChangeFocus       = 8  // 1024: monster ScriptChangeFocus (event 6)
	ScriptNameIsHit             = 9  // 1152: monster ScriptIsHit (event 7)
	ScriptNameRetreat           = 10 // 1280: monster ScriptRetreat (event 8)
	ScriptNameCollision         = 11 // 1408: monster ScriptCollision (event 9)
	ScriptNameEnemyHeard        = 12 // 1536: monster ScriptHearEnemy (event 10)
	ScriptNameEndOfWaypoint     = 13 // 1664: monster ScriptEndOfWaypoint (event 11)
	ScriptNameLostEnemy         = 14 // 1792: monster ScriptLostEnemy (event 13)
	ScriptNameGen52             = 15 // 1920: monster generator FuncInd52 (event 15)
	ScriptNameGen60             = 16 // 2048: monster generator FuncInd60 (event 16)
	ScriptNameGen76             = 17 // 2176: monster generator FuncInd76 (event 18)
	ScriptNameGen68             = 18 // 2304: monster generator FuncInd68 (event 17)
)

// ObjectEditorData is the editor-mode (GameFlag22|GameFlag23) data of an object (Object.Field189):
// script callbacks are kept by name instead of being resolved to function indices.
// Name writes are C strcpy-like and are not bounded by the slot size.
type ObjectEditorData struct {
	ScriptNames [19][128]byte // 0, 0: NUL-terminated callback names, see ScriptName* slots
	Light       [140]byte     // 608, 2432: saved light record of invisible lights (386 drawable layout from LightFlags)
}

var _ = [1]struct{}{}[2572-unsafe.Sizeof(ObjectEditorData{})]

// EditorData returns the editor-mode data of the object, or nil if it has none.
func (obj *Object) EditorData() *ObjectEditorData {
	return (*ObjectEditorData)(obj.Field189)
}

// ScriptName returns a pointer to the start of the given script name slot, or nil if d is nil.
func (d *ObjectEditorData) ScriptName(slot int) unsafe.Pointer {
	if d == nil {
		return nil
	}
	return unsafe.Pointer(&d.ScriptNames[slot][0])
}
