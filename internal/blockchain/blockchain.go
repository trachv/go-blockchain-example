package blockchain

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)
import erc20 "blockchain-portfolio/internal/contracts/erc20"

// Client represents a client for interacting with the blockchain
type Client struct {
	eth     *ethclient.Client
	erc20ABI abi.ABI
}

// NewClient initializes a new blockchain client with the provided RPC URL
func NewClient(rpcURL string) (*Client, error) {
	eth, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, err
	}

	erc20ABI, err := loadERC20ABI()
	if err != nil {
		return nil, err
	}

	return &Client{eth: eth, erc20ABI: erc20ABI}, nil
}

// loadERC20ABI loads the ERC20 contract ABI from a JSON file
func loadERC20ABI() (abi.ABI, error) {
	file, err := os.Open("abis/erc20.json")
	if err != nil {
		return abi.ABI{}, err
	}
	defer file.Close()

	var contractABI abi.ABI
	if err := json.NewDecoder(file).Decode(&contractABI); err != nil {
		return abi.ABI{}, err
	}

	return contractABI, nil
}

// GetBalance retrieves the balance of the specified address
func (c *Client) GetBalance(address string) (*big.Int, error) {
	account := common.HexToAddress(address)
	balance, err := c.eth.BalanceAt(context.Background(), account, nil)
	if err != nil {
		return nil, err
	}
	return balance, nil
}

// GetERC20Balance retrieves the balance of an ERC-20 token for a specific address
func (c *Client) GetERC20Balance(tokenAddress, walletAddress string) (*big.Int, error) {
    token := common.HexToAddress(tokenAddress)
    wallet := common.HexToAddress(walletAddress)

    tokenContract, err := erc20.NewErc20(token, c.eth)
    if err != nil {
        return nil, fmt.Errorf("failed to instantiate token contract: %w", err)
    }

    balance, err := tokenContract.BalanceOf(&bind.CallOpts{Context: context.Background()}, wallet)
    if err != nil {
        return nil, fmt.Errorf("failed to retrieve token balance: %w", err)
    }

    return balance, nil
}

