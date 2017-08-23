package main

import (
	_ "github.com/lib/pq"
	"github.com/nlopes/slack"
)

func StartInstructorConversation(userID, name, teamID string) {
	team := GetTeam(teamID)
	appConn := slack.New(team.Token)
	groupID := FindGroupByName("instructors-and-"+name, appConn)
	if groupID == "" {
		group, err := appConn.CreateGroup("instructors-and-" + name)
		instructors := GetInstructors(teamID)
		for _, instructor := range instructors {
			_, _, err = appConn.InviteUserToGroup(group.ID, instructor.ID)
			check(err)
		}
		groupID = group.ID
		_, _, err = appConn.InviteUserToGroup(groupID, userID)
		check(err)
	}

	appConn.OpenGroup(groupID)
}
