package db

import (
    "database/sql"
    "fmt"
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

    // Create database if it doesn’t exist
    _, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName))
    if err != nil && !isDatabaseExistsError(err) {
        return nil, fmt.Errorf("failed to create database: %w", err)
    }

    // Connect to the created or existing database
    db, err = sql.Open("postgres", fmt.Sprintf("%s dbname=%s", connStr, dbName))
    if err != nil {
        return nil, fmt.Errorf("failed to connect to the database: %w", err)
    }

    return &DBManager{db: db}, nil
}

// isDatabaseExistsError checks if an error is due to an already existing database
func isDatabaseExistsError(err error) bool {
    return err.Error() == "ERROR: database already exists (SQLSTATE 42P04)"
}

// DB returns the database connection
func (d *DBManager) DB() *sql.DB {
    return d.db
}
