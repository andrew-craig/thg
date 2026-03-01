package db

import "database/sql"

// Tag represents a Things 3 tag.
type Tag struct {
	UUID     string
	Title    string
	Shortcut string
	Parent   string
	Index    int
}

// ListTags returns all tags.
func ListTags(db *sql.DB) ([]Tag, error) {
	rows, err := db.Query(`SELECT uuid, title, COALESCE(shortcut,''),
		COALESCE(parent,''), COALESCE("index",0)
		FROM TMTag ORDER BY "index"`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []Tag
	for rows.Next() {
		var t Tag
		if err := rows.Scan(&t.UUID, &t.Title, &t.Shortcut, &t.Parent, &t.Index); err != nil {
			return nil, err
		}
		tags = append(tags, t)
	}
	return tags, rows.Err()
}

// TagsForTask returns the tag names for a given task UUID.
func TagsForTask(database *sql.DB, taskUUID string) ([]string, error) {
	rows, err := database.Query(`SELECT tg.title FROM TMTaskTag tt
		JOIN TMTag tg ON tt.tags = tg.uuid
		WHERE tt.tasks = ? ORDER BY tg."index"`, taskUUID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []string
	for rows.Next() {
		var title string
		if err := rows.Scan(&title); err != nil {
			return nil, err
		}
		tags = append(tags, title)
	}
	return tags, rows.Err()
}

// ResolveTag finds a tag by title (case-insensitive).
func ResolveTag(db *sql.DB, name string) (Tag, error) {
	var t Tag
	err := db.QueryRow(`SELECT uuid, title, COALESCE(shortcut,''),
		COALESCE(parent,''), COALESCE("index",0)
		FROM TMTag WHERE title LIKE ? LIMIT 1`, name).
		Scan(&t.UUID, &t.Title, &t.Shortcut, &t.Parent, &t.Index)
	if err != nil {
		return Tag{}, err
	}
	return t, nil
}
