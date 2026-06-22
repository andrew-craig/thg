# Things 3 Database Reference

`thg` reads tasks by querying the Things 3 SQLite database directly (always in
read-only mode). This document describes that database's location, schema, and
the queries `thg` relies on.

> **Reads only.** `thg` never writes to the database. All writes go through the
> [Things URL scheme](url-scheme.md). Let Things manage its own data.

## Location

```
~/Library/Group Containers/JLMPQHK86H.com.culturedcode.ThingsMac/ThingsData-XXXXX/Things Database.thingsdatabase/main.sqlite
```

The `.thingsdatabase` is a macOS bundle (a directory). The database is
`main.sqlite` inside it. WAL files (`main.sqlite-shm`, `main.sqlite-wal`) may
also be present.

The `ThingsData-XXXXX` directory has a unique suffix that differs per
installation, so `thg` scans the group container for a `ThingsData*` directory
rather than hard-coding the path. You can override the location with the
`THG_DB_PATH` environment variable.

> **Always open read-only** (`?mode=ro` in the SQLite URI, or
> `SQLITE_OPEN_READONLY`) to avoid any risk to the live database. Read-only
> access works while Things is running thanks to SQLite's WAL mode.

## Key Tables

### TMTask

The main table. Holds to-dos, projects, and headings in a single table,
distinguished by `type`.

| Column | Type | Description |
|--------|------|-------------|
| `uuid` | TEXT PK | Unique identifier (used in the URL scheme as `id`) |
| `title` | TEXT | Task/project/heading title |
| `notes` | TEXT | Markdown-ish notes body |
| `type` | INTEGER | **0** = to-do, **1** = project, **2** = heading |
| `status` | INTEGER | **0** = open, **2** = cancelled, **3** = completed |
| `trashed` | INTEGER | **0** = not trashed, **1** = trashed |
| `start` | INTEGER | **1** = started (Anytime/Today), **2** = Someday |
| `startDate` | INTEGER | Scheduled start date (encoded, see below) |
| `startBucket` | INTEGER | Sub-ordering within start |
| `deadline` | INTEGER | Deadline date (same encoding as `startDate`) |
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

### TMArea

| Column | Type | Description |
|--------|------|-------------|
| `uuid` | TEXT PK | Area identifier |
| `title` | TEXT | Area name |
| `visible` | INTEGER | Whether shown in sidebar |
| `index` | INTEGER | Sort order |

### TMTag

| Column | Type | Description |
|--------|------|-------------|
| `uuid` | TEXT PK | Tag identifier |
| `title` | TEXT | Tag name |
| `shortcut` | TEXT | Keyboard shortcut |
| `parent` | TEXT | UUID of parent tag (for nested tags) |
| `index` | INTEGER | Sort order |

### TMTaskTag

Join table linking tasks to tags.

| Column | Type | Description |
|--------|------|-------------|
| `tasks` | TEXT | UUID of the task |
| `tags` | TEXT | UUID of the tag |

### TMChecklistItem

| Column | Type | Description |
|--------|------|-------------|
| `uuid` | TEXT PK | Checklist item identifier |
| `title` | TEXT | Item text |
| `status` | INTEGER | **0** = open, **3** = completed |
| `task` | TEXT | UUID of parent task |
| `index` | INTEGER | Sort order |
| `creationDate` | REAL | Unix timestamp |
| `stopDate` | REAL | Unix timestamp when completed |

## Date Encoding (`startDate`, `deadline`)

The `startDate` and `deadline` columns use a bit-packed integer encoding:

```
Encode: (year << 16) | (month << 12) | (day << 7)
Decode: year  = value >> 16
        month = (value >> 12) & 0xF
        day   = (value >> 7)  & 0x1F
```

Example: `132781056` → `0x07EA1400` → year=2026, month=1, day=8 → **2026-01-08**

The `creationDate`, `userModificationDate`, and `stopDate` fields are standard
Unix timestamps (float, seconds since 1970-01-01).

## Common Queries

```sql
-- All open to-dos (not trashed)
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

-- Today list (tasks shown in Today)
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
