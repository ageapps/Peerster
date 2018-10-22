package data

const (
	PACKET_SIMPLE = "SIMPLE"
	PACKET_RUMOR  = "RUMOR"
	PACKET_STATUS = "STATUS"
)

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
func NewStatusPacket(want *[]PeerStatus) *StatusPacket {
	return &StatusPacket{Want: *want}
}

// NewRumorMessage create
func NewRumorMessage(origin string, ID uint32, text string) *RumorMessage {
	return &RumorMessage{origin, ID, text}
}

// GetPacketType function
func (packet *GossipPacket) GetPacketType() string {
	switch {
	case packet.Rumor != nil && packet.Status == nil && packet.Simple == nil:
		return PACKET_RUMOR
	case packet.Rumor == nil && packet.Status != nil && packet.Simple == nil:
		return PACKET_STATUS
	case packet.Rumor == nil && packet.Status == nil && packet.Simple != nil:
		return PACKET_SIMPLE
	default:
		return ""
	}
}
