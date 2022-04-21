package database

import (
	"fmt"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Database struct {
	Conn *gorm.DB
}

type User struct {
	gorm.Model
	Gamertag    string
	IsProcessed bool
}

func NewDevConnection() *Database {
	db, err := gorm.Open(sqlite.Open("./internal/gotana-dev.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&User{})

	return &Database{
		Conn: db,
	}
}

func (db *Database) CheckIn(gamertag string) (User, error) {
	var user User
	db.Conn.FirstOrCreate(&user, User{Gamertag: gamertag, IsProcessed: false})

	if user.IsProcessed {
		fmt.Println("\n * Data has been processed. Check AWS.\n")
		os.Exit(3)
	}

	return user, nil
}

func (db *Database) MarkComplete(user *User) {
	db.Conn.Model(&user).Update("is_processed", true)
}

func (db *Database) ListUsers() error {
	var users []User
	db.Conn.Find(&users)

	fmt.Println(users)
	return nil
}
