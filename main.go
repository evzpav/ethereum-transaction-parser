package main

import (
	"ethereum-parser/api"
	"ethereum-parser/parser"
	"ethereum-parser/storage"
	"log"
	"time"
)

var nodeUrl = "https://cloudflare-eth.com"
var parsingInterval = 10 * time.Second
var apiPort = 8787

func main() {

	storageInst := storage.New() // in memory storage

	parserInst := parser.New(nodeUrl, storageInst, parsingInterval)

	apiInst := api.New(apiPort, parserInst)

	go parserInst.ParseBlockTransactions()

	log.Fatal(apiInst.Serve())
}
