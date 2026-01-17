/*
 * Octra Go SDK - Client & Wallet Library
 *
 * Copyright (c) 2026 Octra & Community Team
 * Licensed under the MIT License
 * * This package provides a high-level interface for interacting with the Octra Network.
 * It includes secure Ed25519 key management, military-grade wallet encryption (AES-256-GCM),
 * and complete transaction lifecycle management.
 *
 * Technical Specifications:
 * - Encryption: AES-256-GCM
 * - KDF: Scrypt (N=32768, r=8, p=1)
 * - Signatures: Ed25519
 * - Hashing: SHA-256
 */

package client

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/mr-tron/base58"
	"golang.org/x/crypto/scrypt"
)

type OctraClient struct {
	BaseURL    string
	HTTPClient *http.Client
}

type Keystore struct {
	Address string `json:"address"`
	Crypto  struct {
		Cipher     string `json:"cipher"`
		CipherText string `json:"ciphertext"`
		Salt       string `json:"salt"`
		Nonce      string `json:"nonce"`
	} `json:"crypto"`
}

type BalanceInfo struct {
	Address      string `json:"address"`
	Balance      string `json:"balance"`
	BalanceRaw   string `json:"balance_raw"`
	HasPublicKey bool   `json:"has_public_key"`
	Nonce        uint64 `json:"nonce"`
}

type Transaction struct {
	From      string      `json:"from"`
	To        string      `json:"to_"`
	Amount    string      `json:"amount"`
	Nonce     uint64      `json:"nonce"`
	OU        string      `json:"ou"`
	Timestamp json.Number `json:"timestamp"`
}

type SignedTransaction struct {
	Signature string      `json:"signature"`
	PublicKey string      `json:"public_key"`
	Tx        Transaction `json:"tx"`
	Raw       string      `json:"raw,omitempty"`
}

func NewClient(url string) *OctraClient {
	return &OctraClient{
		BaseURL:    strings.TrimSuffix(url, "/"),
		HTTPClient: &http.Client{Timeout: 30 * time.Second},
	}
}

func PublicKeyToAddress(publicKey []byte) string {
	hash := sha256.Sum256(publicKey)
	return "oct" + base58.Encode(hash[:])
}

func GenerateNewKeyPair() (string, string, string, error) {
	pub, priv, err := ed25519.GenerateKey(nil)
	if err != nil { return "", "", "", err }
	return PublicKeyToAddress(pub), base64.StdEncoding.EncodeToString(pub), base64.StdEncoding.EncodeToString(priv.Seed()), nil
}

func GenerateNewKeyPairFromPriv(privateKeyB64 string) (string, string, string, error) {
	seed, _ := base64.StdEncoding.DecodeString(privateKeyB64)
	privateKey := ed25519.NewKeyFromSeed(seed)
	pub := privateKey.Public().(ed25519.PublicKey)
	return PublicKeyToAddress(pub), base64.StdEncoding.EncodeToString(pub), privateKeyB64, nil
}

func FromAtoms(atoms *big.Int) string {
	f := new(big.Float).SetInt(atoms)
	f.Quo(f, big.NewFloat(1e6))
	return f.Text('f', 6)
}

func ToAtoms(amount float64) *big.Int {
	val := big.NewFloat(amount)
	multiplier := new(big.Float).SetFloat64(1e6)
	val.Mul(val, multiplier)
	result := new(big.Int)
	val.Int(result)
	return result
}

func EncryptWallet(privateKeyB64, password string) (string, error) {
	seed, _ := base64.StdEncoding.DecodeString(privateKeyB64)
	salt := make([]byte, 16)
	io.ReadFull(rand.Reader, salt)
	key, _ := scrypt.Key([]byte(password), salt, 32768, 8, 1, 32)
	block, _ := aes.NewCipher(key)
	gcm, _ := cipher.NewGCM(block)
	nonce := make([]byte, gcm.NonceSize())
	io.ReadFull(rand.Reader, nonce)
	cipherText := gcm.Seal(nil, nonce, seed, nil)
	var ks Keystore
	ks.Address = PublicKeyToAddress(ed25519.NewKeyFromSeed(seed).Public().(ed25519.PublicKey))
	ks.Crypto.Cipher = "aes-256-gcm"
	ks.Crypto.CipherText = base64.StdEncoding.EncodeToString(cipherText)
	ks.Crypto.Salt = base64.StdEncoding.EncodeToString(salt)
	ks.Crypto.Nonce = base64.StdEncoding.EncodeToString(nonce)
	data, _ := json.Marshal(ks)
	return string(data), nil
}

func DecryptWallet(keystoreJSON, password string) (string, error) {
	var ks Keystore
	json.Unmarshal([]byte(keystoreJSON), &ks)
	salt, _ := base64.StdEncoding.DecodeString(ks.Crypto.Salt)
	cipherText, _ := base64.StdEncoding.DecodeString(ks.Crypto.CipherText)
	nonce, _ := base64.StdEncoding.DecodeString(ks.Crypto.Nonce)
	key, _ := scrypt.Key([]byte(password), salt, 32768, 8, 1, 32)
	block, _ := aes.NewCipher(key)
	gcm, _ := cipher.NewGCM(block)
	seed, err := gcm.Open(nil, nonce, cipherText, nil)
	if err != nil { return "", fmt.Errorf("invalid password") }
	return base64.StdEncoding.EncodeToString(seed), nil
}

func SignTransaction(tx Transaction, privateKeyB64 string) (*SignedTransaction, error) {
	if tx.OU == "" { tx.OU = "1000" }
	tsRaw := tx.Timestamp.String()
	val, _ := strconv.ParseFloat(tsRaw, 64)
	cleanTS := strconv.FormatFloat(val, 'f', -1, 64)
	tx.Timestamp = json.Number(cleanTS)
	seed, _ := base64.StdEncoding.DecodeString(privateKeyB64)
	privateKey := ed25519.NewKeyFromSeed(seed)
	publicKey := privateKey.Public().(ed25519.PublicKey)
	raw := fmt.Sprintf(`{"from":"%s","to_":"%s","amount":"%s","nonce":%d,"ou":"%s","timestamp":%s}`,
		tx.From, tx.To, tx.Amount, tx.Nonce, tx.OU, cleanTS)
	sig := ed25519.Sign(privateKey, []byte(raw))
	return &SignedTransaction{
		Signature: base64.StdEncoding.EncodeToString(sig),
		PublicKey: base64.StdEncoding.EncodeToString(publicKey),
		Tx:        tx,
		Raw:       raw,
	}, nil
}

func (s *SignedTransaction) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"from": s.Tx.From, "to_": s.Tx.To, "amount": s.Tx.Amount,
		"nonce": s.Tx.Nonce, "ou": s.Tx.OU, "timestamp": s.Tx.Timestamp,
		"signature": s.Signature, "public_key": s.PublicKey,
	}
}

func (c *OctraClient) doRequest(ctx context.Context, method, path string, body interface{}) ([]byte, error) {
	var bodyReader io.Reader
	if body != nil {
		jsonData, _ := json.Marshal(body)
		bodyReader = bytes.NewBuffer(jsonData)
	}
	req, _ := http.NewRequestWithContext(ctx, method, c.BaseURL+path, bodyReader)
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.HTTPClient.Do(req)
	if err != nil { return nil, err }
	defer resp.Body.Close()
	data, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 { return nil, fmt.Errorf("rpc error [%d]: %s", resp.StatusCode, string(data)) }
	return data, nil
}

func (c *OctraClient) GetBalance(ctx context.Context, address string) (*BalanceInfo, error) {
	data, err := c.doRequest(ctx, "GET", "/balance/"+address, nil)
	if err != nil { return nil, err }
	var res BalanceInfo
	json.Unmarshal(data, &res)
	return &res, nil
}

func (c *OctraClient) GetNextNonce(ctx context.Context, address string) (uint64, error) {
	info, err := c.GetBalance(ctx, address)
	if err != nil { return 0, err }
	return info.Nonce + 1, nil
}

func (c *OctraClient) SendTransaction(ctx context.Context, signedTx *SignedTransaction) (map[string]interface{}, error) {
	data, err := c.doRequest(ctx, "POST", "/send-tx", signedTx.ToMap())
	if err != nil { return nil, err }
	var res map[string]interface{}
	json.Unmarshal(data, &res)
	return res, nil
}

func (c *OctraClient) GetTransaction(ctx context.Context, hash string) (map[string]interface{}, error) {
	data, err := c.doRequest(ctx, "GET", "/tx/"+hash, nil)
	if err != nil { return nil, err }
	var res map[string]interface{}
	json.Unmarshal(data, &res)
	return res, nil
}

func (c *OctraClient) WaitTransaction(ctx context.Context, hash string, timeout time.Duration) (map[string]interface{}, error) {
	deadline := time.Now().Add(timeout)
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done(): return nil, ctx.Err()
		case <-ticker.C:
			if time.Now().After(deadline) { return nil, fmt.Errorf("timeout") }
			tx, err := c.GetTransaction(ctx, hash)
			if err == nil && tx != nil {
				if status, _ := tx["status"].(string); status == "confirmed" { return tx, nil }
				if epoch, ok := tx["epoch"]; ok && epoch != nil { return tx, nil }
			}
		}
	}
}
