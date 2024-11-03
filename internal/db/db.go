package db

import (
    "database/sql"
    "fmt"
    "strings"

    _ "github.com/lib/pq"
)

// DBManager manages the database connection
type DBManager struct {
    db *sql.DB
}

// NewDBManager connects to PostgreSQL and creates a new database if it doesn’t exist
func NewDBManager(host, port, user, password, dbName string) (*DBManager, error) {
    connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s sslmode=disable", host, port, user, password)
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        return nil, fmt.Errorf("failed to connect to PostgreSQL: %w", err)
    }

    // Check the connection
    if err := db.Ping(); err != nil {
        db.Close()
        return nil, fmt.Errorf("failed to ping PostgreSQL: %w", err)
    }

    // Try to create the database if it doesn't exist
    _, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName))
    if err != nil && !isDatabaseExistsError(err) {
        db.Close()
        return nil, fmt.Errorf("failed to create database: %w", err)
    }

    // Connect to the specific database
    db, err = sql.Open("postgres", fmt.Sprintf("%s dbname=%s", connStr, dbName))
    if err != nil {
        return nil, fmt.Errorf("failed to connect to the database: %w", err)
    }

    // Check the connection to the specific database
    if err := db.Ping(); err != nil {
        db.Close()
        return nil, fmt.Errorf("failed to ping database %s: %w", dbName, err)
    }

    return &DBManager{db: db}, nil
}

// isDatabaseExistsError checks if the error indicates the database already exists
func isDatabaseExistsError(err error) bool {
    return strings.Contains(err.Error(), "already exists")
}

// DB returns the database connection
func (d *DBManager) DB() *sql.DB {
    return d.db
}

// Close closes the database connection
func (d *DBManager) Close() error {
    if d.db != nil {
        return d.db.Close()
    }
    return nil
}
