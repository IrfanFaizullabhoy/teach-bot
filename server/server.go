package main

import (
	"fmt"
	"os"

	//"github.com/gorilla/schema"
	_ "github.com/lib/pq"
	"github.com/nlopes/slack"
)

func Run() {
	InitializeStudentMap()
	//RegisterAll()

	//SendTestMessage(api, "#teacher-test", "Here to help...")
	//EventChannel = make(chan Event)
}

func GetSlackClient() *slack.Client {
	token := os.Getenv("SLACK_TOKEN")
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

func ChannelExists(channelName string) (bool, string) {
	api := GetSlackClient()
	channels, err := api.GetChannels(false)
	check(err)
	for _, channel := range channels {
		if channel.Name == channelName {
			return true, channel.ID
		}
	}
	return false, ""
}

func StartInstructorConversation(userID, name string) {
	api := GetSlackClient()
	groupID := FindGroupByName("instructors-and-" + name)
	if groupID == "" {
		group, err := api.CreateGroup("instructors-and-" + name)
		instructors := GetInstructors()
		for _, instructor := range instructors {
			_, _, err = api.InviteUserToGroup(group.ID, instructor.ID)
			check(err)
		}
		groupID = group.ID
		_, _, err = api.InviteUserToGroup(groupID, userID)
		check(err)
	}
	api.OpenGroup(groupID)
}
