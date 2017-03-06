package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/nlopes/slack"
)

func AcknowledgeCallback(attachmentAcknowledgeAction slack.AttachmentActionCallback) {
	api := GetSlackClient()
	userID := attachmentAcknowledgeAction.User.ID
	value := attachmentAcknowledgeAction.Actions[0].Value
	ts := attachmentAcknowledgeAction.OriginalMessage.Timestamp
	var acknowledges []AcknowledgeMessage
	db.Where("user_id = ? AND timestamp = ?", userID, ts).Find(&acknowledges)
	if len(acknowledges) == 1 {
		acknowledge := acknowledges[0]
		var acknowledgeactions []AcknowledgeAction
		db.Model(&acknowledge).Related(&acknowledgeactions, "AcknowledgeActions")
		acknowledge.AcknowledgeActions = acknowledgeactions
		for _, action := range acknowledge.AcknowledgeActions {
			fmt.Println(action.ID)
		}
		acknowledgeAction := AcknowledgeAction{UserID: userID, Value: value, AckID: acknowledges[0].ID}
		AddAcknowledgement(acknowledge, acknowledgeAction)
		if FullyAcknowledged(acknowledge, attachmentAcknowledgeAction.Channel) &&
			!strings.Contains(attachmentAcknowledgeAction.OriginalMessage.Text, "[FULLY_ACKNOWLEDGED]") {
			api.UpdateMessage(attachmentAcknowledgeAction.Channel.ID, attachmentAcknowledgeAction.OriginalMessage.Timestamp, " [FULLY_ACKNOWLEDGED] "+attachmentAcknowledgeAction.OriginalMessage.Text)
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
	api := GetSlackClient()
	channelPtr, err := api.GetGroupInfo(channel.ID)
	check(err)
	fmt.Println(strconv.Itoa(len(channelPtr.Members)))
	if len(acknowledge.AcknowledgeActions) >= len(channelPtr.Members)-1 {
		return true
	} else {
		return false
	}
}

func AddAcknowledgement(acknwoledge AcknowledgeMessage, acknowledgeAction AcknowledgeAction) {
	fmt.Println("adding acknowledgement")
	if !ContainsAcknowledge(acknwoledge.AcknowledgeActions, acknowledgeAction) {
		db.Create(&acknowledgeAction)
		acknwoledge.AcknowledgeActions = append(acknwoledge.AcknowledgeActions, acknowledgeAction)
		db.Save(&acknwoledge)
	}

}

func ContainsAcknowledge(acknowledgeActions []AcknowledgeAction, acknowledgeAction AcknowledgeAction) bool {
	for _, AAction := range acknowledgeActions {
		if acknowledgeAction.UserID == AAction.UserID &&
			acknowledgeAction.Value == AAction.Value {
			return true
		}
	}
	//db.Model(&AcknowledgeAction).
	return false
}
