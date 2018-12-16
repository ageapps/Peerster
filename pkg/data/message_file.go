package data

import "github.com/ageapps/Peerster/pkg/utils"

// DataRequest struct
type DataRequest struct {
	Origin      string
	Destination string
	HopLimit    uint32
	HashValue   utils.HashValue
}

// DataReply struct
type DataReply struct {
	Origin      string
	Destination string
	HopLimit    uint32
	HashValue   utils.HashValue
	Data        []byte
}

// NewDataRequest create
func NewDataRequest(origin, destination string, hops uint32, hash utils.HashValue) *DataRequest {
	return &DataRequest{origin, destination, hops, hash}
}

// NewDataReply create
func NewDataReply(origin, destination string, hops uint32, hash utils.HashValue, data []byte) *DataReply {
	return &DataReply{origin, destination, hops, hash, data}
}
