package restaurant

import (
	"log"
	"time"

	"github.com/dev-hyunsang/waiting-services/config"
	"github.com/dev-hyunsang/waiting-services/database"
	"github.com/gofiber/fiber/v2"
	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type RestaurantINFO struct {
	RestaurantUUID           uuid.UUID `json:"restaurant_uuid"`     // 각 레스토랑마다 고유 UUID
	RestaurantName           string    `json:"restaurant_name"`     // 레스토랑 가게명
	RestaurantLocation       string    `json:"restaurant_location"` // 레스토랑 위치
	RestaurantOwnerName      string    `json:"restaurant_owner_name"`
	RestaurantBusinessNumber string    `json:"restaurant_business_number"`
	RestaurantPassword       string    `json:"restaurant_password"`
	CreatedAt                time.Time `json:"created__at"`
	EditedAt                 time.Time `json:"edited_at"`
	DeletedAt                time.Time `json:"deleted_at"`
}

type RequestRestaurantLogin struct {
	RestaurantBusinessNumber string    `json:"restaurant_business_number"`
	RestaurantPassword       string    `json:"restaurant_password"`
	LoginedAt                time.Time `json:"logined_at"`
}

func NewRestaurantHandler(c *fiber.Ctx) error {
	req := new(RestaurantINFO)
	if err := c.BodyParser(req); err != nil {
		log.Println("[ERROR] NewRestaurant | Failed to BodyParser")
		log.Println(err)
	}

	db, err := database.ConntectionSQLite()
	if err != nil {
		log.Println("[ERROR] NewRestaurant | Failed to DataBase Connection")
		log.Println(err)
	}

	restaurantUUID, err := uuid.NewV4()
	if err != nil {
		log.Println("[ERROR] NewRestaurant | Failed to UUID New4")
		log.Println(err)
	}

	pwHash, err := bcrypt.GenerateFromPassword([]byte(req.RestaurantPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Println("[ERROR] NewRestaurant | Failed Generate From Password")
		log.Println(err)
	}

	db.Table("restaurant_infos").Create(RestaurantINFO{
		RestaurantUUID:           restaurantUUID,
		RestaurantName:           req.RestaurantName,
		RestaurantLocation:       req.RestaurantLocation,
		RestaurantOwnerName:      req.RestaurantOwnerName,
		RestaurantBusinessNumber: req.RestaurantBusinessNumber,
		RestaurantPassword:       string(pwHash),
		CreatedAt:                time.Now(),
		EditedAt:                 time.Now(),
		DeletedAt:                time.Now(),
	})

	return c.Status(200).JSON(fiber.Map{
		"status":  200,
		"message": "성공적으로 가게를 등록했어요. 사업자 번호와 패스워드를 통해서 실시간으로 현황을 파악할 수 있어요.",
		"time":    time.Now(),
	})

}

func LoginRestaurantHandler(c *fiber.Ctx) error {
	req := new(RequestRestaurantLogin)
	if err := c.BodyParser(req); err != nil {
		log.Println("[ERROR] NewRestaurant | Failed to BodyParser")
		log.Println(err)
	}

	db, err := database.ConntectionSQLite()
	if err != nil {
		log.Println("[ERROR] NewRestaurant | Failed to DataBase Connection")
		log.Println(err)
	}

	var restaurantInfo RestaurantINFO
	db.Table("").Where("restaurant_business_number = ?", req.RestaurantBusinessNumber).Find(restaurantInfo)

	if req.RestaurantBusinessNumber != restaurantInfo.RestaurantBusinessNumber {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  fiber.StatusUnauthorized,
			"message": "입력하신 정보가 일치하지 않아요. 다시 입력해 보시겠어요?",
			"time":    time.Now(),
		})
	}

	err = bcrypt.CompareHashAndPassword([]byte(restaurantInfo.RestaurantPassword), []byte(req.RestaurantPassword))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  fiber.StatusUnauthorized,
			"message": "입력하신 정보가 일치하지 않아요. 다시 입력해 보시겠어요?",
			"time":    time.Now(),
		})
	}

	claims := jwt.MapClaims{
		"restaurant_uuid": restaurantInfo.RestaurantUUID,
		"restaurant_name": restaurantInfo.RestaurantName,
		"bussines_number": restaurantInfo.RestaurantBusinessNumber,
		"exp":             time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(config.GetEnv("JWT_KEY")))
	if err != nil {
		log.Println("Failed to JWT SignedString")
		log.Println(err)
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  200,
		"message": "정상적으로 로그인이 되었습니다.",
		"token":   t,
		"time":    time.Now(),
	})
}
