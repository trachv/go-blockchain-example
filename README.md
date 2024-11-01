# Blockchain Portfolio

A project for blockchain interaction and ERC-20 contract integration. Includes functions for retrieving Ethereum balances and ERC-20 token balances.

## Contents

- [Technologies](#technologies)
- [Project Structure](#project-structure)
- [Installation and Setup](#installation-and-setup)
- [Usage](#usage)
- [Examples](#examples)

## Technologies

- Go
- Ethereum Go SDK (go-ethereum)
- JSON for configuration management

## Project Structure

```plaintext
blockchain-portfolio/
├── cmd/                    # Main entry file (main.go) for launching the project
├── internal/
│   ├── blockchain/         # Blockchain interaction logic
│   │   ├── blockchain.go
│   ├── config/             # Project configuration
│   │   ├── config.go
├── abis/                   # ABI files for contracts
├── config.json             # Primary configuration file
├── config.example.json     # Example configuration file
└── README.md               # Project documentation
```

## Installation and Setup

1. Clone the repository:
   ```bash
   git clone your-repository-url
   cd blockchain-portfolio
   ```
2. Install project dependencies:
   ```bash
   go mod download
   ```
3. Create a `config.json` file, using `config.example.json` as a reference:
   ```json
   {
       "rpcURL": "https://mainnet.infura.io/v3/YOUR_INFURA_KEY",
       "walletAddress": "YOUR_WALLET_ADDRESS",
       "tokenAddress": "YOUR_TOKEN_ADDRESS"
   }
   ```
4. Make sure `config.json` is added to `.gitignore` to avoid accidental upload to the repository.

## Usage

Run the main project file:

```bash
go run cmd/main.go
```

## Examples

Example output when running the program:

```
ETH Balance of 0x...: 1000000000000000000
ERC20 Token Balance of 0x... for token 0x...: 5000000
```

## Important Notes

- **API keys and sensitive data:** Ensure that all sensitive information, such as API keys and private keys, is protected and not pushed to the repository.
- **License:** Specify a license for your project (e.g., MIT, Apache 2.0) and add a LICENSE file if needed.
