package data

// SearchRequest struct
type SearchRequest struct {
	Origin   string
	Budget   uint64
	Keywords []string
}

// SearchReply struct
type SearchReply struct {
	Origin      string
	Destination string
	HopLimit    uint32
	Results     []*SearchResult
}

// SearchResult struct
type SearchResult struct {
	FileName     string
	MetafileHash []byte
	ChunkMap     []uint64
	ChunkCount   uint64
}

// IsDuplicate func
func (req *SearchRequest) IsDuplicate(r *SearchRequest) bool {
	sameKeywords := true
	if len(req.Keywords) != len(r.Keywords) {
		sameKeywords = false
	}
	for i, v := range req.Keywords {
		if v != r.Keywords[i] {
			sameKeywords = false
		}
	}
	return sameKeywords && req.Origin == r.Origin
}
