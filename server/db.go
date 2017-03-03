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
	db, err := gorm.Open("postgres", fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", os.Getenv("DB_PORT_5432_TCP_ADDR"), os.Getenv("DB_USER"), os.Getenv("DB_NAME"), os.Getenv("DB_PASS")))
	check(err)
	return db
}

func ConnectToRDS() *gorm.DB {
	db, err := gorm.Open("postgres", fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", "usersdev.cv3awqdhwmuz.us-west-2.rds.amazonaws.com", os.Getenv("DB_USER"), os.Getenv("DB_NAME"), os.Getenv("DB_PASS")))
	check(err)
	return db
}

func SetupDB(db *gorm.DB) {
	db.AutoMigrate(&User{}, &Assignment{})
}
