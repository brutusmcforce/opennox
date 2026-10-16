#!/bin/sh
# Builds build/opennox-bundle-<arch>.AppImage (x86_64 or aarch64, the architecture of the system it runs on) from
# build/opennox and build/opennox-server.
# Runs in a 64-bit (amd64) Debian container from the repository root, see .github/workflows/build-and-release.yml:
#   docker run --rm --platform linux/amd64 -v "$PWD:/w" -w /w debian:bookworm scripts/appimage/build.sh
# The AppImage bundles SDL3 (built from source here: Debian 12 has no SDL3 package, and a newer Debian would raise the
# glibc version the AppImage needs). OpenNox plays audio through SDL3 too.
set -e
ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
BUILD="$ROOT/build"
APPDIR="$ROOT/opennox_bundle.AppDir"
CACHE="$ROOT/gocache/appimage"

# update the versions and checksums together
APPIMAGETOOL_VERSION=1.9.1
RUNTIME_VERSION=20251108
ARCH="$(uname -m)"
case "$ARCH" in
x86_64)
	APPIMAGETOOL_SHA256=ed4ce84f0d9caff66f50bcca6ff6f35aae54ce8135408b3fa33abfc3cb384eb0
	RUNTIME_SHA256=2fca8b443c92510f1483a883f60061ad09b46b978b2631c807cd873a47ec260d
	;;
aarch64)
	APPIMAGETOOL_SHA256=f0837e7448a0c1e4e650a93bb3e85802546e60654ef287576f46c71c126a9158
	RUNTIME_SHA256=00cbdfcf917cc6c0ff6d3347d59e0ca1f7f45a6df1a428a0d6d8a78664d87444
	;;
*)
	echo "unsupported architecture: $ARCH"
	exit 1
	;;
esac
SDL3_VERSION=3.4.18
SDL3_SHA256=9c75cf16330322c217dedd2e0609f1124f1b54b8633e763467b4684d0f4334a3

# paths relative to the XDG directories
OPENNOX_CFG="opennox/opennox.yml"
OPENNOX_DATA="opennox"
OPENNOX_STATEFILES="opennox"

apt-get update -qq
# SDL3 is built with its X11, Wayland, PipeWire, PulseAudio and ALSA backends; it loads their libraries at runtime,
# so it works on systems with any of them. OpenNox doesn't use SDL's renderer or GPU API.
apt-get install -y -qq --no-install-recommends ca-certificates curl file \
	cmake ninja-build gcc libc6-dev pkg-config \
	libx11-dev libxext-dev libxrandr-dev libxcursor-dev libxi-dev libxfixes-dev libxss-dev libxtst-dev libxkbcommon-dev \
	libwayland-dev libwayland-bin wayland-protocols libdecor-0-dev libegl-dev libgl-dev libdbus-1-dev libudev-dev \
	libpipewire-0.3-dev libpulse-dev libasound2-dev >/dev/null

# download <url> <sha256> <file>: downloads a file once and checks it
download() {
	if [ ! -f "$3" ]; then
		curl -fsSL -o "$3.tmp" "$1"
		mv "$3.tmp" "$3"
	fi
	echo "$2  $3" | sha256sum -c -
}
mkdir -p "$CACHE"
TOOL="$CACHE/appimagetool-$APPIMAGETOOL_VERSION-$ARCH.AppImage"
RUNTIME="$CACHE/runtime-$RUNTIME_VERSION-$ARCH"
download "https://github.com/AppImage/appimagetool/releases/download/$APPIMAGETOOL_VERSION/appimagetool-$ARCH.AppImage" \
	"$APPIMAGETOOL_SHA256" "$TOOL"
download "https://github.com/AppImage/type2-runtime/releases/download/$RUNTIME_VERSION/runtime-$ARCH" \
	"$RUNTIME_SHA256" "$RUNTIME"
chmod +x "$TOOL"

SDL3_LIB="$CACHE/sdl3-$SDL3_VERSION-audio-$ARCH/libSDL3.so.0"
if [ ! -f "$SDL3_LIB" ]; then
	echo "Building SDL3 $SDL3_VERSION..."
	download "https://github.com/libsdl-org/SDL/releases/download/release-$SDL3_VERSION/SDL3-$SDL3_VERSION.tar.gz" \
		"$SDL3_SHA256" "$CACHE/SDL3-$SDL3_VERSION.tar.gz"
	src="$CACHE/sdl3-src"
	rm -rf "$src"
	mkdir -p "$src"
	tar -xzf "$CACHE/SDL3-$SDL3_VERSION.tar.gz" -C "$src" --strip-components=1
	cmake -S "$src" -B "$src/build" -G Ninja -DCMAKE_BUILD_TYPE=Release \
		-DSDL_SHARED=ON -DSDL_STATIC=OFF -DSDL_TESTS=OFF -DSDL_EXAMPLES=OFF \
		-DSDL_AUDIO=ON -DSDL_PIPEWIRE=ON -DSDL_PULSEAUDIO=ON -DSDL_ALSA=ON -DSDL_RENDER=OFF -DSDL_GPU=OFF -DSDL_CAMERA=OFF -DSDL_KMSDRM=OFF -DSDL_VULKAN=OFF \
		-DSDL_X11=ON -DSDL_WAYLAND=ON -DSDL_OPENGL=ON >/dev/null
	cmake --build "$src/build" >/dev/null
	mkdir -p "$(dirname "$SDL3_LIB")"
	cp -L "$src/build/libSDL3.so.0" "$SDL3_LIB"
	rm -rf "$src"
fi

echo "Preparing the AppDir..."
rm -rf "$APPDIR"
mkdir -p "$APPDIR/usr/bin" "$APPDIR/usr/lib"
for bin in opennox opennox-server; do
	install -m 755 "$BUILD/$bin" "$APPDIR/usr/bin/$bin"
done

# SDL3 and everything it links to (ldd lists all of it), except the libraries that must come from the
# user's system: glibc, and the graphics libraries that belong to the driver.
for path in "$SDL3_LIB"; do
	echo "$path"
	ldd "$path" | sed -n 's/.* => \(\/[^ ]*\) .*/\1/p'
done | sort -u | while read -r path; do
	case "$(basename "$path")" in
	ld-linux.so.* | libc.so.* | libm.so.* | libdl.so.* | libpthread.so.* | librt.so.* | libresolv.so.* | \
		libGL.so.* | libGLX.so.* | libGLdispatch.so.* | libEGL.so.* | libdrm.so.* | libgbm.so.*)
		continue
		;;
	esac
	cp -L "$path" "$APPDIR/usr/lib/"
	echo "bundled $(basename "$path")"
done

cat > "$APPDIR/AppRun" <<EOF
#!/bin/sh
if [ -z "\$APPDIR" ]; then
	APPDIR="\$(dirname "\$(readlink -f "\$0")")"
fi
PATH="\$APPDIR/usr/bin\${PATH:+:\$PATH}"
LD_LIBRARY_PATH="\$APPDIR/usr/lib\${LD_LIBRARY_PATH:+:\$LD_LIBRARY_PATH}"
export APPDIR PATH LD_LIBRARY_PATH

cfgpath="\${XDG_CONFIG_HOME:-\$HOME/.config}/$OPENNOX_CFG"
datapath="\${XDG_DATA_HOME:-\$HOME/.local/share}/$OPENNOX_DATA"
statepath="\${XDG_STATE_HOME:-\$HOME/.local/state}/$OPENNOX_STATEFILES"

mkdir -p "\$(dirname "\$cfgpath")" "\$statepath" || {
	echo "Cannot create '\$(dirname "\$cfgpath")' or '\$statepath'. Check the permissions, or use portable mode (--appimage-portable-home)."
	exit 1
}
if [ ! -d "\$datapath" ]; then
	mkdir -p "\$datapath"
	echo "Please put the Nox game files here and start the AppImage again: \$datapath"
	exit 2
fi
cd "\$statepath"

case "\$1" in
help)
	echo "OpenNox AppImage:"
	echo "  \$ARGV0 [options]          run OpenNox"
	echo "  \$ARGV0 server [options]   run the OpenNox dedicated server"
	echo ""
	echo "Following the XDG Base Directory Specification:"
	echo "  config file: \$cfgpath"
	echo "  game files:  \$datapath"
	echo "  log files:   \$statepath"
	echo ""
	echo "To keep everything next to the AppImage instead, run it once with --appimage-portable-home."
	exit 0
	;;
server)
	shift
	exec opennox-server --config="\$cfgpath" --data="\$datapath" "\$@"
	;;
client)
	shift
	;;
esac
exec opennox --config="\$cfgpath" --data="\$datapath" "\$@"
EOF
chmod +x "$APPDIR/AppRun"

cat > "$APPDIR/opennox.desktop" <<EOF
[Desktop Entry]
Name=OpenNox
Comment=Modern implementation of the Nox game engine
Exec=AppRun
Icon=opennox
Terminal=false
Type=Application
Categories=Game;
EOF
cp "$ROOT/res/opennox_256.png" "$APPDIR/opennox.png"

echo "Building the AppImage..."
# appimagetool is an AppImage itself, and containers have no FUSE
ARCH=$ARCH "$TOOL" --appimage-extract-and-run -n --runtime-file "$RUNTIME" \
	"$APPDIR" "$BUILD/opennox-bundle-$ARCH.AppImage"
rm -rf "$APPDIR"
chmod -R a+rX "$CACHE" "$BUILD"
echo "Done: $BUILD/opennox-bundle-$ARCH.AppImage"
