#!/bin/sh
# Builds the Linux binaries (opennox, opennox-server) into build/.
# Runs in a 64-bit (amd64) Debian container from the repository root, see .github/workflows/build-and-release.yml:
#   docker run --rm --platform linux/amd64 -v "$PWD:/w" -w /w debian:bookworm scripts/ci/build-linux.sh
# Debian 12 is used on purpose: the binaries need the glibc version they were built with, or newer.
set -e
root="$(cd "$(dirname "$0")/../.." && pwd)"
cache="$root/gocache"

# no C compiler or libraries: SDL3 is loaded at runtime (the AppImage bundles it)
apt-get update -qq
apt-get install -y -qq --no-install-recommends ca-certificates curl jq git >/dev/null

goroot="$("$root/scripts/ci/install-go.sh" "$cache/toolchain")"
export PATH="$goroot/bin:$PATH" GOTOOLCHAIN=local
export GOMODCACHE="$cache/mod" GOCACHE="$cache/build"
# the architecture of the system it runs on: amd64, or arm64 on an ARM runner
case "$(uname -m)" in
aarch64) goarch=arm64 ;;
*) goarch=amd64 ;;
esac
export CGO_ENABLED=0 GOARCH=$goarch
git config --global --add safe.directory "$root"

mkdir -p "$root/build"
cd "$root/src"
go run ./internal/noxbuild -o "$root/build"
# the container runs as root; let the CI cache steps read the cache
chmod -R a+rX "$cache" "$root/build"
