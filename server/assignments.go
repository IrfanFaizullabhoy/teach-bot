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

func DownloadFile(fileSharedEvent Event) {
	//GET FILE INFO
	api := GetSlackClient()
	file, _, _, err := api.GetFileInfo(fileSharedEvent.File.ID, 1, 1) //returns file with one comment/onepage
	check(err)

	//FIND IT IN THE DB
	assignment := FindFile(file.User, file.Name)

	// IF ASSIGNMENT = TO ASSIGN

	/*if isAssignment(fileSharedEvent) {
		fmt.Println("is an assignment")
		AssignFile(fileSharedEvent)
	} else if isSubmission() {
		// ATTACHMENT IS A SUBMISSION

		fmt.Println("is a submission")

	}*/

	// ELSE IF ASSIGNMENT = TO COLLECT

	// ELSE RANDOM ASSIGNMENT [make sure its not to assign or collect]

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
}

func IsSubmission(fileSharedChannel string) bool {
	return true
}

func IsAssignment(fileSharedEvent Event) bool {
	/*api := GetSlackClient()
	if isTeacher(fileSharedEvent.User) {
		user := GetUser(fileSharedEvent.User)
		//api.PostMessage(user.ChannelID, "Is t, params)
	}*/
	return true

}

func FindFile(userID, fileName string) Assignment {
	api := GetSlackClient()
	var assignments []Assignment
	params := slack.PostMessageParameters{}
	db.Where("user_id = ? AND downloaded = ?", userID, false).Find(&assignments)
	if len(assignments) == 1 {
		user := GetUser(userID)
		api.PostMessage(user.ChannelID, "Downloading the file: `"+fileName+"`", params)
		return assignments[0]
	} else {
		fmt.Println(userID)
		fmt.Println("more than one assigment not downloaded... " + strconv.Itoa(len(assignments)))
		return assignments[0]
	}
}

func DateInteractive(userID string) {
	api := GetSlackClient()
	params := slack.NewPostMessageParameters()
	attachment := slack.Attachment{CallbackID: "assignment_due", Fallback: "service not working properly"}
	attachmentMondayAction := slack.AttachmentAction{Name: "monday", Text: "Next Monday at 4pm", Type: "button"}
	attachmentOtherAction := slack.AttachmentAction{Name: "other", Text: "Other", Type: "button"}
	attachment.Actions = append(attachment.Actions, attachmentMondayAction)
	attachment.Actions = append(attachment.Actions, attachmentOtherAction)
	params.Attachments = append(params.Attachments, attachment)
	_, _, channel, err := api.OpenIMChannel(userID)
	check(err)
	_, _, err = api.PostMessage(channel, "When would you like the assignment to be due?", params)
	check(err)
}

func HandleDate(attachmentDateAction slack.AttachmentActionCallback) {
	api := GetSlackClient()
	if len(attachmentDateAction.Actions) == 1 {
		if attachmentDateAction.Actions[0].Name == "monday" {
			monday := GetNextMonday().Format(time.ANSIC)
			params := slack.PostMessageParameters{}
			api.PostMessage(attachmentDateAction.Channel.ID, "Great, it'll be due on `"+monday+"`", params)
			fmt.Println(monday)
			assignment := Assignment{DueDate: monday, UserID: attachmentDateAction.User.ID, Downloaded: false}
			CreateAssignment(assignment, db)
			api.PostMessage(attachmentDateAction.Channel.ID, "Can you please post the Google Drive link or File for the homework?", params)
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
