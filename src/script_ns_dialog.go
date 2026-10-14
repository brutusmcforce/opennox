package opennox

import (
	"github.com/brutusmcforce/opennox/libs/strman"
	"github.com/brutusmcforce/opennox/noxscript/ns/v4"
	"github.com/brutusmcforce/opennox/noxscript/ns/v4/audio"

	"github.com/brutusmcforce/opennox/common/sound"
)

func (s noxScriptNS) SetShopkeeperText(obj ns.Obj, text ns.StringID) {
	if obj != nil {
		obj.SetShopText(text)
	}
}

func (s noxScriptNS) SetShopkeeperTextStr(obj ns.Obj, text string) {
	// clients only get string IDs for shop texts
	ScriptLog.Printf("noxscript: SetShopkeeperTextStr is not supported, use SetShopkeeperText (%q)", text)
}

func (s noxScriptNS) TellStory(aud audio.Name, story ns.StringID) {
	caller := s.s.noxScript.Caller()
	trigger := s.s.noxScript.Trigger()
	snd := sound.ByName(string(aud))
	nox_xxx_startShopDialog_548DE0(caller, trigger, snd, strman.ID(story))
}

func (s noxScriptNS) TellStoryStr(audio audio.Name, story string) {
	// clients only get string IDs for stories
	ScriptLog.Printf("noxscript: TellStoryStr is not supported, use TellStory (%q)", story)
}

func (s noxScriptNS) StartDialog(cobj ns.Obj, cother ns.Obj) {
	if cobj == nil || cother == nil {
		return
	}
	obj := s.asObj(cobj)
	other := s.asObj(cother)
	nox_xxx_script_forcedialog_548CD0(other, obj)
}
