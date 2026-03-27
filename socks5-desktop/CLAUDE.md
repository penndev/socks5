# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a SOCKS5 proxy desktop application built with Wails3 (Go backend + Vue 3 frontend). The application provides a system tray interface for managing local SOCKS5 proxy servers that can forward traffic to remote SOCKS5 servers (including SOCKS5-over-TLS).

## Development Commands

### Running the application
```bash
# Development mode with hot-reload
wails3 dev

# Or using Task
task dev
```

### Building
```bash
# Build for current platform
task build

# Package for distribution
task package

# Run the built binary
task run
```

### Frontend development
```bash
cd frontend

# Install dependencies
npm install

# Development server (standalone)
npm run dev

# Build frontend
npm run build

# Build for development (unminified)
npm run build:dev
```

## Architecture

### Backend (Go)

**Core Services** (registered in main.go):
- `Storage`: BoltDB-based key-value storage in user config directory (`~/.config/Socks5Desktop/data.db`)
- `App`: Application constants and configuration
- `Proxy`: SOCKS5 proxy management (local server + remote forwarding)

**Key Components**:
- `main.go`: Application entry point, window management, system tray setup
- `proxy.go`: SOCKS5 proxy logic with support for:
  - Local SOCKS5 server (TCP/UDP)
  - Remote forwarding via `socks5://` or `socks5overtls://` schemes
  - Connection logging via Wails events
- `storage.go`: Persistent storage using BoltDB
- `app.go`: Application constants for event names

**Dependencies**:
- Uses local packages via `replace` directives:
  - `github.com/penndev/gopkg` → `../../gopkg`
  - `github.com/penndev/socks5/core` → `../core`
- These must be available at the specified relative paths

**Event System**:
- `logServerStatus`: Server status updates
- `logProxyList`: Connection logs

### Frontend (Vue 3)

**Tech Stack**:
- Vue 3 with Composition API
- Ant Design Vue for UI components
- Pinia for state management
- Vite for build tooling
- Wails runtime for Go ↔ JS communication

**Structure**:
- `frontend/src/App.vue`: Main application component
- `frontend/src/components/`: UI components
  - `ActionPanel.vue`: Action controls
  - `ServePanel.vue`: Server management
  - `SettingPanel.vue`: Settings interface
  - `BottomBar.vue`: Bottom status bar
  - `bottombar/LogPanel.vue`: Connection logs
  - `bottombar/StatusPanel.vue`: Status display
- `frontend/src/stores/`: Pinia stores
  - `server.js`: Server state management
  - `settings.js`: Settings state management

**Wails Bindings**:
- Auto-generated TypeScript bindings in `frontend/bindings/`
- Call Go methods via `@wailsio/runtime`

## Important Notes

### Local Dependencies
The project depends on two local packages that must exist:
- `../../gopkg` (relative to project root)
- `../core` (relative to project root)

If these are missing, Go commands will fail. Check that the parent directory structure matches the expected layout.

### Platform-Specific Builds
The project includes platform-specific build configurations in `build/`:
- `build/darwin/Taskfile.yml`: macOS builds
- `build/windows/Taskfile.yml`: Windows builds
- `build/linux/Taskfile.yml`: Linux builds

### System Tray Behavior
- Window hides to system tray on close (doesn't quit)
- Click tray icon to toggle window visibility
- Window position is preserved between hide/show
- macOS: App doesn't terminate when last window closes

### Storage Location
Application data is stored in:
- Windows: `%APPDATA%\Socks5Desktop\data.db`
- macOS: `~/Library/Application Support/Socks5Desktop/data.db`
- Linux: `~/.config/Socks5Desktop/data.db`
