#!/bin/sh
# Cross-compiles the Windows binaries (opennox.exe, opennox-server.exe, opennox-debug.exe) into build/,
# together with the DLLs they need.
# Runs in a Debian container from the repository root, see .github/workflows/build-and-release.yml:
#   docker run --rm --platform linux/amd64 -v "$PWD:/w" -w /w debian:bookworm scripts/ci/build-windows.sh
set -e
root="$(cd "$(dirname "$0")/../.." && pwd)"
cache="$root/gocache"
deps="$cache/windows-deps"

# official 32-bit MinGW builds; update the versions and checksums together
SDL2_VERSION=2.32.10
SDL2_SHA256=83a5d74012311edc3c0d40ea6faecbe57ad692aa033fa5dc273cc937e3938ff2
OPENAL_VERSION=1.25.2
OPENAL_SHA256=67a0c4b800bd860c93c04f38caf8cbe4875f9c84700ac430efc451f70e265434

apt-get update -qq
apt-get install -y -qq --no-install-recommends ca-certificates curl jq git unzip gcc-mingw-w64-i686 binutils-mingw-w64-i686 >/dev/null

# download <url> <sha256> <file>: downloads a file once and checks it
download() {
	if [ ! -f "$3" ]; then
		curl -fsSL -o "$3.tmp" "$1"
		mv "$3.tmp" "$3"
	fi
	echo "$2  $3" | sha256sum -c -
}
mkdir -p "$deps"
sdl="$deps/SDL2-$SDL2_VERSION/i686-w64-mingw32"
download "https://github.com/libsdl-org/SDL/releases/download/release-$SDL2_VERSION/SDL2-devel-$SDL2_VERSION-mingw.tar.gz" \
	"$SDL2_SHA256" "$deps/SDL2-devel-$SDL2_VERSION-mingw.tar.gz"
tar -xzf "$deps/SDL2-devel-$SDL2_VERSION-mingw.tar.gz" -C "$deps" "SDL2-$SDL2_VERSION/i686-w64-mingw32"
al="$deps/openal-soft-$OPENAL_VERSION-bin"
download "https://github.com/kcat/openal-soft/releases/download/$OPENAL_VERSION/openal-soft-$OPENAL_VERSION-bin.zip" \
	"$OPENAL_SHA256" "$deps/openal-soft-$OPENAL_VERSION-bin.zip"
rm -rf "$al"
unzip -q "$deps/openal-soft-$OPENAL_VERSION-bin.zip" -d "$deps"

goroot="$("$root/scripts/ci/install-go.sh" "$cache/toolchain")"
export PATH="$goroot/bin:$PATH" GOTOOLCHAIN=local
export GOMODCACHE="$cache/mod" GOCACHE="$cache/build"
# for the MinGW compiler and linker: SDL2/SDL.h, AL/al.h, -lSDL2, -lOpenAL32 (the cross linker ignores LIBRARY_PATH)
export C_INCLUDE_PATH="$sdl/include:$al/include"
export CGO_LDFLAGS="-L$sdl/lib -L$al/libs/Win32"
git config --global --add safe.directory "$root"

out="$root/build"
mkdir -p "$out"
cd "$root/src"
go run ./internal/noxbuild --os windows -o "$out"

# Ship the DLLs the executables need, except the ones that come with Windows.
# OpenAL Soft's own DLL is distributed as OpenAL32.dll, which is what the game links to.
cp "$sdl/bin/SDL2.dll" "$out/"
cp "$al/bin/Win32/soft_oal.dll" "$out/OpenAL32.dll"
gccdir="$(dirname "$(i686-w64-mingw32-gcc -print-libgcc-file-name)")"
search="$out $gccdir /usr/i686-w64-mingw32/lib"
todo="$out/opennox.exe $out/opennox-server.exe $out/opennox-debug.exe $out/SDL2.dll $out/OpenAL32.dll"
while [ -n "$todo" ]; do
	next=""
	for f in $todo; do
		for dll in $(i686-w64-mingw32-objdump -p "$f" | sed -n 's/^\s*DLL Name: //p'); do
			[ -f "$out/$dll" ] && continue
			for d in $search; do
				if [ -f "$d/$dll" ]; then
					cp "$d/$dll" "$out/"
					next="$next $out/$dll"
					break
				fi
			done
		done
	done
	todo="$next"
done
ls -l "$out"
chmod -R a+rX "$cache" "$out"
