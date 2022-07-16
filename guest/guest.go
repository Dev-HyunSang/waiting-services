package guest

import (
	"fmt"
	"log"
	"time"

	"github.com/dev-hyunsang/waiting-services/database"
	"github.com/dev-hyunsang/waiting-services/mail"
	"github.com/gofiber/fiber/v2"
	"github.com/gofrs/uuid"
)

type Waiting struct {
	WaitingUUID uuid.UUID `json:"waiting_uuid" gorm:"column:waiting_uuid;"`
	StoreUUID   uuid.UUID `json:"store_uuid" gorm:"column:store_uuid"`
	StoreName   string    `json:"store_name" gorm:"column:store_name"`
	Name        string    `json:"name" gorm:"column: name;"`
	Email       string    `json:"email" gorm:"column:email"`
	Status      string    `json:"status" gorm:"column:status;"`
	CreatedAt   time.Time `json:"created_at" gorm:"column:created_at;"`
	DeletedAt   time.Time `json:"deleted_at" gorm:"column:deleted_at;"`
}

/* HOW TO WORKING SERVICES?
1. 유저가 새로운 웨이팅을 가게를 검색한 후 등록해요.
2. 유저가 새로운 웨이팅을 가게 등록한 것을 가게 측으로 보내요.
   => 1. 가게 측에서 웨이팅 서비스를 키고 끄고 할 수 있도록 해요.
	  2. n개 이상의 순번이 되면 더 이상 받지 않도록 해요.
3. 가게 측에서 웨이팅 순번이 되면 유저에게 메일로 노티 해요.
   => 10분 이내에 안 오면 순번에서 제외하고 다른 손님을 불러요.
      웨이팅 등록 시 주의사항을 꼭 유저에게 알릴 수 있도록 제작해요.
4. 가게 측에 도착 했으면 가게 측에서 확인하고 순번을 삭제해요.
*/

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

	db.Table("waitings").Create(Waiting{
		WaitingUUID: waitingUUID,
		Name:        req.Name,
		Email:       req.Email,
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

	mail.SendEmail(wating.Email, wating.Name, wating.StoreName)

	return c.Status(200).JSON(fiber.Map{})
}
