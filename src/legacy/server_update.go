package legacy

import "github.com/brutusmcforce/opennox/server"

func Nox_xxx_findObjectAtCursor_54AF40(a1 *server.Object) *server.Object {
	return spatialCursor(a1)
}
func Nox_xxx_updateFallLogic_51B870(a1 *server.Object) {
	motionFall(a1)
}
func Sub_51B810(a1 *server.Object) {
	motionProjectileStep(a1)
}
func Sub_537770(a1 *server.Object) {
	motionProjectileDispatch(a1)
}
