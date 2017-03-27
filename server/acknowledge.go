package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/nlopes/slack"
)

func AcknowledgeCallback(attachmentAcknowledgeAction slack.AttachmentActionCallback) {
	teamID := attachmentAcknowledgeAction.Team.ID
	team := GetTeam(teamID)
	botConn := slack.New(team.BotToken)
	//fmt.Println(userID)
	value := attachmentAcknowledgeAction.Actions[0].Value
	ts := attachmentAcknowledgeAction.OriginalMessage.Timestamp
	var acknowledges []AcknowledgeMessage
	db.Where("timestamp = ?", ts).Find(&acknowledges)
	if len(acknowledges) == 1 {
		acknowledge := acknowledges[0]
		var acknowledgeactions []AcknowledgeAction
		db.Model(&acknowledge).Related(&acknowledgeactions, "AcknowledgeActions")
		acknowledge.AcknowledgeActions = acknowledgeactions
		for _, action := range acknowledge.AcknowledgeActions {
			fmt.Println(action.ID)
		}
		acknowledgeAction := AcknowledgeAction{Value: value, AckID: acknowledges[0].ID}
		AddAcknowledgement(acknowledge, acknowledgeAction)
		if FullyAcknowledged(acknowledge, attachmentAcknowledgeAction.Channel, team) &&
			!strings.Contains(attachmentAcknowledgeAction.OriginalMessage.Text, "~") {
			botConn.UpdateMessage(attachmentAcknowledgeAction.Channel.ID, attachmentAcknowledgeAction.OriginalMessage.Timestamp, "~"+attachmentAcknowledgeAction.OriginalMessage.Text+"~")
		}
	} else {
		//fmt.Println(userID)
		fmt.Println("more than one acknowledgement... " + strconv.Itoa(len(acknowledges)))
	}
}

/*func RemindCallback(attachmentAcknowledgeAction slack.AttachmentActionCallback) {
	api := GetSlackClient()
	userID := attachmentAcknowledgeAction.User.ID
	ts := attachmentAcknowledgeAction.OriginalMessage.Timestamp
	//attachmentAcknowledgeAction.OriginalMessage
	var acknowledges []AcknowledgeMessage
	db.Where("user_id = ? AND timestamp = ?", userID, ts).Find(&acknowledges)
	if len(acknowledges) == 1 {
		acknowledge := acknowledges[0]
		var acknowledgeactions []AcknowledgeAction
		db.Model(&acknowledge).Related(&acknowledgeactions, "AcknowledgeActions")
		acknowledge.AcknowledgeActions = acknowledgeactions
		var userIDs []string
		var allUserIds []string
		for _, user := range GetUsers() {
			allUserIds = append(allUserIds, user.ID)
		}
		for _, action := range acknowledge.AcknowledgeActions {
			userIDs = append(userIDs, action.UserID)
		}
		hasntAck := HasntAcknowledged(userIDs, allUserIds)
		for _, user := range hasntAck {
			api.PostMessage(GetUser(user).ChannelID, "Please acknowledge the message in: #"+attachmentAcknowledgeAction.Channel.Name, slack.PostMessageParameters{})
		}
	} else {
		fmt.Println(userID)
		fmt.Println("more than one acknowledgement... " + strconv.Itoa(len(acknowledges)))
	}
	//attachmentAcknowledgeAction
}*/

func FullyAcknowledged(acknowledge AcknowledgeMessage, channel slack.Channel, team Team) bool {
	/*activeMembers := 0
	for _, member := range channel.Members {
		if IsActiveUser(member) {
			activeMembers++
		}
	}*/
	botConn := slack.New(team.BotToken)
	channelPtr, err := botConn.GetChannelInfo(channel.ID)
	if err != nil {
		groupPtr, err1 := botConn.GetGroupInfo(channel.ID)
		check(err1)
		fmt.Println(strconv.Itoa(len(groupPtr.Members)))
		if len(acknowledge.AcknowledgeActions) >= len(groupPtr.Members)-1 {
			return true
		} else {
			return false
		}
	}

	check(err)

	fmt.Println(strconv.Itoa(len(channelPtr.Members)))
	if len(acknowledge.AcknowledgeActions) >= len(channelPtr.Members)-1 {
		return true
	} else {
		return false
	}
}

func HasntAcknowledged(slice1 []string, slice2 []string) []string {
	diffStr := []string{}
	m := map[string]int{}

	for _, s1Val := range slice1 {
		m[s1Val] = 1
	}
	for _, s2Val := range slice2 {
		m[s2Val] = m[s2Val] + 1
	}

	for mKey, mVal := range m {
		if mVal == 1 {
			diffStr = append(diffStr, mKey)
		}
	}

	return diffStr
}

func AddAcknowledgement(acknwoledge AcknowledgeMessage, acknowledgeAction AcknowledgeAction) {
	//fmt.Println("adding acknowledgement")
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
