package main

import (
	"net/http"
)

// Route struct
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

var routes = Routes{
	// Route{"Index", "GET", "/", Index},
	Route{"Messages", "GET", "/message", GetMessages},
	Route{"Routes", "GET", "/routes", GetRoutes},
	Route{"Routes", "GET", "/files", GetFiles},
	Route{"Private Messages", "GET", "/private", GetPrivateMessages},
	Route{"Nodes", "GET", "/node", GetNodes},
	Route{"ID", "GET", "/id", GetID},
	Route{"Health", "GET", "/health", Health},
	Route{"Messages", "POST", "/message", PostMessage},
	Route{"Nodes", "POST", "/node", PostNode},
	Route{"Private Message", "POST", "/private", PostPrivateMessage},
	Route{"Start", "POST", "/start", Start},
	Route{"Delete", "POST", "/delete", Delete},
	Route{"Upload", "POST", "/upload", Upload},
	Route{"Upload", "POST", "/request", PostRequest},
	Route{"Upload", "POST", "/search", PostSearch},
}
