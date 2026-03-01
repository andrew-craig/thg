package db

import (
	"database/sql"
	"fmt"
)

// Area represents a Things 3 area.
type Area struct {
	UUID    string
	Title   string
	Visible int
	Index   int
}

// ListAreas returns all visible areas.
func ListAreas(db *sql.DB) ([]Area, error) {
	rows, err := db.Query(`SELECT uuid, title, COALESCE(visible,1), COALESCE("index",0)
		FROM TMArea ORDER BY "index"`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var areas []Area
	for rows.Next() {
		var a Area
		if err := rows.Scan(&a.UUID, &a.Title, &a.Visible, &a.Index); err != nil {
			return nil, err
		}
		areas = append(areas, a)
	}
	return areas, rows.Err()
}

// ResolveArea finds an area by UUID or title substring.
func ResolveArea(db *sql.DB, identifier string) (Area, error) {
	// Try exact UUID
	var a Area
	err := db.QueryRow(`SELECT uuid, title, COALESCE(visible,1), COALESCE("index",0)
		FROM TMArea WHERE uuid = ?`, identifier).Scan(&a.UUID, &a.Title, &a.Visible, &a.Index)
	if err == nil {
		return a, nil
	}

	// Try title match
	rows, err := db.Query(`SELECT uuid, title, COALESCE(visible,1), COALESCE("index",0)
		FROM TMArea WHERE title LIKE ? ORDER BY "index"`, "%"+identifier+"%")
	if err != nil {
		return Area{}, err
	}
	defer rows.Close()

	var areas []Area
	for rows.Next() {
		var a Area
		if err := rows.Scan(&a.UUID, &a.Title, &a.Visible, &a.Index); err != nil {
			return Area{}, err
		}
		areas = append(areas, a)
	}
	if len(areas) == 0 {
		return Area{}, fmt.Errorf("no area found matching %q", identifier)
	}
	return areas[0], nil
}
