package data

// SimpleMessage struct
type SimpleMessage struct {
	OriginalName  string
	RelayPeerAddr string
	Contents      string
}

// Message to send
type Message struct {
	Text string
	ID   uint32
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
