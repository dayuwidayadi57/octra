// client/history.go
package client

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
	"strings"
)

type TransactionHistory struct {
	Hash      string      `json:"hash"`
	Epoch     int         `json:"epoch"`
	From      string      `json:"from"`
	To        string      `json:"to"`
	Amount    string      `json:"amount"`
	Timestamp json.Number `json:"timestamp"`
	Status    string      `json:"status"`
}

type WalletStats struct {
	TotalIn  *big.Int
	TotalOut *big.Int
	TxCount  int
}

func (c *OctraClient) GetHistory(ctx context.Context, address string, limit int) ([]TransactionHistory, error) {
	path := fmt.Sprintf("/address/%s?limit=%d", address, limit)
	data, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var wrapper struct {
		RecentTransactions []struct {
			Hash  string `json:"hash"`
			Epoch int    `json:"epoch"`
		} `json:"recent_transactions"`
	}

	if err := json.Unmarshal(data, &wrapper); err != nil {
		return nil, err
	}

	var history []TransactionHistory
	for _, item := range wrapper.RecentTransactions {
		txData, err := c.GetTransaction(ctx, item.Hash)
		if err != nil {
			continue
		}

		parsed, ok := txData["parsed_tx"].(map[string]interface{})
		if !ok {
			continue
		}

		history = append(history, TransactionHistory{
			Hash:      item.Hash,
			Epoch:     item.Epoch,
			From:      fmt.Sprintf("%v", parsed["from"]),
			To:        fmt.Sprintf("%v", parsed["to"]),
			Amount:    fmt.Sprintf("%v", parsed["amount"]),
			Timestamp: json.Number(fmt.Sprintf("%v", parsed["timestamp"])),
			Status:    "confirmed",
		})
	}

	return history, nil
}

func (c *OctraClient) GetStats(ctx context.Context, address string) (*WalletStats, error) {
	history, err := c.GetHistory(ctx, address, 50)
	if err != nil {
		return nil, err
	}

	stats := &WalletStats{
		TotalIn:  big.NewInt(0),
		TotalOut: big.NewInt(0),
		TxCount:  len(history),
	}

	for _, tx := range history {
		amtStr := strings.Fields(tx.Amount)[0]
		amtFloat, _ := strconv.ParseFloat(amtStr, 64)
		amtAtoms := ToAtoms(amtFloat)

		if strings.EqualFold(tx.From, address) {
			stats.TotalOut.Add(stats.TotalOut, amtAtoms)
		} else if strings.EqualFold(tx.To, address) {
			stats.TotalIn.Add(stats.TotalIn, amtAtoms)
		}
	}

	return stats, nil
}
