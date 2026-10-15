# Building for Windows on Linux

This is the main build procedure for the Windows version cross-compiled on Linux.
No C compiler is needed: OpenNox loads SDL3 and OpenAL when it starts.

For building on Windows, see [this page](./build-windows.md).

## Installing Go compiler

See the [official docs](https://go.dev/doc/install). Go 1.25 or newer is needed.

## Building

```bash
cd src
go run ./internal/noxbuild --os windows
```

This should produce `opennox.exe`, `opennox-server.exe` and `opennox-debug.exe`.

`make build-client-win` builds only `opennox.exe`.

## Running

Put `SDL3.dll` and `OpenAL32.dll` next to the executables, see [building on Windows](./build-windows.md#running).

The release build does the same in a container: `scripts/ci/build-windows.sh` downloads the DLLs and checks them.
