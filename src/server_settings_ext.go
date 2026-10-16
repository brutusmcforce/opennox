package opennox

import (
	"context"
	"strconv"
	"strings"

	"github.com/spf13/viper"

	"github.com/brutusmcforce/opennox/libs/console"

	"github.com/brutusmcforce/opennox/legacy"
)

// Server settings added by OpenNox, in the config and as "set" console commands
// (the console changes them until the game is restarted).

const configBookMaxSpellLevel = "server.spells.book_max_level"

func init() {
	configBoolPtr("server.spells.friendly_fire", "NOX_FRIENDLY_FIRE", false, &legacy.SpellFriendlyFire)
	configBoolPtr("server.afk_observer", "NOX_AFK_OBSERVER", true, &legacy.AFKObserver)
	configBoolPtr("server.quest.allow_any_items", "NOX_QUEST_ALLOW_ANY_ITEMS", false, &legacy.QuestAllowAnyItems)
	viper.SetDefault(configBookMaxSpellLevel, 0)
	setLvl := func() { legacy.BookMaxSpellLevel = min(max(viper.GetInt(configBookMaxSpellLevel), 0), 5) }
	setLvl()
	registerOnConfigRead(setLvl)

	noxCmdSet.Register(&console.Command{
		Token: "friendlyfire",
		Help:  "set friendlyfire [on|off]: spells and ranged attacks of players hurt themselves and their allies too, in every game mode",
		Flags: console.Server,
		Func:  consoleSetBool("friendly fire", &legacy.SpellFriendlyFire),
	})
	noxCmdSet.Register(&console.Command{
		Token: "afkobserver",
		Help:  "set afkobserver [on|off]: idle players go to observer mode automatically after 90 seconds",
		Flags: console.Server,
		Func:  consoleSetBool("observer mode for idle players", &legacy.AFKObserver),
	})
	noxCmdSet.Register(&console.Command{
		Token: "questanyitems",
		Help:  "set questanyitems [on|off]: Quest characters keep saved items that can't be found in quests (by default they are dropped on load)",
		Flags: console.Server,
		Func:  consoleSetBool("any items in quest saves", &legacy.QuestAllowAnyItems),
	})
	noxCmdSet.Register(&console.Command{
		Token: "spelllevel",
		Help:  "set spelllevel [1-5|default]: the highest spell level that spell books give (default: 3 in Quest and Coop, 5 otherwise)",
		Flags: console.Server,
		Func: func(ctx context.Context, c *console.Console, tokens []string) bool {
			switch len(tokens) {
			case 0:
			case 1:
				if strings.EqualFold(tokens[0], "default") {
					legacy.BookMaxSpellLevel = 0
					break
				}
				v, err := strconv.Atoi(tokens[0])
				if err != nil || v < 1 || v > 5 {
					return false
				}
				legacy.BookMaxSpellLevel = v
			default:
				return false
			}
			if v := legacy.BookMaxSpellLevel; v > 0 {
				c.Printf(console.ColorLightYellow, "Spell books give up to level %d", v)
			} else {
				c.Printf(console.ColorLightYellow, "Spell books give up to level 3 in Quest and Coop, 5 otherwise")
			}
			return true
		},
	})
}

// consoleSetBool returns a "set" command that shows or changes a setting.
func consoleSetBool(name string, ptr *bool) console.CommandFunc {
	return func(ctx context.Context, c *console.Console, tokens []string) bool {
		switch len(tokens) {
		case 0:
		case 1:
			switch strings.ToLower(tokens[0]) {
			case "on", "1", "true":
				*ptr = true
			case "off", "0", "false":
				*ptr = false
			default:
				return false
			}
		default:
			return false
		}
		state := "off"
		if *ptr {
			state = "on"
		}
		c.Printf(console.ColorLightYellow, "%s: %s", name, state)
		return true
	}
}
