#!/bin/sh
# Two-instance multiplayer test: a host (e2e/mp-host.yaml) and a client (e2e/mp-client.yaml) on loopback.
# usage: scripts/e2e-multiplayer.sh <host binary> <Nox data dir> [work dir] [client binary]
# The client uses the host binary unless another one is given, e.g. a Windows build to test Linux and Windows
# players together. A binary ending in .exe runs under Wine (set WINEPREFIX to use a separate prefix).
# It checks the logs (the client joins and leaves, nothing panics); the screens of both are kept in the work
# directory for review. Uses game ports 18700 (host) and 18720 (client).
set -u
hostbin="$(realpath "$1")"
data="$(realpath "$2")"
root="$(cd "$(dirname "$0")/.." && pwd)"
work="${3:-$(mktemp -d)}"
clientbin="$(realpath "${4:-$1}")"
mkdir -p "$work"
work="$(realpath "$work")"
echo "work dir: $work"

for side in host client; do
	rm -rf "$work/$side"
	mkdir -p "$work/$side/testdata" "$work/$side/home" "$work/$side/data"
	cp "$root/e2e/mp-$side.yaml" "$work/$side/"
	# MP_HOST_YAML / MP_CLIENT_YAML: run other scenarios (e.g. to reproduce a bug), with the same setup
	if [ "$side" = host ] && [ -n "${MP_HOST_YAML:-}" ]; then cp "$MP_HOST_YAML" "$work/$side/mp-$side.yaml"; fi
	if [ "$side" = client ] && [ -n "${MP_CLIENT_YAML:-}" ]; then cp "$MP_CLIENT_YAML" "$work/$side/mp-$side.yaml"; fi
	# e2e mode keeps saves and characters in the data dir: link the game files, keep Save separate
	for f in "$data"/*; do
		[ "$(basename "$f")" = Save ] || ln -s "$f" "$work/$side/data/"
	done
done
# MP_DOWNLOAD=1: the client doesn't have the host's map (BluDeath), so it downloads it from the host
# (opennox#673); the downloaded map is written into the client's own maps dir, not into the game data
if [ -n "${MP_DOWNLOAD:-}" ]; then
	rm "$work/client/data/maps"
	mkdir "$work/client/data/maps"
	for f in "$data"/maps/*; do
		case "$(basename "$f" | tr '[:upper:]' '[:lower:]')" in bludeath) ;; *) ln -s "$f" "$work/client/data/maps/" ;; esac
	done
fi

run() { # side binary port [env...]
	side=$1 bin=$2 port=$3
	shift 3
	# a Windows build runs under Wine, which sees the file system as drive Z:
	wine="" z=""
	case "$bin" in *.exe) wine=wine z="Z:" ;; esac
	echo "$side: $bin"
	(cd "$work/$side" && env HOME="$work/$side/home" XDG_CONFIG_HOME="$work/$side/home/.config" \
		XDG_DATA_HOME="$work/$side/home/.local/share" XDG_CACHE_HOME="$work/$side/home/.cache" \
		XDG_STATE_HOME="$work/$side/home/.local/state" NOX_E2E="$z$work/$side/mp-$side.yaml" NOX_E2E_FAILFAST=false \
		NOX_E2E_OVERRIDE=true WINEDEBUG=-all "$@" timeout 300 $wine "$bin" -port "$port" -data "$z$work/$side/data" \
		>"$work/$side.log" 2>&1)
}

run host "$hostbin" 18700 &
host=$!
i=0
while ! grep -q 'loading map "bludeath.map"' "$work/host.log" 2>/dev/null; do
	i=$((i + 1))
	if [ $i -gt 120 ] || ! kill -0 $host 2>/dev/null; then
		echo "FAIL: the host did not start the arena game"
		exit 1
	fi
	sleep 1
done
run client "$clientbin" 18720 NOX_E2E_SERVERS=127.0.0.1:18700
client_rc=$?
wait $host
host_rc=$?

fail=0
check() { # description file pattern
	if grep -q "$3" "$2"; then echo "ok: $1"; else echo "FAIL: $1"; fail=1; fi
}
[ $host_rc -eq 0 ] && echo "ok: host exited" || { echo "FAIL: host exit code $host_rc"; fail=1; }
[ $client_rc -eq 0 ] && echo "ok: client exited" || { echo "FAIL: client exit code $client_rc"; fail=1; }
if [ -n "${MP_DOWNLOAD:-}" ]; then
	check "client downloaded the map" "$work/client.log" 'fade after map download'
	[ -f "$work/client/data/maps/bludeath/bludeath.map" ] || [ -f "$work/client/data/maps/BluDeath/BluDeath.map" ] &&
		echo "ok: the map file is there" || { echo "FAIL: no downloaded map file"; fail=1; }
fi
check "client found the server" "$work/client.log" 'server response info={"127.0.0.1:18700"'
check "client joined (host side)" "$work/host.log" 'has joined the game'
check "client in the game" "$work/client.log" 'You are in Observation Mode'
check "client saw the host player" "$work/client.log" 'Jack has joined the game'
check "client left (host side)" "$work/host.log" 'has left the game'
# the fight: the host is "Jack", the client "Jack 2" (default names); a hit shows as lost health, a kill as
# "Victim: <name> , Attacker: <name>" in the console log
after="$(sed -n '/STATE: Health after/,$p' "$work/host.log")"
hurt() { # name: the player's health after the fight is below the maximum
	echo "$after" | grep -oE "PLAYER \"$1\" health [0-9]+/[0-9]+" | grep -oE '[0-9]+/[0-9]+' |
		awk -F/ '$1 < $2 { found = 1 } END { exit !found }'
}
if hurt "Jack 2" || grep -q 'Victim: Jack 2 , Attacker: Jack$' "$work/host.log"; then
	echo "ok: the host hit the client"
else
	echo "FAIL: the host did not hit the client"; fail=1
fi
if hurt "Jack" || grep -q 'Victim: Jack , Attacker: Jack 2' "$work/host.log"; then
	echo "ok: the client hit the host"
else
	echo "FAIL: the client did not hit the host"; fail=1
fi
# kills are announced to every player: if the host saw one, the client must have seen it too
if grep -q 'Victim: ' "$work/host.log"; then
	check "client saw the kill" "$work/client.log" 'Victim: '
fi
for side in host client; do
	if grep -q '^panic\|fatal error' "$work/$side.log"; then echo "FAIL: $side panicked"; fail=1; fi
done
[ $fail -eq 0 ] && echo "PASS" || echo "FAILED (logs and screens in $work)"
exit $fail
