package main

import (
    "testing"
    "database/sql"
    "fmt"
    "os"
    "io"
    _ "github.com/mattn/go-sqlite3"
)

// Utility function to set up an in-memory test database.
func setupTestDB(t *testing.T) (*sql.DB, func()) {
    db, err := sql.Open("sqlite3", ":memory:")
    if err != nil {
        t.Fatalf("Failed to create in-memory database: %v", err)
    }

    createTableSQL := `
    CREATE TABLE IF NOT EXISTS aliases (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        folder_path TEXT NOT NULL,
        alias TEXT NOT NULL,
        absolute_path TEXT NOT NULL,
        created_at DATETIME NOT NULL,
        updated_at DATETIME NOT NULL,
        invocation_count INTEGER NOT NULL,
        UNIQUE(folder_path, alias)
    )`
    _, err = db.Exec(createTableSQL)
    if err != nil {
        t.Fatalf("Failed to create table: %v", err)
    }

    // Cleanup function to close the database after the test.
    cleanup := func() {
        db.Close()
    }

    return db, cleanup
}

// Test creating a new alias.
func TestCreateAlias(t *testing.T) {
    db, cleanup := setupTestDB(t)
    defer cleanup()

    folderPath := "code/projects"
    alias := "myapp"
    absolutePath := "/home/user/code/myapp"

    handleInsertOrUpdate(db, folderPath, alias, absolutePath, false)

    var storedPath string
    err := db.QueryRow(`SELECT absolute_path FROM aliases WHERE folder_path = ? AND alias = ?`, folderPath, alias).Scan(&storedPath)
    if err != nil {
        t.Fatalf("Expected alias to be created, got error: %v", err)
    }

    if storedPath != absolutePath {
        t.Fatalf("Expected path %s, got %s", absolutePath, storedPath)
    }
}

// Test updating an existing alias.
func TestUpdateAlias(t *testing.T) {
    db, cleanup := setupTestDB(t)
    defer cleanup()

    folderPath := "code/projects"
    alias := "myapp"
    initialPath := "/home/user/code/myapp"
    updatedPath := "/home/user/newpath/myapp"

    handleInsertOrUpdate(db, folderPath, alias, initialPath, false)
    handleInsertOrUpdate(db, folderPath, alias, updatedPath, true)

    var storedPath string
    err := db.QueryRow(`SELECT absolute_path FROM aliases WHERE folder_path = ? AND alias = ?`, folderPath, alias).Scan(&storedPath)
    if err != nil {
        t.Fatalf("Expected alias to be updated, got error: %v", err)
    }

    if storedPath != updatedPath {
        t.Fatalf("Expected path %s, got %s", updatedPath, storedPath)
    }
}

// Test deleting an alias.
func TestDeleteAlias(t *testing.T) {
    db, cleanup := setupTestDB(t)
    defer cleanup()

    folderPath := "code/projects"
    alias := "myapp"
    absolutePath := "/home/user/code/myapp"

    handleInsertOrUpdate(db, folderPath, alias, absolutePath, false)
    handleDelete(db, folderPath, alias)

    var storedPath string
    err := db.QueryRow(`SELECT absolute_path FROM aliases WHERE folder_path = ? AND alias = ?`, folderPath, alias).Scan(&storedPath)
    if err == nil {
        t.Fatalf("Expected alias to be deleted, but it still exists")
    }
    if err != sql.ErrNoRows {
        t.Fatalf("Unexpected error: %v", err)
    }
}

// Test recalling an alias.
func TestRecallAlias(t *testing.T) {
    db, cleanup := setupTestDB(t)
    defer cleanup()

    folderPath := "code/projects"
    alias := "myapp"
    absolutePath := "/home/user/code/myapp"

    handleInsertOrUpdate(db, folderPath, alias, absolutePath, false)

    // Capturing the output of handleRecall
    output := captureOutput(func() {
        handleRecall(db, folderPath, alias)
    })

    expectedOutput := fmt.Sprintf("cd %s\n", absolutePath)
    if output != expectedOutput {
        t.Fatalf("Expected output %q, got %q", expectedOutput, output)
    }
}

// Test recalling a non-existent alias.
func TestRecallNonExistentAlias(t *testing.T) {
    db, cleanup := setupTestDB(t)
    defer cleanup()

    folderPath := "code/projects"
    alias := "nonexistent"

    output := captureOutput(func() {
        handleRecall(db, folderPath, alias)
    })

    expectedOutput := fmt.Sprintf("Error: Alias '%s' under folder '%s' does not exist.\n", alias, folderPath)
    if output != expectedOutput {
        t.Fatalf("Expected output %q, got %q", expectedOutput, output)
    }
}

// Utility function to capture the output of a function.
func captureOutput(f func()) string {
    rescueStdout := os.Stdout
    r, w, _ := os.Pipe()
    os.Stdout = w

    f()

    w.Close()
    out, _ := io.ReadAll(r)
    os.Stdout = rescueStdout

    return string(out)
}

