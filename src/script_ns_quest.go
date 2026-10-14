package opennox

import (
	"math"

	"github.com/brutusmcforce/opennox/noxscript/ns/v4"

	"github.com/brutusmcforce/opennox/common/sound"
	"github.com/brutusmcforce/opennox/legacy"
	"github.com/brutusmcforce/opennox/server"
)

func (s noxScriptNS) GetQuestStatus(name string) int {
	return int(int32(legacy.QuestStatus(name)))
}

func (s noxScriptNS) GetQuestStatusFloat(name string) float32 {
	return legacy.QuestStatusFloat(name)
}

func (s noxScriptNS) SetQuestStatus(status int, name string) {
	legacy.QuestStatusSet(name, uint32(int32(status)), false)
}

func (s noxScriptNS) SetQuestStatusFloat(status float32, name string) {
	legacy.QuestStatusSet(name, math.Float32bits(status), true)
}

func (s noxScriptNS) ResetQuestStatus(name string) {
	legacy.QuestStatusReset(name)
}

func (s noxScriptNS) JournalEntry(obj ns.Obj, msg ns.StringID, typ ns.EntryType) {
	if obj == nil {
		for _, it := range s.s.Players.ListUnits() {
			legacy.Nox_xxx_comJournalEntryAdd_427500(it, msg, typ)
		}
	} else {
		legacy.Nox_xxx_comJournalEntryAdd_427500(server.ToObject(obj.(server.Obj)), msg, typ)
		if (typ & 0xB) != 0 {
			s.s.Audio.EventObj(sound.SoundJournalEntryAdd, obj.(server.Obj), 0, 0)
		}
	}
}

func (s noxScriptNS) JournalEdit(obj ns.Obj, message ns.StringID, typ ns.EntryType) {
	s.JournalEditStr(obj, string(message), typ)
}

func (s noxScriptNS) JournalDelete(obj ns.Obj, message ns.StringID) {
	s.JournalDeleteStr(obj, string(message))
}

func (s noxScriptNS) JournalEntryStr(obj ns.Obj, msg string, typ ns.EntryType) {
	s.JournalEntry(obj, ns.StringID(msg), typ)
}

func (s noxScriptNS) JournalEditStr(obj ns.Obj, message string, typ ns.EntryType) {
	legacy.ScriptJournal(s.asObj(obj), message, uint16(typ), true)
}

func (s noxScriptNS) JournalDeleteStr(obj ns.Obj, message string) {
	legacy.ScriptJournal(s.asObj(obj), message, 0, false)
}
