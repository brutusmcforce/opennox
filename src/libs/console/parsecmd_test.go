package console

import (
	"testing"

	"github.com/shoenig/test/must"
)

func TestDecodeSecretToken(t *testing.T) {
	got := EncodeSecret("racoiaws")
	must.EqOp(t, "0YAKikQs", got)
}

// Commands match case-insensitively: map rule files have "SET SPELL ..." (opennox: map rules were ignored).
func TestCmdForTokenCase(t *testing.T) {
	set := &Command{Token: "set"}
	secret := &Command{Token: EncodeSecret("racoiaws"), Flags: Secret}
	cmds := []*Command{set, secret}
	must.EqOp(t, set, cmdForToken(cmds, "SET"))
	must.EqOp(t, set, cmdForToken(cmds, "Set"))
	must.EqOp(t, secret, cmdForToken(cmds, "racoiaws"))
	must.Nil(t, cmdForToken(cmds, "settings"))
}
