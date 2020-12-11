package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lighthouse-p2p/hub/internal/models"
)

// ResolveNickName returns the public key for a given nickname
func (h *Handlers) ResolveNickName(ctx *fiber.Ctx) error {
	nickName := ctx.Params("nickname")
	if nickName == "" {
		// Redundent as its handled by the router
		// Just here for edge cases
		return ctx.Status(400).SendString("Nickname must not be empty")
	}

	var record models.DHT
	tx := h.Cfg.Database.Model(&models.DHT{}).Where("nick_name = ?", nickName).First(&record)
	if tx.Error != nil {
		return ctx.Status(404).SendString("Not Found")
	}

	return ctx.JSON(&models.ResolveReponse{
		PubKey: record.PubKey,
	})
}
