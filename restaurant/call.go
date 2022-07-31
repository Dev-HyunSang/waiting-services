package restaurant

import (
	"fmt"
	"log"

	"github.com/dev-hyunsang/waiting-services/database"
	"github.com/dev-hyunsang/waiting-services/models"
	"github.com/gofiber/fiber/v2"
)

func CallHandler(c *fiber.Ctx) error {
	req := new(models.RequestRestaurant)
	if err := c.BodyParser(req); err != nil {
		log.Println("Failed to BodyParser")
		log.Println(err)
	}

	db, err := database.ConntectionSQLite()
	if err != nil {
		log.Println("Failed to Connection SQLite")
		log.Println(err)
	}

	var waiting []models.Waiting
	db.Where("restaurant_uuid = ?", req.RestaurantUUID).Find(waiting)

	return c.Status(200).JSON(fiber.Map{
		"message": fmt.Sprintf("성공적으로 %s의 웨이팅 리스트를 가지고 왔어요.", "hello"),
	})
}
