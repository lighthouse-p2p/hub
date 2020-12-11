package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lighthouse-p2p/hub/internal/models"
)

// Register handles the registration of a user
func (h *Handlers) Register(ctx *fiber.Ctx) error {
	body := new(models.RegisterBody)
	if err := ctx.BodyParser(body); err != nil {
		return ctx.Status(400).SendString("Unable to parse body")
	}

	if body.NickName == "" || body.PubKey == "" {
		return ctx.Status(400).SendString("Either the nickname or the public key is empty")
	}

	tx := h.Cfg.Database.Create(&models.DHT{
		PubKey:   body.PubKey,
		NickName: body.NickName,
	})

	if tx.Error != nil {
		return ctx.Status(409).SendString("Either the nickname, or the public key is already used")
	}

	return ctx.SendStatus(201)
}
