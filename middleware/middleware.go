package middleware

import (
	"strings"

	"github.com/dev-hyunsang/waiting-services/guest"
	"github.com/dev-hyunsang/waiting-services/restaurant"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func Middleware(app *fiber.App) {
	app.Static("/", "./public")

	app.Use(cors.New(cors.Config{
		AllowMethods: strings.Join([]string{
			fiber.MethodGet,
			fiber.MethodPost,
			fiber.MethodDelete,
			fiber.MethodPatch,
		}, ","),
		AllowCredentials: true,
	}))

	api := app.Group("/api")

	Guest := api.Group("/guest")
	Guest.Post("/new", guest.NewWaitingHandler)

	Restaurant := api.Group("/restaurant")
	Restaurant.Post("/join", restaurant.RestaurantSignUpHandler)
	Restaurant.Post("/login", restaurant.RestaurantLogOutHandler)
	Restaurant.Post("/home", restaurant.RestaurantHomeHandler)
	Restaurant.Post("/logout", restaurant.RestaurantLogOutHandler)
}
