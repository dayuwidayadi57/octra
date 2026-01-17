package main

import (
	"context"
	"fmt"
	"time"

	"github.com/dayuwidayadi57/octra/client"
)

func main() {
	octraClient := client.NewClient("https://octra.network")
	addr := "octDErj3QmjQsAacwzcFpt1a3EGiuZJBx3vY8Mc5o4XWUjh"

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	stats, err := octraClient.GetStats(ctx, addr)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}

	fmt.Printf("\n📊 WALLET ANALYTICS: %s\n", addr)
	fmt.Println("────────────────────────────────────────────────────────────────")
	fmt.Printf("Total Outbound : %12s OCT  (%s Atoms)\n", client.FromAtoms(stats.TotalOut), stats.TotalOut.String())
	fmt.Printf("Total Inbound  : %12s OCT  (%s Atoms)\n", client.FromAtoms(stats.TotalIn), stats.TotalIn.String())
	fmt.Printf("Total Activity : %12d Transactions\n", stats.TxCount)
	fmt.Println("────────────────────────────────────────────────────────────────")
}

