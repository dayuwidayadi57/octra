package main

import (
	"context"
	"fmt"
	"time"

	"github.com/dayuwidayadi57/octra/client"
)

func main() {
	octraClient := client.NewClient("https://octra.network")
	targetAddr := "octDErj3QmjQsAacwzcFpt1a3EGiuZJBx3vY8Mc5o4XWUjh"

	fmt.Printf("🔍 Fetching history for: %s...\n\n", targetAddr)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	history, err := octraClient.GetHistory(ctx, targetAddr, 5)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}

	if len(history) == 0 {
		fmt.Println("📭 No transactions found.")
		return
	}

	fmt.Println("📜 RECENT TRANSACTIONS")
	fmt.Println("──────────────────────────────────────────────────────────────────")
	for _, tx := range history {
		typeStr := "📥 IN "
		if tx.From == targetAddr {
			typeStr = "📤 OUT"
		}

		fmt.Printf("[%d] %s | %s | %s... | %s OCT\n",
			tx.Epoch,
			typeStr,
			tx.Hash[:12],
			tx.To[:10],
			tx.Amount,
		)
	}
	fmt.Println("──────────────────────────────────────────────────────────────────")
}

