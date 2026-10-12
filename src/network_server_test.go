package opennox

import (
	"testing"

	"github.com/shoenig/test/must"

	"github.com/brutusmcforce/opennox/libs/noxnet/netmsg"
)

func TestSdecodeLen(t *testing.T) {
	must.EqOp(t, -1, sdecodeLen(nil))
	must.EqOp(t, -1, sdecodeLen([]byte{byte(netmsg.MSG_TEXT_MESSAGE)}))
	must.EqOp(t, 10, sdecodeLen([]byte{byte(netmsg.MSG_NEW_ALIAS)}))
	must.EqOp(t, 22, sdecodeLen([]byte{byte(netmsg.MSG_TRY_SPELL)}))
	must.EqOp(t, 52, sdecodeLen([]byte{byte(netmsg.MSG_VOTE), 0}))
	must.EqOp(t, 2, sdecodeLen([]byte{byte(netmsg.MSG_VOTE), 4}))
	must.EqOp(t, 5, sdecodeLen([]byte{byte(netmsg.MSG_TRADE), 0x17}))
	must.EqOp(t, 2, sdecodeLen([]byte{byte(netmsg.MSG_TRADE), 0x11}))
}
