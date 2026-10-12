#!/bin/sh
# Builds build/opennox-bundle-i386.AppImage from build/opennox and build/opennox-server.
# Runs in a 32-bit Debian container from the repository root, see .github/workflows/build-and-release.yml:
#   docker run --rm --platform linux/386 -v "$PWD:/w" -w /w debian:bookworm scripts/appimage/build.sh
# The AppImage bundles SDL2, OpenAL and their dependencies from Debian 12, the same system the binaries are built on.
set -e
ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
BUILD="$ROOT/build"
APPDIR="$ROOT/opennox_bundle.AppDir"
CACHE="$ROOT/gocache/appimage"

# update the versions and checksums together
APPIMAGETOOL_VERSION=1.9.1
APPIMAGETOOL_SHA256=7ad9ff47c203aae0149b18f6df9e3018b2e2f470ea644a0413e3ded39e9e3bdb
RUNTIME_VERSION=20251108
RUNTIME_SHA256=e72ea0b140a0a16e680713238a6f30aad278b62c4ca17919c554864124515498

# paths relative to the XDG directories
OPENNOX_CFG="opennox/opennox.yml"
OPENNOX_DATA="opennox"
OPENNOX_STATEFILES="opennox"

apt-get update -qq
apt-get install -y -qq --no-install-recommends ca-certificates curl file libsdl2-2.0-0 libopenal1 >/dev/null

# download <url> <sha256> <file>: downloads a file once and checks it
download() {
	if [ ! -f "$3" ]; then
		curl -fsSL -o "$3.tmp" "$1"
		mv "$3.tmp" "$3"
	fi
	echo "$2  $3" | sha256sum -c -
}
mkdir -p "$CACHE"
TOOL="$CACHE/appimagetool-$APPIMAGETOOL_VERSION-i686.AppImage"
RUNTIME="$CACHE/runtime-$RUNTIME_VERSION-i686"
download "https://github.com/AppImage/appimagetool/releases/download/$APPIMAGETOOL_VERSION/appimagetool-i686.AppImage" \
	"$APPIMAGETOOL_SHA256" "$TOOL"
download "https://github.com/AppImage/type2-runtime/releases/download/$RUNTIME_VERSION/runtime-i686" \
	"$RUNTIME_SHA256" "$RUNTIME"
chmod +x "$TOOL"

echo "Preparing the AppDir..."
rm -rf "$APPDIR"
mkdir -p "$APPDIR/usr/bin" "$APPDIR/usr/lib"
for bin in opennox opennox-server; do
	install -m 755 "$BUILD/$bin" "$APPDIR/usr/bin/$bin"
done

# SDL2, OpenAL and everything they link to (ldd lists all of it), except the libraries that must come from the
# user's system: glibc, and the graphics libraries that belong to the driver.
for lib in libSDL2-2.0.so.0 libopenal.so.1; do
	path="/usr/lib/i386-linux-gnu/$lib"
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
ARCH=i686 "$TOOL" --appimage-extract-and-run -n --runtime-file "$RUNTIME" \
	"$APPDIR" "$BUILD/opennox-bundle-i386.AppImage"
rm -rf "$APPDIR"
chmod -R a+rX "$CACHE" "$BUILD"
echo "Done: $BUILD/opennox-bundle-i386.AppImage"
