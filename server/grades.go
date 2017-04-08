package main

import (
	//"fmt"
	"github.com/nlopes/slack"
	"strconv"
	"strings"
)

func GetStudentGrades(studentID, teamID, text string) {
	var submissions []Submission
	team := GetTeam(teamID)
	botConn := slack.New(team.BotToken)
	user := GetUser(studentID)

	params := slack.PostMessageParameters{Markdown: true}
	if user.Role == "teacher" {
		if strings.Contains(text, "@") {

		}
		botConn.PostMessage(user.ChannelID, "Oops, looks like you're a *teacher*... let me show you the grades for a sample student (Bryanna)", params)
		studentID = "U3XTGS09X"

	}
	db.Where("user_id = ? AND team_id = ? AND graded = ?", studentID, teamID, true).Find(&submissions)
	for _, submission := range submissions {
		botConn.PostMessage(user.ChannelID, "For *"+submission.AssignmentName+"* you got a *"+strconv.FormatFloat(submission.Score, 'f', 2, 64)+" / "+strconv.FormatFloat(submission.MaxPoints, 'f', 2, 64)+"*", params)
	}

}

func EnterManualGrades(manualGrading ManualGrade) {
	teamID := "T3Z7YKN07"
	GetTeam(teamID)

	if len(manualGrading.AssignmentGrades) == len(manualGrading.UserIDs) {
		for i, grade := range manualGrading.AssignmentGrades {
			submission := Submission{Graded: true,
				AssignmentName: manualGrading.Name,
				MaxPoints:      manualGrading.MaxPoints,
				Score:          grade,
				UserID:         manualGrading.UserIDs[i],
				TeamID:         teamID}
			db.Create(&submission)
		}
	}
}

func UpdateManualGrades(manualGrading ManualGrade) {
	teamID := "T3Z7YKN07"
	GetTeam(teamID)
	var submissions []Submission
	if len(manualGrading.AssignmentGrades) == len(manualGrading.UserIDs) {
		for i, grade := range manualGrading.AssignmentGrades {
			db.Where("assignment_name = ? AND user_id = ? AND team_id = ? AND graded = ?", manualGrading.Name, manualGrading.UserIDs[i], teamID, true).Find(&submissions)
			if len(submissions) == 1 {
				submissions[0].Score = grade
			}
			db.Save(&submission)
		}
	}
}

/* JSON


 */
