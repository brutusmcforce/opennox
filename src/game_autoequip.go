package opennox

import (
	"context"

	"github.com/spf13/viper"

	"github.com/brutusmcforce/opennox/legacy"
	"github.com/brutusmcforce/opennox/libs/console"
	"github.com/brutusmcforce/opennox/server"
)

// configGameAutoEquip turns off equipping picked up weapons and armor automatically (opennox#741).
// It applies to the local player only: the server doesn't know the choice of other players.
const configGameAutoEquip = "game.auto_equip"

var gameAutoEquip bool

func init() {
	configBoolPtr(configGameAutoEquip, "", true, &gameAutoEquip)
	legacy.AutoEquip = func(u *server.Object) bool {
		if gameAutoEquip {
			return true
		}
		pl := u.UpdateDataPlayer().Player
		return pl == nil || pl.Index() != server.HostPlayerIndex
	}
	noxConsole.Register(&console.Command{
		Token: "autoequip",
		Help:  "autoequip: toggle equipping picked up weapons and armor automatically",
		Flags: console.ClientServer,
		Func: func(ctx context.Context, c *console.Console, tokens []string) bool {
			if len(tokens) != 0 {
				return false
			}
			setGameAutoEquip(!gameAutoEquip)
			if gameAutoEquip {
				c.Print(console.ColorLightYellow, "auto-equip is on")
			} else {
				c.Print(console.ColorLightYellow, "auto-equip is off")
			}
			return true
		},
	})
}

func setGameAutoEquip(on bool) {
	gameAutoEquip = on
	viper.Set(configGameAutoEquip, on)
	writeConfigLater()
}
