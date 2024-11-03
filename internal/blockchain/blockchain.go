package blockchain

import (
    "context"
    "fmt"
    "math/big"
    "time"

    "github.com/ethereum/go-ethereum/accounts/abi/bind"
    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/ethclient"
    erc20 "blockchain-portfolio/internal/contracts/erc20"
)

type BlockchainClient interface {
    GetBalance(address string) (*big.Int, error)
    GetERC20Balance(tokenAddress, walletAddress string) (*big.Int, error)
    Close()
}

type Client struct {
    eth     *ethclient.Client
    timeout time.Duration // Timeout for blockchain calls
}

func NewClient(rpcURL string) (*Client, error) {
    eth, err := ethclient.Dial(rpcURL)
    if err != nil {
        return nil, fmt.Errorf("failed to connect to Ethereum client: %w", err)
    }

    return &Client{
        eth:     eth,
        timeout: 10 * time.Second, // Set a default timeout
    }, nil
}

func (c *Client) GetBalance(address string) (*big.Int, error) {
    account := common.HexToAddress(address)
    ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
    defer cancel()

    balance, err := c.eth.BalanceAt(ctx, account, nil)
    if err != nil {
        return nil, fmt.Errorf("failed to retrieve balance for address %s: %w", address, err)
    }
    return balance, nil
}

func (c *Client) GetERC20Balance(tokenAddress, walletAddress string) (*big.Int, error) {
    token := common.HexToAddress(tokenAddress)
    wallet := common.HexToAddress(walletAddress)

    // Initialize context with timeout
    ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
    defer cancel()

    // Instantiate the ERC20 token contract
    tokenContract, err := erc20.NewErc20(token, c.eth)
    if err != nil {
        return nil, fmt.Errorf("failed to instantiate token contract at address %s: %w", tokenAddress, err)
    }

    // Retrieve balance of the wallet address
    balance, err := tokenContract.BalanceOf(&bind.CallOpts{Context: ctx}, wallet)
    if err != nil {
        return nil, fmt.Errorf("failed to retrieve ERC20 token balance for address %s: %w", walletAddress, err)
    }

    return balance, nil
}

func (c *Client) Close() {
    if c.eth != nil {
        c.eth.Close()
    }
}

// New struct for network client
type NetworkClient struct {
    Name   string
    Client *Client
}

func NewNetworkClient(name, rpcURL string) (*NetworkClient, error) {
    client, err := NewClient(rpcURL)
    if err != nil {
        return nil, err
    }
    return &NetworkClient{
        Name:   name,
        Client: client,
    }, nil
}
