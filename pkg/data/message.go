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

// PrivateMessage to send
type PrivateMessage struct {
	Origin      string
	ID          uint32
	Destination string
	Text        string
	HopLimit    uint32
}

// NewSimpleMessage create
func NewSimpleMessage(ogname, msg, relay string) *SimpleMessage {
	return &SimpleMessage{
		OriginalName:  ogname,
		RelayPeerAddr: relay,
		Contents:      msg,
	}
}

// NewPrivateMessage create
func NewPrivateMessage(origin string, ID uint32, destination, text string, hops uint32) *PrivateMessage {
	return &PrivateMessage{origin, ID, destination, text, hops}
}
