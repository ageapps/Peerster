package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"

	"github.com/ageapps/Peerster/utils"
)

// Health message
func Health(w http.ResponseWriter, r *http.Request) {
	// A very simple health check.
	status := "HEALTHY"
	send(&w, &status)
}

// Index page
func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to Peerster!")
}

// GetMessages func
func GetMessages(w http.ResponseWriter, r *http.Request) {
	name, ok := getNameFromRequest(r)
	if !ok {
		sendError(&w, errors.New("Error: no peer requested"))
		return
	}
	send(&w, getGossiperMessages(name))
}

// GetPrivateMessages func
func GetPrivateMessages(w http.ResponseWriter, r *http.Request) {
	name, ok := getNameFromRequest(r)
	if !ok {
		sendError(&w, errors.New("Error: no peer requested"))
		return
	}
	send(&w, getGossiperPrivateMessages(name))
}

// GetRoutes func
func GetRoutes(w http.ResponseWriter, r *http.Request) {
	name, ok := getNameFromRequest(r)
	if !ok {
		sendError(&w, errors.New("Error: no peer requested"))
		return
	}
	send(&w, getGossiperRoutes(name))
}

// GetNodes func
func GetNodes(w http.ResponseWriter, r *http.Request) {
	name, ok := getNameFromRequest(r)
	if !ok {
		sendError(&w, errors.New("Error: no peer requested"))
		return
	}
	send(&w, getGossiperPeers(name))
}

// PostMessage func
func PostMessage(w http.ResponseWriter, r *http.Request) {
	params := *readBody(&w, r)
	name, ok := params["name"].(string)
	if !ok {
		sendError(&w, errors.New("Error: no peer requested"))
		return
	}
	msg, ok := params["msg"].(string)
	if !ok || !sendMessage(name, msg) {
		sendError(&w, errors.New("Error while sending new message"))
		return
	}
	send(&w, getGossiperMessages(name))
}

// PostPrivateMessage func
func PostPrivateMessage(w http.ResponseWriter, r *http.Request) {
	params := *readBody(&w, r)
	name, ok := params["name"].(string)
	if !ok {
		sendError(&w, errors.New("Error: no peer requested"))
		return
	}
	dest, ok := params["destination"].(string)
	if !ok {
		sendError(&w, errors.New("Error: no destination requested"))
		return
	}
	msg, ok := params["msg"].(string)
	if !ok || !sendPrivateMessage(name, dest, msg) {
		sendError(&w, errors.New("Error while sending new message"))
		return
	}
	send(&w, getGossiperMessages(name))
}

// PostNode func
func PostNode(w http.ResponseWriter, r *http.Request) {
	params := *readBody(&w, r)
	name, ok := params["name"].(string)
	if !ok {
		sendError(&w, errors.New("Error: no peer requested"))
		return
	}
	peer, ok := params["node"].(string)
	if !ok || !addPeer(name, peer) {
		sendError(&w, errors.New("Error while adding new peer"))
		return
	}
	send(&w, getGossiperPeers(name))
}

// GetID func
func GetID(w http.ResponseWriter, r *http.Request) {
	name, ok := getNameFromRequest(r)
	if !ok {
		sendError(&w, errors.New("Error: no peer requested"))
		return
	}

	send(&w, getStatusResponse(name))
}

// Delete gossiper
func Delete(w http.ResponseWriter, r *http.Request) {
	params := *readBody(&w, r)
	name, ok := params["name"].(string)
	if !ok {
		sendError(&w, errors.New("Error: no peer requested"))
		return
	}
	deleteGossiper(name)
	sendOk(&w)
}

// Start gossiper
func Start(w http.ResponseWriter, r *http.Request) {

	params := *readBody(&w, r)

	name := params["name"].(string)
	var peers = utils.PeerAddresses{}
	var gossipAddr = utils.PeerAddress{}

	if address, ok := params["address"]; ok {
		if err := gossipAddr.Set(address.(string)); err != nil {
			sendError(&w, err)
			return
		}
	}
	if peerParams, ok := params["peers"]; ok {
		if err := peers.Set(peerParams.(string)); err != nil {
			sendError(&w, err)
			return
		}
	}
	gossiperName := startGossiper(name, gossipAddr.String(), &peers)
	if gossiperName == "" {
		sendError(&w, errors.New("Error starting gossiper"))
		return
	}
	send(&w, getStatusResponse(gossiperName))
}

func send(w *http.ResponseWriter, v interface{}) {
	if reflect.ValueOf(v).IsNil() {
		sendError(w, errors.New("Error sending response"))
		return
	}
	(*w).Header().Set("Content-Type", "application/json; charset=UTF-8")
	(*w).WriteHeader(http.StatusOK)
	if err := json.NewEncoder(*w).Encode(v); err != nil {
		panic(err)
	}
}
func sendOk(w *http.ResponseWriter) {
	(*w).WriteHeader(http.StatusOK)
	fmt.Fprintf(*w, "OK")
}
func sendError(w *http.ResponseWriter, msg error) {
	(*w).WriteHeader(http.StatusBadRequest)
	fmt.Fprintf(*w, "There was an error processing the request: %v\n", msg.Error())
}
func readBody(w *http.ResponseWriter, r *http.Request) *map[string]interface{} {
	var params map[string]interface{}
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		sendError(w, err)
		return nil
	}
	if err := r.Body.Close(); err != nil {
		sendError(w, err)
		return nil
	}
	if err := json.Unmarshal(body, &params); err != nil {
		sendError(w, err)
		return nil
	}
	return &params
}

func getNameFromRequest(r *http.Request) (string, bool) {
	name, ok := r.URL.Query()["name"]
	if !ok || len(name[0]) < 1 {
		return "", false
	}
	return name[0], true
}
