package main

import (
	"encoding/json"

	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	//"os"

	"github.com/gorilla/schema"
	"github.com/nlopes/slack"
)

// SLASH COMMANDS: /instructors /anonymousQuestion /acknowledge...

func Instructors(w http.ResponseWriter, r *http.Request) {

	var slashPayload SlashPayload

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if origin := r.Header.Get("Origin"); origin != "" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers",
			"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	}
	// Stop here if its Preflighted OPTIONS request
	if r.Method == "OPTIONS" {
		return
	}

	err := r.ParseForm()
	check(err)
	decoder := schema.NewDecoder()
	err = decoder.Decode(&slashPayload, r.PostForm)
	check(err)

	if slashPayload.SSLCheck == 1 {
		fmt.Println("ssl check")
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(403) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}
	StartInstructorConversation(slashPayload.UserID, slashPayload.UserName, slashPayload.TeamID)

	w.WriteHeader(http.StatusOK)
}

func ThisWay(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Disposition", "attachment; filename=thisway.mobileconfig")
	w.Header().Set("Content-Type", "application/x-apple-aspen-config")
	http.ServeFile(w, r, "../this-way.mobileconfig")
}

func RegisterEveryone(w http.ResponseWriter, r *http.Request) {
	var slashPayload SlashPayload

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if origin := r.Header.Get("Origin"); origin != "" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers",
			"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	}
	// Stop here if its Preflighted OPTIONS request
	if r.Method == "OPTIONS" {
		return
	}

	err := r.ParseForm()
	check(err)
	decoder := schema.NewDecoder()
	err = decoder.Decode(&slashPayload, r.PostForm)
	check(err)

	if slashPayload.SSLCheck == 1 {
		fmt.Println("ssl check")
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(403) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}
	team := GetTeam(slashPayload.TeamID)
	if slashPayload.UserID == team.InstallerID {
		go RegisterAll(team)
	}
	w.WriteHeader(http.StatusOK)
}

func EnterGrades(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if origin := r.Header.Get("Origin"); origin != "" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers",
			"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	}
	// Stop here if its Preflighted OPTIONS request
	if r.Method == "OPTIONS" {
		return
	}

	var manualGrade ManualGrade

	body, err := ioutil.ReadAll(r.Body)
	check(err)
	if err = json.Unmarshal(body, &manualGrade); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}

	EnterManualGrades(manualGrade)

	w.WriteHeader(http.StatusOK)
}

func GetGrades(w http.ResponseWriter, r *http.Request) {

	var slashPayload SlashPayload
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if origin := r.Header.Get("Origin"); origin != "" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers",
			"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	}
	// Stop here if its Preflighted OPTIONS request
	if r.Method == "OPTIONS" {
		return
	}

	err := r.ParseForm()
	check(err)
	decoder := schema.NewDecoder()
	err = decoder.Decode(&slashPayload, r.PostForm)
	check(err)

	if slashPayload.SSLCheck == 1 {
		fmt.Println("ssl check")
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(403) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}
	//team := GetTeam(slashPayload.TeamID)
	go GetStudentGrades(slashPayload.UserID, slashPayload.TeamID, slashPayload.Text)

	w.WriteHeader(http.StatusOK)
}

func AnonymousQuestion(w http.ResponseWriter, r *http.Request) {

	var slashPayload SlashPayload

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if origin := r.Header.Get("Origin"); origin != "" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers",
			"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	}
	// Stop here if its Preflighted OPTIONS request
	if r.Method == "OPTIONS" {
		return
	}

	err := r.ParseForm()
	check(err)
	decoder := schema.NewDecoder()
	err = decoder.Decode(&slashPayload, r.PostForm)
	check(err)

	fmt.Println("checkpt 1")

	if slashPayload.Token != "mDMyhrbIMX1k0U6YTsPsw3ca" {
		fmt.Println("wrong token")
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(403) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}

	if slashPayload.SSLCheck == 1 {
		fmt.Println("ssl check")
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(403) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}
	teamID := slashPayload.TeamID
	messageText := slashPayload.Text
	PostAnonymousQuestion(messageText, teamID)
	w.WriteHeader(http.StatusOK)
}

func Acknowledge(w http.ResponseWriter, r *http.Request) {

	var slashPayload SlashPayload

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if origin := r.Header.Get("Origin"); origin != "" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers",
			"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	}
	// Stop here if its Preflighted OPTIONS request
	if r.Method == "OPTIONS" {
		return
	}

	err := r.ParseForm()
	check(err)
	decoder := schema.NewDecoder()
	err = decoder.Decode(&slashPayload, r.PostForm)
	check(err)

	if slashPayload.SSLCheck == 1 {
		fmt.Println("ssl check")
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(403) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}

	teamID := slashPayload.TeamID
	if IsDemoTeam(teamID) {
		DemoAcknowledgePost(teamID, slashPayload.UserID, slashPayload.ChannelID, slashPayload.Text)
		return
	}

	AcknowledgePost(teamID, slashPayload.UserID, slashPayload.ChannelID, slashPayload.Text)

	team := GetTeam(teamID)
	botConn := slack.New(team.BotToken)

	//Respond with acknowledge button
	params := slack.PostMessageParameters{}
	attachment := slack.Attachment{CallbackID: "acknowledge", Fallback: "acknowledge service not working properly"}
	attachment.Actions = append(attachment.Actions, slack.AttachmentAction{Name: "acknowledge", Text: "Acknowledge", Type: "button"})
	params.Attachments = append(params.Attachments, attachment)
	channelID, ts, _ := botConn.PostMessage(slashPayload.ChannelID, slashPayload.Text+" - @"+slashPayload.UserName, params)
	acknowledgeMsg := AcknowledgeMessage{UserID: slashPayload.UserID, Timestamp: ts, ChannelID: channelID}
	db.Create(&acknowledgeMsg)

	/*
		params1 := slack.PostMessageParameters{}
		attachment1 := slack.Attachment{CallbackID: "remind", Fallback: "remind service not working properly"}
		attachment1.Actions = append(attachment1.Actions, slack.AttachmentAction{Name: "remind", Text: "Remind", Type: "button"})
		params1.Attachments = append(params1.Attachments, attachment1)
		api.PostMessage(GetUser(slashPayload.UserID).ChannelID, "Click the `Remind` Button to remind your friends about the following announcement ```"+acknowledgeMsg.Text+"```", params1)
	*/

	acknowledgeAction := AcknowledgeAction{AckID: acknowledgeMsg.ID, UserID: slashPayload.UserID, TeamID: team.TeamID, Value: ""}
	db.Create(&acknowledgeAction)

	w.WriteHeader(http.StatusOK)
}

// Handles an SSL Check for slash command
func SSLCheck(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if origin := r.Header.Get("Origin"); origin != "" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers",
			"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(403) // unprocessable entity
}

// EVENTS API
// Handles event listening and routes further to other actions
func Events(w http.ResponseWriter, r *http.Request) {

	var event OuterEvent

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if origin := r.Header.Get("Origin"); origin != "" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers",
			"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	}
	// Stop here if its Preflighted OPTIONS request
	if r.Method == "OPTIONS" {
		return
	}
	fmt.Println("in events")

	body, err := ioutil.ReadAll(r.Body)
	check(err)
	if err = json.Unmarshal(body, &event); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}

	fmt.Println(event.Event.Type)

	switch event.Type {
	case "url_verification":
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(&ChallengeResponse{Token: event.Token, Challenge: event.Challenge}); err != nil {
			panic(err)
		}
	case "event_callback":
		switch event.Event.Type {
		case "file_shared": //files:read
			HandleFileShared(event.Event, event.TeamID)
		case "team_join": //users:read
			WelcomeToTeam(event.Event, event.TeamID)
		case "message":
			IsInMidterms(event.Event, event.TeamID)
		}
	}

}

// Takes a file and re-uploads as an assignment
func Assign(w http.ResponseWriter, r *http.Request) {

	var slashPayload SlashPayload

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if origin := r.Header.Get("Origin"); origin != "" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers",
			"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	}
	// Stop here if its Preflighted OPTIONS request
	if r.Method == "OPTIONS" {
		return
	}

	err := r.ParseForm()
	check(err)
	decoder := schema.NewDecoder()
	err = decoder.Decode(&slashPayload, r.PostForm)
	check(err)
	DateInteractive(slashPayload.UserID, slashPayload.ChannelID, slashPayload.TeamID)
	w.WriteHeader(http.StatusOK)
}

// Takes in a Student object and adds it to the Database
func OAuth(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if origin := r.Header.Get("Origin"); origin != "" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers",
			"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	}
	// Stop here if its Preflighted OPTIONS request
	if r.Method == "OPTIONS" {
		return
	}

	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")

	GetOAuthToken("135270668007.135692085812", "5bc0dc4bba1567dbf09015375cfbd373", code, "https://teach-bot-api.com/oauth")

	//fmt.Println(scope)
	//os.Setenv("SLACK_TOKEN", accessToken)
	//check(err)
	//fmt.Println("set the slack token to " + accessToken)
	fmt.Println("code is" + code)
	fmt.Println(state)
	// if state != remembered state
	w.WriteHeader(http.StatusOK)
	//if err := json.NewEncoder(w).Encode(); err != nil {
	//	panic(err)
	//}
}

//INTERACTIVE MESSAGES

func Interactive(w http.ResponseWriter, r *http.Request) {
	//api := GetSlackClient()
	var actionResponse ActionResponse

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if origin := r.Header.Get("Origin"); origin != "" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers",
			"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	}
	// Stop here if its Preflighted OPTIONS request
	if r.Method == "OPTIONS" {
		return
	}

	err := r.ParseForm()
	check(err)
	decoder := schema.NewDecoder()
	err = decoder.Decode(&actionResponse, r.PostForm)
	check(err)

	var attachmentActionCallback slack.AttachmentActionCallback
	if err = json.Unmarshal([]byte(actionResponse.Payload), &attachmentActionCallback); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
		return
	}
	teamID := attachmentActionCallback.Team.ID
	team := GetTeam(teamID)
	botConn := slack.New(team.BotToken)

	switch attachmentActionCallback.CallbackID {
	case "student_or_teacher":
		AddToDatabase(attachmentActionCallback)
	case "assignment_due":
		HandleDate(attachmentActionCallback)
	case "acknowledge":
		AcknowledgeCallback(attachmentActionCallback)
	case "reminder":
		botConn.PostMessage(attachmentActionCallback.Channel.ID, "Great I'll send you a few reminders as the due date approaches!", slack.PostMessageParameters{})
	case "reminderAssignment":
		botConn.PostMessage(attachmentActionCallback.Channel.ID, "Got it -- I'll let them know", slack.PostMessageParameters{})
	case "submission_type":
		DemoHandleViewSubmission(attachmentActionCallback)
		//case "remind":
		//	RemindCallback(attachmentActionCallback)
	}

	//fmt.Println(attachmentActionCallback)
}

// GetOAuthToken retrieves an AccessToken
func GetOAuthToken(clientID, clientSecret, code, redirectURI string) (accessToken string, scope string, err error) {
	response, err := GetOAuthResponse(clientID, clientSecret, code, redirectURI)
	if err != nil {
		return "", "", err
	}
	var team Team
	team.TeamID = response.TeamID
	team.InstallerID = response.UserID
	team.Token = response.AccessToken
	team.BotToken = response.Bot.BotAccessToken
	team.BotID = response.Bot.BotUserID
	if IsDemoTeam(team.TeamID) {
		InitializeDemo(team.TeamID)
	}
	db.Create(&team)
	fmt.Println(response.Scope)
	//fmt.Println("bot token is:" + response.Bot.BotAccessToken + " " + response.Bot.BotUserID)
	return response.AccessToken, response.Scope, nil
}

func GetOAuthResponse(clientID, clientSecret, code, redirectURI string) (resp *OAuthResponse, err error) {
	values := url.Values{
		"client_id":     {clientID},
		"client_secret": {clientSecret},
		"code":          {code},
		"scope":         {"bot"},
		"redirect_uri":  {redirectURI},
	}

	var HTTPClient http.Client
	response, err := HTTPClient.PostForm("https://slack.com/api/oauth.access", values)
	check(err)

	body, err := ioutil.ReadAll(response.Body)
	check(err)
	var oAuthResponse OAuthResponse
	if err = json.Unmarshal(body, &oAuthResponse); err != nil {
		panic(err)
	}

	return &oAuthResponse, nil
}
