package main

import (
    "fmt"
    "log"
    "blockchain-portfolio/internal/blockchain"
    "blockchain-portfolio/internal/config"
)

func main() {
    cfg, err := config.LoadConfig()
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }

    client, err := blockchain.NewClient(cfg.RpcURL)
    if err != nil {
        log.Fatalf("Failed to connect to the blockchain: %v", err)
    }

    balance, err := client.GetBalance(cfg.WalletAddress)
    if err != nil {
        log.Fatalf("Failed to get balance: %v", err)
    }
    fmt.Printf("ETH Balance of %s: %s\n", cfg.WalletAddress, balance.String())

    erc20Balance, err := client.GetERC20Balance(cfg.TokenAddress, cfg.WalletAddress)
    if err != nil {
        log.Fatalf("Failed to get ERC20 balance: %v", err)
    }
    fmt.Printf("ERC20 Token Balance of %s for token %s: %s\n", cfg.WalletAddress, cfg.TokenAddress, erc20Balance.String())
}
