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

func SetupDB(db *gorm.DB) {
	db.AutoMigrate(&Student{}, &Teacher{})
}

func CreateStudent(student Student, db *gorm.DB) {
	if db.NewRecord(&student) {
		db.Create(&student)
	} else {
		fmt.Println("error, primary key already exists for user")
	}
}

func CreateTeacher(teacher Teacher, db *gorm.DB) {
	if db.NewRecord(&teacher) {
		db.Create(&teacher)
	} else {
		fmt.Println("error, primary key already exists for user")
	}
}

func GetStudentChannels(db *gorm.DB) []string {
	var channelIDs []string
	db.Table("students").Select("channel_id").Find(&channelIDs)
	return channelIDs
}

func GetInstructorIDs(db *gorm.DB) []string {
	var instructorIDs []string
	db.Table("teachers").Select("user_id").Find(&instructorIDs)
	return instructorIDs
}
