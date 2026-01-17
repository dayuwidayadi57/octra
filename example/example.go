package main

import (
	"fmt"
	"github.com/dayuwidayadi57/octra/client"
)

func main() {
	password := "your-password"

	fmt.Println("🆕 Generating New Wallet...")
	addr, pub, priv, err := client.GenerateNewKeyPair()
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}

	fmt.Printf("✅ Address: %s\n", addr)
	fmt.Printf("✅ Public Key: %s\n", pub)
	fmt.Printf("✅ Private Key (Seed): %s\n", priv)
	fmt.Println("--------------------------------")

	fmt.Println("🔒 Encrypting Wallet to Keystore...")
	keystoreJSON, err := client.EncryptWallet(priv, password)
	if err != nil {
		fmt.Printf("❌ Encryption Error: %v\n", err)
		return
	}
	fmt.Printf("📄 Keystore Result:\n%s\n", keystoreJSON)
	fmt.Println("--------------------------------")

	fmt.Println("🔓 Decrypting Wallet...")
	decryptedPriv, err := client.DecryptWallet(keystoreJSON, password)
	if err != nil {
		fmt.Printf("❌ Decryption Error: %v\n", err)
		return
	}

	if priv == decryptedPriv {
		fmt.Println("✅ Success: Private key matches after decryption!")
		fmt.Printf("🔑 Recovered Priv: %s\n", decryptedPriv)
	} else {
		fmt.Println("❌ Mismatch: Recovery failed.")
	}
}
