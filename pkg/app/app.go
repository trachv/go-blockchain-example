package app

import (
    "blockchain-portfolio/internal/balance_checker"
    "blockchain-portfolio/internal/blockchain"
    "blockchain-portfolio/internal/config"
    "blockchain-portfolio/internal/db"
    "blockchain-portfolio/internal/tracker"
    "fmt"
    "log"
    "os"
    "os/signal"
    "syscall"
    "time"
)

type App struct {
    config         *config.Config
    dbManager      *db.DBManager
    networks       []blockchain.NetworkClient
    balanceChecker *balance_checker.BalanceChecker
}

func NewApp(cfg *config.Config, dbManager *db.DBManager) (*App, error) {
    // Initialize network clients
    var networks []blockchain.NetworkClient
    for _, netCfg := range cfg.Networks {
        netClient, err := blockchain.NewNetworkClient(netCfg.Name, netCfg.RpcURL)
        if err != nil {
            return nil, fmt.Errorf("failed to initialize network client for %s: %w", netCfg.Name, err)
        }
        networks = append(networks, *netClient)
    }

    tracker := tracker.NewTracker(dbManager.DB())
    balanceChecker := balance_checker.NewBalanceChecker(networks, tracker)

    return &App{
        config:         cfg,
        dbManager:      dbManager,
        networks:       networks,
        balanceChecker: balanceChecker,
    }, nil
}

func (a *App) Start() {
    log.Println("Starting application...")

    interval := time.Minute * 1
    ticker := time.NewTicker(interval)
    defer ticker.Stop()

    done := make(chan os.Signal, 1)
    signal.Notify(done, os.Interrupt, syscall.SIGTERM)

    go func() {
        for range ticker.C {
            a.balanceChecker.CheckAndRecordBalances(a.config.WalletAddresses, a.config.TokenAddresses)
        }
    }()

    <-done
    log.Println("Shutting down application...")

    if err := a.dbManager.Close(); err != nil {
        log.Printf("Error closing database connection: %v", err)
    }

    // Close blockchain clients
    for _, network := range a.networks {
        network.Client.Close()
    }

    log.Println("Application stopped successfully.")
}
