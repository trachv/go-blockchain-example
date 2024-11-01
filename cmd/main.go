package main

import (
    "blockchain-portfolio/internal/blockchain"
    "blockchain-portfolio/internal/config"
    "blockchain-portfolio/internal/db"
    "blockchain-portfolio/internal/tracker"
    "fmt"
    "log"
    "time"
)

func main() {
    // Load config
    cfg, err := config.LoadConfig()
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }

    // Connect to PostgreSQL and create the database if it doesn’t exist
    dbManager, err := db.NewDBManager("localhost", cfg.DbPort, cfg.DbUser, cfg.DbPassword, cfg.DbName)
    if err != nil {
        log.Fatalf("Failed to connect or create the database: %v", err)
    }
    defer dbManager.DB().Close()

    // Connect to blockchain client
    client, err := blockchain.NewClient(cfg.RpcURL)
    if err != nil {
        log.Fatalf("Failed to connect to blockchain: %v", err)
    }

    // Initialize Tracker with the database connection
    track := tracker.NewTracker(dbManager.DB())

    // Define the interval for balance checks (e.g., every 10 minutes)
    interval := time.Minute * 1
    ticker := time.NewTicker(interval)
    defer ticker.Stop()

    // Channel to handle graceful shutdown
    done := make(chan bool)

    // Goroutine to execute balance checks periodically
    go func() {
        for {
            select {
            case <-done:
                return
            case <-ticker.C:
                // Fetch balance and record it in the database
                checkAndRecordBalance(client, track, cfg.WalletAddress, cfg.TokenAddress)
            }
        }
    }()

    // Block the main thread until an interrupt signal is received
    select {}
}

// checkAndRecordBalance fetches and records the balance of the specified address and token
func checkAndRecordBalance(client *blockchain.Client, track *tracker.Tracker, address, token string) {
    // Get balance
    balance, err := client.GetERC20Balance(token, address)
    if err != nil {
        log.Printf("Failed to get ERC20 balance: %v", err)
        return
    }

    // Record balance in PostgreSQL
    err = track.RecordBalance(address, token, balance.Int64())
    if err != nil {
        log.Printf("Failed to record balance: %v", err)
        return
    }

    fmt.Printf("Balance recorded: %s for address %s\n", balance.String(), address)
}
