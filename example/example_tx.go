package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/dayuwidayadi57/octra/client"
	"github.com/dayuwidayadi57/osm15"	
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("❌ Error loading .env file")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 150*time.Second)
	defer cancel()

	myPrivateKeyB64 := os.Getenv("OCTRA_PRIVATE_KEY")
	rpcURL := os.Getenv("OCTRA_RPC_URL")
	domainName := os.Getenv("OCTRA_DOMAIN")

	octraClient := client.NewClient(rpcURL)
	addr, pubB64, _, _ := client.GenerateNewKeyPairFromPriv(myPrivateKeyB64)

	fmt.Printf("\033[1;34m[DEBUG]\033[0m Node RPC: %s\n", rpcURL)
	fmt.Printf("\033[1;34m[DEBUG]\033[0m Wallet Address: %s\n", addr)
	fmt.Printf("\033[1;34m[DEBUG]\033[0m Public Key (B64): %s\n", pubB64)

	// --- STEP 1: AUTHENTICATION VIA OSM-15 ---
	authData := osm15.TypedData{
		Domain: osm15.TypedDomain{
			Name:    domainName,
			Version: "1",
			ChainID: 1,
		},
		Types: map[string][]osm15.TypedMember{
			"Login": {{Name: "action", Type: "string"}, {Name: "user", Type: "string"}},
		},
		PrimaryType: "Login",
		Message: map[string]interface{}{
			"action": "Session Login",
			"user":   addr,
		},
	}

	signature, err := osm15.SignTypedData(authData, myPrivateKeyB64)
	if err != nil {
		fmt.Printf("❌ Signing Error: %v\n", err)
		return
	}
	fmt.Printf("\033[1;34m[DEBUG]\033[0m OSM-15 Signature: %s\n", signature)

	recoveredAddr, err := osm15.GetSignerAddress(authData, signature, pubB64)
	fmt.Printf("\033[1;34m[DEBUG]\033[0m Recovered Address: %s\n", recoveredAddr)

	if err != nil || recoveredAddr != addr {
		fmt.Println("❌ Auth Failed: Invalid Identity Signature")
		return
	}
	fmt.Println("✅ Identity Verified via OSM-15 Structured Data")

	// --- STEP 2: CHECK BALANCE ---
	bal, err := octraClient.GetBalance(ctx, addr)
	if err != nil {
		fmt.Printf("❌ RPC Balance Error: %v\n", err)
		return
	}
	fmt.Printf("📊 Balance: %s OCT (Raw: %s Atoms)\n", bal.Balance, bal.BalanceRaw)
	fmt.Printf("\033[1;34m[DEBUG]\033[0m Current Nonce: %d\n", bal.Nonce)

	var destinationAddr string
	var amountInOCT float64

	fmt.Println("-----------------------------")
	fmt.Print("📝 Input Destination Address: ")
	fmt.Scan(&destinationAddr)
	fmt.Print("💰 Input Amount (OCT): ")
	fmt.Scan(&amountInOCT)
	fmt.Println("-----------------------------")

	currentBal, _ := strconv.ParseFloat(bal.Balance, 64)
	if amountInOCT > currentBal {
		fmt.Printf("❌ Insufficient Balance: You have %.6f but trying to send %.6f\n", currentBal, amountInOCT)
		return
	}

	// --- STEP 3: TRANSACTION VIA CLIENT SDK ---
	nonce := bal.Nonce + 1
	atoms := client.ToAtoms(amountInOCT)
	ts := float64(time.Now().UnixNano()) / 1e9
	
	tx := client.Transaction{
		From:      addr,
		To:        destinationAddr,
		Amount:    atoms.String(),
		Nonce:     nonce,
		Timestamp: json.Number(strconv.FormatFloat(ts, 'f', -1, 64)),
	}

	signedTx, _ := client.SignTransaction(tx, myPrivateKeyB64)
	fmt.Printf("\033[1;34m[DEBUG]\033[0m Canonical Payload (OTX-1): %s\n", signedTx.Raw)
	fmt.Printf("\033[1;34m[DEBUG]\033[0m Tx Signature: %s\n", signedTx.Signature)

	fmt.Printf("🚀 Sending %.6f OCT to %s...\n", amountInOCT, destinationAddr)
	
	res, err := octraClient.SendTransaction(ctx, signedTx)
	if err != nil {
		fmt.Printf("❌ Broadcast Error: %v\n", err)
		return
	}

	// Debug Full Response dari Node
	resJSON, _ := json.MarshalIndent(res, "", "  ")
	fmt.Printf("\033[1;34m[DEBUG]\033[0m Node Response: %s\n", string(resJSON))

	txHash, ok := res["tx_hash"].(string)
	if !ok {
		fmt.Println("❌ Node Error: Transaction not accepted")
		return
	}

	fmt.Printf("🔗 Tx Hash: %s\n", txHash)
	fmt.Println("⏳ Waiting for confirmation (polling)...")

	result, err := octraClient.WaitTransaction(ctx, txHash, 120*time.Second)
	if err != nil {
		fmt.Printf("⚠️  Confirmation Alert: %v\n", err)
	} else {
		resultJSON, _ := json.MarshalIndent(result, "", "  ")
		fmt.Printf("\033[1;32m[SUCCESS]\033[0m Transaction Confirmed!\n%s\n", string(resultJSON))
	}
}

