package main

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	"github.com/nlopes/slack"

	_ "github.com/lib/pq"
)

var RegisteredChannels []string
var pg *sql.DB

func ChannelExists(channelName string) bool {
	for _, channel := range RegisteredChannels {
		if channelName == channel {
			return true
		}
	}
	return false
}

func RegisterChannels(api *slack.Client) {
	RegisteredChannels = make([]string, 0)
	IMs, _ := api.GetIMChannels()
	for _, IM := range IMs {
		RegisteredChannels = append(RegisteredChannels, IM.ID)
	}
}

func Run() {
	api := ConnectToSlack()
	fmt.Println("connected to slack")
	RegisterChannels(api)
	fmt.Println("registered channels")
	SendTestMessage(api, "#teacher-test", "Here to help...")
	RegisterCronJob(api)
	Respond(api, "")
}

func ConnectToSlack() *slack.Client {
	token := os.Getenv("SLACK_TOKEN")
	api := slack.New(token)
	return api
}

func SendTestMessage(api *slack.Client, channelName string, messageText string) {
	params := slack.PostMessageParameters{}
	channelID, timestamp, err := api.PostMessage(channelName, messageText, params)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	fmt.Printf("Message successfully sent to channel %s at %s", channelID, timestamp)
}

func PostAnonymousQuestion(api *slack.Client, channelName string, messageText string) {
	params := slack.PostMessageParameters{}
	messageText = "Someone posted an anonymous question: ```" + messageText + "```"
	_, _, err := api.PostMessage(channelName, messageText, params)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
}

func Respond(api *slack.Client, atBot string) {
	rtm := api.NewRTM()
	go rtm.ManageConnection()
	for {
		select {
		case msg := <-rtm.IncomingEvents:
			switch ev := msg.Data.(type) {
			case *slack.MessageEvent:
				RegisterChannels(api) //switch to change to a more efficient version
				params := slack.PostMessageParameters{}
				if ev.Msg.BotID == "" {
					switch {
					/* --------------- ANONYMOUS QUESTION ---------------*/
					case strings.Contains(strings.ToLower(ev.Text), "anonymous"):
						if !ChannelExists(ev.Channel) {
							api.PostMessage(ev.Channel, "Please speak with `@teach-bot` in your direct message channel with `teach-bot`", params)
						} else {
							PostAnonymousQuestion(api, "#teacher-test", ev.Text)
							api.PostMessage(ev.Channel, "Hi - I've posted your anonymous question", params)
						}

					/* -------------- REGISTER STUDENT CONVERSATION ---------------*/
					case strings.Contains(strings.ToLower(ev.Text), "register"):
						if !ChannelExists(ev.Channel) {
							params := slack.PostMessageParameters{}
							api.PostMessage(ev.Channel, "Please register with `@teach-bot` in your direct message channel with `teach-bot`", params)
						} else {
							api.PostMessage(ev.Channel, "Enter your `Firstname Lastname` as shown", params)
							firstName, lastName, err := ParseName(ev.Text)
							student := Student{UserID: ev.User, UserName: ev.Username, FirstName: firstName, LastName: lastName, ChannelID: ev.Channel}
							CreateStudent(student, db)
						}
					/*  ------------------ NONE OF THE ABOVE ---------------- */
					default:
						if !ChannelExists(ev.Channel) {
							params := slack.PostMessageParameters{}
							api.PostMessage(ev.Channel, "You must register", params)
						} else {
							params := slack.PostMessageParameters{}
							api.PostMessage(ev.Channel, "Hi", params)
						}
					}
				}
			case *slack.RTMError:
				fmt.Printf("Error: %s\n", ev.Error())
			default:
				// Ignore other events..
				// fmt.Printf("Unexpected: %v\n", msg.Data)
			}
		}
	}
}
