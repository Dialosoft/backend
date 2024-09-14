package router

import "github.com/gofiber/fiber/v3"

func SetRoutes(app *fiber.App) {

	// api identity
	api := app.Group("/dialosoft-api/v1")

	api.Get("/data", func(ctx fiber.Ctx) error {
		return ctx.SendString("hi")
	})

}
