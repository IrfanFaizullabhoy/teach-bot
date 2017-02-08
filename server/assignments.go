package main

import (
	//"github.com/nlopes/slack"
	"strings"
)

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
