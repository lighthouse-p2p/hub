package models

import (
	"gorm.io/gorm"
)

// DHT is the DB model for the hash table
type DHT struct {
	gorm.Model

	// PubKey is the base64 Curve25519 public key of the peer
	PubKey string `gorm:"uniqueIndex"`

	// NickName is the peer's short name used to make reaching it easier
	NickName string `gorm:"uniqueIndex"`
}
