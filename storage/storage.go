package storage

import "ethereum-parser/models"

type storage struct {
	currentBlock        int
	addressTransactions map[string][]models.Transaction
}

func New() models.Storage {
	return &storage{
		currentBlock:        0,
		addressTransactions: make(map[string][]models.Transaction),
	}
}

func (s *storage) SetCurrentBlock(blockNumber int) {
	s.currentBlock = blockNumber
}

func (s *storage) GetCurrentBlock() int {
	return s.currentBlock
}

func (s *storage) AddAddress(address string) bool {
	if _, ok := s.addressTransactions[address]; ok {
		return false
	}
	s.addressTransactions[address] = []models.Transaction{}
	return true
}

func (s *storage) GetTransactions(address string) []models.Transaction {
	transactions, ok := s.addressTransactions[address]
	if !ok {
		return []models.Transaction{}
	}
	
	return transactions
}

func (s *storage) GetSubscriptions() map[string][]models.Transaction {
	return s.addressTransactions
}
