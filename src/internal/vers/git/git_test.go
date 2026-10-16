package git

import (
	"os/exec"
	"testing"

	"github.com/brutusmcforce/opennox/internal/version"
)

func gitCmd(t *testing.T, args ...string) {
	t.Helper()
	if out, err := exec.Command("git", args...).CombinedOutput(); err != nil {
		t.Fatalf("git %v: %v: %s", args, err, out)
	}
}

// A release commit also has the latest-dev tag of the development build, and the version must be the release
// tag: the v1.10.0-alpha3 binaries reported "latest-dev", because git name-rev returns that one.
func TestVersion(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git is not installed")
	}
	t.Setenv("GIT_TAG", "")
	t.Setenv("GIT_SHA", "")
	t.Chdir(t.TempDir())
	gitCmd(t, "init", "-q")
	gitCmd(t, "-c", "user.email=test@test", "-c", "user.name=test", "commit", "-q", "--allow-empty", "-m", "release")

	t.Run("no tag", func(t *testing.T) {
		vers, sha := Version()
		if vers != version.DefVersion || sha == "" {
			t.Fatalf("Version() = %q, %q; want %q and a commit", vers, sha, version.DefVersion)
		}
	})

	gitCmd(t, "tag", "latest-dev")
	gitCmd(t, "tag", "v9.9.9")

	t.Run("release tag and latest-dev", func(t *testing.T) {
		vers, sha := Version()
		if vers != "v9.9.9" || sha != "" {
			t.Fatalf("Version() = %q, %q; want v9.9.9 and no commit", vers, sha)
		}
	})

	t.Run("commit after the tag", func(t *testing.T) {
		gitCmd(t, "-c", "user.email=test@test", "-c", "user.name=test", "commit", "-q", "--allow-empty", "-m", "next")
		vers, sha := Version()
		if vers != version.DefVersion || sha == "" {
			t.Fatalf("Version() = %q, %q; want %q and a commit", vers, sha, version.DefVersion)
		}
	})
}
