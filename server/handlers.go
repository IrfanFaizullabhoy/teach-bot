package main

import (
	"encoding/json"

	"fmt"
	"io/ioutil"
	"net/http"

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
	fmt.Println("checkpt 2")
	fmt.Println(slashPayload.UserID)
	StartInstructorConversation(slashPayload.UserID)

	w.WriteHeader(http.StatusOK)
	//if err := json.NewEncoder(w).Encode(); err != nil {
	//	panic(err)
	//}
}
