# Changelog

## v1.10.0-alpha3

- The engine no longer uses cgo, and the game moved from SDL2 to SDL3; SDL3, OpenGL and OpenAL are loaded when the game starts.
- Experimental macOS build: a universal `OpenNox.app` for Apple Silicon and Intel.
- The Linux archive needs SDL3 to be installed (Debian 12 and Ubuntu 24.04 do not package it yet; the AppImage bundles SDL3 and OpenAL and needs nothing). The Windows archive ships `SDL3.dll` and `OpenAL32.dll` next to the executable.
- Errors shown before the game window opens, such as missing game data, now use a native message box on Windows, macOS and Linux.
- Fix: the game could crash when opening the multiplayer server list if the socket used to ping servers could not be opened.
- Fix: builds that were not made from a release tag reported version v1.9.0-dev.
- Dedicated servers and the lobby: `/metrics` no longer reports the `go_*` and `process_*` Go runtime metrics; the `nox_*` metrics are unchanged.

### Internal

- No cgo anywhere (`CGO_ENABLED=0` for every target). New packages load the native libraries at runtime with [purego](https://github.com/ebitengine/purego): `libs/dynlib` (dlopen/LoadLibrary), `libs/sdl` (SDL3), `libs/gl` (OpenGL 3.3) and `libs/openal`. `go-gl/gl`, `veandco/go-sdl2` and the vendored `go-openal` are gone.
- The wrappers that call the native functions are marked `//go:uintptrescapes`. `purego.SyscallN` only promises to keep its own arguments alive, and that promise does not reach through a wrapper, so the Go pointers handed to SDL, OpenGL and OpenAL — window size out-parameters, OpenAL buffer and source IDs, uniform values, texture and vertex data — could be collected or moved during the call.
- Dependencies dropped: `julienschmidt/httprouter` (the map server uses `net/http` routing), `tawesoft/golib` (message boxes come from SDL), `golang.org/x/exp` (stdlib `slices`/`maps`), and `prometheus/client_golang`, which took eight more modules with it: the counters and gauges moved to the new `libs/metrics`, which writes the Prometheus text format. The binaries are about 3 MB smaller.
- Map server: `net/http` URL-unescapes a path wildcard, unlike the router it replaced, so a map name from the path is now checked to be a plain file name before it reaches a file path.
- `go vet` runs with `-unsafeptr=false`: the native bindings turn addresses that come from SDL, OpenGL and OpenAL into pointers, which is what the check warns about. It cannot be scoped to those packages, because `go vet` also reports a package's dependencies.
- The e2e screenshots are compared with a tolerance (`NOX_E2E_MAX_DELTA`, `NOX_E2E_MAX_DIFF_PCT`, `NOX_E2E_EXACT` for the old exact comparison). They are rasterized by the GPU, and OpenGL does not promise bit-identical results across drivers: against references recorded on AMD, 102 of the 123 screens differ on NVIDIA, Intel and llvmpipe, nearly always by one.
- Windows: the socket `ioctl` is plain Go now. The raw mmap and memguard files are `unix` rather than linux-only, and `netstr` has a BSD socket path, so macOS and linux/arm64 build.
- CI: a macOS job, a job that cross-builds all five targets, and an AppImage that builds SDL3 from source on Debian 12 (which has no SDL3 package) and bundles it with OpenAL.

## v1.10.0-alpha2

- Multiplayer without port forwarding: the lobby helps players behind routers connect directly, and relays the game when that fails (`network.nat_punch`).
- Fix: Clients got the previous or a broken map when the server's map had no `.nxz` file; the server now compresses it.
- Fix: Clients downloaded the map again on every join (#639); interrupted downloads no longer leave a broken map.
- Map names can be up to 28 characters long.
- Fix: Multiplayer maps starting with "con", "war" or "wiz" were missing from the map list (#651).
- Large map view: the map key cycles small, large and closed, with an opacity option (#716).
- Option to turn off auto-equipping picked-up items (`game.auto_equip`, #741).
- Options menu: OpenNox options on a second page; sliders for UI Scale, Frame Rate Limit and Map Opacity.
- Server settings: friendly fire for spells (`set friendlyfire`, #734), observer mode for idle players (`set afkobserver`, #691), the highest spell book level (`set spelllevel`, #731).
- `set mode` changes the game mode of the next map (#674).
- Spell book pages show the phonemes (#565).
- The mouse wheel over a shop's stock scrolls the stock (#675).
- A clear message when joining a server that runs another version (#419).
- NS4: lives (#669), hiding objects (#724), object lifetime (#619), playing WAV files from the map folder (#620), NPC team chat (#645), item modifiers and colors (#647), trigger activators (#652), player and monster actions (#723), magic walls (#725), quest stage across maps (#727).
- NS4 fixes: the killer in death events (#678), scripted Magic Missile targets (#689), `SetBaseSpeed()` on NPCs after loading (#730), Go scripts in map folders with capital letters.
- Dedicated servers: clear messages for client-only commands, and fixes for crashes in several console commands.
- Fix: A joined player could rename objects on the host.

## v1.10.0-alpha1

- The engine is written entirely in Go and runs as a native 64-bit program (amd64 Linux, x86_64 Windows); 32-bit, Snap and Windows 7 builds were dropped.
- XWIS is no longer supported; the default lobby is `https://nox.thunraz.se`.
- Older OpenNox versions can't play with this one; the software renderer was removed (OpenGL 3.3 is required).
- UI scaling, camera zoom, smooth fonts, vsync, an FPS limit and counter, hardware cursors, and graphics options in the options menu.
- Experimental gamepad support, the original random map generator, map rendering tools and many new NS4 functions.
- Many multiplayer, gameplay, sound and rendering fixes; see the GitHub release for the full list.

## v1.9.0 alphas

See the GitHub releases.

## v1.8.11

- Preparations to make more spells moddable.
- Hidden Death spell now correctly draws spinning skull on monsters.
- Allow using two more effects from scripts: `ENERGY_BOLT` and `PLASMA`.
- Lots of internal changes related to color handling.
- Allow force-switching maps in more cases.
- Stretch image option is now persisted.
- Fix: Potential crash when loading certain maps.
- Fix: Downloading maps via Nox protocol.
- Fix: Incorrect rendering of the observer effect.
- Fix: Rendering of abilities panel when player dies.
- Fix: Obliteration spell effect.
- Fix: Minimap rendering.
- Fix: Parsing of network filters from the config.
- Fix: Transparent color not rendering correctly in some GUI elements.
- Fix: Correctly override background color in GUI.
- Fix: Crash when reading gamma on the server.
- Fix: Server requiring audio/video libs.
- Fix: Text message when switching maps.
- Fix: `GREEN_BOLT` effect when called from a script.
- Fix: Build now doesn't require Git.

## v1.8.10

- First release as an open source project!
- Significant cleanup of the codebase.
- Fix long-standing bug for rendering lights in high-res.
- Long-awaited support for movie playback.
- Experimental support loading campaign in multiplayer.
- Reworked most of the rendering code.
- Allow console commands in Quest.
- Use shaders for gamma control. Should work more reliably.
- Send server resolution and player list to the lobby.
- Allow overriding and initial modding support for spells.
- Allow getting all spells, including hidden ones via `cheat spells all`.
- Allow removing summon limit via `cheat summon.nolimit`.
- Support XDG paths on Linux.
- Fix: Correctly handle `goto` console command.
- Fix: Minor physics deviation compared to vanilla.
- Fix: Teams creation in Quest.
- Fix: Reading gamma settings from config.

## v1.8.9

- OpenNox will now use the new lobby server and fallback to XWIS if no servers were found.
- Ping will now be shown in servers list for servers that support it.
- Server discovery was optimized and is now faster.
- Implemented a tiny builtin SSH remote console. Telnet support removed.
- Game load/save was optimized and works faster now.
- Fix: Not being able to connect to non-OpenNox servers.
- Fix: Spawned creatures being stuck when damaged.
- Fix: Magic missiles targeting dead entities.
- Fix: Energy bolt not working properly in traps.
- Fix: UI no longer freezes when discovering game servers.
- Fix: Player names not being shown in scores table.
- Fix: Sorting in the character list.

## v1.8.8

- Notify if OpenNox update is available.
- Server: add metrics for monitoring.
- Fix: Jumping no longer prevents from taking damage from lava.
- Fix: Do not announce Solo Quest games online.
- Fix: Team colors (regression from v1.8.7).
- Fix: Correctly set config dir for Snap.

## v1.8.7

- Support loading of TrueType and OpenType fonts.
- Automatically fix Russian translation and font files.
- A ton of internal changes and cleanup.
- Fix: Creatures and NPCs with ranging attack trying to fight melee only.
- Fix: Server API being stuck than there are no players on the server.
- Fix: Crash when falling into a pit.
- Fix: AI debug logging.

## v1.8.6

- Release OpenNox for Linux on [Snapcraft](https://snapcraft.io/opennox).
- Reading game files is now safer and more extendable.
- Preparations for saving and replaying Nox matches.
- Preparations for rendering TTF fonts in game engine.
- E2E mode now allows recording player input directly (to reproduce bugs).
- Safer and faster sprite rendering.
- Initial work on rendering map previews.
- A lot of internal refactorings and improvements.
- Allow configuring manual spell cast timeout.
- Allow binding any keys with `bind` command.
- New `cheat spawn` command to spawn items and monsters.
- Fix: Safer GUI code, preventing a major memory corruption.
- Fix: Races in the audio playback.
- Fix: Potential fix for native map downloads.
- Fix: Crashes when playing on broken maps.

## v1.8.5

- Add more resolution options for both regular and HD version.
- Move screen resolution to `opennox.yml` config.
- Fix: Instant victory if player commits a suicide in Arena.
- Fix: Properly propagate closed and private game flags to XWIS.
- Fix: It was possible to join closed games.
- Fix: Map filtering by game type.
- Fix: `cheat health` and `cheat mana` now work properly without arguments.

## v1.8.4

- Better performance for particle effects.
- Allow using emotes in campaign.
- New `cheat sage` command that gives all spells and scrolls, but doesn't make you invincible.
- New `cheat spells` command that only gives all spells for your class.
- New `cheat scrolls` command that only gives all best scrolls.
- New `cheat god` (as opposed to `set god`) that only makes you invincible, but won't give any spells.
- Fix: Set time limit for new key bindings to prevent the same action from executing to fast.
- Fix: Lightning spells were not hitting multiple targets as they should.
- Fix: Controlled creatures disappearing on map switch.
- Fix: Effects from Protection spells were played incorrectly.
- Fix: Random "you cannot wear this" messages when switching equipment.
- Fix: Allow really large amount of gold in `cheat gold`.
- Fix: Announcing Solo games to online lobby.
- Fix: Using relative paths in config.
- Fix: Better protection from potential memory leaks.
- Fix: Crash when reading maps that use the new Panic's script compiler.
- Fix: Multiple other crashes.

## v1.8.3

- Integrated the new (simple) server control panel. 
- Fix: Saving complex input bindings to the config.
- Fix: Wrong interpretation of XWIS option in the new config.
- Fix: Possible crashes for maps using new Panic's compiler.

## v1.8.2

- Write `opennox.yml` config with all extended options.
- Allow changing game data dir via the `opennox.yml` config.
- Fix: Player names being stuck as online in dedicated server.
- Fix: Mouse movement on low sensitivity settings.
- Fix: Assigning keyboard keys instead of mouse for movement.

## v1.8.1

- Fix: Back walls rendering in high-resolution mode.
- Fix: Center in-game menu in high-resolution mode.
- Fix: Possible crash after running for some time.

## v1.8.0

### New

- High-resolution support (up to 4K)! You'll need to run `opennox-hd.exe` for it.
- A lot faster OpenGL-based rendering (partially offloaded to GPU).
- Add `cheat charm.all` to charm any creature (including humanoids).
- Automated end-to-end testing mode, see [opennox-e2e](https://github.com/opennox/opennox-e2e) for details.
- Simple HTTP-based [Server API](./docs/server-api.md) for controlling game servers.
- Allow disabling image smoothing in the video options.
- Allow enabling image stretching in the video options.
- Other minor optimizations.

### Fixes

- Bug with creatures being "deaf" (not reacting to footsteps, etc).
- Properly show the last character in long server names.
- Error when manually saving player in online games.
- Banish of captured creatures (like vampire bats) leading to scripts being stuck.
- Properly forward ports when server port is changed via a flag.
- Disabling soft shadows via `-soft` flag.
- LUA scripts when reloading the same map from the console.
- Nox Reloaded failing to join OpenNox Quest games.
- OpenNox changing `user.rul` file permissions.
- Potential crash on hosting Quest games.
- Potential crash when leaving the game.
- Crash when opening advanced server setting without admin rights.

## v1.7.1

- Fix: Hide OS mouse cursor in windowed mode as well.
- Fix: Rename "Back" button in input config menu to "Accept".
- Fix: Prevent crash on audio errors.
- Fix: Avoid collision of new health and mana cheats with old ones.

## v1.7.0

- Completely reworked input and present pipelines.
- Manual spell casting is now supported. Bind keys to individual spell gestures and become a real wizard/conjurer!
- Gamma and mouse sensitivity sliders were added to the options menu.
- Wide-screen 16:9 resolutions are now available. Older 4:3 resolutions are still supported via config.
- Fix: Input sensitivity should now work correctly in fullscreen mode.
- Fix: Gamma setting in the config is now respected.
- Fix: Stretched video setting in the config is now respected.
- Added key bindings for:
  - Switching to a specific spell slot.
  - Switching to a specific trap slot.
  - Dropping a trap from your inventory.
  - Controlling spawned creature behavior.
  - Accepting item buy/sell/drop.
- Breaking change: `set allow.all` was renamed to `cheat equip.all`
- Breaking change: `set mana` was renamed to `cheat mana`
- Breaking change: `set health` was renamed to `cheat health`

## v1.6.1

- Fix: Prevent freeze when closing the game on Windows. 
- Fix: Prevent crash in `set mana`.
- Fix: Prevent crash when switching shield and bow with `set allow.all` enabled.
- Fix: Ignore item strength requirement when using `allow.all`.
- Fix: Remove restrictions from `cheat gold`, so it's accepted in Quest games.

## v1.6.0

- Rename project to OpenNox to avoid confusion.
- Automatically set Nox data folder. Supports Origin, GoG and Nox Reloaded.
- Keep 4:3 aspect ratio when resizing game window.
- Toggle fullscreen mode with Alt+Enter.
- Add `set health` and `set mana` cheats.
- Add `set allow.all` cheat to remove item class restrictions.
- Allow customizing `gamedata.bin` values without encoding it. See [gamedata.yml](./docs/gamedata-sample.yml).
- LUA now supports trigger events as well as player join/leave events.
- Fix: Detonate traps spell will no longer freeze the game.
- Fix: Fullscreen mode will now be saved correctly.
- Fix: Warriors should now be able to scroll weapons with mouse wheel in multiplayer as well.

## v1.5.0

- Automatically open/forward ports when hosting a game (no need to configure firewall).
- Fixed the bug that prevented warriors from scrolling weapons with mouse wheel.
- Blazing-fast map downloads. Both client and the server should use this version.
- Experimental support for [LUA map scripts](./docs/maps-lua.md).
- Maps will now be transferred with other related files when possible (`.txt`, `.rul`, `.lua`, `.png`).

## v1.4.0

- First public release!
- Compatible with Nox Reloaded.
- Native SDL and OpenAL support (no workarounds required for Win10, streaming, etc).
- Builtin XWIS integration. No need for account in order to host or join games.
- More [Nox Quest](./docs/game-quest.md) options (skip levels, keep portal forever, etc).
- Native Linux build (no Wine).
- Experimental gamepad support.
- Dedicated "headless" server.
