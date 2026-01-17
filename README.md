🧬 Octra Go Client SDK
The Octra Client is a high-performance Go library designed to interface with the Octra Network RPC. It provides comprehensive tools for key management, secure wallet storage, and full transaction lifecycle management (from creation to confirmation).

✨ Key Features
 * Key Management: Generate new Ed25519 keypairs and derive standard oct addresses.
 * Secure Keystores: Industry-standard wallet encryption using Scrypt (KDF) and AES-256-GCM.
 * Transaction Lifecycle: Tools to sign, broadcast, and monitor transactions until they are confirmed in an Epoch.
 * Precision Handling: Built-in utilities to handle OCT atoms (6-decimal precision) using math/big.
 * Async Operations: Full support for context.Context for managing timeouts and cancellations during RPC calls.

🚀 Quick Start
1. Initialize Client
import "local/osm/client"

oc := client.NewClient("https://rpc.octra.network")

2. Wallet Operations
Generate a new wallet or encrypt an existing private key for safe storage.
// Generate new keypair
address, pubKey, privKey, _ := client.GenerateNewKeyPair()

// Encrypt for storage
jsonKeystore, _ := client.EncryptWallet(privKey, "strong_password")

// Decrypt later
originalPrivKey, _ := client.DecryptWallet(jsonKeystore, "strong_password")

3. Send Transaction
A complete workflow to send OCT across the network.
// 1. Get current balance and next nonce
balance, _ := oc.GetBalance(ctx, senderAddr)
nonce, _ := oc.GetNextNonce(ctx, senderAddr)

// 2. Construct TX
tx := client.Transaction{
    From:      senderAddr,
    To:        "octDestinationAddress...",
    Amount:    client.ToAtoms(1.25).String(), // Converts 1.25 OCT to Atoms
    Nonce:     nonce,
    Timestamp: json.Number(fmt.Sprintf("%d", time.Now().Unix())),
}

// 3. Sign and Broadcast
signedTx, _ := client.SignTransaction(tx, privKeyB64)
res, _ := oc.SendTransaction(ctx, signedTx)

// 4. Wait for network confirmation
txHash := res["tx_hash"].(string)
confirmedTx, _ := oc.WaitTransaction(ctx, txHash, 1 * time.Minute)
fmt.Printf("Transaction confirmed in Epoch: %v\n", confirmedTx["epoch"])

🛠 API Reference
Cryptography & Utility
| Function | Description |
|---|---|
| PublicKeyToAddress | Converts raw public key bytes to oct Base58 address. |
| GenerateNewKeyPair | Creates a fresh set of Address, PubKey, and Seed. |
| ToAtoms | Converts float64 OCT values to big.Int Atoms. |
Network RPC
| Method | Description |
|---|---|
| GetBalance | Retrieves balance and nonce info for an address. |
| SendTransaction | Broadcasts a signed transaction to the network. |
| GetTransaction | Fetches details of a specific transaction by hash. |
| WaitTransaction | Polls the network until a transaction is confirmed or timed out. |

🔒 Security Specifications
 * Signature: Ed25519 (Edwards-curve Digital Signature Algorithm).
 * Encryption: AES-256-GCM (Authenticated Encryption).
 * KDF: Scrypt (N=32768, r=8, p=1).
 * Address: oct prefix + Base58 (SHA-256 of Public Key).
Maintained by Octra Community Team
License: MIT
