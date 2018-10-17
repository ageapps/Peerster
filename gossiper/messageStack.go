package gossiper

import (
	"fmt"
	"sync"

	"github.com/ageapps/Peerster/data"
	"github.com/ageapps/Peerster/logger"
)

// MessageStack struct
// that contains as keys the origin
// and as value an array of the rumor
// messages received by that origin
type MessageStack struct {
	Messages map[string][]data.RumorMessage
	mux      sync.Mutex
}

// AreMessagesMissing func
func (stack *MessageStack) AreMessagesMissing(origin string, id uint32) bool {
	stack.mux.Lock()
	defer stack.mux.Unlock()
	messages := stack.Messages[origin]
	if len(messages) <= 0 {
		return false
	}
	if id+1 == messages[len(messages)-1].ID {
		return false
	}
	return true
}

//AddMessage func
func (stack *MessageStack) AddMessage(msg data.RumorMessage) {
	stack.mux.Lock()
	defer stack.mux.Unlock()
	id := msg.ID
	messages := stack.Messages[msg.Origin]
	lastMessageID := messages[len(messages)-1].ID
	if id == uint32(lastMessageID+1) {
		stack.Messages[msg.Origin] = append(messages, msg)
		logger.Log(fmt.Sprintf("Adding message <%v> from <%v>", id, msg.Origin))
	}
}

// PrintStack func
func (stack *MessageStack) PrintStack() {
	stack.mux.Lock()
	defer stack.mux.Unlock()
	for address := range stack.Messages {
		logger.Log(fmt.Sprintf("Sender <%v>, last message %v", address, stack.Messages[address]))
	}
}

func (stack *MessageStack) GetStackMap() map[string]uint32 {
	stack.mux.Lock()
	defer stack.mux.Unlock()
	var stackMap = make(map[string]uint32)
	for address := range stack.Messages {
		messages := stack.Messages[address]
		lastID := uint32(messages[len(messages)-1].ID + 1)
		stackMap[address] = lastID
	}
	return stackMap
}

func (stack *MessageStack) getStatusMessage() *data.StatusPacket {
	stack.mux.Lock()
	defer stack.mux.Unlock()
	var vector []data.PeerStatus
	for address := range stack.Messages {
		messages := stack.Messages[address]
		peerStatus := data.PeerStatus{Identifier: address, NextID: uint32(messages[len(messages)-1].ID + 1)}
		vector = append(vector, peerStatus)
	}
	return data.NewStatusPacket(vector)
}
