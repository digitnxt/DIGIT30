package database

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

// InitDB initializes the PostgreSQL database connection.
func InitDB(connStr string) {
	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	if err = DB.Ping(); err != nil {
		log.Fatalf("Cannot reach database: %v", err)
	}
}

// CreateTables creates all required tables if they do not already exist.
func CreateTables(db *sql.DB) {
	createAccountsTableSQL := `
	CREATE TABLE IF NOT EXISTS accounts (
		id SERIAL PRIMARY KEY,
		accountname VARCHAR(255) UNIQUE NOT NULL,
		admin_email VARCHAR(255),
		admin_phone VARCHAR(255),
		config JSONB NOT NULL,
		created_at TIMESTAMPTZ DEFAULT NOW()
	);
	`
	_, err := db.Exec(createAccountsTableSQL)
	if err != nil {
		log.Fatalf("Failed to create accounts table: %v", err)
	}
	log.Println("Accounts table ensured.")
}
