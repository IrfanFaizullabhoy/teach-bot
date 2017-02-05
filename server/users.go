package main

import (
	"errors"
	//"github.com/nlopes/slack"
	"strings"
)

func ParseName(userInput string) (string, string, error) {
	names := strings.SplitAfter(userInput, " ")
	if len(names) == 2 {
		return names[0], names[1], nil
	} else {
		return "", "", errors.New("Did not input 2 items")
	}

}
