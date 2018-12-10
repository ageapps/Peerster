package data

// UDPMessage struct
type UDPMessage struct {
	Address string
	Packet  GossipPacket
	Message Message
}

// DataRequest struct
type DataRequest struct {
	Origin      string
	Destination string
	HopLimit    uint32
	HashValue   HashValue
}

// DataReply struct
type DataReply struct {
	Origin      string
	Destination string
	HopLimit    uint32
	HashValue   HashValue
	Data        []byte
}

// SimpleMessage struct
type SimpleMessage struct {
	OriginalName  string
	RelayPeerAddr string
	Contents      string
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
