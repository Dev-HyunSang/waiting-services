package guest

import (
	"fmt"
	"log"
	"time"

	"github.com/dev-hyunsang/waiting-services/database"
	"github.com/dev-hyunsang/waiting-services/models"
	"github.com/gofiber/fiber/v2"
	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt/v4"
)

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

func WaitingHomeHandler(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")

	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  fiber.StatusUnauthorized,
			"message": "unauthenticated",
			"time":    time.Now(),
		})
	}

	clamis := token.Claims.(*jwt.StandardClaims)

	db, err := database.ConntectionSQLite()
	if err != nil {
		log.Println("Failed to Connection SQLite")
		log.Println(err)
	}

	var (
		user   models.User
		wating models.Waiting
	)
	db.Table("").Where("user_uuid = ?", clamis.Issuer).First(&user)

	db.Table("").Where("user_uuid = ?", user.UserUUID).First(&wating)

	return c.Status(200).JSON(fiber.Map{
		"status":  200,
		"message": "성공적으로 웨이팅 정보를 불러 왔습니다.",
		"datas":   wating,
		"time":    time.Now(),
	})
}

// 등록 되어 있는 가게를 검색하는 기능
func RestaurantSearch(search string) models.RestaurantINFO {
	db, err := database.ConntectionSQLite()
	if err != nil {
		log.Println("Failed to DataBase Connection SQLite")
		log.Println(err)
	}

	var restaurantInfo models.RestaurantINFO
	db.Table("restaurant_infos").Where("restaurant_name = ?", search).First(&restaurantInfo)

	return restaurantInfo
}

func NewWaitingHandler(c *fiber.Ctx) error {
	req := new(models.Waiting)
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

	var user models.User
	db.Table("user").Where("user_uuid =?", req.UserUUID).Find(&user)

	db.Table("waitings").Create(models.Waiting{
		WaitingUUID:    waitingUUID,
		RestaurantUUID: req.RestaurantUUID,
		UserUUID:       user.UserUUID,
		Status:         false,
		CreatedAt:      time.Now(),
		DeletedAt:      time.Now(),
	})

	return c.Status(200).JSON(fiber.Map{
		"status":  200,
		"message": fmt.Sprintf("%s님! 성공적으로 등록되었어요. 순서가 되면 문자로 알려드릴께요!", user.Name),
		"time":    time.Now(),
	})
}
