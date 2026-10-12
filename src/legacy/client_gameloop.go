//go:build !server

package legacy

import "github.com/brutusmcforce/opennox/client/noxrender"

func Sub_43CEB0() {
	clientFrameAverage()
}
func Sub_40A710(a1 int) uint32 {
	return uint32(serverConfigConnectionRate(int32(a1)))
}
func Nox_client_screenParticlesDraw_431720(vp *noxrender.Viewport) {
	screenParticlesDraw(vp)
}
func Nox_client_newScreenParticle_431540(a1, a2, a3, a4, a5, a6, a7, a8, a9, a10 int) {
	screenParticleCreate(a1, a2, a3, a4, a5, a6, byte(a7), byte(a8), byte(a9), byte(a10))
}
