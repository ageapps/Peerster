package gossiper

import (
	"fmt"
	"sync"

	"github.com/ageapps/Peerster/pkg/data"
	"github.com/ageapps/Peerster/pkg/logger"
)

// PrivateStack struct
// that contains as keys the origin
// and as value an array of the rumor
// messages received by that origin
type PrivateStack struct {
	Messages map[string][]data.PrivateMessage
	mux      sync.Mutex
}

// CompareMessage func
func (stack *PrivateStack) CompareMessage(origin string, id uint32) string {
	stack.mux.Lock()
	defer stack.mux.Unlock()
	messages, ok := stack.Messages[origin]
	if !ok || len(messages) <= 0 {
		return NEW_MESSAGE
	}
	lastMessageID := messages[len(messages)-1].ID
	logger.Log(fmt.Sprintf("Comparing messages %v/%v", lastMessageID, id))
	switch {
	case id == lastMessageID:
		return IN_SYNC
	case id > lastMessageID:
		return NEW_MESSAGE
	case id < lastMessageID:
		return OLD_MESSAGE
	}
	return ""
}

// GetPrivateMessage func
func (stack *PrivateStack) GetPrivateMessage(origin string, id uint32) *data.PrivateMessage {
	stack.mux.Lock()
	defer stack.mux.Unlock()
	messages, ok := stack.Messages[origin]
	if !ok || len(messages) <= 0 {
		return nil
	}
	for _, msg := range messages {
		if msg.ID == id {
			return &msg
		}
	}
	return nil
}

//AddMessage func
func (stack *PrivateStack) AddMessage(msg data.PrivateMessage) {
	stack.mux.Lock()
	defer stack.mux.Unlock()
	id := msg.ID
	messages, ok := stack.Messages[msg.Origin]
	if !ok {
		stack.Messages[msg.Origin] = []data.PrivateMessage{msg}
	} else {
		lastMessageID := messages[len(messages)-1].ID
		if id == uint32(lastMessageID+1) {
			stack.Messages[msg.Origin] = append(messages, msg)
		}
	}
	logger.Log(fmt.Sprintf("Message appended to private stack Origin:%v ID:%v", msg.Origin, id))
}

// PrintStack func
func (stack *PrivateStack) PrintStack() {
	stack.mux.Lock()
	defer stack.mux.Unlock()
	for address := range stack.Messages {
		logger.Log(fmt.Sprintf("Sender <%v>, last message %v", address, stack.Messages[address]))
	}
}

// GetStackMap to get a map
// with latest ids saved from each origin
func (stack *PrivateStack) GetStackMap() *map[string]uint32 {
	stack.mux.Lock()
	defer stack.mux.Unlock()
	var stackMap = make(map[string]uint32)
	for origin := range stack.Messages {
		stackMap[origin] = stack.getLatestMessageID(origin)
	}
	return &stackMap
}

func (stack *PrivateStack) getLatestMessageID(origin string) uint32 {
	messages := stack.Messages[origin]
	lastID := uint32(messages[len(messages)-1].ID)
	return lastID
}

// GetLatestMessages function
// returns an array with the latest rumor messages
func (stack *PrivateStack) GetLatestMessages() *[]data.PrivateMessage {
	stack.mux.Lock()
	defer stack.mux.Unlock()
	var latestMessages = []data.PrivateMessage{}
	for address := range stack.Messages {
		messages := stack.Messages[address]
		latestMessages = append(latestMessages, messages[len(messages)-1])
	}
	return &latestMessages
}

func (stack *PrivateStack) getStatusMessage() *data.StatusPacket {
	stack.mux.Lock()
	defer stack.mux.Unlock()
	var vector []data.PeerStatus
	for address := range stack.Messages {
		messages := stack.Messages[address]
		peerStatus := data.PeerStatus{Identifier: address, NextID: uint32(messages[len(messages)-1].ID + 1)}
		vector = append(vector, peerStatus)
	}
	return data.NewStatusPacket(&vector, "")
}

func (stack *PrivateStack) getPrivateStack() *map[string][]data.PrivateMessage {
	stack.mux.Lock()
	defer stack.mux.Unlock()
	return &stack.Messages
}
