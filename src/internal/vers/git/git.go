package git

import (
	"bytes"
	"os"
	"os/exec"
	"strings"

	"github.com/brutusmcforce/opennox/internal/version"
)

const defVersion = version.DefVersion

func do(cmd ...string) (string, error) {
	c := exec.Command(cmd[0], cmd[1:]...)
	var buf bytes.Buffer
	c.Stdout = &buf
	c.Stderr = os.Stderr
	c.Env = os.Environ()
	err := c.Run()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(buf.String()), nil
}

func FullVersion() string {
	vers, sha := Version()
	if sha == "" {
		return vers
	}
	return vers + "-" + sha
}

func Version() (string, string) {
	sha := ShortSHA()
	if sha == "" {
		return defVersion, ""
	}
	var vers string
	if v := os.Getenv("GIT_TAG"); v != "" {
		vers = v
	} else {
		vers = versionTag(sha)
	}
	if vers == "" || vers == "undefined" {
		return defVersion, sha
	}
	return vers, ""
}

// versionTag returns the version tag that a commit has, if any.
//
// A commit can have more than one tag: the development build moves the latest-dev tag onto the commit it was
// built from, and git name-rev (used here before) returns that one for a release commit, so the v1.10.0-alpha3
// binaries reported their version as "latest-dev". Only a tag of this commit is used now: name-rev also named
// commits BEFORE a tag after it (v1.10.0-alpha3~2).
func versionTag(sha string) string {
	out, err := do("git", "tag", "--points-at", sha)
	if err != nil {
		return ""
	}
	for _, tag := range strings.Fields(out) {
		if isVersionTag(tag) {
			return tag
		}
	}
	return ""
}

// isVersionTag reports whether a tag names a version, like v1.10.0 or v1.10.0-alpha4, and not latest-dev.
func isVersionTag(tag string) bool {
	return len(tag) >= 2 && tag[0] == 'v' && tag[1] >= '0' && tag[1] <= '9'
}

func SHA() string {
	if v := os.Getenv("GIT_SHA"); v != "" {
		return v
	}
	v, _ := do("git", "rev-parse", "HEAD")
	return v
}

func ShortSHA() string {
	if v := os.Getenv("GIT_SHA"); v != "" {
		return v
	}
	v, _ := do("git", "rev-parse", "--short", "HEAD")
	return v
}
