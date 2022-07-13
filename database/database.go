package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func ConntectionSQLite() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("waiting-services.db"), &gorm.Config{})

	return db, err
}
