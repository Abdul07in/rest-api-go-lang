package database

import (
	"database/sql"
	"fmt"
	"log"
	"student-api/internal/config"
)

// Initialize sets up the database and required tables
func Initialize(cfg *config.Config) (*sql.DB, error) {
	// First connect without database name to create the database if it doesn't exist
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort)
	
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("error connecting to MySQL: %w", err)
	}
	defer db.Close()

	// Create database if it doesn't exist
	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", cfg.DBName))
	if err != nil {
		return nil, fmt.Errorf("error creating database: %w", err)
	}

	// Connect to the specific database
	db, err = sql.Open("mysql", cfg.GetDSN())
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	// Create students table if it doesn't exist
	err = createTables(db)
	if err != nil {
		return nil, fmt.Errorf("error creating tables: %w", err)
	}

	log.Println("Database and tables initialized successfully")
	return db, nil
}

func createTables(db *sql.DB) error {
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS students (
		id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
		first_name VARCHAR(50) NOT NULL,
		last_name VARCHAR(50) NOT NULL,
		email VARCHAR(100) NOT NULL UNIQUE,
		age INT NOT NULL,
		grade DECIMAL(4,2) NOT NULL,
		created_at TIMESTAMP NOT NULL,
		updated_at TIMESTAMP NOT NULL,
		INDEX idx_email (email)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;`

	_, err := db.Exec(createTableSQL)
	return err
}
