#!/bin/sh
# Builds build/OpenNox-macos.zip: OpenNox.app (universal: Apple Silicon and Intel) with SDL3 inside, and the
# universal opennox-server.
# Runs on a macOS runner from the repository root, see .github/workflows/build-and-release.yml. Needs Go and the
# macOS tools (lipo, hdiutil, sips, iconutil, codesign, ditto); no C compiler.
set -e
root="$(cd "$(dirname "$0")/../.." && pwd)"
cache="$root/gocache/macos-deps"

# the official SDL3 build for Apple platforms; update the version and checksum together
SDL3_VERSION=3.4.18
SDL3_SHA256=c1adf3da91108dbc3e4205da5dacebc203085d283fa89e86eb2c889abd03a853

mkdir -p "$cache"
dmg="$cache/SDL3-$SDL3_VERSION.dmg"
if [ ! -f "$dmg" ]; then
	curl -fsSL -o "$dmg.tmp" "https://github.com/libsdl-org/SDL/releases/download/release-$SDL3_VERSION/SDL3-$SDL3_VERSION.dmg"
	mv "$dmg.tmp" "$dmg"
fi
echo "$SDL3_SHA256  $dmg" | shasum -a 256 -c -

out="$root/build"
tmp="$root/build/macos-tmp"
rm -rf "$tmp"
mkdir -p "$tmp/arm64" "$tmp/amd64" "$out"

# Go builds both architectures without cgo; lipo joins them into universal binaries
cd "$root/src"
for arch in arm64 amd64; do
	GOARCH=$arch go run ./internal/noxbuild --os darwin -o "$tmp/$arch" opennox opennox-server
done
for bin in opennox opennox-server; do
	lipo -create -output "$tmp/$bin" "$tmp/arm64/$bin" "$tmp/amd64/$bin"
done

app="$tmp/OpenNox.app"
mkdir -p "$app/Contents/MacOS" "$app/Contents/Frameworks" "$app/Contents/Resources"
cp "$tmp/opennox" "$app/Contents/MacOS/opennox"

# SDL3.framework from the disk image (the macOS part of the xcframework is universal);
# the game looks for it in Contents/Frameworks (see libs/dynlib)
mnt="$tmp/sdl3-mount"
mkdir -p "$mnt"
hdiutil attach -nobrowse -readonly -mountpoint "$mnt" "$dmg" >/dev/null
ditto "$mnt/SDL3.xcframework/macos-arm64_x86_64/SDL3.framework" "$app/Contents/Frameworks/SDL3.framework"
hdiutil detach "$mnt" >/dev/null

# icon
iconset="$tmp/opennox.iconset"
mkdir -p "$iconset"
for size in 16 32 128 256 512; do
	sips -z $size $size "$root/res/opennox_512.png" --out "$iconset/icon_${size}x${size}.png" >/dev/null
	double=$((size * 2))
	if [ $double -le 512 ]; then
		sips -z $double $double "$root/res/opennox_512.png" --out "$iconset/icon_${size}x${size}@2x.png" >/dev/null
	fi
done
iconutil -c icns -o "$app/Contents/Resources/opennox.icns" "$iconset"

version="$(git -C "$root" describe --tags --always 2>/dev/null || echo dev)"
short="$(echo "$version" | sed -E 's/^v?([0-9]+(\.[0-9]+)*).*/\1/')"
cat >"$app/Contents/Info.plist" <<EOF
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>CFBundleName</key>
	<string>OpenNox</string>
	<key>CFBundleDisplayName</key>
	<string>OpenNox</string>
	<key>CFBundleIdentifier</key>
	<string>io.github.opennox</string>
	<key>CFBundleExecutable</key>
	<string>opennox</string>
	<key>CFBundleIconFile</key>
	<string>opennox</string>
	<key>CFBundlePackageType</key>
	<string>APPL</string>
	<key>CFBundleShortVersionString</key>
	<string>$short</string>
	<key>CFBundleVersion</key>
	<string>$version</string>
	<key>LSMinimumSystemVersion</key>
	<string>11.0</string>
	<key>LSApplicationCategoryType</key>
	<string>public.app-category.role-playing-games</string>
	<key>NSHighResolutionCapable</key>
	<true/>
</dict>
</plist>
EOF

# Ad-hoc signature: Apple Silicon only runs signed code, and lipo and the added files invalidate the Go linker's
# signatures. Without a Developer ID, Gatekeeper still asks on the first start (right-click > Open).
codesign --force --sign - "$app/Contents/Frameworks/SDL3.framework"
codesign --force --sign - "$app"
codesign --force --sign - "$tmp/opennox-server"
codesign --verify --deep --strict "$app"

# in a directory of its own: the macOS file system is case insensitive, so "OpenNox" next to the "opennox"
# binary that lipo wrote above is the same path, and the mkdir fails
pkg="$tmp/zip/OpenNox"
mkdir -p "$pkg"
ditto "$app" "$pkg/OpenNox.app"
cp "$tmp/opennox-server" "$root/CHANGELOG.md" "$pkg/"
rm -f "$out/OpenNox-macos.zip"
ditto -c -k --keepParent "$pkg" "$out/OpenNox-macos.zip"
rm -rf "$tmp"
ls -l "$out/OpenNox-macos.zip"
