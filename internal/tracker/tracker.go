package tracker

import (
    "database/sql"
    "fmt"
    "math/big"
)

type Tracker struct {
    db *sql.DB
}

func NewTracker(db *sql.DB) *Tracker {
    tracker := &Tracker{db: db}

    if err := tracker.initializeTable(); err != nil {
        fmt.Printf("Failed to initialize table: %v\n", err)
    }

    return tracker
}

func (t *Tracker) initializeTable() error {
    query := `
        CREATE TABLE IF NOT EXISTS balances (
            id SERIAL PRIMARY KEY,
            network VARCHAR(50) NOT NULL,
            address VARCHAR(42) NOT NULL,
            token VARCHAR(42) NOT NULL,
            balance NUMERIC,
            timestamp TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
        );
    `
    _, err := t.db.Exec(query)
    if err != nil {
        return fmt.Errorf("failed to initialize table: %w", err)
    }
    return nil
}

func (t *Tracker) RecordBalance(network, address, token string, balance *big.Int) error {
    balanceStr := balance.String()

    query := "INSERT INTO balances (network, address, token, balance) VALUES ($1, $2, $3, $4)"
    _, err := t.db.Exec(query, network, address, token, balanceStr)
    if err != nil {
        return fmt.Errorf("failed to record balance: %w", err)
    }
    return nil
}
