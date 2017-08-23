package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/nlopes/slack"
)


func HandleFileShared(fileSharedEvent Event, teamID string) {
	if true {
		go DownloadAssignment(fileSharedEvent, teamID)
	} else {
		go DownloadSubmission(fileSharedEvent, teamID)
	}
}

func DownloadAssignment(fileSharedEvent Event, teamID string) {
	team := GetTeam(teamID)

	//DEMO
	if IsDemoTeam(teamID) {
		go DemoDownloadAssignment(fileSharedEvent, teamID)
		return
	}

	botConn := slack.New(team.BotToken)
	file, _, _, err := botConn.GetFileInfo(fileSharedEvent.File.ID, 1, 1) //returns file with one comment/onepage
	check(err)

	//FIND IT IN THE DB
	assignment := FindAssignment(file.User, file.Name, team)
	if assignment.Submissions == nil {
		return
	}
	// GOOGLE DOC
	if strings.Contains(file.URLPrivate, "google.com") {
		fmt.Println("Google Doc")
		return
	}

	//DOWNLOAD FILE
	client := &http.Client{}
	req, _ := http.NewRequest("GET", file.URLPrivate, nil)
	req.Header.Add("Authorization", "Bearer "+team.BotToken)
	response, err := client.Do(req)
	defer response.Body.Close()
	if err != nil {
		panic("Request making error")
	}
	if response.StatusCode != 200 {
		fmt.Println("download error")
		panic(response.Status)
	}

	//WRITE TO DISK
	filePath := "/mounted-volume/" + file.Name
	tmpfile, createErr := os.Create(filePath)
	check(createErr)
	defer tmpfile.Close()
	file_content, readErr := ioutil.ReadAll(response.Body)
	check(readErr)
	size, writeErr := tmpfile.Write(file_content)
	check(writeErr)
	if size != file.Size {
		fmt.Println("file size unequal")
	}

	assignment.FilePath = filePath
	assignment.Downloaded = true
	assignment.FileName = file.Name
	db.Save(&assignment)
	botConn.PostMessage(GetUser(fileSharedEvent.User).ChannelID, "Awesome, sharing with students now!", slack.PostMessageParameters{})
	ConfirmAndSendAssignment(assignment, botConn)
}

func ConfirmAndSendAssignment(assignment Assignment, botConn *slack.Client) {
	/*params := slack.NewPostMessageParameters()
	attachment := slack.Attachment{CallbackID: "confirm_assignment", Fallback: "service not working properly"}
	attachmentSendAction := slack.AttachmentAction{Name: "send", Text: "Send Assignment", Type: "button"}
	attachmentDeleteAction := slack.AttachmentAction{Name: "delete", Text: "This is Incorrect", Type: "button"}
	//AttachmentGroupAction
	attachment.Actions = append(attachment.Actions, attachmentSendAction)
	attachment.Actions = append(attachment.Actions, attachmentDeleteAction)
	params.Attachments = append(params.Attachments, attachment)
	_, _, err := api.PostMessage(channel, "Confirming, I should send out `"+
		assignment.FileName+"` due on *"+assignment.DueDate+"*", params)
	check(err)*/

	students := GetStudents()
	var channels []string
	for _, student := range students {
		channels = append(channels, student.ChannelID)
		botConn.PostMessage(student.ChannelID, "The following assignment is due on *"+assignment.DueDate+"*", slack.PostMessageParameters{})
	}
	params := slack.FileUploadParameters{File: assignment.FilePath, Filename: assignment.FileName, Channels: channels}
	_, err := botConn.UploadFile(params)
	check(err)

	botConn.PostMessage(GetUser(assignment.UserID).ChannelID, "Successfully posted the assignment!", slack.PostMessageParameters{})
}

func FindAssignment(userID, fileName string, team Team) Assignment {
	var assignments []Assignment
	params := slack.PostMessageParameters{}
	botConn := slack.New(team.BotToken)
	db.Where("user_id = ? AND team_id = ? AND downloaded = ?", userID, team.TeamID, false).Find(&assignments)
	if len(assignments) == 1 {
		user := GetUser(userID)
		botConn.PostMessage(user.ChannelID, "Got it! Will distribute: `"+fileName+"`", params)
		return assignments[0]
	} else if len(assignments) == 0 {
		return assignments[0]
	} else {
		fmt.Println(userID)
		fmt.Println("more than one assigment not downloaded... " + strconv.Itoa(len(assignments)))
		return WhichAssignment(userID, assignments)
	}
}

func WhichAssignment(userID string, assignments []Assignment) Assignment {
	return assignments[0]
}

func DateInteractive(userID, channelID, teamID string) {
	team := GetTeam(teamID)
	if IsDemoTeam(team.TeamID) {
		DemoDateInteractive(userID, channelID, teamID)
		return
	}
	botConn := slack.New(team.BotToken)
	params := slack.NewPostMessageParameters()
	attachment := slack.Attachment{CallbackID: "assignment_due", Fallback: "service not working properly"}
	attachmentMondayAction := slack.AttachmentAction{Name: "monday", Text: "Next Monday at 4pm", Type: "button"}
	attachmentOtherAction := slack.AttachmentAction{Name: "other", Text: "Other", Type: "button"}
	attachment.Actions = append(attachment.Actions, attachmentMondayAction)
	attachment.Actions = append(attachment.Actions, attachmentOtherAction)
	params.Attachments = append(params.Attachments, attachment)
	_, _, err := botConn.PostMessage(channelID, "When would you like the assignment to be due?", params)
	check(err)
}

func HandleDate(attachmentDateAction slack.AttachmentActionCallback) {
	team := GetTeam(attachmentDateAction.Team.ID)
	//DEMO
	if IsDemoTeam(team.TeamID) {
		DemoHandleDate(attachmentDateAction)
		return
	}
	botConn := slack.New(team.BotToken)
	if len(attachmentDateAction.Actions) == 1 {
		if attachmentDateAction.Actions[0].Name == "monday" {
			monday := GetNextMonday().Format(time.ANSIC)
			params := slack.PostMessageParameters{}
			assignment := Assignment{DueDate: monday, UserID: attachmentDateAction.User.ID, TeamID: team.TeamID, Downloaded: false}
			go CreateAssignment(assignment, db)
			botConn.PostMessage(attachmentDateAction.Channel.ID, "Great, it'll be due on *"+monday+"*\nPlease press the `+` sign below/to my left to attach the homework?", params)
			// go routine that cleans up db if there is no file that gets uploaded in 5 minutes
		} else {
			//TODO
		}
	}
}

func NoFileUploaded(user User) {

}

func checkFileChannel(channelID string, file slack.File) bool {
	for _, channel := range file.Channels {
		if channel == channelID {
			return true
		}
	}
	return false
}

func CopyGoogleDriveFile() {}

func GetNextMonday() time.Time {
	currentTime := time.Now()
	day := currentTime.Weekday()
	daysToAdd := 8 - int(day)
	currentTime = currentTime.AddDate(0, 0, daysToAdd)
	dueTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 16, 0, 0, 0, currentTime.Location())
	return dueTime
}

func CreateAssignment(assignment Assignment, db *gorm.DB) {
	if db.NewRecord(&assignment) {
		db.Create(&assignment)
	} else {
		fmt.Println("error, primary key already exists for assignment")
		//format.errorf
	}
}

func ParseAssignment(userInput string) string {
	names := strings.Split(userInput, " ")
	if len(names) == 3 &&
		strings.ToLower(names[0]) == "assignment" &&
		strings.ToLower(names[1]) == "due" {
		return names[2]
	} else {
		return ""
	}
}

func ParseDate(dateString string) string {
	return "date"
}

// func DownloadPDF()

// func getStudentChannels()

// func Assign(file, channels, due date)
