package tests

import (
	"database/sql"
	"os"
	"testing"

	_ "modernc.org/sqlite"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open test DB: %v", err)
	}

	schema := `
	CREATE TABLE users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL,
		real_name TEXT,
		role TEXT DEFAULT 'user',
		status TEXT DEFAULT 'pending',
		avatar TEXT,
		team_id INTEGER,
		team_name TEXT,
		team_color TEXT,
		create_time DATETIME DEFAULT CURRENT_TIMESTAMP,
		update_time DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE blogger (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_name TEXT,
		real_name TEXT,
		platform TEXT,
		category TEXT,
		follower_count INTEGER DEFAULT 0,
		avg_view_count INTEGER DEFAULT 0,
		status TEXT DEFAULT 'pending',
		is_deleted INTEGER DEFAULT 0,
		user_id INTEGER,
		create_time DATETIME DEFAULT CURRENT_TIMESTAMP,
		update_time DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`

	_, err = db.Exec(schema)
	if err != nil {
		t.Fatalf("Failed to create schema: %v", err)
	}

	return db
}

func cleanupTestDB(t *testing.T, db *sql.DB) {
	db.Close()
}

func TestDatabaseConnection(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	err := db.Ping()
	if err != nil {
		t.Errorf("Failed to ping database: %v", err)
	}
}

func TestUserInsert(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	result, err := db.Exec(`
		INSERT INTO users (username, password, real_name, role, status)
		VALUES (?, ?, ?, ?, ?)
	`, "testuser", "hashedpass", "Test User", "user", "active")

	if err != nil {
		t.Errorf("Failed to insert user: %v", err)
	}

	id, _ := result.LastInsertId()
	if id == 0 {
		t.Error("Expected non-zero user ID")
	}
}

func TestBloggerInsert(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	result, err := db.Exec(`
		INSERT INTO blogger (user_name, real_name, platform, category, follower_count, status)
		VALUES (?, ?, ?, ?, ?, ?)
	`, "testblogger", "Test Blogger", "抖音", "美妆", 10000, "active")

	if err != nil {
		t.Errorf("Failed to insert blogger: %v", err)
	}

	id, _ := result.LastInsertId()
	if id == 0 {
		t.Error("Expected non-zero blogger ID")
	}
}
