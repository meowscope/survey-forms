package repository

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

const initSchema = `
CREATE TABLE IF NOT EXISTS surveys (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS questions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    survey_id TEXT,
    content TEXT,
    FOREIGN KEY(survey_id) REFERENCES surveys(id)
);`

func OpenDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./my.db")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	fmt.Println("established connection to db")
	return db, nil
}

func InitSchema(db *sql.DB) error {
	_, err := db.Exec(initSchema)
	if err != nil {
		return fmt.Errorf("failed to initialize tables %w", err)
	}

	_, err = db.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		return fmt.Errorf("failed to turn on fkeys at %w", err)
	}
	return nil
}
