package script

import "github.com/brutusmcforce/opennox/noxscript/ns/v3"

func OpenNorthEastWall() {
	ns.WallOpen(ns.Wall(110, 98))
}
func OpenNorthWestWall() {
	ns.WallOpen(ns.Wall(99, 109))
}
func OpenSouthWestWall() {
	ns.WallOpen(ns.Wall(128, 140))
}
func OpenSouthEastWall() {
	ns.WallOpen(ns.Wall(141, 127))
}
