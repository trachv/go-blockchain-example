package balance_checker

import (
    "blockchain-portfolio/internal/blockchain"
    "blockchain-portfolio/internal/tracker"
    "fmt"
    "log"
    "sync"
)

type BalanceChecker struct {
    networks []blockchain.NetworkClient
    tracker  *tracker.Tracker
}

func NewBalanceChecker(networks []blockchain.NetworkClient, tracker *tracker.Tracker) *BalanceChecker {
    return &BalanceChecker{
        networks: networks,
        tracker:  tracker,
    }
}

func (bc *BalanceChecker) CheckAndRecordBalances(addresses, tokens []string) {
    var wg sync.WaitGroup

    for _, network := range bc.networks {
        for _, address := range addresses {
            for _, token := range tokens {
                wg.Add(1)
                go func(nw blockchain.NetworkClient, addr, tok string) {
                    defer wg.Done()
                    bc.checkAndRecordBalance(nw, addr, tok)
                }(network, address, token)
            }
        }
    }

    wg.Wait()
}

func (bc *BalanceChecker) checkAndRecordBalance(network blockchain.NetworkClient, address, token string) {
    balance, err := network.Client.GetERC20Balance(token, address)
    if err != nil {
        log.Printf("[%s] Failed to get ERC20 balance for address %s and token %s: %v", network.Name, address, token, err)
        return
    }

    err = bc.tracker.RecordBalance(network.Name, address, token, balance)
    if err != nil {
        log.Printf("[%s] Failed to record balance for address %s and token %s: %v", network.Name, address, token, err)
        return
    }

    fmt.Printf("[%s] Balance recorded: %s for address %s and token %s\n", network.Name, balance.String(), address, token)
}
