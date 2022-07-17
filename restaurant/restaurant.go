package restaurant

import (
	"github.com/gofrs/uuid"
)

type restaurant struct {
	RestaurantUUID uuid.UUID `json:"restaurant_uuid"`
	RestaurantName string    `json:"restaurant_name"`
}
