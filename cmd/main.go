package main

import (
    "blockchain-portfolio/internal/config"
    "blockchain-portfolio/internal/db"
    "blockchain-portfolio/pkg/app"
    "log"
)

func main() {
    cfg, err := config.LoadConfig()
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }

    dbManager, err := db.NewDBManager(cfg.DbHost, cfg.DbPort, cfg.DbUser, cfg.DbPassword, cfg.DbName)
    if err != nil {
        log.Fatalf("Failed to initialize database manager: %v", err)
    }

    application, err := app.NewApp(cfg, dbManager)
    if err != nil {
        log.Fatalf("Failed to initialize application: %v", err)
    }

    application.Start()
}
