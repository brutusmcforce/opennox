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

## Known issue: floating point on Apple Silicon

arm64 fuses a multiply followed by an add into a single FMA instruction, which rounds once where amd64 rounds
twice, so the same code gives slightly different results on Apple Silicon than on Linux, Windows and Intel
Macs. Go permits this and has no flag to turn it off; the fix is per expression, an explicit conversion that
forces the intermediate rounding:

```go
sum := a + float64(b*c) // not a + b*c, which arm64 fuses
```

`libs/prand` is written that way already, because the generator has to produce the same numbers everywhere.
The rest of the engine is not: the arm64 build has about 660 of these sites, mostly in the geometry, AI and
map code (`legacy/placement_trig.go`, `server/object_ai_path.go`, `server/wall.go`, `client/sight.go`...),
plus the ones inside Go's own `math` package. Whether any of it is visible in the game is untested, which is
why the macOS build is experimental. To list the sites:

```bash
GOOS=darwin GOARCH=arm64 go build -gcflags='all=-S' ./... 2>&1 | grep -E 'FMADDD|FMSUBD|FNMADDD|FNMSUBD'
```

If Apple Silicon ever has to match the other platforms exactly, converting those expressions is the way to do
it; building the release for amd64 only and letting Rosetta 2 run it is the blunt alternative.

## Running

OpenNox looks for the Nox game files in a folder named `Nox` in your home folder or in
`~/Library/Application Support`, or use `-data <path>`.
It keeps its settings and saves in `~/Library/Application Support/OpenNox`, and logs in `~/Library/Logs/OpenNox`.

The release app is not signed with an Apple Developer ID: the first time, open it with right-click > Open.
