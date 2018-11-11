package data

const (
	PACKET_SIMPLE = "SIMPLE"
	PACKET_RUMOR  = "RUMOR"
	PACKET_STATUS = "STATUS"
	PACKET_PRIVATE = "PRIVATE"
)

// GossipPacket struct
type GossipPacket struct {
	Simple *SimpleMessage
	Rumor *RumorMessage
	Status *StatusPacket
	Private *PrivateMessage
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
func NewStatusPacket(want *[]PeerStatus) *StatusPacket {
	return &StatusPacket{Want: *want}
}

// NewRumorMessage create
func NewRumorMessage(origin string, ID uint32, text string) *RumorMessage {
	return &RumorMessage{origin, ID, text}
}

// NewPrivateMessage create
func NewPrivateMessage(origin string, ID uint32, destination, text string, hops uint32) *PrivateMessage {
	return &PrivateMessage{origin, ID, destination, text, hops}
}

// GetPacketType function
func (packet *GossipPacket) GetPacketType() string {
	switch {
	case packet.Rumor != nil && packet.Status == nil && packet.Simple == nil  && packet.Private == nil:
		return PACKET_RUMOR
	case packet.Rumor == nil && packet.Status != nil && packet.Simple == nil  && packet.Private == nil:
		return PACKET_STATUS
	case packet.Rumor == nil && packet.Status == nil && packet.Simple != nil  && packet.Private == nil:
		return PACKET_SIMPLE
	case packet.Rumor == nil && packet.Status == nil && packet.Simple == nil  && packet.Private != nil:
		return PACKET_PRIVATE
	default:
		return ""
	}
}
