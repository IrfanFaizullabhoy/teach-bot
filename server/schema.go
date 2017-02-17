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
	Token       string `schema:"token"`
	TeamID      string `schema:"team_id"`
	TeamDomain  string `schema:"team_domain"`
	ChannelID   string `schema:"channel_id"`
	ChannelName string `schema:"channel_name"`
	UserID      string `schema:"user_id"`
	UserName    string `schema:"user_name"`
	Command     string `schema:"command"`
	Text        string `schema:"text"`
	ResponseURL string `schema:"response_url"`
}

type GroupCreateResponse struct {
	OK    bool  `json:"ok"`
	Group Group `json:"group"`
}

type Group struct {
	ID   string `json:"id"`
	name string `json:"name"`
}
