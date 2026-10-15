package script

import "github.com/brutusmcforce/opennox/noxscript/ns/v3"

func OpenHammerWall() {
	ns.WallOpen(ns.Wall(115, 151))
}
