package main

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/ageapps/Peerster/pkg/file"

	"github.com/ageapps/Peerster/gossiper"
	"github.com/ageapps/Peerster/pkg/data"
	"github.com/ageapps/Peerster/pkg/logger"
	"github.com/ageapps/Peerster/pkg/router"
	"github.com/ageapps/Peerster/pkg/utils"
)

func NewGossiperPool() *GossiperPool {
	return &GossiperPool{
		gossipers: make(map[string]*gossiper.Gossiper),
	}
}

// GossiperPool struct cointaining gossipers
type GossiperPool struct {
	gossipers map[string]*gossiper.Gossiper
	mux       sync.Mutex
}

func (pool *GossiperPool) addGossiper(gossiper *gossiper.Gossiper) {
	pool.mux.Lock()
	pool.gossipers[gossiper.Name] = gossiper
	pool.mux.Unlock()
}
func (pool *GossiperPool) deleteGossiper(name string) {
	pool.mux.Lock()
	delete(pool.gossipers, name)
	pool.mux.Unlock()
}
func (pool *GossiperPool) getGossiper(name string) (foundGossiper *gossiper.Gossiper, found bool) {
	pool.mux.Lock()
	defer pool.mux.Unlock()
	foundGossiper, found = pool.gossipers[name]
	return
}
func (pool *GossiperPool) findGossiper(name, address string) (*gossiper.Gossiper, bool) {
	pool.mux.Lock()
	defer pool.mux.Unlock()
	for _, gossiper := range pool.gossipers {
		if gossiper.Address.String() == address || gossiper.Name == name {
			logger.Log(fmt.Sprintf("Running gossiper found Name:%v Address:%v", name, address))
			return gossiper, true
		}
	}
	return nil, false
}

var (
	gossiperPool = NewGossiperPool()
)

// StatusResponse struct
type StatusResponse struct {
	Address string `json:"address"`
	Name    string `json:"name"`
}

func startGossiper(name, address string, peers *utils.PeerAddresses) string {
	logger.CreateLogger(name, address, true)
	targetGossiper, found := gossiperPool.findGossiper(name, address)

	if !found {
		newGossiper, err := gossiper.NewGossiper(address, name, false, 5)
		if err != nil {
			logger.Log(fmt.Sprintln("Error creating new Gossiper ", err))
			return ""
		}
		targetGossiper = newGossiper
		if peers != nil && len(peers.GetAdresses()) > 0 {
			go targetGossiper.AddPeers(peers)
		}
		go func() {
			if err := targetGossiper.ListenToPeers(); err != nil {
				log.Fatal(err)
			}
		}()
		gossiperPool.addGossiper(targetGossiper)
	}
	return targetGossiper.Name
}

func getGossiperRoutes(name string) *router.RoutingTable {
	targetGossiper, found := gossiperPool.getGossiper(name)
	if !found {
		return nil
	}

	return targetGossiper.GetRoutes()
}

func indexFileInGossiper(name, file string) map[string]file.File {
	targetGossiper, found := gossiperPool.getGossiper(name)
	if !found {
		return nil
	}

	blob, localFile := gossiper.SaveLocalFile(file)
	targetGossiper.IndexAndPublishBundle(localFile, blob, uint32(10))

	return getGossiperFiles(name)
}

func getGossiperFiles(name string) map[string]file.File {
	targetGossiper, found := gossiperPool.getGossiper(name)
	if !found {
		return nil
	}
	files := targetGossiper.GetFiles()
	return files
}

func getGossiperMessages(name string) *[]data.RumorMessage {
	targetGossiper, found := gossiperPool.getGossiper(name)
	if !found {
		return nil
	}
	return targetGossiper.GetLatestMessages()
}
func getGossiperPrivateMessages(name string) *map[string][]data.PrivateMessage {
	targetGossiper, found := gossiperPool.getGossiper(name)
	if !found {
		return nil
	}
	return targetGossiper.GetPrivateMessages()
}

func getGossiperPeers(name string) *[]string {
	targetGossiper, found := gossiperPool.getGossiper(name)
	if !found {
		return nil
	}
	return targetGossiper.GetPeerArray()
}

func getStatusResponse(name string) *StatusResponse {
	targetGossiper, found := gossiperPool.getGossiper(name)
	if !found {
		return nil
	}
	return &StatusResponse{
		Name:    targetGossiper.Name,
		Address: targetGossiper.Address.String(),
	}
}

func deleteGossiper(name string) {
	targetGossiper, found := gossiperPool.getGossiper(name)
	if !found {
		return
	}
	go targetGossiper.Kill()
	gossiperPool.deleteGossiper(targetGossiper.Name)
}

func addPeer(name, peer string) bool {
	targetGossiper, found := gossiperPool.getGossiper(name)
	if !found {
		return false
	}
	targetGossiper.AddAndNotifyPeer(peer)
	return true
}

func sendMessage(name, msg string) bool {
	targetGossiper, found := gossiperPool.getGossiper(name)
	if !found {
		return false
	}
	newMsg := &data.Message{
		Text: msg,
	}
	targetGossiper.HandleClientMessage(newMsg)
	return true
}

func sendSearchMessage(name, keyboards string) bool {
	targetGossiper, found := gossiperPool.getGossiper(name)
	if !found {
		return false
	}
	newMsg := &data.Message{
		Keywords: strings.Split(keyboards, ","),
		Budget:   uint64(2),
	}
	targetGossiper.HandleClientMessage(newMsg)
	return true
}

func sendPrivateMessage(name, destination, msg string) bool {
	targetGossiper, found := gossiperPool.getGossiper(name)
	if !found {
		return false
	}
	newMsg := &data.Message{
		Text:        msg,
		Destination: destination,
	}
	targetGossiper.HandleClientMessage(newMsg)
	return true
}
func sendFileRequest(name, destination, fileName, hash string) bool {
	targetGossiper, found := gossiperPool.getGossiper(name)
	if !found {
		return false
	}
	newMsg := &data.Message{
		Destination: destination,
		FileName:    fileName,
		RequestHash: hash,
	}
	targetGossiper.HandleClientMessage(newMsg)
	return true
}
