package main

import (
	"fmt"
	"reflect"

	"github.com/ageapps/Peerster/data"
	"github.com/ageapps/Peerster/gossiper"
	"github.com/ageapps/Peerster/logger"
	"github.com/ageapps/Peerster/utils"
)

var (
	serverGossiper = make(map[string]*gossiper.Gossiper)
)

// StatusResponse struct
type StatusResponse struct {
	Address string `json:"address"`
	Name    string `json:"name"`
}

func startGossiper(name, address string, peers *utils.PeerAddresses) string {
	logger.CreateLogger(name, address, true)
	newGossiper, err := gossiper.NewGossiper(address, name, false)
	if err != nil {
		logger.Log(fmt.Sprintln("Error creating new Gossiper ", err))
		for _, gossiper := range serverGossiper {
			if gossiper.Address.String() == address || gossiper.Name == name {
				logger.Log(fmt.Sprintf("Running gossiper found Name:%v Address:%v", gossiper.Name, gossiper.Address))
				logger.Log(fmt.Sprintf("Running gossiper found Name:%v Address:%v", name, address))
				return gossiper.Name
			}
		}
		return ""
	}
	serverGossiper[name] = newGossiper
	serverGossiper[name].SetPeers(peers)
	go serverGossiper[name].ListenToPeers()
	serverGossiper[name].StartEntropyTimer()
	return name
}

func getGossiperMessages(name string) *[]data.RumorMessage {
	if reflect.ValueOf(serverGossiper).IsNil() {
		return nil
	}
	return serverGossiper[name].GetLatestMessages()
}

func getGossiperPeers(name string) *[]string {
	if reflect.ValueOf(serverGossiper).IsNil() {
		return nil
	}
	return serverGossiper[name].GetPeers()
}

func getStatusResponse(name string) *StatusResponse {
	if reflect.ValueOf(serverGossiper).IsNil() {
		return nil
	}
	return &StatusResponse{
		Name:    serverGossiper[name].Name,
		Address: serverGossiper[name].Address.String(),
	}
}

func deleteGossiper(name string) {
	if len(serverGossiper) > 0 {
		go serverGossiper[name].Kill()
		delete(serverGossiper, name)
	}
}

func addPeer(name, peer string) bool {
	if reflect.ValueOf(serverGossiper).IsNil() {
		return false
	}
	serverGossiper[name].AddPeer(peer)
	return true
}

func sendMessage(name, msg string) bool {
	if reflect.ValueOf(serverGossiper).IsNil() {
		return false
	}
	newMsg := &data.Message{
		Text: msg,
	}
	serverGossiper[name].HandleClientMessage(newMsg)
	return true
}
