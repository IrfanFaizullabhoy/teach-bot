package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/nlopes/slack"
	"fmt"
)

var StudentMap map[string]User
var TeamMap map[string]Team

func RegisterAll(team Team) {
	userIDs := FindAllUserIDs(team)
	for _, id := range userIDs {
		WelcomeToTeamByID(id, team)
	}
}

func CreateUser(user User, db *gorm.DB) {
	if db.NewRecord(&user) {
		db.Create(&user)
	} else {
		fmt.Println("error, primary key already exists for user")
	}
}

func FillUserInfo(user slack.User, role string, db *gorm.DB, botConn *slack.Client) {
	var dbUser User
	db.Where("user_id = ?", user.ID).First(&dbUser)
	dbUser.Profile = user.Profile
	dbUser.Name = user.Name
	dbUser.Role = role
	dbUser.IsBot = user.IsBot
	db.Save(&dbUser)

	//CONFIRM ENROLLMENT
	params := slack.NewPostMessageParameters()
	_, _, err := botConn.PostMessage(dbUser.ChannelID, "Awesome -- "+user.Name+" you're all registered. I'll be contacting you in the future for when assignments are posted and collected!", params)
	check(err)

}

func GetUser(userID string) User {
	dbUser, ok := StudentMap[userID]
	if !ok {
		db.Where("user_id = ?", userID).First(&dbUser)
		StudentMap[userID] = dbUser
	}
	return dbUser
}

func WelcomeToTeam(TeamJoinEvent Event, TeamID string) {
	userID := TeamJoinEvent.User
	if len(userID) != 9 {
		fmt.Println("user ID is not length 9") //ERROR
		return
	}
	team := GetTeam(TeamID)
	appConn := slack.New(team.Token)
	botConn := slack.New(team.BotToken)
	userInfo, err := appConn.GetUserInfo(userID)
	if err != nil {
		fmt.Println("can't find user") //ERROR
		return
	}

	//FORMING MESSAGE
	params := slack.NewPostMessageParameters()
	attachment := slack.Attachment{CallbackID: "student_or_teacher", Fallback: "service not working properly"}
	attachmentStudentAction := slack.AttachmentAction{Name: "student", Text: "Student", Type: "button"}
	attachmentTeacherAction := slack.AttachmentAction{Name: "teacher", Text: "Teacher", Type: "button"}
	attachment.Actions = append(attachment.Actions, attachmentStudentAction)
	attachment.Actions = append(attachment.Actions, attachmentTeacherAction)
	params.Attachments = append(params.Attachments, attachment)
	groupID := FindGroupByName(userInfo.Name+"-and-teachbot", appConn)
	if groupID == "" {
		group, err := appConn.CreateGroup(userInfo.Name + "-and-teachbot")
		check(err)
		groupID = group.ID
		_, _, err = appConn.InviteUserToGroup(groupID, userID)
		check(err)
		_, _, err = appConn.InviteUserToGroup(groupID, team.BotID) //Invite Teachbot
		check(err)
		err = appConn.KickUserFromGroup(groupID, team.InstallerID)
		check(err)
		user := User{ID: userID, ChannelID: groupID, TeamID: team.TeamID}
		CreateUser(user, db)
		_, _, err = botConn.PostMessage(user.ChannelID, "Welcome! Are you a student or a teacher?", params)
		check(err)
	} else {
		//
	}
}

func FindAllUserIDs(team Team) []string {
	appConn := slack.New(team.Token)
	users, err := appConn.GetUsers()
	check(err)
	var allUserIDs []string
	for _, user := range users {
		if !user.IsBot {
			allUserIDs = append(allUserIDs, user.ID)
		}
	}
	return allUserIDs
}

func WelcomeToTeamByID(userID string, team Team) {
	if userID == team.BotID {
		return
	}

	if userID == "USLACKBOT" {
		return
	}
	if len(userID) != 9 {
		fmt.Println("user ID is not length 9") //ERROR
		return
	}
	appConn := slack.New(team.Token)
	info, err := appConn.GetUserInfo(userID)
	if info.IsBot || info.IsRestricted || info.IsUltraRestricted {
		return
	}

	botConn := slack.New(team.BotToken)
	userInfo, err := appConn.GetUserInfo(userID)
	if err != nil {
		fmt.Println("can't find user") //ERROR
		return
	}

	//FORMING MESSAGE
	params := slack.NewPostMessageParameters()
	attachment := slack.Attachment{CallbackID: "student_or_teacher", Fallback: "service not working properly"}
	attachmentStudentAction := slack.AttachmentAction{Name: "student", Text: "Student", Type: "button"}
	attachmentTeacherAction := slack.AttachmentAction{Name: "teacher", Text: "Teacher", Type: "button"}
	attachment.Actions = append(attachment.Actions, attachmentStudentAction)
	attachment.Actions = append(attachment.Actions, attachmentTeacherAction)
	params.Attachments = append(params.Attachments, attachment)
	groupID := FindGroupByName(userInfo.Name+"+bot", appConn)
	if groupID == "" {
		group, err := appConn.CreateGroup(userInfo.Name + "+bot")
		check(err)
		groupID = group.ID
		if userID != team.InstallerID {
			fmt.Println(userID)
			_, _, err = appConn.InviteUserToGroup(groupID, userID)
			check(err)
		}
		_, _, err = appConn.InviteUserToGroup(groupID, team.BotID)
		check(err)
		if userID != team.InstallerID {
			err = appConn.LeaveGroup(groupID)
			check(err)
		}
		user := User{ID: userID, ChannelID: groupID, TeamID: team.TeamID}
		CreateUser(user, db)
		_, _, err = botConn.PostMessage(user.ChannelID, "Welcome! Are you a student or a teacher?", params)
		check(err)
	} else {
		//
	}
}

func FindGroupByName(groupName string, appConn *slack.Client) string {
	groups, err := appConn.GetGroups(false)
	check(err)
	for _, group := range groups {
		if group.Name == groupName {
			return group.ID
		}
	}
	return ""
}

func AddToDatabase(StudentOrTeacherAction slack.AttachmentActionCallback) {
	team := GetTeam(StudentOrTeacherAction.Team.ID)
	botConn := slack.New(team.BotToken)
	user := StudentOrTeacherAction.User
	if len(StudentOrTeacherAction.Actions) == 1 {
		if StudentOrTeacherAction.Actions[0].Name == "student" {
			FillUserInfo(user, "student", db, botConn)
		} else {
			installer := GetUser(team.InstallerID)
			params := slack.PostMessageParameters{}
			botConn.PostMessage(installer.ChannelID, "user: "+StudentOrTeacherAction.User.Name+"registered as a teacher", params)
			FillUserInfo(user, "teacher", db, botConn)
		}
	}
}

func GetInstructors(teamID string) []User {
	var teachers []User
	db.Where("role = ? AND team_id = ?", "teacher", teamID).Find(&teachers)
	return teachers
}

func GetStudents() []User {
	var students []User
	db.Where("role = ?", "student").Find(&students)
	return students
}

func GetUsers(teamID string) []User {
	var users []User
	db.Where("team_id = ?", teamID).Find(&users)
	return users
}

func InitializeStudentMap() {
	StudentMap = make(map[string]User)
	students := GetStudents()
	for _, student := range students {
		StudentMap[student.ID] = student
	}
}

func GetTeams() []Team {
	var teams []Team
	db.Find(&teams)
	return teams
}

func InitializeTeamMap() {
	TeamMap = make(map[string]Team)
	teams := GetTeams()
	for _, team := range teams {
		TeamMap[team.TeamID] = team
	}
}

func GetTeam(id string) Team {
	dbTeam, ok := TeamMap[id]
	if !ok {
		db.Where("team_id = ?", id).First(&dbTeam)
		TeamMap[id] = dbTeam
	}
	return dbTeam
}
