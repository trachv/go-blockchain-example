package config

import (
    "encoding/json"
    "fmt"
    "os"
)

type NetworkConfig struct {
    Name   string `json:"name"`
    RpcURL string `json:"rpcURL"`
}

type Config struct {
    Networks        []NetworkConfig `json:"networks"`
    WalletAddresses []string        `json:"walletAddresses"`
    TokenAddresses  []string        `json:"tokenAddresses"`
    DbHost          string          `json:"dbHost"`
    DbUser          string          `json:"dbUser"`
    DbPassword      string          `json:"dbPassword"`
    DbName          string          `json:"dbName"`
    DbPort          string          `json:"dbPort"`
}

func LoadConfig() (*Config, error) {
    file, err := os.Open("config.json")
    if err != nil {
        return nil, fmt.Errorf("failed to open config file: %w", err)
    }
    defer file.Close()

    config := &Config{}
    if err := json.NewDecoder(file).Decode(config); err != nil {
        return nil, fmt.Errorf("failed to decode config JSON: %w", err)
    }
    return config, nil
}
