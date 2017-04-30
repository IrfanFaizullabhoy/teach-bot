package main

import (
	"fmt"
	"github.com/nlopes/slack"
	"time"
)

//Create User & Team
//Welcome to Slack ?

//Install teachbot to team
// registerAll
//Same username and password
// create field-trips channel
// add teach bot to midterm

//Create Google Drive & Link for Homework
//allow google drive to import file

//Add Team to DemoTeams map
var teamIDtoTS map[string]string
var teamIDtoChannelID map[string]string

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
	screenshotFileName := "assignment1name"

	file, _, _, err := botConn.GetFileInfo(fileSharedEvent.File.ID, 1, 1) //returns file with one comment/onepage
	check(err)

	if file.Name == screenshotFileName {
		return
	}

	fmt.Println("in demo download assignment")
	fileChannelID := GetUser(file.User).ChannelID
	fmt.Println(fileChannelID)
	botConn.PostMessage(fileChannelID, "Awesome, sharing with students now!", slack.PostMessageParameters{})
	fmt.Println("done")
	screenshotFilePath := "../mounted-volume/sample.png"
	params := slack.FileUploadParameters{File: screenshotFilePath, Filename: screenshotFileName, Channels: []string{fileChannelID}, InitialComment: "This is what students are seeing right now!", Title: "Student's Perspective"}
	_, err = botConn.UploadFile(params)
	check(err)

	go func() {
		time.Sleep(7 * time.Second)
		DemoAcknowledge(teamID, team.InstallerID)
	}()
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

func DemoAcknowledge(teamID, userID string) {
	// IN TEACH-BOT PRIV CHANNEL
	fmt.Println("in this one")
	team := GetTeam(teamID)
	botConn := slack.New(team.BotToken)
	user := GetUser(userID)
	fieldTripID := GetOrCreateChannel("field-trips", teamID)
	fieldTrip := "Next we will use the `/acknowledge` feature to make an announcement and see which students have seen the announcement in real time. \n *Go to the* #field-trips *channel to try it out.*"
	botConn.PostMessage(user.ChannelID, fieldTrip, slack.PostMessageParameters{LinkNames: 1})

	// IN #FIELD-TRIPS
	hogwarts := "Type `/acknowledge` followed by an announcement you want to make. \n This will post the announcement in the current channel. *For example:* \n `/acknowledge the dress code for next week's trip to Hogwarts is wizard-casual`."
	botConn.PostMessage(fieldTripID, hogwarts, slack.PostMessageParameters{})
}

func GetOrCreateChannel(channelName, teamID string) string {
	team := GetTeam(teamID)
	appConn := slack.New(team.Token)
	channels, err := appConn.GetChannels(true)
	check(err)
	for _, channel := range channels {
		if channel.Name == channelName {
			return channel.ID
		}
	}
	channel, err := appConn.CreateChannel(channelName)
	check(err)
	for _, user := range GetUsers(teamID) {
		appConn.InviteUserToChannel(channel.ID, user.ID)
	}
	return channel.ID
}

func DemoAcknowledgePost(teamID, userID, channelID, text string) {
	team := GetTeam(teamID)
	botConn := slack.New(team.BotToken)
	user := GetUser(userID)
	fmt.Println("demo ack post")

	//Respond with acknowledge button
	params := slack.PostMessageParameters{EscapeText: true}
	attachment := slack.Attachment{CallbackID: "acknowledge", Fallback: "acknowledge service not working properly"}
	attachment.Actions = append(attachment.Actions, slack.AttachmentAction{Name: "acknowledge", Text: "Acknowledge", Type: "button"})
	params.Attachments = append(params.Attachments, attachment)
	botConn.PostMessage(channelID, text+" - @"+user.Name, params)

	go func() {
		channel, ts, err := botConn.PostMessage(channelID, "Responses so far: 0/5", slack.PostMessageParameters{})
		check(err)
		teamIDtoTS[teamID] = ts
		teamIDtoChannelID[teamID] = channelID
		time.Sleep(1 * time.Second)
		botConn.UpdateMessage(channel, ts, "Responses so far: 1/5")
		time.Sleep(500 * time.Millisecond)
		botConn.UpdateMessage(channel, ts, "Responses so far: 2/5")
		time.Sleep(900 * time.Millisecond)
		botConn.UpdateMessage(channel, ts, "Responses so far: 3/5")
		time.Sleep(2 * time.Second)
		botConn.PostMessage(channelID, "*Let's look at another feature* while we're waiting for your students to finish their assignments and acknowledge your message.", slack.PostMessageParameters{})
		time.Sleep(5 * time.Second)
		midtermID := GetOrCreateChannel("midterm", teamID)
		messageText := "i'm kinda embarrassed that I can't find this information, but what's the room number for the midterm?"
		botConn.PostMessage(midtermID, "Someone posted an anonymous question: \n ```"+messageText+"```", slack.PostMessageParameters{})
		time.Sleep(500 * time.Millisecond)
		botConn.PostMessage(channelID, "It looks like a student has asked an anonymous question by using teach-bot's `/anonymousQuestion` feature. *Go to the* #midterm *channel* see what it is, and then answer. \n (hint: the answer is 42)", slack.PostMessageParameters{LinkNames: 1})
	}()
	//acknowledgeMsg := AcknowledgeMessage{UserID: userID, Timestamp: ts, ChannelID: channelID}
}

func IsInMidterms(MessageChannelEvent Event, teamID string) {
	fmt.Println("hi")
	team := GetTeam(teamID)
	if MessageChannelEvent.User == team.InstallerID &&
		MessageChannelEvent.Channel == GetOrCreateChannel("midterm", teamID) {
		fmt.Println("done!")
	} else {
		return
	}
}

/* Acknowledge flow
faizulla bot - now follow me into the Midterm channel I just created!
Notify all the students that there is a midterm with the following - see how the `/acknowledge` feature works!
```
/acknowledge There will be a midterm in 2 weeks on Units 1 - 5 @channel !
Click acknowledge when you get this!
```
*/ //ttrttrrrerererr4ewë44

//https://drive.google.com/open?id=0B38oEsv5Mt0-cC04NjJURWVvaDg
