# thg

A fast command-line interface for [Things 3](https://culturedcode.com/things/)
on macOS. List, search, add, and complete tasks without leaving the terminal.

`thg` reads directly from the Things 3 database (read-only) and writes through
the official `things:///` URL scheme — so reads are instant and writes go
through Things itself, keeping your data safe.

```console
$ thg
ID      Title                     Project        When     Deadline
a1b2c3  Review pull request       thg            Today
d4e5f6  Buy milk                                 Today
...

$ thg add "Email the contractor" --when today --tags errand
Added: Email the contractor

$ thg done a1b2c3
Completed: Review pull request
```

## Features

- **Read instantly** — queries the local Things SQLite database directly, in
  read-only mode.
- **Write safely** — adds and updates go through the Things URL scheme, never
  by writing to the database.
- **Familiar views** — Today, Inbox, Someday, Upcoming, plus filtering by
  project, area, or tag.
- **Scriptable** — `--json` on read commands for piping into other tools.

## Requirements

- **macOS** with [Things 3](https://culturedcode.com/things/) installed.
- **[Go](https://go.dev/dl/) 1.25 or newer** (only needed to build from source).
- An **auth token** for commands that modify existing tasks (`done`, `update`).
  Found in **Things → Settings → General → Enable Things URLs → Manage**.

## Installation

### Build from source (recommended)

```bash
git clone https://github.com/andrew-craig/thg.git
cd thg
./install.sh
```

`install.sh` builds the binary and installs it to `/usr/local/bin` (override
with `PREFIX=~/.local ./install.sh` to install to `~/.local/bin`).

To build without installing:

```bash
go build -o thg .
./thg --help
```

### With `go install`

```bash
go install github.com/andrewstuart/thg@latest
```

This installs `thg` into `$(go env GOPATH)/bin` — make sure that directory is
on your `PATH`.

## Configuration

Both settings are optional.

| Setting | How to set | Purpose |
|---------|-----------|---------|
| Auth token | `--auth-token`, `THG_AUTH_TOKEN`, or `~/.config/thg/config.json` | Required by `done` and `update` |
| Database path | `THG_DB_PATH` | Override automatic database discovery |

The auth token is resolved in this order: the `--auth-token` flag, then the
`THG_AUTH_TOKEN` environment variable, then a config file:

```bash
mkdir -p ~/.config/thg
echo '{"auth_token":"YOUR_TOKEN"}' > ~/.config/thg/config.json
```

`thg` finds the Things database automatically by scanning the standard group
container. If you keep it elsewhere, set `THG_DB_PATH` to the full path of
`main.sqlite`.

## Usage

Running `thg` with no arguments shows your Today list.

### Listing and viewing

```bash
thg                       # Today (default)
thg list inbox            # Inbox
thg list someday          # Someday
thg list upcoming         # Scheduled / upcoming
thg list "Project Name"   # To-dos in a project
thg list --area Work      # Filter by area
thg list --tag errand     # Filter by tag
thg list --all            # All open to-dos

thg show a1b2c3           # Task details (by short ID or title)
thg projects              # Open projects (--area, --all)
thg areas                 # Areas
thg tags                  # Tags

# Machine-readable output
thg list --json
```

### Adding tasks

```bash
thg add "Buy milk"
thg add "Prep slides" --when today --deadline 2026-07-01 --list Work
thg add "Plan trip" --tags travel,personal --notes "Check passport"
thg add "Groceries" --checklist "Milk,Eggs,Bread"
```

`add` flags: `--when` (today/tomorrow/evening/anytime/someday/`YYYY-MM-DD`),
`--deadline`, `--list` (project or area), `--heading`, `--tags`, `--notes`,
`--checklist`, `--reveal`.

### Completing and updating

These commands modify existing tasks and require an auth token (see
[Configuration](#configuration)).

```bash
thg done a1b2c3
thg update a1b2c3 --title "New title" --when tomorrow
thg update a1b2c3 --append-notes "Follow-up needed"
thg update a1b2c3 --canceled
```

`update` flags: `--title`, `--when`, `--deadline`, `--tags`, `--notes`,
`--append-notes`, `--list`, `--canceled`.

Use `thg <command> --help` for the full list of flags on any command.

## How it works

- **Reads** query the Things 3 SQLite database directly over a read-only
  connection. See [docs/database.md](docs/database.md).
- **Writes** are sent through the `things:///` URL scheme via `open`, so Things
  applies every change itself. See [docs/url-scheme.md](docs/url-scheme.md).

`thg` never writes to the database. More detail on the architecture lives in
[CLAUDE.md](CLAUDE.md).

## Documentation

- [docs/database.md](docs/database.md) — Things 3 database schema and queries
- [docs/url-scheme.md](docs/url-scheme.md) — full `things:///` URL-scheme reference
- [CLAUDE.md](CLAUDE.md) — architecture and contributor notes

## License

[MIT](LICENSE) © Andrew Craig

`thg` is an independent project and is not affiliated with or endorsed by
Cultured Code, the makers of Things.
