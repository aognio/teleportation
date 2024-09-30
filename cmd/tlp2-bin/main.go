package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var sourced bool // Global variable to determine if commands should be output for sourcing

func main() {
	args := os.Args

	// Default action flags
	update := false
	list := false
	deleteFlag := false
	renameFlag := false
	renameFolderFlag := false
	sqliteFlag := false
	var newAliasName string
	var newFolderName string

	if len(args) < 2 {
		handleFatalErrorOutput("Invalid usage. Please check the instructions.")
		printUsage(args[0])
		return
	}

	// Extract the flag if the first argument is a switch
	for len(args) > 1 && strings.HasPrefix(args[1], "--") {
		switch args[1] {
		case "--update":
			update = true
		case "--list":
			list = true
		case "--delete":
			deleteFlag = true
		case "--rename":
			renameFlag = true
		case "--rename-folder":
			renameFolderFlag = true
		case "--sqlite":
			sqliteFlag = true
		case "--sourced":
			sourced = true
		default:
			handleFatalErrorOutput("Invalid flag provided. Please check the instructions.")
			printUsage(args[0])
			return
		}
		args = args[1:]
	}

	if sqliteFlag {
		handleOpenSQLite()
		return
	}

	db := initDatabase()
	defer db.Close()

	if list {
		if len(args) == 1 {
			handleList(db)
		} else if len(args) == 2 {
			partialFolder := args[1]
			handleListWithFolder(db, partialFolder)
		} else {
			printUsage(args[0])
		}
		return
	}

	if deleteFlag {
		if len(args) != 3 {
			printUsage(args[0])
			return
		}
		folderPath := args[1]
		alias := args[2]
		handleDelete(db, folderPath, alias)
		return
	}

	if renameFlag {
		if len(args) != 4 {
			printUsage(args[0])
			return
		}
		folderPath := args[1]
		oldAlias := args[2]
		newAliasName = args[3]
		handleRename(db, folderPath, oldAlias, newAliasName)
		return
	}

	if renameFolderFlag {
		if len(args) != 3 {
			printUsage(args[0])
			return
		}
		oldFolderPath := args[1]
		newFolderName = args[2]
		handleRenameFolder(db, oldFolderPath, newFolderName)
		return
	}

	// Handle insert/update and recall cases with no initial switch
	if len(args) == 2 {
		alias := args[1]
		handleAliasSearch(db, alias)
	} else if len(args) == 3 || len(args) == 4 {
		folderPath := args[1]
		alias := args[2]

		if len(args) == 4 {
			absolutePath := args[3]
			handleInsertOrUpdate(db, folderPath, alias, absolutePath, update)
		} else {
			handleRecall(db, folderPath, alias)
		}
	} else {
		printUsage(args[0])
	}
}

// handleFatalErrorOutput handles printing fatal error messages and exits
func handleFatalErrorOutput(message string) {
	if sourced {
		// Print the error with "echo" and exit with an error code
		fmt.Printf("echo %q\n", message)
		os.Exit(1)
	} else {
		// Use Go's default logging for non-sourced cases
		log.Fatalf("%s", message)
	}
}

// handleOpenSQLite opens the SQLite database file using the sqlite3 command-line client
func handleOpenSQLite() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		handleFatalErrorOutput(fmt.Sprintf("Error getting home directory: %v", err))
	}

	dbPath := filepath.Join(homeDir, ".config", "teleportation", "teleportation.sqlite3")
	if sourced {
		// When sourced, output the actual sqlite3 command without echo
		handleCommandOutput(fmt.Sprintf("sqlite3 %s", dbPath))
	} else {
		// When not sourced, provide an informative message
		handleInformativeOutput(fmt.Sprintf("Use the \"sqlite3\" command to manually edit the entries at the SQLite database at '%s'", dbPath))
	}
}

// initDatabase initializes the database and creates necessary tables
func initDatabase() *sql.DB {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		handleFatalErrorOutput(fmt.Sprintf("Error getting home directory: %v", err))
	}

	configDir := filepath.Join(homeDir, ".config", "teleportation")
	if err := os.MkdirAll(configDir, os.ModePerm); err != nil {
		handleFatalErrorOutput(fmt.Sprintf("Error creating config directory: %v", err))
	}

	dbPath := filepath.Join(configDir, "teleportation.sqlite3")
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		handleFatalErrorOutput(fmt.Sprintf("Error opening database: %v", err))
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
		handleFatalErrorOutput(fmt.Sprintf("Error creating table: %v", err))
	}

	return db
}

// handleInformativeOutput handles printing informational messages to the user
func handleInformativeOutput(message string) {
    if sourced {
        fmt.Printf("echo %s\n", escapeForShell(message))
    } else {
        fmt.Println(message)
    }
}

// handleCommandOutput handles printing commands to the output for sourcing
func handleCommandOutput(command string) {
    fmt.Println(command)
}

func escapeForShell(input string) string {
    replacer := strings.NewReplacer(
        `"`, `\"`,
        `'`, `\'`,
        "`", "\\`",
        "&", `\&`,
        "|", `\|`,
        "<", `\<`,
        ">", `\>`,
        ";", `\;`,
        "(", `\(`,
        ")", `\)`,
        "$", `\$`,
    )
    return replacer.Replace(input)
}

// handleErrorOutput handles printing error messages
func handleErrorOutput(message string) {
	if sourced {
		fmt.Printf("echo %q\n", message)
	} else {
		fmt.Println(message)
	}
}

// handleInsertOrUpdate handles inserting a new alias or updating an existing one
func handleInsertOrUpdate(db *sql.DB, folderPath, alias, absolutePath string, update bool) {
	now := time.Now().Format(time.RFC3339)

	// Check if alias already exists
	var existingPath string
	selectSQL := `SELECT absolute_path FROM aliases WHERE folder_path = ? AND alias = ?`
	err := db.QueryRow(selectSQL, folderPath, alias).Scan(&existingPath)

	if err == nil && !update {
		handleErrorOutput(fmt.Sprintf("Error: Alias '%s' under folder '%s' already exists. Use --update to modify it.", alias, folderPath))
		return
	}

	if err == nil && update {
		// Update existing alias without changing invocation count
		updateSQL := `
            UPDATE aliases
            SET absolute_path = ?, updated_at = ?
            WHERE folder_path = ? AND alias = ?`
		_, err := db.Exec(updateSQL, absolutePath, now, folderPath, alias)
		if err != nil {
			handleFatalErrorOutput(fmt.Sprintf("Error updating alias: %v", err))
		}
		handleInformativeOutput(fmt.Sprintf("Alias '%s' updated.", alias))
	} else if err == sql.ErrNoRows {
		// Insert new alias with invocation count set to 0
		insertSQL := `
            INSERT INTO aliases (folder_path, alias, absolute_path, created_at, updated_at, invocation_count)
            VALUES (?, ?, ?, ?, ?, 0)`
		_, err := db.Exec(insertSQL, folderPath, alias, absolutePath, now, now)
		if err != nil {
			handleFatalErrorOutput(fmt.Sprintf("Error inserting alias: %v", err))
		}
		handleInformativeOutput(fmt.Sprintf("Alias '%s' created.", alias))
	} else if err != nil {
		handleFatalErrorOutput(fmt.Sprintf("Error retrieving alias: %v", err))
	}
}

// handleRecall handles recalling an alias from a specific folder and printing the appropriate output
func handleRecall(db *sql.DB, folderPath, alias string) {
    var absolutePath string
    selectSQL := `SELECT absolute_path FROM aliases WHERE folder_path = ? AND alias = ?`
    err := db.QueryRow(selectSQL, folderPath, alias).Scan(&absolutePath)
    if err == sql.ErrNoRows {
        handleErrorOutput(fmt.Sprintf("Error: Alias '%s' under folder '%s' does not exist.", alias, folderPath))
    } else if err != nil {
        handleErrorOutput(fmt.Sprintf("Error retrieving alias: %v", err))
    } else {
        // Update the invocation count
        updateSQL := `UPDATE aliases SET invocation_count = invocation_count + 1 WHERE folder_path = ? AND alias = ?`
        _, err = db.Exec(updateSQL, folderPath, alias)
        if err != nil {
            handleErrorOutput(fmt.Sprintf("Error updating invocation count: %v", err))
            return
        }
        
        // For recall, assume the parent shell changes directories
        fmt.Printf("cd %s\n", absolutePath)
    }
}

// handleList lists all aliases in the database
func handleList(db *sql.DB) {
	rows, err := db.Query(`SELECT folder_path, alias, absolute_path FROM aliases ORDER BY folder_path, alias`)
	if err != nil {
		handleFatalErrorOutput(fmt.Sprintf("Error listing aliases: %v", err))
	}
	defer rows.Close()

	var aliases []string
	for rows.Next() {
		var folderPath, alias, absolutePath string
		if err := rows.Scan(&folderPath, &alias, &absolutePath); err != nil {
			handleFatalErrorOutput(fmt.Sprintf("Error reading alias: %v", err))
		}
		aliases = append(aliases, fmt.Sprintf("%s %s %s", alias, folderPath, absolutePath))
	}

	if len(aliases) == 0 {
		handleInformativeOutput("No aliases found.")
	} else {
		handleInformativeOutput(fmt.Sprintf("%d aliases found:", len(aliases)))
		for _, alias := range aliases {
			handleInformativeOutput(alias)
		}
	}
}

// handleListWithFolder lists all aliases under a specific folder path
func handleListWithFolder(db *sql.DB, partialFolder string) {
	rows, err := db.Query(`SELECT folder_path, alias, absolute_path FROM aliases WHERE folder_path LIKE ? ORDER BY folder_path, alias`, "%"+partialFolder+"%")
	if err != nil {
		handleFatalErrorOutput(fmt.Sprintf("Error listing aliases: %v", err))
	}
	defer rows.Close()

	var aliases []string
	for rows.Next() {
		var folderPath, alias, absolutePath string
		if err := rows.Scan(&folderPath, &alias, &absolutePath); err != nil {
			handleFatalErrorOutput(fmt.Sprintf("Error reading alias: %v", err))
		}
		aliases = append(aliases, fmt.Sprintf("%s %s %s", alias, folderPath, absolutePath))
	}

	if len(aliases) == 0 {
		handleInformativeOutput(fmt.Sprintf("No aliases found under the folder path '%s'.", partialFolder))
	} else {
		handleInformativeOutput(fmt.Sprintf("%d aliases found under the folder path '%s':", len(aliases), partialFolder))
		for _, alias := range aliases {
			handleInformativeOutput(alias)
		}
	}
}

// handleDelete deletes an alias from the database
func handleDelete(db *sql.DB, folderPath, alias string) {
	deleteSQL := `DELETE FROM aliases WHERE folder_path = ? AND alias = ?`
	result, err := db.Exec(deleteSQL, folderPath, alias)
	if err != nil {
		handleFatalErrorOutput(fmt.Sprintf("Error deleting alias: %v", err))
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		handleFatalErrorOutput(fmt.Sprintf("Error fetching delete result: %v", err))
	}

	if rowsAffected == 0 {
		handleErrorOutput(fmt.Sprintf("Alias '%s' under folder '%s' not found.", alias, folderPath))
	} else {
		handleInformativeOutput(fmt.Sprintf("Alias '%s' under folder '%s' deleted successfully.", alias, folderPath))
	}
}

// handleRename renames an existing alias
func handleRename(db *sql.DB, folderPath, oldAlias, newAliasName string) {
	updateSQL := `UPDATE aliases SET alias = ?, updated_at = ? WHERE folder_path = ? AND alias = ?`
	now := time.Now().Format(time.RFC3339)
	result, err := db.Exec(updateSQL, newAliasName, now, folderPath, oldAlias)
	if err != nil {
		handleFatalErrorOutput(fmt.Sprintf("Error renaming alias: %v", err))
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		handleFatalErrorOutput(fmt.Sprintf("Error fetching rename result: %v", err))
	}

	if rowsAffected == 0 {
		handleErrorOutput(fmt.Sprintf("Alias '%s' under folder '%s' not found.", oldAlias, folderPath))
	} else {
		handleInformativeOutput(fmt.Sprintf("Alias '%s' under folder '%s' renamed to '%s'.", oldAlias, folderPath, newAliasName))
	}
}

// handleRenameFolder renames an existing folder path for all aliases
func handleRenameFolder(db *sql.DB, oldFolderPath, newFolderPath string) {
	updateSQL := `UPDATE aliases SET folder_path = ?, updated_at = ? WHERE folder_path = ?`
	now := time.Now().Format(time.RFC3339)
	result, err := db.Exec(updateSQL, newFolderPath, now, oldFolderPath)
	if err != nil {
		handleFatalErrorOutput(fmt.Sprintf("Error renaming folder path: %v", err))
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		handleFatalErrorOutput(fmt.Sprintf("Error fetching rename result: %v", err))
	}

	if rowsAffected == 0 {
		handleErrorOutput(fmt.Sprintf("No aliases found under folder path '%s'.", oldFolderPath))
	} else {
		handleInformativeOutput(fmt.Sprintf("Folder path '%s' renamed to '%s' for %d aliases.", oldFolderPath, newFolderPath, rowsAffected))
	}
}

// handleAliasSearch searches for an alias across all folders and teleports if only one match is found
func handleAliasSearch(db *sql.DB, alias string) {
    rows, err := db.Query(`SELECT folder_path, alias, absolute_path FROM aliases WHERE alias = ? ORDER BY folder_path, alias`, alias)
    if err != nil {
        handleFatalErrorOutput(fmt.Sprintf("Error searching for alias: %v", err))
    }
    defer rows.Close()

    type AliasRecord struct {
        FolderPath   string
        Alias        string
        AbsolutePath string
    }

    var aliases []AliasRecord
    for rows.Next() {
        var record AliasRecord
        if err := rows.Scan(&record.FolderPath, &record.Alias, &record.AbsolutePath); err != nil {
            handleFatalErrorOutput(fmt.Sprintf("Error reading alias: %v", err))
        }
        aliases = append(aliases, record)
    }

    if len(aliases) == 0 {
        handleInformativeOutput(fmt.Sprintf("Alias '%s' not found.", alias))
    } else if len(aliases) == 1 {
        // Only one alias found, teleport directly there
        handleCommandOutput(fmt.Sprintf("cd %s", aliases[0].AbsolutePath))
    } else {
        // Multiple aliases found, print informative output
        handleInformativeOutput(fmt.Sprintf("There are %d aliases under different folders. Please specify the folder too:", len(aliases)))
        for _, aliasRecord := range aliases {
            handleInformativeOutput(fmt.Sprintf("%s %s %s", aliasRecord.Alias, aliasRecord.FolderPath, aliasRecord.AbsolutePath))
        }
    }
}

// updateInvocationCount increments the invocation count for a given alias
func updateInvocationCount(db *sql.DB, folderPath, alias string) {
	updateSQL := `UPDATE aliases SET invocation_count = invocation_count + 1 WHERE folder_path = ? AND alias = ?`
	_, err := db.Exec(updateSQL, folderPath, alias)
	if err != nil {
		handleFatalErrorOutput(fmt.Sprintf("Error updating invocation count: %v", err))
	}
}

// printUsage prints usage information for the tlp2 command
func printUsage(command string) {
	handleInformativeOutput(fmt.Sprintf("Usage for storing: %s <folder_path> <alias> <absolute_path> [--update] [--sourced]", command))
	handleInformativeOutput(fmt.Sprintf("Usage for recalling: %s <folder_path> <alias> [--sourced]", command))
	handleInformativeOutput(fmt.Sprintf("Usage for searching: %s <alias> [--sourced]", command))
	handleInformativeOutput(fmt.Sprintf("Usage for listing all aliases: %s --list [<partial_folder>] [--sourced]", command))
	handleInformativeOutput(fmt.Sprintf("Usage for deleting: %s --delete <folder_path> <alias> [--sourced]", command))
	handleInformativeOutput(fmt.Sprintf("Usage for renaming an alias: %s --rename <folder_path> <alias> <new_alias_name> [--sourced]", command))
	handleInformativeOutput(fmt.Sprintf("Usage for renaming a folder: %s --rename-folder <old_folder_path> <new_folder_path> [--sourced]", command))
	handleInformativeOutput(fmt.Sprintf("Usage for opening SQLite: %s --sqlite", command))
}

