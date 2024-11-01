package tracker

import (
    "database/sql"
    "fmt"
    "time"
)

type Tracker struct {
    db *sql.DB
}

func NewTracker(db *sql.DB) *Tracker {
    tracker := &Tracker{db: db}

    // Initialize table if it doesn’t exist
    err := tracker.initializeTable()
    if err != nil {
        fmt.Printf("Failed to initialize table: %v\n", err)
    }

    return tracker
}

func (t *Tracker) initializeTable() error {
    query := `
        CREATE TABLE IF NOT EXISTS balances (
            id SERIAL PRIMARY KEY,
            address VARCHAR(42),
            token VARCHAR(42),
            balance NUMERIC,
            timestamp TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
        );
    `
    _, err := t.db.Exec(query)
    return err
}

func (t *Tracker) RecordBalance(address, token string, balance int64) error {
    query := "INSERT INTO balances (address, token, balance) VALUES ($1, $2, $3)"
    _, err := t.db.Exec(query, address, token, balance)
    if err != nil {
        return fmt.Errorf("failed to record balance: %w", err)
    }
    return nil
}

func (t *Tracker) GetBalanceHistory(address, token string) ([]BalanceRecord, error) {
    query := "SELECT balance, timestamp FROM balances WHERE address = $1 AND token = $2 ORDER BY timestamp DESC"
    rows, err := t.db.Query(query, address, token)
    if err != nil {
        return nil, fmt.Errorf("failed to query balance history: %w", err)
    }
    defer rows.Close()

    var records []BalanceRecord
    for rows.Next() {
        var record BalanceRecord
        if err := rows.Scan(&record.Balance, &record.Timestamp); err != nil {
            return nil, err
        }
        records = append(records, record)
    }

    return records, nil
}

type BalanceRecord struct {
    Balance   int64
    Timestamp time.Time
}
