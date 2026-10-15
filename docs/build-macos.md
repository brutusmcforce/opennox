# Building for macOS

No C compiler is needed: OpenNox loads SDL3 and OpenAL when it starts. macOS includes OpenAL.

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

## Running

OpenNox looks for the Nox game files in a folder named `Nox` in your home folder or in
`~/Library/Application Support`, or use `-data <path>`.
It keeps its settings and saves in `~/Library/Application Support/OpenNox`, and logs in `~/Library/Logs/OpenNox`.

The release app is not signed with an Apple Developer ID: the first time, open it with right-click > Open.
