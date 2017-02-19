package main

import (
	"encoding/json"

	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/gorilla/schema"
)

// Takes in an email and password, returns Student Object
func Login(w http.ResponseWriter, r *http.Request) {

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
	body, err := ioutil.ReadAll(r.Body)
	check(err)

	var login LoginMsg
	err = json.Unmarshal(body, &login)
	check(err)
	hashedPassword := string(hash(login.Password))

	var Student Student
	db.Where("email = ?", login.Email).First(&Student)
	if hashedPassword == Student.HashedPassword {
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(Student); err != nil {
			panic(err)
		}
	} else { // login failed
		w.WriteHeader(http.StatusOK)
	}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(Student); err != nil {
		panic(err)
	}
}

// Takes in a Student object and adds it to the Database
func SignUp(w http.ResponseWriter, r *http.Request) {

	var Student Student

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

	body, err := ioutil.ReadAll(r.Body)
	check(err)
	if err = json.Unmarshal(body, &Student); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}

	//CreateStudent(Student, db)

	w.WriteHeader(http.StatusOK)
	//if err := json.NewEncoder(w).Encode(); err != nil {
	//	panic(err)
	//}
}

// Takes in a Student object and adds it to the Database
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

	fmt.Println("checkpt 1")

	if slashPayload.Token != "qZNXELMoLQhPLLiae2ih7yER" {
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
	fmt.Println("checkpt 2")
	fmt.Println(slashPayload.UserID)
	StartInstructorConversation(slashPayload.UserID)

	w.WriteHeader(http.StatusOK)
	//if err := json.NewEncoder(w).Encode(); err != nil {
	//	panic(err)
	//}
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

// Takes in a Student object and adds it to the Database
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

	body, err := ioutil.ReadAll(r.Body)
	check(err)
	if err = json.Unmarshal(body, &event); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}
	switch event.Type {
	case "url_verification":
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(&ChallengeResponse{Token: event.Token}); err != nil {
			panic(err)
		}
	case "event_callback":
		if event.Event.Type == "file_shared" {
			DownloadFile(event.Event)
		}
	}

}

// Takes in a Student object and adds it to the Database
func Assignment(w http.ResponseWriter, r *http.Request) {

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

	if slashPayload.Token != "qZNXELMoLQhPLLiae2ih7yER" {
		fmt.Println("wrong token")
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(403) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}

	w.WriteHeader(http.StatusOK)
	//if err := json.NewEncoder(w).Encode(); err != nil {
	//	panic(err)
	//}
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

	accessToken, scope, err := GetOAuthToken("135270668007.135692085812", "5bc0dc4bba1567dbf09015375cfbd373", code, "https://2a84aaec.ngrok.io/oauth")

	fmt.Println(scope)
	os.Setenv("SLACK_TOKEN", accessToken)
	check(err)
	fmt.Println("set the slack token to " + accessToken)
	fmt.Println(code)
	fmt.Println(state)
	// if state != remembered state
	w.WriteHeader(http.StatusOK)
	//if err := json.NewEncoder(w).Encode(); err != nil {
	//	panic(err)
	//}
}

// GetOAuthToken retrieves an AccessToken
func GetOAuthToken(clientID, clientSecret, code, redirectURI string) (accessToken string, scope string, err error) {
	response, err := GetOAuthResponse(clientID, clientSecret, code, redirectURI)
	if err != nil {
		return "", "", err
	}
	return response.AccessToken, response.Scope, nil
}

func GetOAuthResponse(clientID, clientSecret, code, redirectURI string) (resp *OAuthResponse, err error) {
	values := url.Values{
		"client_id":     {clientID},
		"client_secret": {clientSecret},
		"code":          {code},
		"redirect_uri":  {redirectURI},
	}

	var HTTPClient http.Client
	response, err := HTTPClient.PostForm("https://slack.com/api/oauth.access", values)

	body, err := ioutil.ReadAll(response.Body)
	check(err)

	var oAuthResponse OAuthResponse
	if err = json.Unmarshal(body, &oAuthResponse); err != nil {
		panic(err)
	}

	return &oAuthResponse, nil
}
