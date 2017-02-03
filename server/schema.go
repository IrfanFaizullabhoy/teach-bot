package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"time"
)

type User struct {
	gorm.Model
	FirstName       string //`json:"first_name"`
	LastName        string //`json:"last_name"`
	PhoneNumber     string //`json:"phone_number"`
	ProfileImage    string //`json:"profile_image"`
	Email           string //`json:"email"`
	HashedPassword  string //`json:"hashed_password"`
	ChannelID       string //`json:"channel_id"`
	MuncherySession string //`json:"munchery-session"`
}

type LoginMsg struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Listing struct {
	gorm.Model
	PicURL   string   `json:"pic_url"`
	Address  string   `json:"address"`
	About    string   `json:"about"`
	Owner    string   `json:"owner"`
	Category string   `json:"category"`
	Reviews  []Review `json:"reviews"`
	//Payments []Payment `json:"payment"`
	//Title    String    `json:"title"`
}

type Booking struct {
	gorm.Model
	Listing   Listing   `json:"listing"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	User      User      `json:"user"`
}

type Review struct {
	gorm.Model
	Rating      int     `json:"rating"`
	Listing     Listing `json:"listing"`
	CreatedBy   User    `json:"created_by"`
	Description string  `json:"description"`
	Owner       string  `json:"owner"`
}
