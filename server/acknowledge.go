package main

import (
	"fmt"
	"strconv"

	"github.com/nlopes/slack"
)

func AcknowledgeCallback(attachmentAcknowledgeAction slack.AttachmentActionCallback) {
	api := GetSlackClient()
	userID := attachmentAcknowledgeAction.User.ID
	ts := attachmentAcknowledgeAction.OriginalMessage.Timestamp
	var acknowledges []AcknowledgeMessage
	db.Where("user_id = ? AND timestamp = ? AND ", userID, ts).Find(&acknowledges)
	if len(acknowledges) == 1 {
		AddAcknowledgement(userID, acknowledges[0])
		if FullyAcknowledged(acknowledges[0], attachmentAcknowledgeAction.Channel) {
			api.UpdateMessage(attachmentAcknowledgeAction.Channel.ID, attachmentAcknowledgeAction.OriginalMessage.Timestamp, "*[FULLY_ACKNOWLEDGED]*"+attachmentAcknowledgeAction.OriginalMessage.Text)
		}
	} else {
		fmt.Println(userID)
		fmt.Println("more than one acknowledgement... " + strconv.Itoa(len(acknowledges)))
	}
	//attachmentAcknowledgeAction
}

func FullyAcknowledged(acknowledge AcknowledgeMessage, channel slack.Channel) bool {
	/*activeMembers := 0
	for _, member := range channel.Members {
		if IsActiveUser(member) {
			activeMembers++
		}
	}*/
	if len(acknowledge.UsersAcknowledged) == len(channel.Members)-1 {
		return true
	} else {
		return false
	}
}

func AddAcknowledgement(userID string, acknwoledge AcknowledgeMessage) {
	if !contains(acknwoledge.UsersAcknowledged, userID) {
		acknwoledge.UsersAcknowledged = append(acknwoledge.UsersAcknowledged, userID)
		db.Save(&acknwoledge)
	}
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
