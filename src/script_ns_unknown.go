package opennox

import (
	"github.com/brutusmcforce/opennox/libs/strman"

	"github.com/brutusmcforce/opennox/legacy"
)

func (s noxScriptNS) Unused5e(id string) int {
	str := s.s.Strings().GetStringInFile(strman.ID(id), "CScrFunc.c")
	return legacy.Sub_512E80(str)
}

func (s noxScriptNS) Unused74(arg1 int, arg2 int) {}

func (s noxScriptNS) Unknownb8(id int) bool {
	obj := s.s.noxScript.ScriptToObject(id)
	return obj != nil && uint32(obj.ObjSubClass)&0x100 != 0
}

func (s noxScriptNS) Unknownb9(id int) bool {
	obj := s.s.noxScript.ScriptToObject(id)
	return obj != nil && uint32(obj.ObjSubClass)&0x80 != 0
}

func (s noxScriptNS) Unknownc4() {
	noxClient.Nox_video_stopAllFades44E040()
}
