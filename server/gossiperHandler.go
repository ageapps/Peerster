package main

import (
	"log"
	"reflect"

	"github.com/ageapps/Peerster/data"
	"github.com/ageapps/Peerster/gossiper"
	"github.com/ageapps/Peerster/logger"
	"github.com/ageapps/Peerster/utils"
)

var (
	serverGossiper *gossiper.Gossiper
)

// StatusResponse struct
type StatusResponse struct {
	Address string `json:"address"`
	Name    string `json:"name"`
}

func startGossiper(name, address string, peers *utils.PeerAddresses) {
	logger.CreateLogger(name, address, true)
	newGossiper, err := gossiper.NewGossiper(address, name, false)
	if err != nil {
		log.Fatal(err)
	}
	serverGossiper = newGossiper
	serverGossiper.SetPeers(peers)
	go serverGossiper.ListenToPeers()
	serverGossiper.StartEntropyTimer()
}

func getGossiperMessages() *[]data.RumorMessage {
	if reflect.ValueOf(serverGossiper).IsNil() {
		return nil
	}
	return serverGossiper.GetLatestMessages()
}

func getGossiperPeers() *[]string {
	if reflect.ValueOf(serverGossiper).IsNil() {
		return nil
	}
	return serverGossiper.GetPeers()
}

func getStatusResponse() *StatusResponse {
	if reflect.ValueOf(serverGossiper).IsNil() {
		return nil
	}
	return &StatusResponse{
		Name:    serverGossiper.Name,
		Address: serverGossiper.Address.String(),
	}
}

func addPeer(peer string) bool {
	if reflect.ValueOf(serverGossiper).IsNil() {
		return false
	}
	serverGossiper.AddPeer(peer)
	return true
}

func sendMessage(msg string) bool {
	if reflect.ValueOf(serverGossiper).IsNil() {
		return false
	}
	newMsg := &data.Message{
		Text: msg,
	}
	serverGossiper.HandleClientMessage(newMsg)
	return true
}
