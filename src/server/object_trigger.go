package server

import (
	"unsafe"

	"github.com/brutusmcforce/opennox/common/memmap"
)

// TriggerUpdateData is the update data of triggers, pressure plates and switches (TriggerUpdate, ToggleUpdate).
type TriggerUpdateData struct {
	Flags            uint32         // 0, 0: 1 pressed this frame, 2 trigger once, 4 pressed, 8 updated
	Activator        uint32         // 1, 4: object that pressed the trigger (as a memmap word)
	State            byte           // 2, 8
	Field9           byte           // 2, 9
	Field10          [2]byte        // 2, 10
	ScriptCollide    ScriptCallback // 3, 12: script event 0; filters objects that can press it
	ScriptActivate   ScriptCallback // 5, 20: script event 1
	ScriptDeactivate ScriptCallback // 7, 28: script event 2
	SoundOn          uint32         // 9, 36
	SoundOff         uint32         // 10, 40
	AllowClass       uint32         // 11, 44: if set, only objects of these classes can press it
	IgnoreClass      uint32         // 12, 48: objects of these classes cannot press it
	AllowTeam        int8           // 13, 52: if set, only objects of this team can press it
	IgnoreTeam       int8           // 13, 53: objects of this team cannot press it
	Colors           [6]byte        // 13, 54: two RGB edge colors of a pressure plate
}

var _ = [1]struct{}{}[60-unsafe.Sizeof(TriggerUpdateData{})]

// ActivatorObj returns the object stored in Activator.
func (d *TriggerUpdateData) ActivatorObj() *Object {
	return (*Object)(memmap.WordPtr(d.Activator))
}

// ScriptAt returns the script callback at a 386 offset of the data (12: collide, 20: activate, 28: deactivate).
func (d *TriggerUpdateData) ScriptAt(off int) *ScriptCallback {
	switch off {
	case 12:
		return &d.ScriptCollide
	case 20:
		return &d.ScriptActivate
	case 28:
		return &d.ScriptDeactivate
	}
	panic("invalid trigger script offset")
}

// UpdateDataTrigger returns the update data of triggers and switches (TriggerUpdate, ToggleUpdate).
func (obj *Object) UpdateDataTrigger() *TriggerUpdateData {
	// TODO: verify this conversion by checking ObjectType
	return updateDataAs[TriggerUpdateData](obj)
}

// UpdateDataTrigger returns the update data of trigger and switch types (TriggerUpdate, ToggleUpdate).
func (t *ObjectType) UpdateDataTrigger() *TriggerUpdateData {
	// TODO: verify this conversion by checking the update function
	return (*TriggerUpdateData)(t.UpdateData)
}
