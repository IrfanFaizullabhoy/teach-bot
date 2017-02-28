package main

import (
	"fmt"
	"github.com/nlopes/slack"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

func GrabAssignMessage(channelID, userID string) {

}

func DownloadFile(fileSharedEvent Event) (*slack.File, string) {
	//GET FILE INFO
	api := GetSlackClient()
	file, _, _, err := api.GetFileInfo(fileSharedEvent.File.ID, 1, 1) //returns file with one comment/onepage
	check(err)

	//DOWNLOAD FILE
	client := &http.Client{}
	req, _ := http.NewRequest("GET", file.URLPrivate, nil)
	req.Header.Add("Authorization", "Bearer "+os.Getenv("SLACK_TOKEN"))
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
		panic("file size unequal in DownloadFile")
	}
	return file, filePath
}

func DateInteractive(userID string) {
	api := GetSlackClient()
	params := slack.NewPostMessageParameters()
	attachment := slack.Attachment{CallbackID: "assignment_due", Fallback: "service not working properly"}
	attachmentMondayAction := slack.AttachmentAction{Name: "monday", Text: "Next Monday at 4pm", Type: "button"}
	attachmentOtherAction := slack.AttachmentAction{Name: "other", Text: "Other", Type: "button"}
	attachment.Actions = append(attachment.Actions, attachmentMondayAction)
	attachment.Actions = append(attachment.Actions, attachmentOtherAction)
	_, _, channel, err := api.OpenIMChannel(userID)
	check(err)
	_, _, err = api.PostMessage(channel, "When would you like the assignment to be due?", params)
	check(err)
}

func HandleDate(attachmentDateAction slack.AttachmentActionCallback) {
	if len(attachmentDateAction.Actions) == 1 {
		if attachmentDateAction.Actions[0].Name == "monday" {
			monday := GetNextMonday().Format(time.ANSIC)
			fmt.Println(monday)
			assignment := Assignment{DueDate: monday}
			assignment = assignment
		} else {

		}
	}
}

func GetNextMonday() time.Time {
	currentTime := time.Now()
	day := currentTime.Weekday()
	daysToAdd := 8 - int(day)
	currentTime = currentTime.AddDate(0, 0, daysToAdd)
	dueTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 16, 0, 0, 0, currentTime.Location())
	return dueTime
}

func PostAssignmentInStudents(file *slack.File, filePath string) {
	api := GetSlackClient()
	//GET STUDENT CHANNELS
	//channels := GetStudentChannels(db)

	//UPLOAD FILE
	var channels []string
	fileParams := slack.FileUploadParameters{Filename: file.Name, File: filePath, Filetype: file.Filetype, Channels: channels}
	api.UploadFile(fileParams)
	//channels = GetStudentChannels(db)
	fileParams.Channels = channels
	api.UploadFile(fileParams)

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
