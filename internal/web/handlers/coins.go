package handlers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/lighthouse-p2p/hub/internal/models"
)

// Coins returns the coins of a pub key
func (h *Handlers) Coins(c *fiber.Ctx) error {
	pubKey := c.Query("pub_key")
	if pubKey == "" {
		return c.SendStatus(400)
	}

	db := h.Cfg.Database

	var lastBlockForPubKey models.CoinBlock
	tx := db.Model(&models.CoinBlock{}).Where("pub_key = ?", pubKey).Last(&lastBlockForPubKey)
	if tx.Error != nil {
		return c.SendStatus(404)
	}

	return c.SendString(fmt.Sprintf("%f", lastBlockForPubKey.TotalCoins))
}
