package models

import (
	"time"

	"github.com/gofrs/uuid"
)

type User struct {
	UserUUID    uuid.UUID `json:"user_uuid" gorm:"user_uuid;"`
	Name        string    `json:"name" gorm:"column: name;"`
	Email       string    `json:"email" gorm:"column:email;"`
	PhoneNumber string    `json:"phone_number" gorm:"column:phone_number;"`
	CreatedAt   time.Time `json:"created_at" gorm:"column:created_at;"`
	DeletedAt   time.Time `json:"deleted_at" gorm:"column:deleted_at;"`
}

type Waiting struct {
	WaitingUUID    uuid.UUID `json:"waiting_uuid" gorm:"column:waiting_uuid"`
	RestaurantUUID uuid.UUID `json:"restaurant_uuid" gorm:"column:restaurant_uuid"`
	UserUUID       uuid.UUID `json:"user_uuid" gorm:"column:user_uuid;"`
	Status         bool      `json:"status" gorm:"colum:status"`
	CreatedAt      time.Time `json:"created_at" gorm:"coloumn:created_at"`
	DonedAt        time.Time `json:"doned_at" gorm:"coloumn:doned_at"`
	DeletedAt      time.Time `json:"deleted_at" gorm:"coloumn:deleted_at"`
}
