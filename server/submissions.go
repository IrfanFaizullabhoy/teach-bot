package main

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	//"strconv"
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
	SubmissionReport(assignment)
}

func SubmissionReport(assignment Assignment) {
	//	students := GetStudents()
	zipFilename := assignment.UserID + ".zip"
	newfile, err := os.Create(zipFilename)
	check(err)
	zipit := zip.NewWriter(newfile)

	for _, submission := range assignment.Submissions {
		if submission.Submitted {
			submissionFile, err := os.Open(submission.FilePath)
			check(err)
			defer submissionFile.Close()

			info, err := submissionFile.Stat()
			check(err)

			header, err := zip.FileInfoHeader(info)
			check(err)

			writer, err := zipit.CreateHeader(header)
			check(err)

			_, err = io.Copy(writer, submissionFile)
			check(err)
		}
	}

	err = zipit.Close()
	check(err)
	err = newfile.Close()
	check(err)

	info, err := newfile.Stat()
	check(err)

	path, err := filepath.Abs(filepath.Dir(info.Name()))
	path = path + newfile.Name()

	channels := GetInstructorChannels()
	fileParams := slack.FileUploadParameters{File: path, Filetype: ".zip", Filename: info.Name(), Channels: channels}
	team := GetTeam(assignment.TeamID)
	botConn := slack.New(team.BotToken)
	botConn.UploadFile(fileParams)
}

func DownloadSubmission(fileSharedEvent Event, teamID string) {
	team := GetTeam(teamID)
	botConn := slack.New(team.BotToken)
	file, _, _, err := botConn.GetFileInfo(fileSharedEvent.File.ID, 1, 1) //returns file with one comment/onepage
	check(err)

	//FIND IT IN THE DB
	user := GetUser(file.User)
	submission := FindSubmission(file.User, user.ChannelID, file.Name, team)

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

func FindSubmission(userID, channelID, fileName string, team Team) Submission {
	var submissions []Submission
	params := slack.PostMessageParameters{}
	botConn := slack.New(team.BotToken)
	db.Where("user_id = ? AND channel_id = ? AND team_id = ? AND submitted = ?", userID, channelID, team.TeamID, false).Find(&submissions)
	if len(submissions) == 1 {
		user := GetUser(userID)
		botConn.PostMessage(user.ChannelID, "Downloading the file: `"+fileName+"`", params)
		return submissions[0]
	} else if len(submissions) == 0 {
		return submissions[0]
	} else {
		//fmt.Println(userID)
		//fmt.Println("more than one submission not updated... " + strconv.Itoa(len(submissions)))
		//return WhichAssignment(userID, submissions)
		return submissions[0]
	}
}
