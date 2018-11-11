package data

import (
)
// SimpleMessage struct
type SimpleMessage struct {
	OriginalName  string
	RelayPeerAddr string
	Contents      string
}

// Message to send
type Message struct {
	Text string
	Destination string
	ID   uint32
}

// PrivateMessage to send
type PrivateMessage struct {
	Origin string
	ID   uint32
	Destination string
	Text string
	HopLimit   uint32
}

// RumorMessage to send
type RumorMessage struct {
	Origin string `json:"origin"`
	ID     uint32 `json:"id"`
	Text   string `json:"text"`
}

// PeerStatus to send
type PeerStatus struct {
	Identifier string
	NextID     uint32
}

func (rumor *RumorMessage) IsRouteRumor() bool{
	return rumor.Text == ""
}
func (msg *Message) IsPrivate() bool{
	return msg.Destination != ""
}
