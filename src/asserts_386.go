//go:build 386

package opennox

// Layout assertions: these structs keep the layout of the original 32-bit binary.
// On other architectures pointers are larger and the layout differs.

import (
	"github.com/brutusmcforce/opennox/client"
	"github.com/brutusmcforce/opennox/client/noxrender"
	"github.com/brutusmcforce/opennox/server"
	"unsafe"
)

var _ = [1]struct{}{}[288-unsafe.Sizeof(audioStructXxx{})]
var _ = [1]struct{}{}[36-unsafe.Sizeof(audioStructYyy{})]
var _ = [1]struct{}{}[1056-unsafe.Sizeof(noxrender.RenderData{})]
var _ = [1]struct{}{}[48-unsafe.Sizeof(noxrender.RenderMat{})]
var _ = [1]struct{}{}[52-unsafe.Sizeof(noxrender.Viewport{})]
var _ = [1]struct{}{}[4-unsafe.Sizeof(int(0))]
var _ = [1]struct{}{}[80-unsafe.Sizeof(client.DrawableFX{})]
var _ = [1]struct{}{}[20-unsafe.Sizeof(server.Settings3{})]
var _ = [1]struct{}{}[60-unsafe.Sizeof(server.Settings2{})]
var _ = [1]struct{}{}[168-unsafe.Sizeof(server.Settings{})]
var _ = [1]struct{}{}[96-unsafe.Sizeof(server.MapGroup{})]
var _ = [1]struct{}{}[unsafe.Sizeof(int(0))-4]
var _ = [1]struct{}{}[780-unsafe.Sizeof(Object{})]
var _ = [1]struct{}{}[20-unsafe.Sizeof(server.HealthData{})]
var _ = [1]struct{}{}[24-unsafe.Sizeof(harpoonPtr{})]
var _ = [1]struct{}{}[16-unsafe.Sizeof(server.MagicWallSave{})]
var _ = [1]struct{}{}[40-unsafe.Sizeof(server.PhonemeLeaf{})]
var _ = [1]struct{}{}[24-unsafe.Sizeof(MonsterListen{})]
var _ = [1]struct{}{}[248-unsafe.Sizeof(server.MonsterDef{})]
var _ = [1]struct{}{}[36-unsafe.Sizeof(Wall{})]
