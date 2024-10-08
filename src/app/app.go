package app

import (
	"line/src/configs/db"
	"line/src/configs/env"
	"line/src/configs/rabbitmq"
	"line/src/handlers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"github.com/gofiber/websocket/v2"
)

func SetupApp() *fiber.App {
	app := fiber.New()

	env.LoadConfig()
	db.ConnectDB()
	rabbitmq.Connect()

	app.Use(logger.New())

	app.Post("/graphql", handlers.GraphQLHandler())

	app.Get("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			return websocket.New(handlers.HandleWebSocket)(c)
		}
		return fiber.ErrBadRequest
	})
	return app
}
