#!/bin/sh
# Builds the Linux binaries (opennox, opennox-server) into build/.
# Runs in a 32-bit Debian container from the repository root, see .github/workflows/build-and-release.yml:
#   docker run --rm --platform linux/386 -v "$PWD:/w" -w /w debian:bookworm scripts/ci/build-linux.sh
# Debian 12 is used on purpose: the binaries need the glibc version they were built with, or newer.
set -e
root="$(cd "$(dirname "$0")/../.." && pwd)"
cache="$root/gocache"

apt-get update -qq
apt-get install -y -qq --no-install-recommends ca-certificates curl jq git gcc libc6-dev pkg-config \
	libsdl2-dev libopenal-dev >/dev/null

goroot="$("$root/scripts/ci/install-go.sh" "$cache/toolchain")"
export PATH="$goroot/bin:$PATH" GOTOOLCHAIN=local
export GOMODCACHE="$cache/mod" GOCACHE="$cache/build"
export CGO_ENABLED=1 GOARCH=386
git config --global --add safe.directory "$root"

mkdir -p "$root/build"
cd "$root/src"
go run ./internal/noxbuild -o "$root/build"
# the container runs as root; let the CI cache steps read the cache
chmod -R a+rX "$cache" "$root/build"
