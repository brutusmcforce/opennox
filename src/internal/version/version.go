package version

import (
	"github.com/brutusmcforce/opennox/libs/env"
	"github.com/brutusmcforce/opennox/libs/log"
)

var Log = log.New("version")

func LogVersion() {
	Log.Info("engine", "version", Version(), "commit", Commit())
}

const (
	// DefVersion is the version of builds that are not on a tag. Bump it right after a release, or every
	// local and latest-dev build claims to be older than what was already released.
	DefVersion = "v1.10.0-dev"
)

const (
	devCommit = "<dev>"
)

var (
	version = DefVersion
	commit  = devCommit
)

func Version() string {
	return version
}

func Commit() string {
	return commit
}

func ClientVersion() string {
	if env.IsE2E() {
		// The version is drawn on screen, and e2e screenshots must not depend on the commit.
		return "e2e"
	}
	if IsDev() {
		return version + " (" + commit + ")"
	}
	return version
}

func IsDev() bool {
	return commit == devCommit || semverIsDev(version)
}
