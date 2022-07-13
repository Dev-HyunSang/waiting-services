package guest

import (
	"fmt"
	"log"
	"time"

	"github.com/dev-hyunsang/waiting-services/database"
	"github.com/gofiber/fiber/v2"
	"github.com/gofrs/uuid"
)

type Waiting struct {
	WaitingUUID uuid.UUID `json:"waiting_uuid" gorm:"colum:"`
	Name        string    `json:"name"`
	PhoneNumber string    `json:"phone_number"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	DeletedAt   time.Time `json:"deleted_at"`
}

func NewWaitingHandler(c *fiber.Ctx) error {
	req := new(Waiting)
	if err := c.BodyParser(req); err != nil {
		log.Println("Failed to BodyParser")
		log.Println(err)
	}

	waitingUUID, err := uuid.NewV4()
	if err != nil {
		log.Println("Failed Create UUID v4")
		log.Fatalln(err)
	}

	db, err := database.ConntectionSQLite()
	if err != nil {
		log.Println("Failed to Connection SQLite")
		log.Fatalln(err)
	}

	db.Table("").Create(Waiting{
		WaitingUUID: waitingUUID,
		Name:        req.Name,
		PhoneNumber: req.PhoneNumber,
		Status:      "waiting",
		CreatedAt:   time.Now(),
		DeletedAt:   time.Now(),
	})

	return c.Status(200).JSON(fiber.Map{
		"status":  200,
		"message": fmt.Sprintf("%s님! 성공적으로 등록되었어요. 순서가 되면 문자로 알려드릴께요!", req.Name),
		"time":    time.Now(),
	})
}

func CallWaitingHandler(c *fiber.Ctx) error {
	var wating Waiting
	req := new(Waiting)
	if err := c.BodyParser(req); err != nil {
		log.Println("Failed to BodyParser")
		log.Println(err)
	}

	db, err := database.ConntectionSQLite()
	if err != nil {
		log.Println("Failed to Connection SQLite")
		log.Fatalln(err)
	}

	db.Where("wating_uuid = ?", req.WaitingUUID).Find(wating)

	return c.Status(200).JSON(fiber.Map{})
}
