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
	UserID         string `gorm:"unique", json:"user_id"`
	HashedPassword string `json:"hashed_password"`
	ChannelID      string `json:"channel_id"`
}

type Teacher struct {
	gorm.Model
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	UserName       string `json:"username"`
	UserID         string `gorm:"unique", json:"user_id"`
	Email          string `json:"email"`
	HashedPassword string `json:"hashed_password"`
	ChannelID      string `json:"channel_id"`
}

type LoginMsg struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SlashPayload struct {
	Token       string `json:"token"`
	TeamID      string `json:"team_id"`
	TeamDomain  string `json:"team_domain"`
	ChannelID   string `json:"channel_id"`
	ChannelName string `json:"channel_name"`
	UserID      string `json:"user_id"`
	UserName    string `json:"user_name"`
	Command     string `json:"command"`
	Text        string `json:"text"`
	ResponseURL string `json:"response_url"`
}

type GroupCreateResponse struct {
	OK    bool  `json:"ok"`
	Group Group `json:"group"`
}

type Group struct {
	ID   string `json:"id"`
	name string `json:"name"`
}
