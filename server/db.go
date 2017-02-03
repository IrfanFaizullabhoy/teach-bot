package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"os"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func ConnectToPG() *gorm.DB {
	db, err := gorm.Open("postgres", fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", os.Getenv("DB_SERVICE"), os.Getenv("DB_USER"), os.Getenv("DB_NAME"), os.Getenv("DB_PASS")))
	check(err)
	return db
}

func SetupDB(db *gorm.DB) {
	db.AutoMigrate(&User{})
}

func CreateUser(user User, db *gorm.DB) {
	if db.NewRecord(&user) {
		db.Create(&user)
	} else {
		fmt.Println("error, primary key already exists for user")
	}
}
