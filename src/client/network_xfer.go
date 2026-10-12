package client

import (
	"github.com/brutusmcforce/opennox/common/ntype"
	"github.com/brutusmcforce/opennox/server"
)

func (c *Client) handleCBORXfer(ind ntype.PlayerInd, m server.NetXferExt) {
	switch m := m.(type) {
	case *server.XferObjectSetLabel:
		if data := c.Objs.Ext.SetByNetCode(m.Object); data != nil {
			data.DisplayName = m.Label
			data.DisplayNameColor = m.Color.Color()
		}
	}
}
