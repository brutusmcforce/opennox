package ns

import (
	"math/rand"
)

type Func = any

// MapEvent is an event type for OnMapEvent function.
type MapEvent int

const (
	MapInitialize = MapEvent(iota + 1)
	MapEntry
	MapExit
	MapShutdown
)

// MapEventFunc is a callback type for OnMapEvent.
type MapEventFunc func()

// OnMapEvent registers a callback function that will be called each server map event.
//
// See MapEvent for a list of events and MapEventFunc for callback type.
func OnMapEvent(typ MapEvent, fnc MapEventFunc) {
	if impl == nil {
		return
	}
	impl.OnMapEvent(typ, fnc)
}

// RandomFloat generates random float.
func RandomFloat(min float32, max float32) float32 {
	if impl == nil {
		return min + rand.Float32()*(max-min)
	}
	return impl.RandomFloat(min, max)
}

// Random generates random int.
func Random(min int, max int) int {
	if impl == nil {
		return min + rand.Intn(max-min)
	}
	return impl.Random(min, max)
}

// StopScript exits current script function.
func StopScript(value any) {
	if impl == nil {
		return
	}
	impl.StopScript(value)
}

// AutoSave triggers an autosave. Only solo games.
func AutoSave() {
	if impl == nil {
		return
	}
	impl.AutoSave()
}

// StartupScreen shows startup screen to host.
func StartupScreen(which int) {
	if impl == nil {
		return
	}
	impl.StartupScreen(which)
}

// DeathScreen shows death screen to host.
func DeathScreen(which int) {
	if impl == nil {
		return
	}
	impl.DeathScreen(which)
}

// LoadMapOptions is a set of options for LoadMap.
type LoadMapOptions struct {
	// IgnoreMapType allows bypassing regular map loading rules.
	// Setting this allows loading maps with a different game mode, switching to/from solo maps, etc.
	IgnoreMapType bool
	// HideScores hides scores screen when switching from Quest or other similar map.
	HideScores bool
	// HideTitleScreen hides boot/title screen of the map when loading it.
	HideTitleScreen bool
	// KeepQuestStage loads the next Quest map at the current stage. By default, each map is one stage
	// further, and the monsters get stronger.
	KeepQuestStage bool
	// QuestStage loads the next Quest map at the given stage, if it's above 0.
	QuestStage int
}

// SetSaveData keeps data of the map's scripts in saved games: when the game is saved and loaded again, the map's
// objects are restored, but its scripts start again, so they read it back with SaveData. The data belongs to the
// current map. A nil data removes the key.
func SetSaveData(key string, data []byte) {
	if impl == nil {
		return
	}
	impl.SetSaveData(key, data)
}

// SaveData returns the data kept with SetSaveData in the saved game that was loaded (or set since), or nil.
func SaveData(key string) []byte {
	if impl == nil {
		return nil
	}
	return impl.SaveData(key)
}

// HideQuestScreens hides the screens that Quest shows when a map starts, from the next map on until it's called
// again: the stage screen (the stage, the map's title and picture) and the scores of the previous map.
// For maps linked with exits that don't start a new stage (e.g. an open world), and for dungeons that still show
// them, call it again. LoadMapOptions.HideTitleScreen and HideScores hide them for one map.
func HideQuestScreens(title, scores bool) {
	if impl == nil {
		return
	}
	impl.HideQuestScreens(title, scores)
}

// LoadMap loads a map with a given name and options.
// Name should not contain .map suffix or any other path elements.
func LoadMap(name string, opts *LoadMapOptions) {
	if impl == nil {
		return
	}
	if opts == nil {
		opts = new(LoadMapOptions)
	}
	impl.LoadMap(name, opts)
}
