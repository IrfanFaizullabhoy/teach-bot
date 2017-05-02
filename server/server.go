package main

import (
	"fmt"

	//"github.com/gorilla/schema"
	_ "github.com/lib/pq"
	"github.com/nlopes/slack"
)

var TeamDemoMap map[string]bool

func Run() {
	InitializeStudentMap()
	InitializeTeamMap()
	InitializeDemoMaps()
	//SendTestMessage(api, "#teacher-test", "Here to help...")
	//EventChannel = make(chan Event)
}

/*
func GetSlackClient() *slack.Client {
	token := os.Getenv("SLACK_TOKEN")
	api := slack.New(token)
	return api
}

func GetSlackBotClient() *slack.Client {
	token := os.Getenv("SLACKBOT_TOKEN")
	api := slack.New(token)
	return api
}
*/

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

/*
//U4LC8A0TH
func FakeAssignmentAssign() {
	api := GetSlackClient()
	fileparams := slack.FileUploadParameters{File: "/mounted-volume/Assignment1.pdf", Filename: "Assignment1.pdf", Channels: []string{"G4K1FBBK3"}}
	_, err := api.UploadFile(fileparams)
	check(err)
	api.PostMessage("G4K1FBBK3", "@irfan assigned the above assignment: `Assignment1.pdf`", slack.PostMessageParameters{})
	params := slack.NewPostMessageParameters()
	attachment := slack.Attachment{CallbackID: "reminder", Fallback: "service not working properly"}
	attachmentMondayAction := slack.AttachmentAction{Name: "yes", Text: "Yes Please!", Type: "button"}
	attachmentOtherAction := slack.AttachmentAction{Name: "no", Text: "No thanks", Type: "button"}
	attachment.Actions = append(attachment.Actions, attachmentMondayAction)
	attachment.Actions = append(attachment.Actions, attachmentOtherAction)
	params.Attachments = append(params.Attachments, attachment)
	params.LinkNames = 1
	api.PostMessage("G4K1FBBK3", "It's due on *Mon Mar 20 16:00:00 2017*, would you like reminders?", params)
}

//U4LC8A0TH
func FakeAssignmentZip() {
	api := GetSlackClient()
	fileparams := slack.FileUploadParameters{File: "/mounted-volume/Assignment1.zip", Filename: "Assignment1.zip", Channels: []string{"G4K0E8WJU"}}
	_, err := api.UploadFile(fileparams)
	check(err)
	api.PostMessage("G4K0E8WJU", "Here is the zip folder for the following assignment: `Assignment1.pdf`", slack.PostMessageParameters{})
	params := slack.NewPostMessageParameters()
	attachment := slack.Attachment{CallbackID: "reminderAssignment", Fallback: "service not working properly"}
	attachmentMondayAction := slack.AttachmentAction{Name: "yes", Text: "Yes", Type: "button"}
	attachmentOtherAction := slack.AttachmentAction{Name: "no", Text: "No", Type: "button"}
	attachment.Actions = append(attachment.Actions, attachmentMondayAction)
	attachment.Actions = append(attachment.Actions, attachmentOtherAction)
	params.Attachments = append(params.Attachments, attachment)
	api.PostMessage("G4K0E8WJU", "The following students did not submit yet: `@martin` `@adena` -- Should I let them know?", params)
}
*/
func StartInstructorConversation(userID, name, teamID string) {
	team := GetTeam(teamID)
	appConn := slack.New(team.Token)
	groupID := FindGroupByName("instructors-and-"+name, appConn)
	if groupID == "" {
		group, err := appConn.CreateGroup("instructors-and-" + name)
		instructors := GetInstructors(teamID)
		for _, instructor := range instructors {
			_, _, err = appConn.InviteUserToGroup(group.ID, instructor.ID)
			check(err)
		}
		groupID = group.ID
		_, _, err = appConn.InviteUserToGroup(groupID, userID)
		check(err)
	}
	appConn.OpenGroup(groupID)
}
