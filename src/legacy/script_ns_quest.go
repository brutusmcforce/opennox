package legacy

import (
	"github.com/brutusmcforce/opennox/noxscript/ns/v4"

	"github.com/brutusmcforce/opennox/server"
)

func Nox_xxx_comJournalEntryAdd_427500(a1 *server.Object, msg ns.StringID, a3 ns.EntryType) {
	journalUnitAdd(a1, string(msg), uint16(a3))
}
