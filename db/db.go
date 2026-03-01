package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

const (
	// Default database path components
	groupContainer = "Library/Group Containers/JLMPQHK86H.com.culturedcode.ThingsMac"
	dbBundle       = "Things Database.thingsdatabase"
	dbFile         = "main.sqlite"
)

// Open opens the live Things 3 database in read-only mode.
// It checks THG_DB_PATH env var first, then the standard location.
func Open() (*sql.DB, error) {
	path := os.Getenv("THG_DB_PATH")
	if path == "" {
		var err error
		path, err = findDB()
		if err != nil {
			return nil, err
		}
	}

	dsn := fmt.Sprintf("file:%s?mode=ro", path)
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	// Verify connection
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("connect to Things database at %s: %w", path, err)
	}

	return db, nil
}

func findDB() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("get home directory: %w", err)
	}

	containerDir := filepath.Join(home, groupContainer)

	// Find the ThingsData-* directory (suffix varies per install)
	entries, err := os.ReadDir(containerDir)
	if err != nil {
		return "", fmt.Errorf("read Things container at %s: %w\nIs Things 3 installed?", containerDir, err)
	}

	for _, e := range entries {
		if e.IsDir() && len(e.Name()) > 10 && e.Name()[:10] == "ThingsData" {
			path := filepath.Join(containerDir, e.Name(), dbBundle, dbFile)
			if _, err := os.Stat(path); err == nil {
				return path, nil
			}
		}
	}

	return "", fmt.Errorf("Things 3 database not found in %s\nIs Things 3 installed?", containerDir)
}
