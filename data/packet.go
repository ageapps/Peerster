package data

// GossipPacket struct
type GossipPacket struct {
	Simple *SimpleMessage
	Rumor  *RumorMessage
	Status *StatusPacket
}

// StatusPacket to send
type StatusPacket struct {
	Want []PeerStatus
}

// NewSimpleMessage create
func NewSimpleMessage(ogname, msg, relay string) *SimpleMessage {
	return &SimpleMessage{
		OriginalName:  ogname,
		RelayPeerAddr: relay,
		Contents:      msg,
	}
}

// NewStatusPacket create
func NewStatusPacket(want []PeerStatus) *StatusPacket {
	return &StatusPacket{Want: want}
}

// NewRumorMessage create
func NewRumorMessage(origin string, ID uint32, text string) *RumorMessage {
	return &RumorMessage{origin, ID, text}
}

// IsStatusMessage function
func (packet *GossipPacket) IsStatusMessage() bool {
	return packet.Rumor == nil && packet.Status != nil && packet.Simple == nil
}

// IsRumorMessage function
func (packet *GossipPacket) IsRumorMessage() bool {
	return packet.Rumor != nil && packet.Status == nil && packet.Simple == nil
}
