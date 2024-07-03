package models

import "encoding/json"

type Parser interface {
	// last parsed block
	GetCurrentBlock() int
	// add address to observer
	Subscribe(address string) bool
	// list of inbound or outbound transactions for an address
	GetTransactions(address string) []Transaction

	ParseBlockTransactions()
}

type Storage interface {
	SetCurrentBlock(blockNumber int)
	GetCurrentBlock() int
	AddAddress(address string) bool
	GetTransactions(address string) []Transaction
	GetSubscriptions() map[string][]Transaction
}

type JsonRPCRequest struct {
	Jsonrpc string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	ID      int           `json:"id"`
}

type JsonRPCResponse struct {
	Jsonrpc string          `json:"jsonrpc"`
	ID      int             `json:"id"`
	Result  json.RawMessage `json:"result"`
	Error   *JsonRPCError   `json:"error"`
}

type JsonRPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type Block struct {
	Number       string        `json:"number"`
	Transactions []Transaction `json:"transactions"`
}

type Transaction struct {
	From           string `json:"from"`
	To             string `json:"to"`
	Hash           string `json:"hash"`
	Value          string `json:"value"`
	ValueInt       int    `json:"valueInt"`
	Type           string `json:"type"`
	BlockNumberHex string `json:"blockNumber"`
	BlockNumberInt int    `json:"blockNumberInt"`
	Gas            string `json:"gas"`
	GasPrice       string `json:"gasPrice"`
	Input          string `json:"input"`
}

type AddressTransactions struct {
	Address      string        `json:"address"`
	Transactions []Transaction `json:"transactions"`
	NewestBlock  int           `json:"newestBlock"`
	OldestBlock  int           `json:"oldestBlock"`
}

type SubscriptionResponse struct {
	Jsonrpc string             `json:"jsonrpc"`
	Method  string             `json:"method"`
	Params  SubscriptionParams `json:"params"`
}

type SubscriptionParams struct {
	Subscription string          `json:"subscription"`
	Result       json.RawMessage `json:"result"`
}
