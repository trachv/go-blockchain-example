// internal/config/config.go
package config

import (
    "encoding/json"
    "os"
)

type Config struct {
    RpcURL        string `json:"rpcURL"`
    WalletAddress string `json:"walletAddress"`
    TokenAddress  string `json:"tokenAddress"`
}

func LoadConfig() (*Config, error) {
    file, err := os.Open("config.json")
    if err != nil {
        return nil, err
    }
    defer file.Close()

    config := &Config{}
    if err := json.NewDecoder(file).Decode(config); err != nil {
        return nil, err
    }
    return config, nil
}