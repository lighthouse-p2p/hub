package http

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/lighthouse-p2p/hub/internal/config"
)

// InitHTTP initializes the HTTP server on the given address
func InitHTTP(cfg *config.Config) {
	addr := cfg.HTTPConfig.HTTPAddr

	log.Printf("Starting the HTTP server on %s\n", addr)

	app := fiber.New()

	go log.Fatal(app.Listen(addr))
	log.Println("HTTP server is up!")
}
