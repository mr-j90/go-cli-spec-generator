# go-cli-spec-generator WIP

An interactive CLI tool for generating software specification documents through a guided terminal wizard. Built with Go, Cobra, and Bubble Tea.

## Overview

`specgen` walks you through a multi-step TUI wizard to collect project details — your CLI profile type, optional feature areas, and targeted questions — then exports the resulting specification as a PDF, DOCX, or Markdown document.

## Features

- **Interactive TUI wizard** — 5-step flow powered by Bubble Tea
- **CLI profile selection** — oneshot, daemon, subcommand, or hybrid
- **Feature area selection** — 12 optional areas (auth, storage, API, testing, observability, deployment, security, caching, messaging, search, notifications, configuration)
- **Adaptive question bank** — 27 questions scoped to your chosen profile and features
- **Session persistence** — save and resume interrupted sessions
- **Multiple export formats** — PDF (via GoPDF), DOCX (via UniOffice), Markdown

## Requirements

- Go 1.26.1 or later

## Installation

```bash
git clone https://github.com/zyx-holdings/go-spec.git
cd go-spec
go build -o specgen .
```

## Usage

### Launch the interactive wizard

```bash
specgen new
```

Steps:
1. **Profile** — select your CLI's execution model
2. **Features** — toggle optional feature areas (Space to select, Enter to confirm)
3. **Questions** — answer targeted questions for each selected area
4. **Review** — inspect all answers before export
5. **Export** — choose output format (PDF / DOCX / Markdown)

### Generate a spec from a file (planned)

```bash
specgen generate --input answers.json --output spec.pdf --format pdf
```

### Resume a previous session (planned)

```bash
specgen resume
```

### Show version

```bash
specgen version
```

### Global flags

| Flag | Description |
|------|-------------|
| `--no-color` | Disable color output |

## TUI Keyboard Shortcuts

| Key | Action |
|-----|--------|
| Arrow keys | Navigate options |
| Space | Toggle selection (multi-select) |
| Enter | Confirm / advance |
| Tab | Next question |
| Shift+Tab | Previous question |
| Alt+Enter | Insert newline (textarea inputs) |
| Esc / Backspace | Go back |
| Ctrl+C | Quit |

## Project Structure

```
.
├── cmd/                    # Cobra command implementations
│   ├── root.go             # Root command, version, global flags
│   ├── new.go              # `specgen new` — launch TUI wizard
│   ├── generate.go         # `specgen generate` — file-based generation
│   ├── resume.go           # `specgen resume` — restore saved session
│   └── version.go          # `specgen version`
├── internal/
│   ├── cli/                # CLI utilities
│   ├── export/             # PDF and DOCX export logic
│   ├── questions/          # Question registry, profiles, feature areas
│   ├── render/             # Document rendering
│   ├── session/            # Session state and JSON persistence
│   └── tui/                # Bubble Tea TUI (steps, widgets, styles)
└── main.go                 # Entry point
```

## CLI Profiles

| Profile | Description |
|---------|-------------|
| `oneshot` | Runs once, performs a task, then exits (e.g. `ls`, `curl`) |
| `daemon` | Long-running background service |
| `subcommand` | Multiple subcommands (e.g. `git`, `docker`) |
| `hybrid` | Supports both one-shot and daemon execution modes |

## Feature Areas

| Area | Description |
|------|-------------|
| authentication | Auth mechanisms and user identity |
| storage | Databases and file storage |
| api | HTTP/gRPC API layer |
| testing | Testing strategy and tooling |
| observability | Logging, metrics, tracing |
| deployment | Packaging and deployment targets |
| security | Security controls and hardening |
| caching | In-memory or distributed caching |
| messaging | Event queues and pub/sub |
| search | Full-text or vector search |
| notifications | Alerts and notification delivery |
| configuration | Config loading and management |

## Session Persistence

Sessions are saved as JSON and can be resumed if interrupted. Each session tracks:

- A UUID session ID and timestamps
- Selected profile and feature areas
- All question answers (string or multi-value)

## Dependencies

| Package | Version | Purpose |
|---------|---------|---------|
| [cobra](https://github.com/spf13/cobra) | v1.10.2 | CLI framework |
| [bubbletea](https://github.com/charmbracelet/bubbletea) | v1.3.10 | TUI framework |
| [bubbles](https://github.com/charmbracelet/bubbles) | v1.0.0 | TUI components |
| [lipgloss](https://github.com/charmbracelet/lipgloss) | v1.1.0 | Terminal styling |
| [gopdf](https://github.com/signintech/gopdf) | v0.36.0 | PDF generation |
| [unioffice](https://github.com/unidoc/unioffice) | v1.39.0 | DOCX generation |

## Development

### Run tests

```bash
go test ./...
```

### Run with live reload

```bash
go run . new
```

### Build binary

```bash
go build -o specgen .
```

## Current Status

| Feature | Status |
|---------|--------|
| TUI wizard (steps 1–5) | Implemented |
| Question registry (27 questions) | Implemented |
| Session persistence | Implemented |
| PDF export | Partial (skeleton) |
| DOCX export | Partial (basic) |
| Markdown export | Planned |
| `generate` command | Planned |
| `resume` command | Planned |

## Module

```
github.com/zyx-holdings/go-spec
```
