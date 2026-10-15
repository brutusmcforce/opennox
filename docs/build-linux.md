# Building for Linux

This is the main build procedure for the Linux version.

No C compiler or development libraries are needed: OpenNox loads SDL3 and OpenAL when it starts.

## Installing Go compiler

Your Go installation must be at least version 1.25.

If you have the supported version in your repository, you can install Go from there (`golang` package).
Otherwise, see the [official docs](https://go.dev/doc/install).

## Building

```bash
cd src
go run ./internal/noxbuild
```

This should produce `opennox` and `opennox-server` binaries (64-bit, amd64; 32-bit builds are not supported).

## Running

The game needs the SDL3 and OpenAL libraries:

- Debian 13 / Ubuntu 25.04 or newer: `sudo apt install libsdl3-0 libopenal1`
- Arch Linux: `sudo pacman -S sdl3 openal`
- Fedora: `sudo dnf install SDL3 openal-soft`

On older distributions, use the AppImage from the releases: it includes both.

The dedicated server (`opennox-server`) doesn't need them.
