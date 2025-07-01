package database

import (
	"database/sql"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

// Initialize sets up the SQLite database and creates tables
func Initialize() {
	var err error
	DB, err = sql.Open("sqlite3", "./visits.db")
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	createTables()
	seedData()
	log.Println("Database initialized successfully")
}

// createTables creates the necessary tables
func createTables() {
	scheduleTable := `
	CREATE TABLE IF NOT EXISTS schedules (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		client_name TEXT NOT NULL,
		shift_start DATETIME NOT NULL,
		shift_end DATETIME NOT NULL,
		latitude REAL NOT NULL,
		longitude REAL NOT NULL,
		status TEXT NOT NULL DEFAULT 'upcoming',
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	taskTable := `
	CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		schedule_id INTEGER NOT NULL,
		description TEXT NOT NULL,
		status TEXT NOT NULL DEFAULT 'pending',
		reason TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (schedule_id) REFERENCES schedules (id)
	);`

	visitTable := `
	CREATE TABLE IF NOT EXISTS visits (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		schedule_id INTEGER NOT NULL,
		start_time DATETIME,
		end_time DATETIME,
		start_lat REAL,
		start_lng REAL,
		end_lat REAL,
		end_lng REAL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (schedule_id) REFERENCES schedules (id)
	);`

	activityTable := `
	CREATE TABLE IF NOT EXISTS activities (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		schedule_id INTEGER NOT NULL,
		title TEXT NOT NULL,
		description TEXT NOT NULL,
		is_resolved BOOLEAN NOT NULL DEFAULT 0,
		reason TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (schedule_id) REFERENCES schedules (id)
	);`

	tables := []string{scheduleTable, taskTable, visitTable, activityTable}
	for _, table := range tables {
		if _, err := DB.Exec(table); err != nil {
			log.Fatal("Failed to create table:", err)
		}
	}
}

// seedData loads and executes the comprehensive seed data from SQL file
func seedData() {
	// Check if data already exists
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM schedules").Scan(&count)
	if err != nil {
		log.Fatal("Failed to check existing data:", err)
	}

	if count > 0 {
		log.Println("Data already exists, skipping seed")
		return // Data already exists
	}

	// Get the current directory and construct path to SQL file
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatal("Failed to get current directory:", err)
	}

	sqlFilePath := filepath.Join(currentDir, "database", "seed_data.sql")
	
	// Read the SQL file
	sqlBytes, err := ioutil.ReadFile(sqlFilePath)
	if err != nil {
		log.Printf("Warning: Could not read seed_data.sql file at %s: %v", sqlFilePath, err)
		log.Println("Falling back to minimal seed data...")
		seedMinimalData()
		return
	}

	sqlContent := string(sqlBytes)
	
	// Execute the entire SQL content as a single transaction
	// This is more reliable than splitting by semicolons
	tx, err := DB.Begin()
	if err != nil {
		log.Printf("Failed to begin transaction: %v", err)
		seedMinimalData()
		return
	}
	defer tx.Rollback()

	// Remove comments and empty lines for cleaner execution
	lines := strings.Split(sqlContent, "\n")
	var cleanedLines []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "--") {
			cleanedLines = append(cleanedLines, line)
		}
	}
	cleanedSQL := strings.Join(cleanedLines, "\n")

	// Split by semicolon and execute each statement
	statements := strings.Split(cleanedSQL, ";")
	
	for _, statement := range statements {
		statement = strings.TrimSpace(statement)
		if statement == "" {
			continue
		}
		
		log.Printf("Executing SQL statement: %s", statement[:min(50, len(statement))])
		_, err := tx.Exec(statement)
		if err != nil {
			log.Printf("Error executing SQL statement: %v\nStatement: %s", err, statement)
			tx.Rollback()
			log.Println("Rolling back transaction and falling back to minimal seed data...")
			seedMinimalData()
			return
		}
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		log.Printf("Failed to commit transaction: %v", err)
		seedMinimalData()
		return
	}

	log.Println("Comprehensive seed data loaded successfully from SQL file")
}

// Helper function to get minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// seedMinimalData provides fallback minimal data if SQL file can't be loaded
func seedMinimalData() {
	now := time.Now()
	
	// Create proper time.Time values for different dates
	todayMorning := time.Date(now.Year(), now.Month(), now.Day(), 9, 0, 0, 0, now.Location())
	todayMorningEnd := time.Date(now.Year(), now.Month(), now.Day(), 11, 0, 0, 0, now.Location())
	todayAfternoon := time.Date(now.Year(), now.Month(), now.Day(), 14, 0, 0, 0, now.Location())
	todayAfternoonEnd := time.Date(now.Year(), now.Month(), now.Day(), 16, 0, 0, 0, now.Location())
	
	yesterday := now.AddDate(0, 0, -1)
	yesterdayMorning := time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 10, 0, 0, 0, yesterday.Location())
	yesterdayMorningEnd := time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 12, 0, 0, 0, yesterday.Location())
	
	tomorrow := now.AddDate(0, 0, 1)
	tomorrowMorning := time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 8, 0, 0, 0, tomorrow.Location())
	tomorrowMorningEnd := time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 10, 0, 0, 0, tomorrow.Location())

	// Sample schedules with proper time formatting
	schedules := []map[string]interface{}{
		{
			"client_name": "John Smith",
			"shift_start": todayMorning.Format("2006-01-02 15:04:05"),
			"shift_end":   todayMorningEnd.Format("2006-01-02 15:04:05"),
			"latitude":    40.7128,
			"longitude":   -74.0060,
			"status":      "upcoming",
		},
		{
			"client_name": "Mary Johnson",
			"shift_start": todayAfternoon.Format("2006-01-02 15:04:05"),
			"shift_end":   todayAfternoonEnd.Format("2006-01-02 15:04:05"),
			"latitude":    40.7589,
			"longitude":   -73.9851,
			"status":      "upcoming",
		},
		{
			"client_name": "Robert Davis",
			"shift_start": yesterdayMorning.Format("2006-01-02 15:04:05"),
			"shift_end":   yesterdayMorningEnd.Format("2006-01-02 15:04:05"),
			"latitude":    40.6892,
			"longitude":   -74.0445,
			"status":      "missed",
		},
		{
			"client_name": "Sarah Wilson",
			"shift_start": tomorrowMorning.Format("2006-01-02 15:04:05"),
			"shift_end":   tomorrowMorningEnd.Format("2006-01-02 15:04:05"),
			"latitude":    40.7831,
			"longitude":   -73.9712,
			"status":      "upcoming",
		},
	}

	for _, schedule := range schedules {
		result, err := DB.Exec(`
			INSERT INTO schedules (client_name, shift_start, shift_end, latitude, longitude, status)
			VALUES (?, ?, ?, ?, ?, ?)`,
			schedule["client_name"], schedule["shift_start"], schedule["shift_end"],
			schedule["latitude"], schedule["longitude"], schedule["status"])
		if err != nil {
			log.Printf("Failed to insert schedule: %v", err)
			continue
		}

		scheduleID, _ := result.LastInsertId()

		// Sample tasks for each schedule
		tasks := []string{
			"Assist with morning medication",
			"Help with personal hygiene",
			"Prepare light meal",
			"Check vital signs",
			"Light housekeeping",
		}

		for _, task := range tasks {
			_, err := DB.Exec(`
				INSERT INTO tasks (schedule_id, description, status)
				VALUES (?, ?, 'pending')`,
				scheduleID, task)
			if err != nil {
				log.Printf("Failed to insert task: %v", err)
			}
		}

		// Create visit entry for each schedule
		_, err = DB.Exec(`
			INSERT INTO visits (schedule_id)
			VALUES (?)`, scheduleID)
		if err != nil {
			log.Printf("Failed to create visit entry: %v", err)
		}

		// Sample activities for each schedule
		activities := []map[string]interface{}{
			{
				"title":       "Room Cleaning",
				"description": "Clean and organize the client's living room and bedroom",
				"is_resolved": false,
				"reason":      "",
			},
			{
				"title":       "Medication Check",
				"description": "Verify medication schedule and ensure proper dosage",
				"is_resolved": true,
				"reason":      "",
			},
			{
				"title":       "Meal Preparation",
				"description": "Prepare healthy lunch according to dietary requirements",
				"is_resolved": false,
				"reason":      "Client was not hungry at the time",
			},
		}

		for _, activity := range activities {
			_, err := DB.Exec(`
				INSERT INTO activities (schedule_id, title, description, is_resolved, reason)
				VALUES (?, ?, ?, ?, ?)`,
				scheduleID, activity["title"], activity["description"], 
				activity["is_resolved"], activity["reason"])
			if err != nil {
				log.Printf("Failed to insert activity: %v", err)
			}
		}
	}

	log.Println("Minimal fallback data seeded successfully")
}

// Close closes the database connection
func Close() {
	if DB != nil {
		DB.Close()
	}
} 