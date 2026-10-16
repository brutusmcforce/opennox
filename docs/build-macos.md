# Building for macOS

No C compiler is needed: OpenNox loads SDL3 when it starts.

## Installing Go and SDL3

Install [Go](https://go.dev/dl/) 1.25 or newer, and SDL3 with [Homebrew](https://brew.sh): `brew install sdl3`.

## Building

```bash
cd src
go run ./internal/noxbuild opennox opennox-server
```

This builds for the Mac you're on. For the other architecture, set `GOARCH=amd64` (Intel) or `GOARCH=arm64`
(Apple Silicon); any system can build for macOS this way.

`scripts/ci/build-macos.sh` builds the release: a universal `OpenNox.app` with SDL3 inside.

## Floating point on arm64

arm64 (Apple Silicon, ARM Linux) fuses a multiply followed by an add into a single FMA instruction, which
rounds once where amd64 rounds twice, so the same code computes slightly different floats. The game has to
simulate identically on every platform: multiplayer clients run the same simulation, saved games and the e2e
references come from amd64.

Go permits the fusion and has no language-level switch for it, but the compiler can be told not to do it:

```
GOCOMPILEDEBUG=fmahash=qn
```

`internal/noxbuild` sets this for every build, so the release binaries never fuse - in the game's own code and
in Go's `math` package alike (`math.Sin`, `Cos`, `Atan2`, `Hypot`, `Pow` and `Acos` all fuse without it, and
the game uses them). It does nothing on amd64, which never fuses: the amd64 machine code is byte for byte the
same with and without it. `scripts/ci/lint-fma.sh` checks the result rather than the flag - it compiles the
whole arm64 binary and fails if a single FMA instruction is left, and also checks that the build fuses
*without* the flag, so the check cannot pass because a new toolchain renamed the instructions.

A plain `go build` or `go test` does not read noxbuild's setting. On an arm64 machine, pass it yourself when
the results have to match the other platforms:

```bash
GOCOMPILEDEBUG=fmahash=qn go test ./...
GOOS=darwin GOARCH=arm64 go build -gcflags='all=-S' ./... 2>&1 | grep -E 'FMADD|FMSUB|FNMADD|FNMSUB' # lists what fuses
```

An expression can also be written so that no compiler may fuse it, with a conversion that forces the
intermediate rounding:

```go
sum := a + float64(b*c) // not a + b*c, which arm64 fuses
```

`libs/prand` is written that way, so the random numbers match even in a build without the flag.

What is still untested on arm64: whether the simulation really is identical end to end. The e2e test and the
multiplayer test have only ever run on amd64, so an arm64 build has not yet been checked against the recorded
references or played against an amd64 client.

## Running

OpenNox looks for the Nox game files in a folder named `Nox` in your home folder or in
`~/Library/Application Support`, or use `-data <path>`.
It keeps its settings and saves in `~/Library/Application Support/OpenNox`, and logs in `~/Library/Logs/OpenNox`.

The release app is not signed with an Apple Developer ID: the first time, open it with right-click > Open.
