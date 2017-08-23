package main

import (
	"fmt"

	_ "github.com/lib/pq"
	"github.com/nlopes/slack"
)


func PostAnonymousQuestion(messageText, teamID string) {
	team := GetTeam(teamID)
	appConn := slack.New(team.Token)
	channelID := ""
	channels, err := appConn.GetChannels(false)
	check(err)
	for _, channel := range channels {
		if channel.Name == "questions" {
			channelID = channel.ID
		}
	}
	if channelID == "" { // DOESNT EXIST YET, invite everyone
		channel, _ := appConn.CreateChannel("questions")
		appConn.SetChannelPurpose(channel.ID, "A Channel for you to ask questions that apply to the whole class, to ask anonymously type `/anonymousQuestion` followed by your question")
		for _, user := range GetUsers(teamID) {
			appConn.InviteUserToChannel(channel.ID, user.ID)
		}
		channelID = channel.ID
	}

	params := slack.PostMessageParameters{}
	messageText = "Someone posted an anonymous question: ```" + messageText + "```"
	_, _, err = appConn.PostMessage(channelID, messageText, params)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
}
