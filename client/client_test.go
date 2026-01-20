package client

import (
	"encoding/json"
	"fmt"
	"testing"
	"strings"
)

func TestOctraSDK_FullFlow(t *testing.T) {
	fmt.Println("\n\033[1;36m>>> STARTING FULL DEBUG TEST <<<\033[0m")

	// 1. Key Generation Debug
	addr, pub, priv, err := GenerateNewKeyPair()
	if err != nil {
		t.Fatalf("Failed to gen key: %v", err)
	}
	fmt.Printf("\033[1;34m[DEBUG]\033[0m Address: %s\n", addr)
	fmt.Printf("\033[1;34m[DEBUG]\033[0m PubKey (B64): %s\n", pub)
	fmt.Printf("\033[1;34m[DEBUG]\033[0m PrivKey (Seed B64): %s\n", priv)

	// 2. Wallet Encryption Debug
	password := "octra123"
	keystoreJSON, err := EncryptWallet(priv, password)
	if err != nil {
		t.Fatalf("Encryption failed: %v", err)
	}
	fmt.Printf("\033[1;34m[DEBUG]\033[0m Keystore Created: %s\n", keystoreJSON)

	decryptedPriv, err := DecryptWallet(keystoreJSON, password)
	if err != nil || decryptedPriv != priv {
		t.Fatalf("Decryption mismatch!")
	}
	fmt.Println("✅ Keystore Encryption/Decryption Verified")

	// 3. Transaction Signing Debug (Standard OTX-1)
	tx := Transaction{
		From:      addr,
		To:        "octD4RxTBurSjSUp3mdM3eAH4Qo4GyU3Ay29oTez3eWVuWV",
		Amount:    "5000000",
		Nonce:     10,
		Timestamp: json.Number("1737273600"),
		Message:   "Test Debug Message",
	}

	signed, err := SignTransaction(tx, priv)
	if err != nil {
		t.Fatalf("Signing failed: %v", err)
	}

	fmt.Printf("\033[1;34m[DEBUG]\033[0m Canonical Data (Signed): %s\n", signed.Raw)
	fmt.Printf("\033[1;34m[DEBUG]\033[0m Signature: %s\n", signed.Signature)

	// Verify that Message is NOT in the signed Raw data (OTX-1 Rule)
	if strings.Contains(signed.Raw, "Message") || strings.Contains(signed.Raw, "Test Debug Message") {
		t.Errorf("FAIL: Message should NOT be in the signed canonical payload")
	} else {
		fmt.Println("✅ OTX-1 Compliance: Message excluded from signature payload")
	}

	// 4. ToMap Debug (Checking if message is included for node broadcast)
	finalMap := signed.ToMap()
	mapJSON, _ := json.MarshalIndent(finalMap, "", "  ")
	fmt.Printf("\033[1;34m[DEBUG]\033[0m Final JSON to Node:\n%s\n", string(mapJSON))

	if finalMap["message"] != "Test Debug Message" {
		t.Errorf("FAIL: Message should be in the final map for node")
	} else {
		fmt.Println("✅ Network Compliance: Message included in final broadcast map")
	}

	fmt.Println("\033[1;32m>>> ALL TESTS PASSED <<<\033[0m")
}

func TestAtomsConversion(t *testing.T) {
	amount := 1.25
	atoms := ToAtoms(amount)
	if atoms.String() != "1250000" {
		t.Errorf("ToAtoms failed: expected 1250000, got %s", atoms.String())
	}
	
	backToOCT := FromAtoms(atoms)
	if !strings.HasPrefix(backToOCT, "1.25") {
		t.Errorf("FromAtoms failed: expected 1.25, got %s", backToOCT)
	}
	fmt.Printf("\033[1;34m[DEBUG]\033[0m Conversion: %.2f OCT -> %s Atoms -> %s OCT\n", amount, atoms.String(), backToOCT)
}
