package main

import (
	"log"

	"github.com/dev-hyunsang/waiting-services/database"
	"github.com/dev-hyunsang/waiting-services/guest"
	"github.com/dev-hyunsang/waiting-services/middleware"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	middleware.Middleware(app)

	db, err := database.ConntectionSQLite()
	if err != nil {
		log.Println("Failed to Connection SQLite")
		log.Fatalln(err)
	}

	db.AutoMigrate(&guest.Waiting{})

	if err := app.Listen(":3000"); err != nil {
		log.Println("Failed to Runing Server...!")
		log.Fatalln(err)
	}
}
