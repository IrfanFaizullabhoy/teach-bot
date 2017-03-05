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

func GrabAssignMessage(channelID, userID string) {

}

func HandleFileShared(fileSharedEvent Event) {
	if isTeacher(fileSharedEvent.User) {
		go DownloadAssignment(fileSharedEvent)
	} else {
		go SubmitAssignment(fileSharedEvent)
	}
}

// ELSE IF ASSIGNMENT = TO COLLECT
func SubmitAssignment(fileSharedEvent Event) {
}

// ELSE RANDOM ASSIGNMENT [make sure its not to assign or collect]
func DownloadAssignment(fileSharedEvent Event) {
	api := GetSlackClient()
	file, _, _, err := api.GetFileInfo(fileSharedEvent.File.ID, 1, 1) //returns file with one comment/onepage
	check(err)

	//FIND IT IN THE DB
	assignment := FindAssignment(file.User, file.Name)

	// GOOGLE DOC
	if strings.Contains(file.URLPrivate, "google.com") {
		fmt.Println("Google Doc")
		return
	}

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

	assignment.FilePath = filePath
	assignment.Downloaded = true
	assignment.FileName = file.Name
	db.Save(&assignment)

	ConfirmAndSendAssignment(assignment)
}

func ConfirmAndSendAssignment(assignment Assignment) {
	api := GetSlackClient()
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
	}
	params := slack.FileUploadParameters{File: assignment.FilePath, Filename: assignment.FileName, Channels: channels}
	_, err := api.UploadFile(params)
	check(err)
}

func FindAssignment(userID, fileName string) Assignment {
	api := GetSlackClient()
	var assignments []Assignment
	params := slack.PostMessageParameters{}
	db.Where("user_id = ? AND downloaded = ?", userID, false).Find(&assignments)
	if len(assignments) == 1 {
		user := GetUser(userID)
		api.PostMessage(user.ChannelID, "Downloading the file: `"+fileName+"`", params)
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

func DateInteractive(userID, channel string) {
	api := GetSlackClient()
	params := slack.NewPostMessageParameters()
	attachment := slack.Attachment{CallbackID: "assignment_due", Fallback: "service not working properly"}
	attachmentMondayAction := slack.AttachmentAction{Name: "monday", Text: "Next Monday at 4pm", Type: "button"}
	attachmentOtherAction := slack.AttachmentAction{Name: "other", Text: "Other", Type: "button"}
	attachment.Actions = append(attachment.Actions, attachmentMondayAction)
	attachment.Actions = append(attachment.Actions, attachmentOtherAction)
	params.Attachments = append(params.Attachments, attachment)
	_, _, err := api.PostMessage(channel, "When would you like the assignment to be due?", params)
	check(err)
}

func HandleDate(attachmentDateAction slack.AttachmentActionCallback) {
	api := GetSlackClient()
	if len(attachmentDateAction.Actions) == 1 {
		if attachmentDateAction.Actions[0].Name == "monday" {
			monday := GetNextMonday().Format(time.ANSIC)
			params := slack.PostMessageParameters{}
			assignment := Assignment{DueDate: monday, UserID: attachmentDateAction.User.ID, Downloaded: false}
			go CreateAssignment(assignment, db)
			api.PostMessage(attachmentDateAction.Channel.ID, "Great, it'll be due on *"+monday+"*\nPlease press the `+` sign below/to my left to attach the homework?", params)
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
