package data

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

// StatusPacket to send
type StatusPacket struct {
	Want  []PeerStatus
	Route string
}

// NewStatusPacket create
func NewStatusPacket(want *[]PeerStatus, route string) *StatusPacket {
	if want != nil {
		return &StatusPacket{Want: *want}
	}
	return &StatusPacket{Route: route}
}

// IsRouteRumor check if Rumor is a route message
func (rumor *RumorMessage) IsRouteRumor() bool {
	return rumor.Text == ""
}

// NewRumorMessage create
func NewRumorMessage(origin string, ID uint32, text string) *RumorMessage {
	return &RumorMessage{origin, ID, text}
}

// IsRouteStatus create
func (status *StatusPacket) IsRouteStatus() bool {
	return status.Route != ""
}
