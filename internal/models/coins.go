package models

import (
	"gorm.io/gorm"
)

// CoinChain the the coin blockchain stored onto a psql DB for now
type CoinChain struct {
	gorm.Model

	// PubKey is the public key of the peer
	PubKey string

	// TotalCoins is the number of coins the peer has after the txn
	TotalCoins float64

	// TXN is the amount added/reduced from the old `TotalCoins`
	TNX float64

	// Hash is the blake2 hash of the last node's protobuf in hex format
	Hash string
}
