package main

import (
	"encoding/json"

	"io/ioutil"
	"net/http"
	//"fmt"
)

// Takes in an email and password, returns User Object
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

	var user User
	db.Where("email = ?", login.Email).First(&user)
	if hashedPassword == user.HashedPassword {
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(user); err != nil {
			panic(err)
		}
	} else { // login failed
		w.WriteHeader(http.StatusOK)
	}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		panic(err)
	}
}

// Takes in a User object and adds it to the Database
func SignUp(w http.ResponseWriter, r *http.Request) {

	var user User

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
	if err = json.Unmarshal(body, &user); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}
	// DATABASE CODE //

	CreateUser(user, db)

	// END OF DATABASE CODE //
	w.WriteHeader(http.StatusOK)
	//if err := json.NewEncoder(w).Encode(); err != nil {
	//	panic(err)
	//}
}
