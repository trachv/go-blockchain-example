package main

import (
    "blockchain-portfolio/internal/blockchain"
    "blockchain-portfolio/internal/config"
    "blockchain-portfolio/internal/db"
    "blockchain-portfolio/internal/tracker"
    "fmt"
    "log"
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

    address := cfg.WalletAddress
    token := cfg.TokenAddress

    // Get balance
    balance, err := client.GetERC20Balance(token, address)
    if err != nil {
        log.Fatalf("Failed to get ERC20 balance: %v", err)
    }

    // Record balance in PostgreSQL
    err = track.RecordBalance(address, token, balance.Int64())
    if err != nil {
        log.Fatalf("Failed to record balance: %v", err)
    }

    fmt.Printf("Balance recorded: %s for address %s\n", balance.String(), address)
}
