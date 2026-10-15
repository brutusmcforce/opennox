package opennox

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/brutusmcforce/opennox/common/ntype"
	"github.com/brutusmcforce/opennox/legacy"
	"github.com/brutusmcforce/opennox/server"
)

// scriptAudioFiles are the audio files the server sent (see script_audio.go), by map and name.
var scriptAudioFiles = make(map[string]string)

// handleScriptAudio plays an audio file that a map script plays on the server.
func handleScriptAudio(_ ntype.PlayerInd, m server.NetXferExt) {
	am, ok := m.(*server.XferPlayAudio)
	if !ok {
		return
	}
	name, err := scriptAudioName(am.Name)
	mapName := strings.ToLower(am.Map)
	if err != nil || mapName == "" || !filepath.IsLocal(mapName) || strings.ContainsAny(mapName, "/\\:") {
		gameLog.Printf("script audio: invalid file %q %q", am.Map, am.Name)
		return
	}
	key := mapName + "/" + name
	if len(am.Data) != 0 {
		path, err := saveScriptAudio(mapName, name, am.Data)
		if err != nil {
			gameLog.Printf("script audio: cannot save %q: %v", key, err)
			return
		}
		scriptAudioFiles[key] = path
	}
	path := scriptAudioFiles[key]
	if path == "" {
		gameLog.Printf("script audio: no data for %q", key)
		return
	}
	legacy.Dialogs.PlayFile(path, max(0, min(am.Volume, 100)))
}

// saveScriptAudio writes a received audio file to the cache, so the audio player can stream it.
func saveScriptAudio(mapName, name string, data []byte) (string, error) {
	dir, err := os.UserCacheDir()
	if err != nil {
		dir = os.TempDir()
	}
	path := filepath.Join(dir, "opennox", "audio", mapName, filepath.FromSlash(name))
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return "", err
	}
	tmp := path + ".part"
	if err := os.WriteFile(tmp, data, 0644); err != nil {
		return "", err
	}
	if err := os.Rename(tmp, path); err != nil {
		_ = os.Remove(tmp)
		return "", err
	}
	return path, nil
}
