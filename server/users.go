package main

import (
	"errors"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/nlopes/slack"
	//"os"
	"fmt"
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

func CreateUser(user User, db *gorm.DB) {
	if db.NewRecord(&user) {
		db.Create(&user)
	} else {
		fmt.Println("error, primary key already exists for user")
		//format.errorf
	}
}

func FillUserInfo(user slack.User, role string, db *gorm.DB) {
	var dbUser User
	db.Where("user_id = ?", user.ID).First(&dbUser)
	dbUser.Profile = user.Profile
	dbUser.Name = user.Name
	dbUser.Role = role
	dbUser.IsBot = user.IsBot
	db.Save(&dbUser)

	//CONFIRM ENROLLMENT
	api := GetSlackClient()
	params := slack.NewPostMessageParameters()
	_, _, err := api.PostMessage(dbUser.ChannelID, "Awesome -- "+user.Name+" you're all registered. I'll be contacting you in the future for when assignments are posted and collected!", params)
	check(err)

}

func GetUser(userID string) User {
	var dbUser User
	db.Where("user_id = ?", userID).First(&dbUser)
	return dbUser
}

func WelcomeToTeam(TeamJoinEvent Event) {
	userID := TeamJoinEvent.User
	fmt.Println("userID")
	api := GetSlackClient()
	params := slack.NewPostMessageParameters()
	attachment := slack.Attachment{CallbackID: "student_or_teacher", Fallback: "service not working properly"}
	attachmentStudentAction := slack.AttachmentAction{Name: "student", Text: "Student", Type: "button"}
	attachmentTeacherAction := slack.AttachmentAction{Name: "teacher", Text: "Teacher", Type: "button"}
	attachment.Actions = append(attachment.Actions, attachmentStudentAction)
	attachment.Actions = append(attachment.Actions, attachmentTeacherAction)
	_, _, channel, err := api.OpenIMChannel(userID)
	check(err)
	user := User{ID: userID, ChannelID: channel}
	CreateUser(user, db)
	_, _, err = api.PostMessage(channel, "Welcome! Are you a student or a teacher?", params)
	check(err)
}

func WelcomeToTeamTest(userID string) {
	api := GetSlackClient()
	params := slack.NewPostMessageParameters()
	attachment := slack.Attachment{CallbackID: "student_or_teacher", Fallback: "service not working properly"}
	attachmentStudentAction := slack.AttachmentAction{Name: "student", Text: "Student", Type: "button"}
	attachmentTeacherAction := slack.AttachmentAction{Name: "teacher", Text: "Teacher", Type: "button"}
	attachment.Actions = append(attachment.Actions, attachmentStudentAction)
	attachment.Actions = append(attachment.Actions, attachmentTeacherAction)
	params.Attachments = append(params.Attachments, attachment)
	_, _, channel, err := api.OpenIMChannel(userID)
	check(err)
	user := User{ID: userID, ChannelID: channel}
	CreateUser(user, db)
	_, _, err = api.PostMessage(channel, "Welcome! Are you a student or a teacher?", params)
	check(err)
}

func AddToDatabase(StudentOrTeacherAction slack.AttachmentActionCallback) {
	user := StudentOrTeacherAction.User
	if len(StudentOrTeacherAction.Actions) == 1 {
		if StudentOrTeacherAction.Actions[0].Name == "student" {
			FillUserInfo(user, "student", db)
		} else {
			FillUserInfo(user, "teacher", db)
		}
	}
}

func GetInstructors() []User {
	var teachers []User
	db.Where("role = ?", "teacher").Find(&teachers)
	//instructors := []string{"U3YKBAK1S", "U42EVJF7E", "U3YK6EPV0"}
	return teachers
}

func GetStudents() []User {
	var students []User
	db.Where("role = ?", "teacher").Find(&students)
	return students
}

func isTeacher(userID string) bool {
	var dbUser User
	db.Where("user_id = ?", userID).First(&dbUser)
	if dbUser.Role == "teacher" {
		return true
	} else {
		return false
	}
}
