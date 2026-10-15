package opennox

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"runtime/debug"

	"github.com/brutusmcforce/opennox/libs/env"
	"github.com/brutusmcforce/opennox/libs/log"
)

var (
	// logFile isn't buffered: the log must be complete when the game crashes or is killed.
	logFile    io.Closer
	logFileHnd = log.NewTextHandler(nil)
)

func init() {
	logFileHnd.SetLevel(slog.LevelDebug)
	log.AddHandler(logFileHnd)
}

func closeLog() {
	setLogFile(nil)
}

func setLogFile(w io.WriteCloser) {
	if f := logFile; f != nil {
		defer func() { _ = f.Close() }()
	}
	if w == nil {
		logFile = nil
		logFileHnd.SetWriter(nil)
		return
	}
	logFile = w
	logFileHnd.SetWriter(w)
}

func defaultLogDir() string {
	if sdir := env.AppUserDir(); sdir != "" {
		return filepath.Join(sdir, "logs")
	}
	if runtime.GOOS == "darwin" {
		// not next to the binary: inside an app bundle that breaks its signature, and /Applications is read-only
		if home, err := os.UserHomeDir(); err == nil {
			return filepath.Join(home, "Library", "Logs", "OpenNox")
		}
	}
	return filepath.Join(filepath.Dir(os.Args[0]), "logs")
}

func writeLogsToDir(dir string) error {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	name := "opennox.log"
	if isDedicatedServer {
		name = "opennox-server.log"
	}
	f, err := os.Create(filepath.Join(dir, name))
	if err != nil {
		return err
	}
	setLogFile(f)
	// Panics and fatal errors only go to stderr, which players usually don't see (e.g. the Windows build started
	// from the file manager), so the log wouldn't say why the game crashed.
	if err := debug.SetCrashOutput(f, debug.CrashOptions{}); err != nil {
		slog.Warn("cannot write crashes to the log", "err", err)
	}
	return nil
}
