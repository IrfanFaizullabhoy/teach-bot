package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/nlopes/slack"
)

func SetupSubmissions(assignment Assignment) {
	students := GetStudents()
	for _, student := range students {
		SetupSubmission(assignment, student) // CHANGE TO BATCH CREATE
	}
}

func SetupSubmission(assignment Assignment, user User) {
	submission := Submission{AssigmentID: assignment.ID, UserID: user.ID, ChannelID: user.ChannelID, Submitted: false}
	db.Create(&submission)
}

func CollectSubmissions(assignment Assignment) {
	var submissions []Submission
	db.Model(&assignment).Related(&submissions, "Submissions")
	assignment.Submissions = submissions
}

func SubmissionReport(assignment Assignment) {
	//	students := GetStudents()
	var studentsWithoutSubmission []string
	for _, submission := range assignment.Submissions {
		if submission.Submitted == false {
			studentsWithoutSubmission = append(studentsWithoutSubmission, submission.UserID)
		}
	}

}

func DownloadSubmission(fileSharedEvent Event) {
	api := GetSlackClient()
	file, _, _, err := api.GetFileInfo(fileSharedEvent.File.ID, 1, 1) //returns file with one comment/onepage
	check(err)

	//FIND IT IN THE DB
	user := GetUser(file.User)
	submission := FindSubmission(file.User, user.ChannelID, file.Name)

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

	submission.FilePath = filePath
	submission.Submitted = true
	submission.FileName = file.Name
	db.Save(&submission)

	//ConfirmAndSendAssignment(assignment)
}

func FindSubmission(userID, channelID, fileName string) Submission {
	api := GetSlackClient()
	var submissions []Submission
	params := slack.PostMessageParameters{}
	db.Where("user_id = ? AND channel_id = ? AND submitted = ?", userID, channelID, false).Find(&submissions)
	if len(submissions) == 1 {
		user := GetUser(userID)
		api.PostMessage(user.ChannelID, "Downloading the file: `"+fileName+"`", params)
		return submissions[0]
	} else if len(submissions) == 0 {
		return submissions[0]
	} else {
		fmt.Println(userID)
		fmt.Println("more than one submission not updated... " + strconv.Itoa(len(submissions)))
		//return WhichAssignment(userID, submissions)
		return submissions[0]
	}
}
