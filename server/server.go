package main

import (
	"fmt"
	"os"

	//"github.com/gorilla/schema"
	_ "github.com/lib/pq"
	"github.com/nlopes/slack"
)

//var EventChannel chan Event

func Run() {
	//WelcomeToTeamTest("U3YF3JM35")
	//SendTestMessage(api, "#teacher-test", "Here to help...")
	//EventChannel = make(chan Event)
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

func GroupExists(groupName string) (bool, string) {
	api := GetSlackClient()
	groups, err := api.GetGroups(false)
	check(err)
	for _, group := range groups {
		if group.Name == groupName {
			return true, group.ID
		}
	}
	return false, ""
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
	exists, groupID := GroupExists("instructors-and-" + name)
	if !exists {
		group, err := api.CreateGroup("instructors-and-" + name)
		instructors := GetInstructors()
		for _, instructor := range instructors {
			_, _, err = api.InviteUserToGroup(group.ID, instructor.ID)
			check(err)
			groupID = group.ID
		}
		_, _, err = api.InviteUserToGroup(group.Name, userID)
		check(err)
	}
	api.OpenGroup(groupID)
}
