package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
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
	send(&w, getGossiperMessages())
}

// GetNodes func
func GetNodes(w http.ResponseWriter, r *http.Request) {
	send(&w, getGossiperPeers())
}

// PostMessage func
func PostMessage(w http.ResponseWriter, r *http.Request) {

	var params map[string]interface{}
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		sendError(&w, err)
		return
	}
	if err := r.Body.Close(); err != nil {
		sendError(&w, err)
		return
	}
	if err := json.Unmarshal(body, &params); err != nil {
		sendError(&w, err)
		return
	}
	msg, ok := params["msg"].(string)
	if !ok || !sendMessage(msg) {
		sendError(&w, errors.New("Error while sending new message"))
		return
	}
	send(&w, getGossiperMessages())
}

// PostNode func
func PostNode(w http.ResponseWriter, r *http.Request) {

	var params map[string]interface{}
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		sendError(&w, err)
		return
	}
	if err := r.Body.Close(); err != nil {
		sendError(&w, err)
		return
	}
	if err := json.Unmarshal(body, &params); err != nil {
		sendError(&w, err)
		return
	}
	peer, ok := params["node"].(string)
	if !ok || !addPeer(peer) {
		sendError(&w, errors.New("Error while adding new peer"))
		return
	}
	send(&w, getGossiperPeers())
}

// GetID func
func GetID(w http.ResponseWriter, r *http.Request) {
	send(&w, getStatusResponse())
}

// Start gossiper
func Start(w http.ResponseWriter, r *http.Request) {

	var params map[string]interface{}
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		sendError(&w, err)
		return
	}
	if err := r.Body.Close(); err != nil {
		sendError(&w, err)
		return
	}
	if err := json.Unmarshal(body, &params); err != nil {
		sendError(&w, err)
		return
	}

	name := params["name"].(string)
	var peers = utils.PeerAddresses{}
	var gossipAddr = utils.PeerAddress{IP: net.ParseIP("127.0.0.1"), Port: 5000}

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
	startGossiper(name, gossipAddr.String(), &peers)
	send(&w, getStatusResponse())
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
