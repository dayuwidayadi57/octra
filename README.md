# 🧬 Octra Go Client SDK

[![Network](https://img.shields.io/badge/Network-Octra-blueviolet)](https://octra.network)
[![Protocol](https://img.shields.io/badge/Protocol-OTX--1-success)](https://github.com/dayuwidayadi57/octra)
[![Team](https://img.shields.io/badge/Developed%20By-QiubitLabs-blue)](https://github.com/QiubitLabs)
[![Status](https://img.shields.io/badge/Status-Stable-blue)](https://github.com/dayuwidayadi57/octra)
[![License](https://img.shields.io/badge/License-MIT-yellow)](https://github.com/dayuwidayadi57/octra/blob/main/LICENSE)

The Octra Client is a high-performance Go library designed to interface with the Octra Network RPC. It provides comprehensive tools for key management, secure wallet storage, and full transaction lifecycle management (from creation to confirmation).

✨ Key Features
- **Key Management**: Generate new Ed25519 keypairs and derive standard oct addresses.
- **Secure Keystores**: Industry-standard wallet encryption using Scrypt (KDF) and AES-256-GCM.
- **Transaction Lifecycle**: Tools to sign, broadcast, and monitor transactions until they are confirmed in an Epoch.
- **Precision Handling**: Built-in utilities to handle OCT atoms (6-decimal precision) using math/big.
- **OSM-15 Support**: Standardized structured data signing for secure authentication.
- **Metadata Support**: Ability to include custom messages in transactions (OTX-1 Compliant).

🚀 Installation
```bash
go get github.com/dayuwidayadi57/octra
```

🚀 Quick Start

### 1. Initialize Client
```go
import "github.com/dayuwidayadi57/octra/client"

oc := client.NewClient("https://rpc.octra.network")
```

### 2. Wallet Operations
```go
// Generate new keypair
address, pubKey, privKey, _ := client.GenerateNewKeyPair()

// Encrypt for storage
jsonKeystore, _ := client.EncryptWallet(privKey, "strong_password")
```

### 3. Send Transaction
```go
// Get current balance and next nonce
balance, _ := oc.GetBalance(ctx, senderAddr)

// Construct TX with optional message
tx := client.Transaction{
    From:      senderAddr,
    To:        "octDestinationAddress...",
    Amount:    client.ToAtoms(1.25).String(),
    Nonce:     balance.Nonce + 1,
    Timestamp: json.Number(strconv.FormatFloat(float64(time.Now().Unix()), 'f', -1, 64)),
    Message:   "Payment for services",
}

// Sign and Broadcast
signedTx, _ := client.SignTransaction(tx, privKeyB64)
res, _ := oc.SendTransaction(ctx, signedTx)

// Wait for network confirmation
txHash := res["tx_hash"].(string)
confirmedTx, _ := oc.WaitTransaction(ctx, txHash, 2 * time.Minute)
```

🛠 API Reference

#### Cryptography & Utility
- **PublicKeyToAddress**: Converts raw public key bytes to oct Base58 address.
- **GenerateNewKeyPair**: Creates a fresh set of Address, PubKey, and Seed.
- **ToAtoms**: Converts float64 OCT values to big.Int Atoms.
- **SignTransaction**: Signs OTX-1 compliant transactions (isolates message from payload).

#### Network RPC
- **GetBalance**: Retrieves balance and nonce info for an address.
- **SendTransaction**: Broadcasts a signed transaction to the network.
- **WaitTransaction**: Polls the network until a transaction is confirmed or timed out.

🔒 Security Specifications
- **Signature**: Ed25519 (Edwards-curve Digital Signature Algorithm).
- **Encryption**: AES-256-GCM (Authenticated Encryption).
- **KDF**: Scrypt (N=32768, r=8, p=1).
- **Address**: oct prefix + Base58 (SHA-256 of Public Key).
- **Compliance**: Fully supports **OTX-1** and **OSM-15** standards.

---
Maintained by **dayuwidayadi57** & **Qiubit Team**
License: **MIT**
