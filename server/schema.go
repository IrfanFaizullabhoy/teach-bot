package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/nlopes/slack"
)

type User struct {
	gorm.Model
	ID        string            `json:"id" gorm:"column:user_id" gorm:"unique"`
	Name      string            `json:"name"`
	ChannelID string            `json:"channel_id"`
	IsBot     bool              `json:"is_bot"`
	Profile   slack.UserProfile `json:"profile"`
	Role      string            `json:"role"`
}

type Assignment struct {
	gorm.Model
	DueDate  string `json:"due_date"`
	Link     string `json:"link"`
	FileName string `json:"file_name"`
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
	SSLCheck    int    `schema:"ssl_check"`
}

type OuterEvent struct {
	Token       string   `json:"token"`
	Challenge   string   `json:"challenge"`
	Type        string   `json:"type"`
	Event       Event    `json:"event"`
	TeamID      string   `json:"team_id"`
	APIAppID    string   `json:"api_app_id"`
	AuthedUsers []string `json:"authed_users"`
}

type OAuthResponseIncomingWebhook struct {
	URL              string `json:"url"`
	Channel          string `json:"channel"`
	ChannelID        string `json:"channel_id,omitempty"`
	ConfigurationURL string `json:"configuration_url"`
}

type OAuthResponseBot struct {
	BotUserID      string `json:"bot_user_id"`
	BotAccessToken string `json:"bot_access_token"`
}

type OAuthResponse struct {
	AccessToken     string                       `json:"access_token"`
	Scope           string                       `json:"scope"`
	TeamName        string                       `json:"team_name"`
	TeamID          string                       `json:"team_id"`
	IncomingWebhook OAuthResponseIncomingWebhook `json:"incoming_webhook"`
	Bot             OAuthResponseBot             `json:"bot"`
	UserID          string                       `json:"user_id,omitempty"`
	//SlackResponse
}

type ChallengeResponse struct {
	Token     string `json:"token"`
	Challenge string `json:"challenge"`
	Type      string `json:"type"`
}

type Event struct {
	Type    string `json:"type"`
	EventTS string `json:"event_ts"`
	User    string `json:"user"`
	TS      string `json:"ts"`
	File    File   `json:"file"`
}

type File struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Title      string `json:"title"`
	User       string `json:"user"`
	URLPrivate string `json:"url_private"`
	Filetype   string `json:"filetype"`
	Size       string `json:"size"`
}

type GroupCreateResponse struct {
	OK    bool  `json:"ok"`
	Group Group `json:"group"`
}

type ActionResponse struct {
	Payload string `json:"payload"`
}

type Group struct {
	ID   string `json:"id"`
	name string `json:"name"`
}
