package parser

import (
	"ethereum-parser/models"
	"fmt"
	"log"
	"math/big"
	"strings"
	"time"
)

type parser struct {
	nodeUrl         string
	storage         models.Storage
	parsingInterval time.Duration
	retryDelay      time.Duration
}

func New(nodeUrl string, storage models.Storage, parsingInterval time.Duration) models.Parser {
	return &parser{
		nodeUrl:         nodeUrl,
		storage:         storage,
		parsingInterval: parsingInterval,
		retryDelay:      2 * time.Second,
	}
}

func (p *parser) GetCurrentBlock() int {
	return p.storage.GetCurrentBlock()
}

func (p *parser) setCurrentBlock(blockNumber int) {
	p.storage.SetCurrentBlock(blockNumber)
}

func (p *parser) Subscribe(address string) bool {
	address = strings.ToLower(address)
	return p.storage.AddAddress(address)
}

func (p *parser) GetTransactions(address string) []models.Transaction {
	address = strings.ToLower(address)
	return p.storage.GetTransactions(address)
}

func (p *parser) GetSubscriptions() map[string][]models.Transaction {
	return p.storage.GetSubscriptions()
}

func (p *parser) ParseBlockTransactions() {
	p.parseTransactionsPerBlock()

	ticker := time.NewTicker(p.parsingInterval)

	for range ticker.C {
		p.parseTransactionsPerBlock()
	}
}

func (p *parser) parseTransactionsPerBlock() {
	fmt.Println("Parsing...")

	lastestBlock, err := p.getETHLatestBlock()
	if err != nil {
		log.Printf("error getting latest block number: %v", err)
		time.Sleep(p.retryDelay)
		return
	}

	current := big.NewInt(int64(p.GetCurrentBlock()))
	if current.Cmp(big.NewInt(0)) == 0 {
		current = new(big.Int).Sub(lastestBlock, big.NewInt(10))
	}

	if current.Cmp(lastestBlock) >= 0 {
		return
	}

	fmt.Println("Current block: ", current.String())
	fmt.Println("Lastest block: ", lastestBlock.String())

	addressMap := p.GetSubscriptions()

	if len(addressMap) == 0 {
		fmt.Println("no addresses to parse")
		return
	}

	for current.Cmp(lastestBlock) <= 0 {

		block, err := p.getETHBlockByNumber(current)
		if err != nil {
			log.Printf("error getting block: %v\n", err)
			time.Sleep(p.retryDelay)
			continue
		}

		for _, tx := range block.Transactions {

			_, ok := addressMap[tx.From]
			if ok {
				tx.Type = "outbound"
				tx.BlockNumberInt = int(current.Int64())
				valueBigInt, _ := hexToBigInt(tx.Value)
				tx.ValueInt = int(valueBigInt.Int64())
				addressMap[tx.From] = append(addressMap[tx.From], tx)
			}

			_, ok = addressMap[tx.To]
			if ok {
				tx.Type = "inbound"
				addressMap[tx.To] = append(addressMap[tx.To], tx)
			}
		}

		p.setCurrentBlock(int(current.Int64()))
		current = current.Add(current, big.NewInt(1))
	}
}
