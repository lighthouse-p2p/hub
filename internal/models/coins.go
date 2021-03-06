package models

import (
	"encoding/base64"
	"errors"

	"github.com/lighthouse-p2p/hub/internal/config"
	"github.com/vmihailenco/msgpack/v5"
	"golang.org/x/crypto/blake2b"
	"gorm.io/gorm"
)

// CoinBlock the the coin blockchain stored onto a psql DB for now
type CoinBlock struct {
	gorm.Model

	// PubKey is the public key of the peer
	PubKey string `gorm:"index"`

	// TotalCoins is the number of coins the peer has after the txn
	TotalCoins float64

	// TXN is the amount added/reduced from the old `TotalCoins`
	TXN float64

	// Hash is the blake2 hash of the last node's msgpack in base64
	LastHash string
}

// CoinPack is used to pack a record into msgpack
type CoinPack struct {
	PubKey     string  `msgpack:"public_key"`
	TotalCoins float64 `msgpack:"total_coins"`
	TXN        float64 `msgpack:"txn"`
	LastHash   string  `msgpack:"hash"`
}

// AddBlock adds a new block to the coin chain
func AddBlock(cfg *config.Config, pubKey string, txn float64) error {
	db := cfg.Database

	newBlock := &CoinBlock{}

	var lastBlock CoinBlock
	tx := db.Model(&CoinBlock{}).Last(&lastBlock)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			// Handle record not found (first record)
			newBlock.LastHash = "init_block"
		} else {
			return tx.Error
		}
	} else {
		pack, err := msgpack.Marshal(&CoinPack{
			PubKey:     lastBlock.PubKey,
			TotalCoins: lastBlock.TotalCoins,
			TXN:        lastBlock.TXN,
			LastHash:   lastBlock.LastHash,
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

		newBlock.LastHash = base64.StdEncoding.EncodeToString(blake.Sum(nil))
	}

	var lastBlockForPubKey CoinBlock
	tx = db.Model(&CoinBlock{}).Where("pub_key = ?", pubKey).Last(&lastBlockForPubKey)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			newBlock.TotalCoins = txn
		} else {
			return tx.Error
		}
	} else {
		newBlock.TotalCoins = lastBlockForPubKey.TotalCoins + txn
	}

	newBlock.PubKey = pubKey
	newBlock.TXN = txn

	tx = db.Create(&newBlock)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}
