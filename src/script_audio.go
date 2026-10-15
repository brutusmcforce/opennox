package opennox

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/brutusmcforce/opennox/libs/datapath"
	"github.com/brutusmcforce/opennox/libs/ifs"

	"github.com/brutusmcforce/opennox/server"
)

// Audio files that map scripts play (ns4 PlayAudioFile, opennox#620): the server reads them from the map folder and
// sends each file once to every player over the OpenNox extension channel; clients cache and play them.

// scriptAudioMaxSize limits the size of an audio file a script can send to the players.
const scriptAudioMaxSize = 16 << 20

type scriptAudio struct {
	data map[string][]byte // file contents by map and name
	sent [32]struct {      // files each player got
		netCode uint32 // the player that got them (a new player in the slot has none)
		files   map[string]bool
	}
}

// scriptAudioName checks the name of an audio file in the map folder.
func scriptAudioName(name string) (string, error) {
	name = filepath.ToSlash(filepath.Clean(name))
	if name == "" || !filepath.IsLocal(name) || strings.ContainsAny(name, "\\:") {
		return "", fmt.Errorf("invalid audio file name: %q", name)
	}
	if !strings.EqualFold(filepath.Ext(name), ".wav") {
		return "", fmt.Errorf("unsupported audio file (only .wav): %q", name)
	}
	return name, nil
}

func (s *Server) readScriptAudio(mapName, name string) ([]byte, error) {
	key := mapName + "/" + name
	if data, ok := s.audio.data[key]; ok {
		return data, nil
	}
	f, err := ifs.Open(datapath.Maps(mapName, name))
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()
	data, err := io.ReadAll(io.LimitReader(f, scriptAudioMaxSize+1))
	if err != nil {
		return nil, err
	} else if len(data) > scriptAudioMaxSize {
		return nil, fmt.Errorf("audio file %q is larger than %d bytes", name, scriptAudioMaxSize)
	}
	if s.audio.data == nil {
		s.audio.data = make(map[string][]byte)
	}
	s.audio.data[key] = data
	return data, nil
}

// playScriptAudio makes all players play an audio file from the map folder.
func (s *Server) playScriptAudio(name string) error {
	name, err := scriptAudioName(name)
	if err != nil {
		return err
	}
	mapName := strings.TrimSuffix(strings.ToLower(s.nox_server_currentMapGetFilename_409B30()), ".map")
	data, err := s.readScriptAudio(mapName, name)
	if err != nil {
		return err
	}
	key := mapName + "/" + name
	for _, pl := range s.Players.List() {
		i := int(pl.PlayerIndex())
		if i < 0 || i >= len(s.audio.sent) {
			continue
		}
		sent := &s.audio.sent[i]
		if sent.files == nil || sent.netCode != pl.NetCodeVal {
			sent.netCode = pl.NetCodeVal
			sent.files = make(map[string]bool)
		}
		m := &server.XferPlayAudio{Map: mapName, Name: name, Volume: 100}
		if !sent.files[key] {
			m.Data = data
		}
		if s.NetXferSendExt(pl.PlayerIndex(), m) {
			sent.files[key] = true
		}
	}
	return nil
}

func (s noxScriptNS) PlayAudioFile(name string) {
	if err := s.s.playScriptAudio(name); err != nil {
		s.s.Log.Error("PlayAudioFile", "err", err)
	}
}

func (s noxScriptNS) QuestStage() int {
	return s.s.nox_game_getQuestStage_4E3CC0()
}
