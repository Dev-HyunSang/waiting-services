package models

import (
	"time"

	"github.com/gofrs/uuid"
)

type RestaurantINFO struct {
	RestaurantUUID           uuid.UUID `json:"restaurant_uuid"`     // 각 레스토랑마다 고유 UUID
	RestaurantName           string    `json:"restaurant_name"`     // 레스토랑 가게명
	RestaurantLocation       string    `json:"restaurant_location"` // 레스토랑 위치
	RestaurantPhoneNumber    string    `json:"restaurant_phone_number"`
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

type RequestRestaurant struct {
	RestaurantUUID uuid.UUID `json:"restaurant_uuid"`
	Date           time.Time `json:"date"`
}
