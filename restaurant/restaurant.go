package restaurant

import (
	"log"
	"time"

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

func RestaurantSignUpHandler(c *fiber.Ctx) error {
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

func RestaurantLoginHandler(c *fiber.Ctx) error {
	req := new(RequestRestaurantLogin)
	if err := c.BodyParser(req); err != nil {
		log.Println("[ERROR] NewRestaurant | Failed to BodyParser")
		log.Println(err)
	}

	// 사용자가 입력을 하였는지 확인함.
	if req.RestaurantBusinessNumber == "" || req.RestaurantPassword == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": "정보가 제대로 입력되지 않았습니다. 확인 후 다시 시도 해 주세요.",
			"time":    time.Now(),
		})
	}

	db, err := database.ConntectionSQLite()
	if err != nil {
		log.Println("[ERROR] NewRestaurant | Failed to DataBase Connection")
		log.Println(err)
	}

	var restaurantInfo RestaurantINFO
	result := db.Table("restaurant_infos").Where("restaurant_business_number = ?", req.RestaurantBusinessNumber).First(&restaurantInfo)
	if result.RowsAffected == 0 {
		return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"status":  fiber.ErrBadRequest.Code,
			"message": "입력하신 정보를 토대로 사용자 정보를 찾을 수 없어요. 다시 학인해 주세요!",
			"time":    time.Now(),
		})
	}
	err = bcrypt.CompareHashAndPassword([]byte(restaurantInfo.RestaurantPassword), []byte(req.RestaurantPassword))
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  fiber.StatusUnauthorized,
			"message": "입력하신 정보가 일치하지 않아요. 다시 입력해 보시겠어요?",
			"time":    time.Now(),
		})
	}

	// JWT 토큰 발행
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["restaurant_uuid"] = restaurantInfo.RestaurantUUID
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix() // 토큰 만료 30분 설정
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		log.Println("Failed to Publish JWT")
		log.Println(err)
	}

	// 쿠키 설정
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    t,
		Expires:  time.Now().Add(time.Minute * 30), // 쿠키 만료 시간 30분
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	return c.Status(200).JSON(fiber.Map{
		"status":  200,
		"message": "성공적으로 로그인 되었습니다.",
		"time":    time.Now(),
	})
}

func RestaurantHomeHandler(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")

	token, err := jwt.ParseWithClaims(cookie, jwt.StandardClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	if err != nil {
		log.Println("Failed to Get JWT")
		log.Println(err)
	}

	claims := token.Claims.(*jwt.StandardClaims)

	db, err := database.ConntectionSQLite()
	if err != nil {
		log.Println("Failed to Connection SQLite")
		log.Println(err)
	}

	var RestaurantInfo RestaurantINFO
	db.Where("restaurant_uuid = ?", claims.Issuer).First(&RestaurantInfo)

	return c.Status(200).JSON(RestaurantInfo)

}

func RestaurantLogOutHandler(c *fiber.Ctx) error {
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	return c.Status(200).JSON(fiber.Map{
		"status":  200,
		"message": "성공적으로 로그아웃 되셨습니다.",
		"time":    time.Now(),
	})
}
