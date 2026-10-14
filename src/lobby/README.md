# Nox lobby server

This project provides a Nox game lobby which exposes a simple HTTP API for both listing and registering Nox game servers.

The main use case for the lobby is to support OpenNox, but the API can also be used for bots
that want to notify players about currently active Nox games.

## Public lobby

To get a list of games on the public lobby:

```bash
curl 'https://nox.thunraz.se/api/v0/games/list'
```

This package is also the Go client for the HTTP API; OpenNox uses it to list and register games.

## Running locally

The simplest way to run the lobby server locally is via Docker. Build the image from the `src` directory:

```bash
docker build -f lobby/Dockerfile -t nox-lobby .
docker run -d --restart always --name nox-lobby -p 8080:80 nox-lobby
```

To get a list of games via local lobby:

```bash
curl 'http://127.0.0.1:8080/api/v0/games/list'
```

## Origin

Originally [github.com/opennox/lobby](https://github.com/opennox/lobby), now maintained here as part of OpenNox.
