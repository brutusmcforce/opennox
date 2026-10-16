package opennox

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/brutusmcforce/opennox/libs/ifs"
)

// Data of the map scripts that is kept in saved games (opennox#618): the map's objects are restored from a save,
// but a script starts again, so it keeps what it needs with ns.SetSaveData and reads it back with ns.SaveData.
// It's stored next to the saved map, as <map>.nsdata (JSON), and read back when that map is loaded, before the
// scripts start. Maps loaded from the game's maps directory have none.

func scriptDataPath(mapPath string) string {
	return strings.TrimSuffix(mapPath, filepath.Ext(mapPath)) + ".nsdata"
}

// saveScriptData writes the scripts' data next to a map that was saved.
func (s *Server) saveScriptData(mapPath string) {
	path := scriptDataPath(mapPath)
	if len(s.scriptData) == 0 {
		_ = ifs.Remove(path) // no stale data from an earlier save
		return
	}
	data, err := json.Marshal(s.scriptData)
	if err != nil {
		ScriptLog.Printf("cannot save script data: %v", err)
		return
	}
	f, err := ifs.Create(path)
	if err != nil {
		ScriptLog.Printf("cannot save script data: %v", err)
		return
	}
	defer f.Close()
	if _, err := f.Write(data); err != nil {
		ScriptLog.Printf("cannot save script data: %v", err)
	}
}

// loadScriptData reads the scripts' data saved with a map, if any.
func (s *Server) loadScriptData(mapPath string) {
	s.scriptData = nil
	f, err := ifs.Open(scriptDataPath(mapPath))
	if err != nil {
		if !os.IsNotExist(err) {
			ScriptLog.Printf("cannot read script data: %v", err)
		}
		return
	}
	defer f.Close()
	data, err := io.ReadAll(f)
	if err == nil {
		err = json.Unmarshal(data, &s.scriptData)
	}
	if err != nil {
		ScriptLog.Printf("cannot read script data: %v", err)
		s.scriptData = nil
	}
}

func (s noxScriptNS) SetSaveData(key string, data []byte) {
	if data == nil {
		delete(s.s.scriptData, key)
		return
	}
	if s.s.scriptData == nil {
		s.s.scriptData = make(map[string][]byte)
	}
	s.s.scriptData[key] = append([]byte(nil), data...)
}

func (s noxScriptNS) SaveData(key string) []byte {
	if v, ok := s.s.scriptData[key]; ok {
		return append([]byte(nil), v...)
	}
	return nil
}
