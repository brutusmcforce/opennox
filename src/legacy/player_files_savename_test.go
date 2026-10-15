package legacy

import (
	"path/filepath"
	"testing"

	"github.com/brutusmcforce/opennox/internal/cryptfile"
	"github.com/brutusmcforce/opennox/server"
)

// The map name of a save survives a write and a read, also when it's longer than the original 8 characters,
// and doesn't overwrite what follows in the player.
func TestPlayerFileSaveName(t *testing.T) {
	path := filepath.Join(t.TempDir(), "save.plr")
	for _, name := range []string{"g_castle", "MyLongQuestCastle", "MyVeryLongQuestMapName123456"} {
		if err := cryptfile.OpenGlobal(path, cryptfile.WriteOnly, 27); err != nil {
			t.Fatal(err)
		}
		var w server.Player
		playerFileSaveName(playerFileStream(), &w, name)
		if got := w.SaveName(); got != name {
			t.Errorf("write: got %q", got)
		}
		_ = cryptfile.Close()

		if err := cryptfile.OpenGlobal(path, cryptfile.ReadOnly, 27); err != nil {
			t.Fatal(err)
		}
		var p server.Player
		p.Ankhs[0] = &server.Object{}
		p.Field4792 = 1
		ankh := p.Ankhs[0]
		playerFileSaveName(playerFileStream(), &p, "")
		_ = cryptfile.Close()
		if got := p.SaveName(); got != name {
			t.Errorf("read: got %q, want %q", got, name)
		}
		if p.Ankhs[0] != ankh || p.Field4792 != 1 {
			t.Errorf("%q: the name overwrote the fields after it", name)
		}
	}
}
