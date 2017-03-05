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
	userInfo, err := api.GetUserInfo(userID)
	check(err)
	params := slack.NewPostMessageParameters()
	attachment := slack.Attachment{CallbackID: "student_or_teacher", Fallback: "service not working properly"}
	attachmentStudentAction := slack.AttachmentAction{Name: "student", Text: "Student", Type: "button"}
	attachmentTeacherAction := slack.AttachmentAction{Name: "teacher", Text: "Teacher", Type: "button"}
	attachment.Actions = append(attachment.Actions, attachmentStudentAction)
	attachment.Actions = append(attachment.Actions, attachmentTeacherAction)
	groupID := FindGroupByName("teachbot-and-" + userInfo.Name)
	if groupID == "" {
		group, err := api.CreateGroup("teachbot-and-" + userInfo.Name)
		check(err)
		groupID = group.ID
	}
	api.InviteUserToGroup(groupID, userID)
	teachBotID := GetTeachBotID()
	api.InviteUserToGroup(groupID, teachBotID)
	user := User{ID: userID, ChannelID: groupID}
	CreateUser(user, db)
	_, _, err = api.PostMessage(user.ChannelID, "Welcome! Are you a student or a teacher?", params)
	check(err)
}

func GetTeachBotID() string {
	api := GetSlackClient()
	users, err := api.GetUsers()
	check(err)
	for _, user := range users {
		if user.IsBot && user.Name == "teachbot2" {
			return user.ID
		}
	}
	return ""
}

func WelcomeToTeamTest(userID string) {
	api := GetSlackClient()
	userInfo, err := api.GetUserInfo(userID)
	check(err)
	params := slack.NewPostMessageParameters()
	attachment := slack.Attachment{CallbackID: "student_or_teacher", Fallback: "service not working properly"}
	attachmentStudentAction := slack.AttachmentAction{Name: "student", Text: "Student", Type: "button"}
	attachmentTeacherAction := slack.AttachmentAction{Name: "teacher", Text: "Teacher", Type: "button"}
	attachment.Actions = append(attachment.Actions, attachmentStudentAction)
	attachment.Actions = append(attachment.Actions, attachmentTeacherAction)
	params.Attachments = append(params.Attachments, attachment)
	groupID := FindGroupByName("teachbot-and-" + userInfo.Name)
	if groupID == "" {
		group, err := api.CreateGroup("teachbot-and-" + userInfo.Name)
		check(err)
		groupID = group.ID
	}
	api.InviteUserToGroup(groupID, userID)
	teachBotID := GetTeachBotID()
	api.InviteUserToGroup(groupID, teachBotID)
	user := User{ID: userID, ChannelID: groupID}
	CreateUser(user, db)
	_, _, err = api.PostMessage(user.ChannelID, "Welcome! Are you a student or a teacher?", params)
	check(err)
}

func FindGroupByName(groupName string) string {
	api := GetSlackClient()
	groups, err := api.GetGroups(false)
	check(err)
	for _, group := range groups {
		if group.Name == groupName {
			return group.ID
		}
	}
	return ""
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
