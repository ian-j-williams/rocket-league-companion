# Rocket League Companion

A lightweight self-contained Go desktop companion for Rocket League match stats.

## Features

- Local Rocket League Stats API connection to port 49123 by default
- Guided configuration for PacketSendRate tuning
- Live match snapshot panel
- Per-session win/loss tracker
- Reset button for the current session state
- Single executable build workflow

## Configuration notes

The Rocket League Stats API must be enabled before the game starts. Update the Rocket League config before launching the client.

Typical config values:

- PacketSendRate = 30
- Port = 49123
- WebPort = 49124

The game must be restarted after the config is changed.

## Run

```bash
go run ./cmd/companion
```

## Build

```bash
go build ./cmd/companion
```
