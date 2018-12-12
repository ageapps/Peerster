package data

import (
	"reflect"
)

const (
	// PACKET_SIMPLE type
	PACKET_SIMPLE = "SIMPLE"
	// PACKET_RUMOR type
	PACKET_RUMOR = "RUMOR"
	// PACKET_STATUS type
	PACKET_STATUS = "STATUS"
	// PACKET_PRIVATE type
	PACKET_PRIVATE = "PRIVATE"
	// PACKET_DATA_REQUEST type
	PACKET_DATA_REQUEST = "DATA_REQUEST"
	// PACKET_DATA_REPLY type
	PACKET_DATA_REPLY = "DATA_REPLY"
	// PACKET_SEARCH_REQUEST type
	PACKET_SEARCH_REQUEST = "SEARCH_REQUEST"
	// PACKET_SEARCH_REPLY type
	PACKET_SEARCH_REPLY = "SEARCH_REPLY"
	// PACKET_SEARCH_RESULT type
	PACKET_SEARCH_RESULT = "SEARCH_RESULT"
	// PACKET_TX_PUBLISH type
	PACKET_TX_PUBLISH = "TX_PUBLISH"
	// PACKET_BLOCK_PUBLISH type
	PACKET_BLOCK_PUBLISH = "BLOCK_PUBLISH"
)

// GossipPacket struct
type GossipPacket struct {
	Simple        *SimpleMessage
	Rumor         *RumorMessage
	Status        *StatusPacket
	Private       *PrivateMessage
	DataRequest   *DataRequest
	DataReply     *DataReply
	SearchRequest *SearchRequest
	SearchReply   *SearchReply
	TxPublish     *TxPublish
	BlockPublish  *BlockPublish
}

// GetPacketType function
func (packet *GossipPacket) GetPacketType() string {
	types := []string{
		PACKET_SIMPLE,
		PACKET_RUMOR,
		PACKET_STATUS,
		PACKET_PRIVATE,
		PACKET_DATA_REPLY,
		PACKET_DATA_REQUEST,
		PACKET_SEARCH_REPLY,
		PACKET_SEARCH_REQUEST,
		PACKET_TX_PUBLISH,
		PACKET_BLOCK_PUBLISH,
	}
	var values []interface{}
	values = append(values, packet.Simple)
	values = append(values, packet.Rumor)
	values = append(values, packet.Status)
	values = append(values, packet.Private)
	values = append(values, packet.DataReply)
	values = append(values, packet.DataRequest)
	values = append(values, packet.SearchReply)
	values = append(values, packet.SearchRequest)
	values = append(values, packet.TxPublish)
	values = append(values, packet.BlockPublish)

	notNull := -1

	for index := 0; index < len(values); index++ {
		if !reflect.ValueOf(values[index]).IsNil() {
			//fmt.Printf("YYYYYY %v - %v - %v\n", types[index], notNull, values[index])
			// 2 or more properties where != null
			if notNull >= 0 {
				return ""
			}
			notNull = index
		}
	}
	if notNull >= 0 {
		return types[notNull]
	}
	return ""
	// switch {
	// case packet.Rumor != nil && packet.Status == nil && packet.Simple == nil && packet.Private == nil && packet.DataReply == nil && packet.DataRequest == nil:
	// 	return PACKET_RUMOR
	// case packet.Rumor == nil && packet.Status != nil && packet.Simple == nil && packet.Private == nil && packet.DataReply == nil && packet.DataRequest == nil:
	// 	return PACKET_STATUS
	// case packet.Rumor == nil && packet.Status == nil && packet.Simple != nil && packet.Private == nil && packet.DataReply == nil && packet.DataRequest == nil:
	// 	return PACKET_SIMPLE
	// case packet.Rumor == nil && packet.Status == nil && packet.Simple == nil && packet.Private != nil && packet.DataReply == nil && packet.DataRequest == nil:
	// 	return PACKET_PRIVATE
	// case packet.Rumor == nil && packet.Status == nil && packet.Simple == nil && packet.Private == nil && packet.DataReply != nil && packet.DataRequest == nil:
	// 	return PACKET_DATA_REPLY
	// case packet.Rumor == nil && packet.Status == nil && packet.Simple == nil && packet.Private == nil && packet.DataReply == nil && packet.DataRequest != nil:
	// 	return PACKET_DATA_REQUEST
	// default:
	// 	return ""
	// }
}
