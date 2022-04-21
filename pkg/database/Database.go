package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Gamertag string
}

func NewDevConnection() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("../../internal/gotana-dev.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&User{})

	return db
}
