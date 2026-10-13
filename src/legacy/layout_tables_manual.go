package legacy

import (
	"github.com/brutusmcforce/opennox/common/layout"
	"github.com/brutusmcforce/opennox/server"
)

// Offset tables for structs whose fields are addressed by 386 offsets kept in data tables
// (see common/layout; layout_tables.go holds the generated ones).
var (
	layoutServerMonsterDef        = layout.Of[server.MonsterDef]()
	layoutServerPlayerUpdateData  = layout.Of[server.PlayerUpdateData]()
	layoutServerMonsterGenUpdData = layout.Of[server.MonsterGenUpdateData]()
)
