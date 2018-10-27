package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Routes arr
type Routes []Route

// NewRouter func
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
			Handler(handler)
	}
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("../app/")))

	return router
}
