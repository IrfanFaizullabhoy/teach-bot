package main

import (
	"fmt"
	"github.com/nlopes/slack"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func GrabAssignMessage(channelID, userID string) {

}

func DownloadFile(fileSharedEvent Event) (*slack.File, string) {
	//GET FILE INFO
	api := GetSlackClient()
	fmt.Println("Printing the File Name" + fileSharedEvent.File.ID)
	file, _, _, err := api.GetFileInfo(fileSharedEvent.File.ID, 1, 1) //returns file with one comment/onepage
	check(err)
	fmt.Println(file.URLPrivate)
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

func PostAssignmentInStudents(file *slack.File, filePath string) {
	api := GetSlackClient()
	//GET STUDENT CHANNELS
	channels := GetStudentChannels(db)

	//UPLOAD FILE
	fileParams := slack.FileUploadParameters{Filename: file.Name, File: filePath, Filetype: file.Filetype, Channels: channels}
	api.UploadFile(fileParams)
	channels = GetStudentChannels(db)
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
