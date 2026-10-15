package opennox

import (
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/spf13/viper"

	"github.com/brutusmcforce/opennox/libs/env"
	"github.com/brutusmcforce/opennox/libs/log"
)

func init() {
	wd, err := os.Getwd()
	if err != nil {
		wd = "."
	}
	viper.SetConfigName(configName)
	// user config directory is checked first, and new config files are created there
	if udir := configUserDir(); udir != "" {
		viper.AddConfigPath(udir)
		configPath = filepath.Join(udir, configName+"."+configExt)
	} else {
		configPath = filepath.Join(wd, configName+"."+configExt)
	}
	viper.AddConfigPath(wd)
	viper.AddConfigPath(filepath.Dir(os.Args[0]))
}

// configUserDir returns the default directory for the config file.
//
// On Linux it's ~/.config/opennox (or $XDG_CONFIG_HOME/opennox), on Windows it's %AppData%\OpenNox.
func configUserDir() string {
	if sdir := env.AppUserDir(); sdir != "" {
		return sdir
	}
	if runtime.GOOS == "windows" {
		dir, err := os.UserConfigDir()
		if err != nil {
			return ""
		}
		return filepath.Join(dir, "OpenNox")
	}
	if dir := os.Getenv("XDG_CONFIG_HOME"); dir != "" {
		return filepath.Join(dir, "opennox")
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	// Linux Snapcraft installation replaces HOME variable
	if rhome := os.Getenv("SNAP_REAL_HOME"); rhome != "" {
		home = rhome
	}
	return filepath.Join(home, ".config", "opennox")
}

const (
	configName        = "opennox"
	configExt         = "yml"
	configNoxDataPath = "game.data"
	configNoxSerial   = "game.serial"
)

var (
	configLog      = log.New("config")
	configPath     string
	configDirty    bool
	configDirtyAt  time.Time
	configReadOnly = isDedicatedServer
	onConfigRead   []func()
)

func registerOnConfigRead(fnc func()) {
	onConfigRead = append(onConfigRead, fnc)
}

func writeConfig() error {
	if configReadOnly || env.IsE2E() {
		return nil
	}
	configLog.Printf("writing to %q", configPath)
	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		configLog.Printf("cannot create config dir: %v", err)
		return err
	}
	if err := viper.WriteConfigAs(configPath); err != nil {
		configLog.Printf("cannot save file: %v", err)
		return err
	}
	return nil
}

func maybeWriteConfig() {
	if !configDirty {
		return
	}
	if writeConfig() == nil {
		configDirty = false
	}
}

func writeConfigLater() {
	configDirty = true
	configDirtyAt = time.Now()
}

// configSaveDelay is how long changed settings wait before they are written.
const configSaveDelay = 2 * time.Second

// flushConfig writes changed settings a moment after the last change (called every frame). Writing them
// only on exit lost them when the game was killed, e.g. a resolution picked on a Steam Deck (opennox#643).
func flushConfig() {
	if configDirty && time.Since(configDirtyAt) >= configSaveDelay {
		maybeWriteConfig()
	}
}

func configStrPtr(key, env string, def string, ptr *string) {
	viper.SetDefault(key, def)
	if env != "" {
		viper.BindEnv(key, env)
	}
	*ptr = viper.GetString(key)
	registerOnConfigRead(func() {
		*ptr = viper.GetString(key)
	})
}

func configBoolPtr(key, env string, def bool, ptr *bool) {
	viper.SetDefault(key, def)
	if env != "" {
		viper.BindEnv(key, env)
	}
	*ptr = viper.GetBool(key)
	registerOnConfigRead(func() {
		*ptr = viper.GetBool(key)
	})
}

func configHiddenBoolPtr(key, env string, ptr *bool) {
	if env != "" {
		viper.BindEnv(key, env)
	}
	*ptr = viper.GetBool(key)
	registerOnConfigRead(func() {
		*ptr = viper.GetBool(key)
	})
}

func readConfig(path string) error {
	defer func() {
		for _, fnc := range onConfigRead {
			fnc()
		}
	}()
	if env.IsE2E() {
		// tests must not depend on the local config
		configLog.Println("e2e: using defaults")
		return nil
	}
	if path != "" {
		if abs, err := filepath.Abs(path); err == nil {
			path = abs
		}
		viper.SetConfigFile(path)
		configPath = path
	}
	err := viper.ReadInConfig()
	if err == nil {
		configPath = viper.ConfigFileUsed()
		configDirty = false
		configLog.Printf("using file: %q", configPath)
		return nil
	}
	if _, ok := err.(viper.ConfigFileNotFoundError); ok || os.IsNotExist(err) {
		writeConfigLater()
		configLog.Println("file not found, using defaults")
		return nil
	}
	// There's a weird behavior in viper that it may try to read config file without extension.
	// It's a problem for us because we have "opennox" binary on Linux, so viper will read the binary itself and fail.
	// So we check this case and ignore it as if the config was not found.
	if filepath.Base(viper.ConfigFileUsed()) == configName {
		writeConfigLater()
		configLog.Println("file not found, using defaults")
		return nil
	}
	return err // error parsing config
}
