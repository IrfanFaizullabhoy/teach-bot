package main

import (
	"errors"
	//"github.com/nlopes/slack"
	"strings"
)

func ParseStudentName(userInput string) (string, string, error) {
	names := strings.Split(userInput, " ")
	if len(names) == 3 &&
		strings.ToLower(names[0]) == "register" {
		return names[1], names[2], nil
	} else {
		return "", "", errors.New("Did not input 3 items")
	}
}

func ParseTeacherName(userInput string) (string, string, error) {
	names := strings.Split(userInput, " ")
	if len(names) == 4 &&
		strings.ToLower(names[0]) == "register" &&
		strings.ToLower(names[1]) == "teacher" {
		return names[2], names[3], nil
	} else {
		return "", "", errors.New("Did not input 4 items")
	}
}
