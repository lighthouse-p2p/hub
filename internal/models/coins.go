package models

import (
	"errors"

	"github.com/lighthouse-p2p/hub/internal/config"
	"github.com/vmihailenco/msgpack/v5"
	"golang.org/x/crypto/blake2b"
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
	TXN float64

	// Hash is the blake2 hash of the last node's msgpack in hex format
	Hash string
}

// CoinPack is used to pack a record into msgpack
type CoinPack struct {
	PubKey     string  `msgpack:"public_key"`
	TotalCoins float64 `msgpack:"total_coins"`
	TXN        float64 `msgpack:"txn"`
	Hash       string  `msgpack:"hash"`
}

// AddBlock adds a new block to the coin chain
func AddBlock(cfg *config.Config, pubKey string, totalCoins, txn float64) error {
	db := cfg.Database

	newBlock := &CoinChain{}

	var lastBlock CoinChain
	tx := db.Model(&CoinChain{}).Last(lastBlock)

	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			// Handle record not found (first record)
			newBlock.Hash = ""
		} else {
			return tx.Error
		}
	} else {
		pack, err := msgpack.Marshal(&CoinPack{
			PubKey:     lastBlock.PubKey,
			TotalCoins: lastBlock.TotalCoins,
			TXN:        lastBlock.TXN,
			Hash:       lastBlock.Hash,
		})
		if err != nil {
			return err
		}

		blake, err := blake2b.New512(nil)
		if err != nil {
			return err
		}

		_, err = blake.Write(pack)
		if err != nil {
			return err
		}

		newBlock.Hash = string(blake.Sum(nil))
	}

	newBlock.PubKey = pubKey
	newBlock.TotalCoins = totalCoins
	newBlock.TXN = txn

	tx = db.Create(&newBlock)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}
