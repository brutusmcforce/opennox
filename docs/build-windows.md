# Building for Windows

This is the main build procedure for Windows version.

No C compiler is needed: OpenNox loads SDL3 and OpenAL when it starts.

For cross-compilation from Linux to Windows (including WSL), see [this page](./build-windows-on-linux.md).

## Installing Go compiler

See [official docs](https://go.dev/dl/).

To check that installation was successful, open `cmd.exe` and write `go version`.
It should print the version you installed.

## Downloading the project

### As a ZIP file

Download [the latest source](https://github.com/brutusmcforce/opennox/archive/refs/heads/dev.zip) and extract it to some directory.

We will assume the project is unpacked to `C:\Dev\opennox` for this guide. Prefer path without spaces and non-English characters.

### Using Git

1. Install [Git for Windows](https://gitforwindows.org).
   Click "next" for all the steps - defaults are good already.
2. Install [TortoiseGit](https://tortoisegit.org/download/).
   When asked for a name and email, either fill it in, or skip by checking the box.
   It's only used when sending changes to the project.
3. Create a directory for projects, let's say `C:\Dev\`.
4. Right click on empty space in the new folder, select "Git Clone...".
5. Paste the project URL: `https://github.com/brutusmcforce/opennox.git`.
6. Check the "Branch" box, write `dev` to the field.
7. Click OK, wait for the project to download.

## Building

Open `cmd.exe` (or PowerShell). Assuming OpenNox was downloaded into `C:\Dev\opennox`:

```bat
cd C:\Dev\opennox\src
go run ./internal/noxbuild
```

You should see that the project dependencies are being downloaded and compiled.
Build may take a while, especially when building for the first time.

Once finished, it should produce `opennox.exe`, `opennox-server.exe` and `opennox-debug.exe` in `C:\Dev\opennox\src`.

To build only the game: `go run ./internal/noxbuild client`.

## Running

The game loads two DLLs, which must be next to `opennox.exe`:

- `SDL3.dll`: from the SDL3 release for Windows (`SDL3-<version>-win32-x64.zip` on the [SDL releases page](https://github.com/libsdl-org/SDL/releases)).
- `OpenAL32.dll`: OpenAL Soft's `bin/Win64/soft_oal.dll` from the [OpenAL Soft releases](https://github.com/kcat/openal-soft/releases), renamed to `OpenAL32.dll`.

Both are also in every OpenNox release for Windows.

1. Create a new folder for OpenNox, let's say `C:\Games\OpenNox`.
2. Copy `opennox.exe` and `opennox-server.exe` from `C:\Dev\opennox\src`, and the two DLLs.
3. Run `opennox.exe`.
   - If it says that the Nox directory was not found, it works!
     Check [this guide](https://mod.io/g/nox/r/how-to-install-opennox-windows) to set the Nox game path.
