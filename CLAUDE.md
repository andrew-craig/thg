# thg - Things 3 CLI

Command-line interface for reading and writing tasks in Things 3 (macOS).

## Architecture

- **Read**: Query the Things 3 SQLite database directly (read-only copy)
- **Write**: Use the `things:///` URL scheme via `open` to add/update tasks (never write to the DB)

## Things 3 Database

### Location

```
~/Library/Group Containers/JLMPQHK86H.com.culturedcode.ThingsMac/ThingsData-HK9N6/Things Database.thingsdatabase/main.sqlite
```

The `.thingsdatabase` is a macOS bundle (directory). The database is `main.sqlite` inside it.
WAL files (`main.sqlite-shm`, `main.sqlite-wal`) may also be present.

**Important**: Always open the database in read-only mode (`?mode=ro` in SQLite URI or `SQLITE_OPEN_READONLY`) to avoid any risk to the live database. Things must handle its own writes.

### Key Tables

#### TMTask

The main table. Holds to-dos, projects, and headings in a single table distinguished by `type`.

| Column | Type | Description |
|--------|------|-------------|
| `uuid` | TEXT PK | Unique identifier (used in URL scheme as `id`) |
| `title` | TEXT | Task/project/heading title |
| `notes` | TEXT | Markdown-ish notes body |
| `type` | INTEGER | **0** = to-do, **1** = project, **2** = heading |
| `status` | INTEGER | **0** = open, **2** = cancelled, **3** = completed |
| `trashed` | INTEGER | **0** = not trashed, **1** = trashed |
| `start` | INTEGER | **1** = started (Anytime/Today), **2** = Someday |
| `startDate` | INTEGER | Scheduled start date (encoded, see below) |
| `startBucket` | INTEGER | Sub-ordering within start |
| `deadline` | INTEGER | Deadline date (same encoding as startDate) |
| `project` | TEXT | UUID of parent project (NULL if not in a project) |
| `heading` | TEXT | UUID of parent heading within a project |
| `area` | TEXT | UUID of parent area |
| `index` | INTEGER | Sort order within its container |
| `todayIndex` | INTEGER | Sort order in Today list (>0 if shown in Today) |
| `creationDate` | REAL | Unix timestamp (seconds since 1970-01-01) |
| `userModificationDate` | REAL | Unix timestamp of last user edit |
| `stopDate` | REAL | Unix timestamp when completed/cancelled |
| `checklistItemsCount` | INTEGER | Total checklist items |
| `openChecklistItemsCount` | INTEGER | Incomplete checklist items |
| `untrashedLeafActionsCount` | INTEGER | Total child to-dos (for projects) |
| `openUntrashedLeafActionsCount` | INTEGER | Open child to-dos (for projects) |
| `rt1_repeatingTemplate` | TEXT | UUID of repeating template (if recurring) |
| `contact` | TEXT | UUID of delegated contact |
| `reminderTime` | INTEGER | Reminder time |

#### TMArea

| Column | Type | Description |
|--------|------|-------------|
| `uuid` | TEXT PK | Area identifier |
| `title` | TEXT | Area name |
| `visible` | INTEGER | Whether shown in sidebar |
| `index` | INTEGER | Sort order |

Current areas:
- `SmmUNQ4iCsQ9NFMC6bAiAB` — Dev & Admin
- `TEqjNkGzLkRCaMNeaUhvi3` — Work
- `V9QTStxBCV9BNuULFsPwhv` — Personal

#### TMTag

| Column | Type | Description |
|--------|------|-------------|
| `uuid` | TEXT PK | Tag identifier |
| `title` | TEXT | Tag name |
| `shortcut` | TEXT | Keyboard shortcut |
| `parent` | TEXT | UUID of parent tag (for nested tags) |
| `index` | INTEGER | Sort order |

#### TMTaskTag

Join table linking tasks to tags.

| Column | Type | Description |
|--------|------|-------------|
| `tasks` | TEXT | UUID of the task |
| `tags` | TEXT | UUID of the tag |

#### TMChecklistItem

| Column | Type | Description |
|--------|------|-------------|
| `uuid` | TEXT PK | Checklist item identifier |
| `title` | TEXT | Item text |
| `status` | INTEGER | **0** = open, **3** = completed |
| `task` | TEXT | UUID of parent task |
| `index` | INTEGER | Sort order |
| `creationDate` | REAL | Unix timestamp |
| `stopDate` | REAL | Unix timestamp when completed |

### Date Encoding (startDate, deadline)

The `startDate` and `deadline` columns use a bit-packed integer encoding:

```
Encode: (year << 16) | (month << 12) | (day << 7)
Decode: year = value >> 16
        month = (value >> 12) & 0xF
        day = (value >> 7) & 0x1F
```

Example: `132781056` → `0x07EA1400` → year=2026, month=1, day=8 → **2026-01-08**

The `creationDate`, `userModificationDate`, and `stopDate` fields are standard Unix timestamps (float, seconds since 1970-01-01).

### Common Queries

```sql
-- All open to-dos (not trashed, not in a project that is completed)
SELECT * FROM TMTask
WHERE type = 0 AND status = 0 AND trashed = 0;

-- All open projects
SELECT * FROM TMTask
WHERE type = 1 AND status = 0 AND trashed = 0;

-- To-dos in a specific project, with headings
SELECT t.title, t.status, h.title AS heading
FROM TMTask t
LEFT JOIN TMTask h ON t.heading = h.uuid
WHERE t.project = '<project-uuid>'
  AND t.type = 0 AND t.trashed = 0
ORDER BY t."index";

-- To-dos grouped by project within an area
SELECT p.title AS project, t.title AS task
FROM TMTask t
JOIN TMTask p ON t.project = p.uuid
WHERE p.area = '<area-uuid>'
  AND t.type = 0 AND t.status = 0 AND t.trashed = 0
ORDER BY p."index", t."index";

-- Today list (tasks with todayIndex > 0 or startDate = today)
SELECT * FROM TMTask
WHERE type = 0 AND status = 0 AND trashed = 0
  AND todayIndex > 0
ORDER BY todayIndex;

-- Someday tasks
SELECT * FROM TMTask
WHERE type = 0 AND status = 0 AND trashed = 0 AND start = 2;

-- Tags for a task
SELECT tg.title FROM TMTaskTag tt
JOIN TMTag tg ON tt.tags = tg.uuid
WHERE tt.tasks = '<task-uuid>';

-- Checklist items for a task
SELECT title, status FROM TMChecklistItem
WHERE task = '<task-uuid>'
ORDER BY "index";
```

## Things URL Scheme (Writing)

All writes go through the `things:///` URL scheme, invoked via `open` on macOS.
Full documentation is in `things-url.md`.

### Adding a to-do

```bash
open "things:///add?title=Buy%20milk&list=Shopping&when=today"
```

Key parameters: `title`, `notes`, `when` (today/tomorrow/evening/anytime/someday/date), `deadline`, `tags`, `list` (project or area name), `list-id`, `heading`, `checklist-items` (newline-separated), `completed`, `reveal`.

### Adding a project

```bash
open "things:///add-project?title=Plan%20Trip&area=Personal&to-dos=Book%20flights%0aBook%20hotel"
```

### Updating (requires auth-token)

```bash
open "things:///update?auth-token=TOKEN&id=UUID&completed=true"
```

The auth-token is found in Things → Settings → General → Enable Things URLs → Manage.

### JSON command (bulk/complex operations)

```bash
open "things:///json?data=$(python3 -c 'import urllib.parse, json; print(urllib.parse.quote(json.dumps([{"type":"to-do","attributes":{"title":"Test","list":"Work"}}])))')"
```

The JSON command accepts an array of to-do, project, heading, and checklist-item objects. Supports both `create` and `update` operations. See `things-url.md` for full schema.

### URL Scheme Limits

- Max 250 items per 10-second window
- String fields max 4,000 chars (notes max 10,000)
- Checklist items max 100 per task
- Update/update-project commands require `auth-token`
- Add commands do not require auth-token

## Development Notes

- A snapshot of the database is at `things.sqlite` in this repo for schema exploration (do not commit — add to `.gitignore`)
- The live database path contains a unique suffix (`ThingsData-HK9N6`) that may differ per installation
- Things must not be running when copying the database for a clean snapshot, but read-only access to the live file works while Things is running (SQLite WAL mode)
