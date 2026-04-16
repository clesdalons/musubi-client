# README

# Musubi Client

Musubi Client is a desktop synchronization tool for Divinity: Original Sin 2 save games. It watches the local
`Story` save directory, uploads new save folders automatically to an Azure-backed API, and can fetch the latest
cloud save on demand.

## Architecture

- `main.go` — lightweight Wails bootstrap and frontend asset configuration.
- `internal/application` — application service exposing Wails bindings and orchestrating watcher, sync, and config.
- `internal/config` — persistent settings management and automatic DOS2 save path detection.
- `internal/watcher` — filesystem watcher monitoring new save folders and emitting upload events.
- `internal/cloud` — Azure API synchronization, including upload and cloud download.
- `internal/storage` — ZIP packaging and extraction utilities.
- `frontend` — React + Vite UI rendering the configuration and sync status.

## Development

1. Start the frontend development server:

```bash
cd frontend
npm install
npm run dev
```

2. Start the Wails desktop application:

```bash
cd ..
wails dev
```

## Build

To compile the application for production:

```bash
wails build
```

## Notes

- Uses a local JSON config file at `~/.musubi/config.json`.
- Automatically detects common DOS2 save directories in `%USERPROFILE%\Documents` and OneDrive.
- The app is now structured with clear internal packages for a portfolio-grade Go project.
