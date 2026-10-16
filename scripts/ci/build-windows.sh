#!/bin/sh
# Cross-compiles the Windows binaries (opennox.exe, opennox-server.exe, opennox-debug.exe) into build/,
# together with the DLL they load (SDL3). No C compiler is needed.
# Runs in a Debian container from the repository root, see .github/workflows/build-and-release.yml:
#   docker run --rm --platform linux/amd64 -v "$PWD:/w" -w /w debian:bookworm scripts/ci/build-windows.sh
set -e
root="$(cd "$(dirname "$0")/../.." && pwd)"
cache="$root/gocache"
deps="$cache/windows-deps"

# official builds of the DLLs (the 64-bit parts are used); update the versions and checksums together
SDL3_VERSION=3.4.18
SDL3_SHA256=049b9dd711e3fb3f8d84feb998dfec54a75e17cf539c4a80ff4d058d839771cf

# no C compiler: the executables load SDL3.dll at runtime, which is shipped next to them
apt-get update -qq
apt-get install -y -qq --no-install-recommends ca-certificates curl jq git >/dev/null

# download <url> <sha256> <file>: downloads a file once and checks it
download() {
	if [ ! -f "$3" ]; then
		curl -fsSL -o "$3.tmp" "$1"
		mv "$3.tmp" "$3"
	fi
	echo "$2  $3" | sha256sum -c -
}
mkdir -p "$deps"
sdl="$deps/SDL3-$SDL3_VERSION/x86_64-w64-mingw32"
download "https://github.com/libsdl-org/SDL/releases/download/release-$SDL3_VERSION/SDL3-devel-$SDL3_VERSION-mingw.tar.gz" \
	"$SDL3_SHA256" "$deps/SDL3-devel-$SDL3_VERSION-mingw.tar.gz"
tar -xzf "$deps/SDL3-devel-$SDL3_VERSION-mingw.tar.gz" -C "$deps" "SDL3-$SDL3_VERSION/x86_64-w64-mingw32"

goroot="$("$root/scripts/ci/install-go.sh" "$cache/toolchain")"
export PATH="$goroot/bin:$PATH" GOTOOLCHAIN=local
export GOMODCACHE="$cache/mod" GOCACHE="$cache/build"
git config --global --add safe.directory "$root"

out="$root/build"
mkdir -p "$out"
cd "$root/src"
go run ./internal/noxbuild --os windows -o "$out"

# Ship the DLL the executables load. The official build is self-contained: it only needs DLLs that come with Windows.
cp "$sdl/bin/SDL3.dll" "$out/"
ls -l "$out"
chmod -R a+rX "$cache" "$out"
