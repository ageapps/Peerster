package main

import (
	"net/http"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

var routes = Routes{
	Route{"Index", "GET", "/", Index},
	Route{"Messages", "GET", "/message", GetMessages},
	Route{"Nodes", "GET", "/node", GetNodes},
	Route{"ID", "GET", "/id", GetID},
	Route{"Health", "GET", "/health", Health},
	Route{"Messages", "POST", "/message", PostMessage},
	Route{"Nodes", "POST", "/node", PostNode},
	Route{"Start", "POST", "/start", Start},
}
