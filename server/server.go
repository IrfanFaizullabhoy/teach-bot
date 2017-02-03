package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/nlopes/slack"
	"github.com/robfig/cron"

	_ "github.com/lib/pq"
)

var RegisteredChannels []string
var pg *sql.DB

func SetupTable(db *sql.DB, tableName string) {
	_, err := db.Exec("CREATE TABLE users (channel text PRIMARY KEY, session text)")
	if err != nil {
		log.Printf("Error inserting into DB: %+v", err)
		return
	}
}

func GetUser(db *sql.DB, channelID string) *User {
	var session string
	row := db.QueryRow("SELECT session FROM users WHERE channel = $1", channelID)
	err := row.Scan(&session)
	if err != nil {
		return nil
	}
	user := new(User)
	user.ChannelID = channelID
	user.MuncherySession = session
	return user
}

func GetUsers(db *sql.DB, api *slack.Client) (users []*User) {
	IMs, _ := api.GetIMChannels()
	for _, IM := range IMs {
		user := GetUser(db, IM.ID)
		if user != nil {
			users = append(users, user)
		}
	}
	return users
}

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

func RegisterCronJob(api *slack.Client, db *sql.DB) {
	c := cron.New()
	// gonna have to figure out timezones
	c.AddFunc("0 0 21 * * MON-FRI", func() { runCronPost(api, db) })
	c.Start()
}

func runCronPost(api *slack.Client, db *sql.DB) {
	users := GetUsers(db, api)
	for _, user := range users {
		user = user
		//go menuPost(user.MuncherySession, api, user.ChannelID)
	}
}

func Run() {
	api := ConnectToSlack()
	fmt.Println("connected to slack")
	RegisterChannels(api)
	fmt.Println("registered channels")
	SendTestMessage(api, "#teacher-test", "Here to help...")
	atMB := GetAtTeachId(api)
	RegisterCronJob(api, pg)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		runCronPost(api, pg)
	})
	go http.ListenAndServe(":8080", nil)

	Respond(api, atMB)
}

func ConnectToSlack() *slack.Client {
	token := os.Getenv("SLACK_TOKEN")
	api := slack.New(token)
	return api
}

func GetAtTeachId(api *slack.Client) string {
	users, _ := api.GetUsers()
	for _, user := range users {
		if user.IsBot && user.Name == "teach" {
			return "<@" + user.ID + ">"
		}
	}
	return "Couldn't find teach"
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

					/* --------------- MENU CONVERSATION ---------------*/
					case strings.Contains(strings.ToLower(ev.Text), "menu") && !strings.Contains(strings.ToLower(ev.Text), "register"):
						if !ChannelExists(ev.Channel) {
							api.PostMessage(ev.Channel, "Please speak with `@teach-bot` in your direct message channel with `teach-bot`", params)
						} else {
							user := GetUser(pg, ev.Channel)
							if user == nil {
								api.PostMessage(ev.Channel, "Hi", params)
							} else {
								api.PostMessage(ev.Channel, "Hey", params)
								//menuPost(user.MuncherySession, api, ev.Channel)
							}
						}

					/* --------------- ORDER CONVERSATION ---------------*/
					case strings.Contains(strings.ToLower(ev.Text), "order") && !strings.Contains(strings.ToLower(ev.Text), "register"):
						if !ChannelExists(ev.Channel) {
							api.PostMessage(ev.Channel, "Please speak with `@teach-bot` in your direct message channel with `@teach-bot`", params)
						} else {
							ids, parseError := ParseOrder(ev.Text)
							if ids == nil || parseError {
								api.PostMessage(ev.Channel, "Sorry, didn't understand `", params)
							} else {
								api.PostMessage(ev.Channel, "Hi", params)
								//user := GetUser(pg, ev.Channel)
								//addToBasket(user.MuncherySession, ids)
								//checkout(user.MuncherySession)
							}
						}

					/* -------------- REGISTER CONVERSATION ---------------*/
					case strings.Contains(strings.ToLower(ev.Text), "register"):
						if !ChannelExists(ev.Channel) {
							params := slack.PostMessageParameters{}
							api.PostMessage(ev.Channel, "Yes", params)
						} else {
							muncherySessionID, skip := ParseRegistration(ev.Text, api, ev.Channel) // TODO
							if skip {
								api.PostMessage(ev.Channel, "Sorry`"+muncherySessionID+"` was not valid", params)
								break
							}
							if true {
								api.PostMessage(ev.Channel, "Sorry`"+muncherySessionID+"` was not valid", params)
							} else {
								api.PostMessage(ev.Channel, "Perfect", params)
								user := new(User)
								user.ChannelID = ev.Channel
								//user.MuncherySession = muncherySessionID
								//RegisterUserToDB(pg, user)
							}
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

func ParseRegistration(messageBody string, api *slack.Client, channel string) (string, bool) {
	params := slack.PostMessageParameters{}
	registrationText := strings.Split(messageBody, " ")
	if len(registrationText) < 2 {
		api.PostMessage(channel, "Looks", params)
		return "", false
	}
	if strings.ToLower(registrationText[0]) != "register" {
		api.PostMessage(channel, "Looks", params)
		return "", true
	}
	var token string
	for i, strings := range registrationText {
		if i >= 1 {
			token = token + strings
		}
	}
	return token, false
}

func ParseOrder(order string) ([]int, bool) {
	orders := strings.Split(order, " ")
	var orderNums []int
	for j, order := range orders {
		if j == 0 {

		} else {
			order = strings.Replace(order, ",", "", -1)
			i, err := strconv.Atoi(order)
			if err != nil {
				return nil, true
			}
			orderNums = append(orderNums, i)
		}
	}
	return orderNums, false
}
