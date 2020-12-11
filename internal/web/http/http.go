package http

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/lighthouse-p2p/hub/internal/config"
	"github.com/lighthouse-p2p/hub/internal/web/http/handlers"
)

// InitHTTP initializes the HTTP server on the given address
func InitHTTP(cfg *config.Config) {
	addr := cfg.HTTPConfig.HTTPAddr
	log.Printf("Starting the HTTP server on %s\n", addr)

	app := fiber.New(fiber.Config{
		ServerHeader: "lighthousehub/fiber/1.0",
	})

	app.Use(logger.New())
	app.Use(recover.New())

	handlersInit := handlers.Handlers{Cfg: cfg}

	v1Group := app.Group("/v1")
	v1Group.Post("/register", handlersInit.Register)

	go log.Fatal(app.Listen(addr))
	log.Println("HTTP server is up!")
}
