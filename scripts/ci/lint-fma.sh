#!/bin/sh
# Checks that the arm64 build contains no fused multiply-add (FMA) instructions.
#
# arm64 fuses "a + b*c" into one FMADD, which rounds once where amd64 rounds twice, so the two
# architectures compute different floats from the same code. The game must simulate identically
# everywhere (multiplayer, saved games, the e2e references), so the builds set GOCOMPILEDEBUG=fmahash=qn
# (src/internal/noxbuild/main.go), which turns the optimization off in every package.
#
# This checks the result rather than the flag: the arm64 assembly of the whole binary, the standard
# library included, must have no FMA instruction. It also counts them WITHOUT the flag, so that the
# check can't pass because the toolchain renamed the instructions or stopped emitting the listed ones.
#
# An intentional math.FMA call would also show up here (it compiles to FMADDD and is the same
# operation on amd64, so it would be correct) - nothing uses it today.
set -e
cd "$(dirname "$0")/../../src"

pattern='[[:space:]](FMADD|FMSUB|FNMADD|FNMSUB)[DS][[:space:]]'
asm() {
	GOOS=linux GOARCH=arm64 go build -gcflags=all=-S -o /dev/null ./cmd/opennox 2>&1
}

echo "== arm64 assembly without fmahash=qn (expected to fuse)"
baseline=$(asm | grep -cE "$pattern" || true)
echo "FMA instructions: $baseline"
if [ "$baseline" -eq 0 ]; then
	echo "FAIL: no FMA instructions even with the optimization on;" >&2
	echo "the toolchain changed, and this check no longer proves anything." >&2
	exit 1
fi

echo "== arm64 assembly as the game is built (fmahash=qn)"
found=$(GOCOMPILEDEBUG=fmahash=qn asm | grep -E "$pattern" || true)
if [ -n "$found" ]; then
	echo "FAIL: the arm64 build still fuses multiply-add, so it computes different floats than amd64:" >&2
	echo "$found" | head -20 >&2
	exit 1
fi
echo "FMA instructions: 0"
echo "OK: the arm64 build computes floats like amd64"
