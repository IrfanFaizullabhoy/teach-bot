package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Student struct {
	gorm.Model
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	Email          string `json:"email"`
	UserName       string `json:"username"`
	UserID         string `json:"user_id"`
	HashedPassword string `json:"hashed_password"`
	ChannelID      string `json:"channel_id"`
}

type Teacher struct {
	gorm.Model
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	UserName       string `json:"username"`
	UserID         string `json:"user_id"`
	Email          string `json:"email"`
	HashedPassword string `json:"hashed_password"`
	ChannelID      string `json:"channel_id"`
}

type LoginMsg struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
