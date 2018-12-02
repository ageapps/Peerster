package data

// UDPMessage struct
type UDPMessage struct {
	Address string
	Packet  GossipPacket
	Message Message
}

// SimpleMessage struct
type SimpleMessage struct {
	OriginalName  string
	RelayPeerAddr string
	Contents      string
}

// Message to send
type Message struct {
	Text        string
	Destination string
	ID          uint32
}

// PrivateMessage to send
type PrivateMessage struct {
	Origin      string
	ID          uint32
	Destination string
	Text        string
	HopLimit    uint32
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

// IsRouteRumor check if Rumor is a route message
func (rumor *RumorMessage) IsRouteRumor() bool {
	return rumor.Text == ""
}

// IsPrivate check if is private message
func (msg *Message) IsPrivate() bool {
	return msg.Destination != ""
}
