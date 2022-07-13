package middleware

import "github.com/gofiber/fiber/v2"

func Middleware(app *fiber.App) {
	api := app.Group("/api")

	guest := api.Group("/guest")
	guest.Post("/new", )
	guest.Get("/list", )
}
