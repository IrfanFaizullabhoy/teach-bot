package main

import (
	"fmt"
	"os"

	//"github.com/gorilla/schema"
	_ "github.com/lib/pq"
	"github.com/nlopes/slack"
)

func Run() {
	//WelcomeToTeamTest("U3YF3JM35")
	//SendTestMessage(api, "#teacher-test", "Here to help...")
}

func GetSlackClient() *slack.Client {
	token := os.Getenv("SLACK_TOKEN")
	fmt.Println("connecting to client w token" + token)
	api := slack.New(token)
	return api
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

func StartInstructorConversation(userID, name string) {
	api := GetSlackClient()
	group, err := api.GetGroupInfo("instructors-and-" + name)
	if group == nil || err != nil {
		group, err = api.CreateGroup("instructors-and-" + name)
		check(err)
		instructors := GetInstructors()
		for _, instructor := range instructors {
			_, _, err = api.InviteUserToGroup(group.Name, instructor.ID)
			check(err)
		}
		_, _, err := api.InviteUserToGroup(group.Name, userID)
		check(err)
	}
}
