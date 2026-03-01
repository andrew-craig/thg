package db

import "database/sql"

// ChecklistItem represents a Things 3 checklist item.
type ChecklistItem struct {
	UUID     string
	Title    string
	Status   int // 0=open, 3=completed
	TaskUUID string
	Index    int
}

// ChecklistForTask returns checklist items for a task.
func ChecklistForTask(db *sql.DB, taskUUID string) ([]ChecklistItem, error) {
	rows, err := db.Query(`SELECT uuid, title, COALESCE(status,0), task, COALESCE("index",0)
		FROM TMChecklistItem WHERE task = ? ORDER BY "index"`, taskUUID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []ChecklistItem
	for rows.Next() {
		var ci ChecklistItem
		if err := rows.Scan(&ci.UUID, &ci.Title, &ci.Status, &ci.TaskUUID, &ci.Index); err != nil {
			return nil, err
		}
		items = append(items, ci)
	}
	return items, rows.Err()
}
