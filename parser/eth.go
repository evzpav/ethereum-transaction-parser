package parser

import (
	"bytes"
	"encoding/json"
	"ethereum-parser/models"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"strings"
	"time"
)

func (p *parser) makeRPCRequest(method string, params []interface{}) (json.RawMessage, error) {

	requestBody := models.JsonRPCRequest{
		Jsonrpc: "2.0",
		Method:  method,
		Params:  params,
		ID:      time.Now().Nanosecond(),
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request body: %v", err)
	}

	resp, err := http.Post(p.nodeUrl, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error sending HTTP request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	var rpcResponse models.JsonRPCResponse
	err = json.Unmarshal(body, &rpcResponse)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %v", err)
	}

	if rpcResponse.Error != nil {
		return nil, fmt.Errorf("JSON-RPC error: %v", rpcResponse.Error)
	}

	return rpcResponse.Result, nil
}

func (p *parser) getETHLatestBlock() (*big.Int, error) {
	result, err := p.makeRPCRequest("eth_blockNumber", []interface{}{})
	if err != nil {
		return nil, fmt.Errorf("error getting latest block number: %v", err)
	}

	blockStr := strings.ReplaceAll(string(result), "\"", "")

	return hexToBigInt(blockStr)
}

func (p *parser) getETHBlockByNumber(blockNumber *big.Int) (*models.Block, error) {
	result, err := p.makeRPCRequest("eth_getBlockByNumber", []interface{}{
		fmt.Sprintf("0x%x", blockNumber), true,
	})
	if err != nil {
		return nil, fmt.Errorf("error getting block by number: %v", err)
	}

	var block models.Block
	err = json.Unmarshal(result, &block)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling block: %v", err)
	}

	return &block, nil
}

func hexToBigInt(hexStr string) (*big.Int, error) {
	n := new(big.Int)
	_, ok := n.SetString(hexStr[2:], 16) // Removing the "0x" prefix and setting the base to 16
	if !ok {
		return nil, fmt.Errorf("invalid hexadecimal string")
	}
	return n, nil
}
