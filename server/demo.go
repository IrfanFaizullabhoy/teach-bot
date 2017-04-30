package main

import (
	"fmt"
	"github.com/nlopes/slack"
)

//Create User & Team
//Welcome to Slack ?

//Install teachbot to team
// registerAll
//Same username and password

//Create Google Drive & Link for Homework
//allow google drive to import file

//Add Team to DemoTeams map

func IsDemoTeam(teamID string) bool {
	return TeamDemoMap[teamID]
}

func CheckPresence(userID, teamID string) {
	userID = "U3YF3JM35"
	teamID = "T3Z7YKN07"
	team := GetTeam(teamID)
	botConn := slack.New(team.BotToken)
	for {
		presence, err := botConn.GetUserPresence(userID)
		check(err)
		fmt.Println("in loop")
		if presence.Presence == "active" {
			user := GetUser(userID)
			botConn.PostMessage(user.ChannelID, "Gotcha!", slack.PostMessageParameters{})
			fmt.Println("posted message")
			return
		}
		time.Sleep(2 * time.Second)
	}
}

func DemoDateInteractive(userID, channelID, teamID string) {
	team := GetTeam(teamID)
	botConn := slack.New(team.BotToken)
	params := slack.NewPostMessageParameters()
	attachment := slack.Attachment{CallbackID: "assignment_due", Fallback: "service not working properly"}
	attachment5MinAction := slack.AttachmentAction{Name: "5min", Text: "5 Minutes from Now", Type: "button"}
	attachment.Actions = append(attachment.Actions, attachment5MinAction)
	params.Attachments = append(params.Attachments, attachment)
	_, _, err := botConn.PostMessage(channelID, "When would you like the assignment to be due?", params)
	check(err)
}

func DemoHandleDate(attachmentDateAction slack.AttachmentActionCallback) {
	team := GetTeam(attachmentDateAction.Team.ID)
	//DEMO
	botConn := slack.New(team.BotToken)
	if len(attachmentDateAction.Actions) == 1 {
		if attachmentDateAction.Actions[0].Name == "5min" {
			params := slack.PostMessageParameters{}
			botConn.PostMessage(attachmentDateAction.Channel.ID, "Great, it'll be due in 5 mintues! \n Paste the following link as a sample assignment ```https://docs.google.com/a/usc.edu/document/d/1JvNwiTnMiqWtZGeWkuqGQQDG0CaZcdFPwpxIjgRoqKw/edit?usp=drive_web``` ", params)
			// go routine that cleans up db if there is no file that gets uploaded in 5 minutes
		} else {
			//TODO
		}
	}
}

func DemoDownloadAssignment(fileSharedEvent Event, teamID string) {
	// Sends Picture of what was sent to users
	team := GetTeam(teamID)
	botConn := slack.New(team.BotToken)
	fmt.Println("in demo download assignment")
	file, _, _, err := botConn.GetFileInfo(fileSharedEvent.File.ID, 1, 1)
	check(err)
	fileChannelID := GetUser(file.User).ChannelID
	fmt.Println(fileChannelID)
	botConn.PostMessage(fileChannelID, "Awesome, sharing with students now!", slack.PostMessageParameters{})
	fmt.Println("done")
	screenshotFilePath := "../mounted-volume/sample.png"
	var screenshotFileName string
	params := slack.FileUploadParameters{File: screenshotFilePath, Filename: screenshotFileName, Channels: []string{fileChannelID}, InitialComment: "This is what students are seeing right now!", Title: "Student's Perspective"}
	_, err = botConn.UploadFile(params)
	check(err)
	//DemoViewSubmissions(fileChannelID, teamID)
	return
}

func DemoViewSubmissions(channelID, teamID string) {
	team := GetTeam(teamID)
	botConn := slack.New(team.BotToken)
	params := slack.NewPostMessageParameters()
	attachment := slack.Attachment{CallbackID: "submission_type", Fallback: "service not working properly"}
	attachmentZipAction := slack.AttachmentAction{Name: "zip", Text: "Zip File", Type: "button"}
	attachmentDriveAction := slack.AttachmentAction{Name: "drive", Text: "Drive Folder", Type: "button"}
	attachment.Actions = append(attachment.Actions, attachmentZipAction)
	attachment.Actions = append(attachment.Actions, attachmentDriveAction)
	params.Attachments = []slack.Attachment{attachment}

	botConn.PostMessage(channelID, "Would you like the submissions as a `.zip` file or in a `Google Drive` Folder?", params)
}

func DemoHandleViewSubmission(attachmentSubmissionViewAction slack.AttachmentActionCallback) {
	team := GetTeam(attachmentSubmissionViewAction.Team.ID)
	//DEMO
	botConn := slack.New(team.BotToken)
	if len(attachmentSubmissionViewAction.Actions) == 1 {
		if attachmentSubmissionViewAction.Actions[0].Name == "drive" {
			botConn.PostMessage(attachmentSubmissionViewAction.Channel.ID, "Here's a Google Drive Folder with the submissions... \n `Assignment 1`: https://drive.google.com/open?id=0B38oEsv5Mt0-cC04NjJURWVvaDg", slack.PostMessageParameters{})
			// go routine that cleans up db if there is no file that gets uploaded in 5 minutes
		} else if attachmentSubmissionViewAction.Actions[0].Name == "zip" {
			zipFilePath := "../mounted-volume/assignment1.zip"
			var screenshotFileName string
			params := slack.FileUploadParameters{File: zipFilePath, Filename: screenshotFileName, Channels: []string{attachmentSubmissionViewAction.Channel.ID}, InitialComment: "Take a look!", Title: "Assignment 1 Submissions"}
			_, err := botConn.UploadFile(params)
			check(err)
			//TODO
		}
	}
}

func DemoAcknowledge() {

	"Next we will use the /acknowledge feature to make an announcement and see which students have seen the announcement in real time. \nType /acknowledge followed by an announcement you want to make. \n For example: /acknowledge the British are coming."
}

func DemoAcknowledgePost(teamID, userID, channelID, text string) {
	team := GetTeam(teamID)
	botConn := slack.New(team.BotToken)
	user := GetUser(userID)

	//Respond with acknowledge button
	params := slack.PostMessageParameters{EscapeText: true}
	attachment := slack.Attachment{CallbackID: "acknowledge", Fallback: "acknowledge service not working properly"}
	attachment.Actions = append(attachment.Actions, slack.AttachmentAction{Name: "acknowledge", Text: "Acknowledge", Type: "button"})
	params.Attachments = append(params.Attachments, attachment)
	botConn.PostMessage(channelID, text+" - @"+user.Name, params)
	//acknowledgeMsg := AcknowledgeMessage{UserID: userID, Timestamp: ts, ChannelID: channelID}
}

/* Acknowledge flow
faizulla bot - now follow me into the Midterm channel I just created!
Notify all the students that there is a midterm with the following - see how the `/acknowledge` feature works!
```
/acknowledge There will be a midterm in 2 weeks on Units 1 - 5 @channel !
Click acknowledge when you get this!
```
*/

//https://drive.google.com/open?id=0B38oEsv5Mt0-cC04NjJURWVvaDg
