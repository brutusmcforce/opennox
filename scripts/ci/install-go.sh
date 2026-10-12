#!/bin/sh
# Installs the Go version required by src/go.mod for the architecture of the current (Debian) system,
# and prints its GOROOT. The download is checked against the checksum published on go.dev.
# usage: install-go.sh <dir>   (installs into <dir>/go<version>-<arch>, reused if it exists)
set -e
dir="$1"
root="$(cd "$(dirname "$0")/../.." && pwd)"
version="go$(sed -n 's/^go \([0-9.]*\)$/\1/p' "$root/src/go.mod")"
case "$(dpkg --print-architecture)" in
i386) arch=386 ;;
amd64) arch=amd64 ;;
*) echo "unsupported architecture: $(dpkg --print-architecture)" >&2; exit 1 ;;
esac
goroot="$dir/$version-$arch/go"
if [ -x "$goroot/bin/go" ]; then
	echo "$goroot"
	exit 0
fi
file="$version.linux-$arch.tar.gz"
sum="$(curl -fsSL 'https://go.dev/dl/?mode=json&include=all' | jq -r --arg f "$file" '.[].files[] | select(.filename == $f) | .sha256')"
if [ -z "$sum" ]; then
	echo "no checksum for $file on go.dev" >&2
	exit 1
fi
tmp="$(mktemp)"
curl -fsSL -o "$tmp" "https://go.dev/dl/$file"
echo "$sum  $tmp" | sha256sum -c - >&2
mkdir -p "$dir/$version-$arch"
tar -xzf "$tmp" -C "$dir/$version-$arch"
rm -f "$tmp"
echo "$goroot"
