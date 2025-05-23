package main

import (
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Get("/spoorwebsite/:domain", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"ok": true,
		})
	})

	app.Listen(":8000")
}
