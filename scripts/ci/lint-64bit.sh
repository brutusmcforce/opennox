#!/bin/sh
# Checks for new 64-bit layout hazards in legacy code:
#   - raw 386 byte offsets into Go structs whose field moved on amd64 (internal/offsets, MOVED sites);
#   - pointers kept in 32-bit integers (internal/ptrint);
#   - blob memory (the original binary's 32-bit memory images) used as a Go type whose layout differs on amd64
#     (internal/blobcast).
# Known, reviewed sites are listed in scripts/ci/lint-64bit.baseline (without line numbers); the check fails
# only on new ones. After reviewing a new site, update the baseline with: scripts/ci/lint-64bit.sh -update
# Runs from anywhere; needs Go (no cgo).
set -e
root="$(cd "$(dirname "$0")/../.." && pwd)"
baseline="$root/scripts/ci/lint-64bit.baseline"
tmp="$(mktemp -d)"
trap 'rm -rf "$tmp"' EXIT
cd "$root/src"
# file:line -> file, collapse spaces
norm() { sed -E 's/^([^ :]+):[0-9]+/\1/; s/[[:space:]]+/ /g; s/ $//'; }
go run ./internal/offsets -v ./... >"$tmp/offsets.txt"
go run ./internal/ptrint -v ./... >"$tmp/ptrint.txt"
go run ./internal/blobcast ./... >"$tmp/blobcast.txt"
{
	grep -E ' MOVED ' "$tmp/offsets.txt" | norm | sed 's/^/offset: /'
	grep -E '^[^ #]+:[0-9]+ +(narrow|widen|arith) ' "$tmp/ptrint.txt" | norm | sed 's/^/ptrint: /'
	norm <"$tmp/blobcast.txt" | sed 's/^/blobcast: /'
} | sort >"$tmp/found.txt"
if [ "$1" = "-update" ]; then
	cp "$tmp/found.txt" "$baseline"
	echo "updated $baseline ($(wc -l <"$baseline") sites)"
	exit 0
fi
touch "$baseline"
sort "$baseline" >"$tmp/baseline.txt"
comm -13 "$tmp/baseline.txt" "$tmp/found.txt" >"$tmp/new.txt"
if [ -s "$tmp/new.txt" ]; then
	echo "New 64-bit layout hazards (check them, then run scripts/ci/lint-64bit.sh -update):"
	cat "$tmp/new.txt"
	exit 1
fi
echo "lint-64bit: no new sites ($(wc -l <"$tmp/found.txt") known)"
