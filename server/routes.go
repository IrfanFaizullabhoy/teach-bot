package main

import (
	"github.com/gorilla/mux"
	"net/http"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler

		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}

	return router
}

var routes = Routes{

	Route{
		"Login",
		"POST",
		"/login",
		Login,
	},
	Route{
		"SignUp",
		"POST",
		"/signup",
		SignUp,
	},
	Route{
		"SSLCheck",
		"GET",
		"/instructors",
		SSLCheck,
	},
	Route{
		"Instructors",
		"POST",
		"/instructors",
		Instructors,
	},
	Route{
		"Events",
		"POST",
		"/events",
		Events,
	},
	Route{
		"OAuth",
		"GET",
		"/oauth",
		OAuth,
	},
}
