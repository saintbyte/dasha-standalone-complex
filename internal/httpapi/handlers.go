package httpapi

import (
	"github.com/gofiber/fiber/v2"
)

func New() *fiber.App {
	app := fiber.New()
	app.Get("/", index)
	return app
}

func index(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"service": "dasha-standalone-complex", "status": "ok"})
}
