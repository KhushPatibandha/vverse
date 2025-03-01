package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func SetupDB() error {
	err := openDB()
	if err != nil {
		return err
	}
	err = setupTables()
	if err != nil {
		return err
	}
	return nil
}

func CloseDB() error {
	return db.Close()
}

func ExecCmd(query string, args ...any) (sql.Result, error) {
	return db.Exec(query, args...)
}

func openDB() error {
	dataB, err := sql.Open("sqlite3", "./sqlite3.db")
	if err != nil {
		return err
	}
	db = dataB
	return nil
}

func setupTables() error {
	videoTable := `
    CREATE TABLE IF NOT EXISTS videos (
        id INTEGER PRIMARY KEY AUTOINCREMENT,           -- Unique ID for the video
        name TEXT NOT NULL,                             -- Name of the video
        size INTEGER NOT NULL,                          -- Size of the video in bytes
        duration INTEGER NOT NULL,                      -- Duration of the video in seconds
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP  -- Time when the video was uploaded
    );`
	linkTable := `
    CREATE TABLE IF NOT EXISTS links (
        video_id INTEGER NOT NULL,                      -- ID of the video
        link TEXT PRIMARY KEY,                          -- Link to the video
        expiry TIMESTAMP NOT NULL                       -- Time when the link will expire (5 mins)
    );`

	_, err := db.Exec(videoTable)
	if err != nil {
		return err
	}
	_, err = db.Exec(linkTable)
	if err != nil {
		return err
	}
	return nil
}
