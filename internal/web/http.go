package web

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/websocket/v2"
	"github.com/lighthouse-p2p/hub/internal/config"
	"github.com/lighthouse-p2p/hub/internal/web/handlers"
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
	v1Group.Get("/coins", handlersInit.Coins)
	v1Group.Post("/register", handlersInit.Register)
	v1Group.Get("/resolve/:nickname", handlersInit.ResolveNickName)

	wsGroup := v1Group.Group("/ws")
	wsGroup.Use("/", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			return c.Next()
		}

		return c.Status(400).SendString("Plain HTTP to websocket endpoint :(")
	})
	wsGroup.Get("/signaling", websocket.New(handlersInit.Signaling))

	go log.Fatal(app.Listen(addr))
	log.Println("HTTP server is up!")
}
