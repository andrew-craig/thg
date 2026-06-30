# thg — Contributor Guide

Command-line interface for reading and writing tasks in Things 3 (macOS).

This file is guidance for anyone (humans or AI agents) working on the `thg`
codebase. For installation and usage, see [README.md](README.md). For the
detailed Things 3 references, see the [docs](docs/) folder.

## Architecture

- **Read**: Query the Things 3 SQLite database directly, always via a read-only
  connection. See [docs/database.md](docs/database.md) for the schema, the date
  encoding, and the queries used.
- **Write**: Use the `things:///` URL scheme (invoked with `open`) to add and
  update tasks — never write to the database. See
  [docs/url-scheme.md](docs/url-scheme.md) for the full scheme.

This split is deliberate: Things owns its database and must handle its own
writes. Writing to the live database directly risks corrupting it.

## Layout

| Path | Purpose |
|------|---------|
| `main.go` | Entry point; calls `cmd.Execute()` |
| `cmd/` | Cobra commands (one file per command) |
| `db/` | Read-only SQLite access, models, date/tag/checklist helpers |
| `things/` | URL-scheme builder and `open` invocation (writes) |
| `config/` | Auth-token and config-file loading |
| `format/` | Table / JSON output printer |
| `docs/` | Things 3 database and URL-scheme references |

## Commands

`thg` (default: today), `list`, `show`, `add`, `done`, `update`, `areas`,
`projects`, `tags`. Each maps to a file in `cmd/`. The default command and
`list` read the database; `add`, `done`, and `update` write via the URL scheme.

## Conventions

- The database is opened read-only (`?mode=ro`). Never change this.
- `THG_DB_PATH` overrides database discovery; the standard path's
  `ThingsData-XXXXX` suffix varies per install, so we scan for it.
- Write commands that modify existing tasks (`done`, `update`) require an
  auth-token, resolved by `config.LoadAuthToken` from flag → `THG_AUTH_TOKEN`
  → `~/.config/thg/config.json`.
- Most read commands accept a `--json` flag for machine-readable output.

## Development Notes

- A snapshot database can be placed at `things.sqlite` for schema exploration;
  it is git-ignored and should not be committed.
- Read-only access to the live database works while Things is running (SQLite
  WAL mode). For a clean *copy*, Things should not be writing during the copy.

## Issue Tracking

This project uses **bd (beads)** for issue tracking — see
[AGENTS.md](AGENTS.md). Do not introduce parallel markdown TODO lists.
