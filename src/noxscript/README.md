# NoxScript

Script APIs for Nox maps, used by the OpenNox script runtimes:

- [ns/asm](./ns/asm) and [ns/nsvm](./ns/nsvm): the original NoxScript bytecode and builtins.
- [ns/v3](./ns/v3): NoxScript 3, the functions used by the original Nox maps, and a decompiler (`ns/v3/cmd/noxscript-decomp`).
- [ns/v4](./ns/v4): NoxScript 4, the object-oriented API for OpenNox maps.
- [eud/v171](./eud/v171): a compatibility layer for Panic's EUD map scripts.

Originally [github.com/opennox/noxscript](https://github.com/opennox/noxscript), now maintained here as part of OpenNox.

Map scripts still import these packages by their original paths (e.g. `github.com/opennox/noxscript/ns/v4`),
which are the keys in the `Symbols` maps of the `eval` packages. Keep those keys when changing the extracted files.
