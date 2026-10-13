# Building for Linux

This is the main build procedure for Linux version.

## Installing Go compiler

Your Go installation must be at least version 1.22.

If you have the supported version in your repository, you can install Go from there (`golang` package).

Otherwise, use Snap version
```bash
sudo snap install --classic go
```
or see [official docs](https://golang.org/doc/install) for manual installation.

## Installing C compiler

```bash
sudo apt install gcc
```

## Installing other dependencies

```bash
sudo apt update
sudo apt install libsdl2-dev libopenal-dev
```

## Building

```bash
cd src
go run ./internal/noxbuild
```

This should produce `opennox` and `opennox-server` binaries (64-bit, amd64; 32-bit builds are not supported).
