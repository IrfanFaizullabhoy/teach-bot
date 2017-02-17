package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	//"github.com/gorilla/schema"
	_ "github.com/lib/pq"
	"github.com/nlopes/slack"
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

//Change to a DB lookup for student / teacher / both channels
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

func StartInstructorConversation(userID string) {
	//instructorIDs := GetInstructorIDs(db) //implement this
	apiUrl := "https://slack.com"
	resource := "/api/groups.create"
	data := url.Values{}
	userID = "U42D42KLG"
	data.Set("token", "xoxp-135270668007-134513633107-139843303798-8a4a0f1918cd7cd754a65f3540777d95")
	data.Add("name", "the_one_and_only")

	u, _ := url.ParseRequestURI(apiUrl)
	u.Path = resource
	u.RawQuery = data.Encode()
	urlStr := fmt.Sprintf("%v", u) // "https://api.com/user/?name=foo&surname=bar"

	client := &http.Client{}
	r, _ := http.NewRequest("POST", urlStr, nil)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	resp, err := client.Do(r)

	check(err)

	var groupCreateResponse GroupCreateResponse
	body, err := ioutil.ReadAll(resp.Body)
	check(err)
	if err := json.Unmarshal(body, &groupCreateResponse); err != nil {
		panic(err)
	}

	fmt.Println("checkpt 3")

	groupID := groupCreateResponse.Group.ID
	instructors := []string{"U3YKBAK1S", "U42EVJF7E", "U3YK6EPV0", userID}
	for _, instructorID := range instructors {
		InviteInstructorToPrivateChannel(groupID, instructorID)
	}
}

func InviteInstructorToPrivateChannel(channelID, instructorID string) {
	apiUrl := "https://slack.com"
	resource := "/api/groups.invite"
	data := url.Values{}
	data.Set("token", "xoxp-135270668007-134513633107-139843303798-8a4a0f1918cd7cd754a65f3540777d95")
	data.Add("channel", channelID)
	data.Add("user", instructorID)

	u, _ := url.ParseRequestURI(apiUrl)
	u.Path = resource
	u.RawQuery = data.Encode()
	urlStr := fmt.Sprintf("%v", u) // "https://api.com/user/?name=foo&surname=bar"

	client := &http.Client{}
	r, _ := http.NewRequest("POST", urlStr, nil)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	_, err := client.Do(r)
	check(err)
}

//Change to be invoked as a go routine
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

					/* -------------- REGISTER TEACHER CONVERSATION ---------------*/
					case strings.Contains(strings.ToLower(ev.Text), "register") && strings.Contains(strings.ToLower(ev.Text), "teacher"):
						if !ChannelExists(ev.Channel) {
							params := slack.PostMessageParameters{}
							api.PostMessage(ev.Channel, "Please register with `@teach-bot` in your direct message channel with `teach-bot`", params)
						} else {
							//api.PostMessage(ev.Channel, "Enter your `Firstname Lastname` as shown", params)
							firstName, lastName, err := ParseTeacherName(ev.Text)
							if err != nil {
								api.PostMessage(ev.Channel, "Oops, didn't register properly! Try again `register fname lname`", params)
							}
							teacher := Teacher{UserID: ev.User, UserName: ev.Username, FirstName: firstName, LastName: lastName, ChannelID: ev.Channel}
							CreateTeacher(teacher, db)
						}

					/* -------------- REGISTER STUDENT CONVERSATION ---------------*/
					case strings.Contains(strings.ToLower(ev.Text), "register"):
						if !ChannelExists(ev.Channel) {
							params := slack.PostMessageParameters{}
							api.PostMessage(ev.Channel, "Please register with `@teach-bot` in your direct message channel with `teach-bot`", params)
						} else {
							//api.PostMessage(ev.Channel, "Enter your `Firstname Lastname` as shown", params)
							firstName, lastName, err := ParseStudentName(ev.Text)
							if err != nil {
								api.PostMessage(ev.Channel, "Oops, didn't register properly! Try again `register fname lname`", params)
							}
							student := Student{UserID: ev.User, UserName: ev.Username, FirstName: firstName, LastName: lastName, ChannelID: ev.Channel}
							CreateStudent(student, db)
						}

					/* -------------- ASSIGNMENT ASSIGNING CONVERSATION ---------------*/
					case strings.Contains(strings.ToLower(ev.Text), "assignment"):
						if !ChannelExists(ev.Channel) {
							params := slack.PostMessageParameters{}
							api.PostMessage(ev.Channel, "Please register with `@teach-bot` in your direct message channel with `teach-bot`", params)
						} else {
							api.PostMessage(ev.Channel, strconv.Itoa(ev.File.Size), params)
							client := &http.Client{}
							req, _ := http.NewRequest("GET", ev.File.URLPrivate, nil)
							req.Header.Add("Authorization", "Bearer xoxb-136568950452-8X180knozh1mI8hYPXqzYDZR")
							response, err := client.Do(req)

							if err != nil {
								panic("Request making error")
							}

							if response.StatusCode != 200 {
								fmt.Println("download error")
								panic(response.Status)
							}

							api.PostMessage(ev.Channel, response.Status, params)

							defer response.Body.Close()
							file := ev.File
							filePath := "/mounted-volume/" + file.Name
							tmpfile, err1 := os.Create(filePath)
							defer tmpfile.Close()

							if err1 != nil {
								panic("error creating the file")
							}
							file_content, err := ioutil.ReadAll(response.Body)
							if err != nil {
								panic("read file error")
							}
							size, err2 := tmpfile.Write(file_content)
							if err2 != nil {
								panic("write file error")
							}
							api.PostMessage(ev.Channel, strconv.Itoa(size), params)

							var channels []string
							channels = append(channels, ev.Channel)
							fileParams := slack.FileUploadParameters{Filename: file.Name, File: filePath, Filetype: file.Filetype, Channels: channels}
							api.UploadFile(fileParams)
							channels = GetStudentChannels(db)
							fileParams.Channels = channels
							api.UploadFile(fileParams)

							// get due date

							// there should be some kind of staging
						}

					case strings.Contains(strings.ToLower(ev.Text), "instructor"):
						if !ChannelExists(ev.Channel) {
							params := slack.PostMessageParameters{}
							api.PostMessage(ev.Channel, "Please register with `@teach-bot` in your direct message channel with `teach-bot`", params)
						} else {
							fmt.Println(ev.User)
							StartInstructorConversation("irf")
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
				// fmt.Printf("Unexpected: %v\n",  msg.Data)
			}
		}
	}
}
