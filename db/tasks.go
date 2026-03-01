package db

import (
	"database/sql"
	"fmt"
	"strings"
)

// Task represents a Things 3 task (to-do, project, or heading).
type Task struct {
	UUID         string
	Title        string
	Notes        string
	Type         int // 0=todo, 1=project, 2=heading
	Status       int // 0=open, 2=cancelled, 3=completed
	Start        int // 1=started, 2=someday
	StartDate    int64
	Deadline     int64
	ProjectUUID  string
	ProjectTitle string
	HeadingUUID  string
	HeadingTitle string
	AreaUUID     string
	AreaTitle     string
	CreationDate float64
	ModDate      float64
	StopDate     float64
	TodayIndex   int64

	// Project-specific
	TotalTasks int
	OpenTasks  int

	// Populated separately
	Tags      []string
	Checklist []ChecklistItem
}

// ShortID returns the first 6 characters of the UUID.
func (t Task) ShortID() string {
	if len(t.UUID) >= 6 {
		return t.UUID[:6]
	}
	return t.UUID
}

// StatusText returns a human-readable status.
func (t Task) StatusText() string {
	switch t.Status {
	case 0:
		return "open"
	case 2:
		return "cancelled"
	case 3:
		return "completed"
	default:
		return "unknown"
	}
}

// WhenText returns a human-readable when value.
func (t Task) WhenText() string {
	if t.TodayIndex > 0 {
		return "today"
	}
	if t.Start == 2 {
		return "someday"
	}
	if t.StartDate > 0 {
		return FormatDate(t.StartDate)
	}
	if t.Start == 1 {
		return "anytime"
	}
	return "inbox"
}

const taskColumns = `t.uuid, t.title, COALESCE(t.notes,''), t.type, t.status,
	t.start, COALESCE(t.startDate,0), COALESCE(t.deadline,0),
	COALESCE(t.project,''), COALESCE(p.title,''),
	COALESCE(t.heading,''), COALESCE(h.title,''),
	COALESCE(t.area,''), COALESCE(a.title,''),
	COALESCE(t.creationDate,0), COALESCE(t.userModificationDate,0),
	COALESCE(t.stopDate,0), COALESCE(t.todayIndex,0),
	COALESCE(t.untrashedLeafActionsCount,0), COALESCE(t.openUntrashedLeafActionsCount,0)`

const taskJoins = `FROM TMTask t
	LEFT JOIN TMTask p ON t.project = p.uuid
	LEFT JOIN TMTask h ON t.heading = h.uuid
	LEFT JOIN TMArea a ON t.area = a.uuid`

func scanTask(row interface{ Scan(...any) error }) (Task, error) {
	var t Task
	err := row.Scan(
		&t.UUID, &t.Title, &t.Notes, &t.Type, &t.Status,
		&t.Start, &t.StartDate, &t.Deadline,
		&t.ProjectUUID, &t.ProjectTitle,
		&t.HeadingUUID, &t.HeadingTitle,
		&t.AreaUUID, &t.AreaTitle,
		&t.CreationDate, &t.ModDate, &t.StopDate, &t.TodayIndex,
		&t.TotalTasks, &t.OpenTasks,
	)
	return t, err
}

func scanTasks(rows *sql.Rows) ([]Task, error) {
	var tasks []Task
	for rows.Next() {
		t, err := scanTask(rows)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, rows.Err()
}

// ListToday returns tasks in the Today list.
func ListToday(db *sql.DB) ([]Task, error) {
	query := fmt.Sprintf(`SELECT %s %s
		WHERE t.type = 0 AND t.status = 0 AND t.trashed = 0
		AND t.todayIndex > 0
		ORDER BY t.todayIndex`, taskColumns, taskJoins)
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanTasks(rows)
}

// ListInbox returns tasks in the Inbox (no project, no area, start=0).
func ListInbox(db *sql.DB) ([]Task, error) {
	query := fmt.Sprintf(`SELECT %s %s
		WHERE t.type = 0 AND t.status = 0 AND t.trashed = 0
		AND (t.project IS NULL OR t.project = '')
		AND (t.area IS NULL OR t.area = '')
		AND t.start = 0
		AND COALESCE(t.startDate, 0) = 0
		ORDER BY t."index"`, taskColumns, taskJoins)
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanTasks(rows)
}

// ListSomeday returns tasks with start=2 (Someday).
func ListSomeday(db *sql.DB) ([]Task, error) {
	query := fmt.Sprintf(`SELECT %s %s
		WHERE t.type = 0 AND t.status = 0 AND t.trashed = 0
		AND t.start = 2
		ORDER BY t."index"`, taskColumns, taskJoins)
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanTasks(rows)
}

// ListUpcoming returns tasks with a future startDate.
func ListUpcoming(db *sql.DB) ([]Task, error) {
	today := TodayEncoded()
	query := fmt.Sprintf(`SELECT %s %s
		WHERE t.type = 0 AND t.status = 0 AND t.trashed = 0
		AND COALESCE(t.startDate, 0) > ?
		ORDER BY t.startDate, t."index"`, taskColumns, taskJoins)
	rows, err := db.Query(query, today)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanTasks(rows)
}

// ListAll returns all open tasks.
func ListAll(db *sql.DB) ([]Task, error) {
	query := fmt.Sprintf(`SELECT %s %s
		WHERE t.type = 0 AND t.status = 0 AND t.trashed = 0
		ORDER BY t."index"`, taskColumns, taskJoins)
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanTasks(rows)
}

// ListByProject returns tasks in a specific project.
func ListByProject(database *sql.DB, projectUUID string) ([]Task, error) {
	query := fmt.Sprintf(`SELECT %s %s
		WHERE t.project = ? AND t.type = 0 AND t.trashed = 0
		AND t.status = 0
		ORDER BY t."index"`, taskColumns, taskJoins)
	rows, err := database.Query(query, projectUUID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanTasks(rows)
}

// ListByArea filters open tasks by area UUID.
func ListByArea(db *sql.DB, areaUUID string) ([]Task, error) {
	query := fmt.Sprintf(`SELECT %s %s
		WHERE t.type = 0 AND t.status = 0 AND t.trashed = 0
		AND (t.area = ? OR p.area = ?)
		ORDER BY t."index"`, taskColumns, taskJoins)
	rows, err := db.Query(query, areaUUID, areaUUID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanTasks(rows)
}

// ListByTag filters open tasks by tag UUID.
func ListByTag(db *sql.DB, tagUUID string) ([]Task, error) {
	query := fmt.Sprintf(`SELECT %s %s
		JOIN TMTaskTag tt ON tt.tasks = t.uuid
		WHERE t.type = 0 AND t.status = 0 AND t.trashed = 0
		AND tt.tags = ?
		ORDER BY t."index"`, taskColumns, taskJoins)
	rows, err := db.Query(query, tagUUID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanTasks(rows)
}

// ResolveTask finds a task by full UUID, UUID prefix (6+ chars), or title substring.
func ResolveTask(db *sql.DB, identifier string) (Task, error) {
	// Try exact UUID first
	query := fmt.Sprintf(`SELECT %s %s WHERE t.uuid = ?`, taskColumns, taskJoins)
	row := db.QueryRow(query, identifier)
	t, err := scanTask(row)
	if err == nil {
		return t, nil
	}

	// Try UUID prefix (6+ chars)
	if len(identifier) >= 6 {
		query = fmt.Sprintf(`SELECT %s %s WHERE t.uuid LIKE ? AND t.trashed = 0`, taskColumns, taskJoins)
		rows, err := db.Query(query, identifier+"%")
		if err != nil {
			return Task{}, err
		}
		defer rows.Close()
		tasks, err := scanTasks(rows)
		if err != nil {
			return Task{}, err
		}
		if len(tasks) == 1 {
			return tasks[0], nil
		}
		if len(tasks) > 1 {
			return Task{}, fmt.Errorf("ambiguous ID prefix %q matches %d tasks", identifier, len(tasks))
		}
	}

	// Try title substring match
	query = fmt.Sprintf(`SELECT %s %s
		WHERE t.trashed = 0 AND t.status = 0
		AND t.title LIKE ?
		ORDER BY t.type, t."index"`, taskColumns, taskJoins)
	rows, err := db.Query(query, "%"+identifier+"%")
	if err != nil {
		return Task{}, err
	}
	defer rows.Close()
	tasks, err := scanTasks(rows)
	if err != nil {
		return Task{}, err
	}
	if len(tasks) == 0 {
		return Task{}, fmt.Errorf("no task found matching %q", identifier)
	}
	if len(tasks) == 1 {
		return tasks[0], nil
	}

	// Multiple matches: show them
	var b strings.Builder
	fmt.Fprintf(&b, "ambiguous match %q — %d results:\n", identifier, len(tasks))
	for _, tk := range tasks {
		typeName := "todo"
		if tk.Type == 1 {
			typeName = "project"
		}
		fmt.Fprintf(&b, "  %s  %s  (%s)\n", tk.ShortID(), tk.Title, typeName)
	}
	return Task{}, fmt.Errorf("%s", b.String())
}

// ResolveProject finds a project by UUID prefix or title substring.
func ResolveProject(db *sql.DB, identifier string) (Task, error) {
	// Try UUID prefix
	if len(identifier) >= 6 {
		query := fmt.Sprintf(`SELECT %s %s
			WHERE t.uuid LIKE ? AND t.type = 1 AND t.trashed = 0`, taskColumns, taskJoins)
		rows, err := db.Query(query, identifier+"%")
		if err != nil {
			return Task{}, err
		}
		defer rows.Close()
		tasks, err := scanTasks(rows)
		if err != nil {
			return Task{}, err
		}
		if len(tasks) == 1 {
			return tasks[0], nil
		}
	}

	// Try title match
	query := fmt.Sprintf(`SELECT %s %s
		WHERE t.type = 1 AND t.trashed = 0 AND t.status = 0
		AND t.title LIKE ?
		ORDER BY t."index"`, taskColumns, taskJoins)
	rows, err := db.Query(query, "%"+identifier+"%")
	if err != nil {
		return Task{}, err
	}
	defer rows.Close()
	tasks, err := scanTasks(rows)
	if err != nil {
		return Task{}, err
	}
	if len(tasks) == 0 {
		return Task{}, fmt.Errorf("no project found matching %q", identifier)
	}
	if len(tasks) == 1 {
		return tasks[0], nil
	}

	var b strings.Builder
	fmt.Fprintf(&b, "ambiguous project match %q:\n", identifier)
	for _, t := range tasks {
		fmt.Fprintf(&b, "  %s  %s\n", t.ShortID(), t.Title)
	}
	return Task{}, fmt.Errorf("%s", b.String())
}

// ListProjects returns all open projects with task counts.
func ListProjects(db *sql.DB, areaFilter string, showAll bool) ([]Task, error) {
	where := "t.type = 1 AND t.trashed = 0"
	var args []any
	if !showAll {
		where += " AND t.status = 0"
	}
	if areaFilter != "" {
		where += " AND (t.area = ? OR a.title LIKE ?)"
		args = append(args, areaFilter, "%"+areaFilter+"%")
	}
	query := fmt.Sprintf(`SELECT %s %s WHERE %s ORDER BY t."index"`, taskColumns, taskJoins, where)
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanTasks(rows)
}
