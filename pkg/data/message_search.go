package data

import (
	"github.com/ageapps/Peerster/pkg/utils"
)

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
	MetafileHash utils.HashValue
	ChunkMap     []uint64
	ChunkCount   uint64
}

// FileResult struct
type FileResult struct {
	FileName     string
	Destination  string
	MetafileHash utils.HashValue
}

// NewSearchRequest create
func NewSearchRequest(ogname string, budget uint64, keywords []string) *SearchRequest {
	return &SearchRequest{
		Origin:   ogname,
		Budget:   budget,
		Keywords: keywords,
	}
}

// NewFileResult create
func NewFileResult(filename, destination string, metahash utils.HashValue) *FileResult {
	return &FileResult{
		FileName:     filename,
		Destination:  destination,
		MetafileHash: metahash,
	}
}

// NewSearchResult create
func NewSearchResult(file string, metafileHash utils.HashValue, chunkMap []uint64, chunkCount uint64) *SearchResult {
	return &SearchResult{
		FileName:     file,
		MetafileHash: metafileHash,
		ChunkMap:     chunkMap,
		ChunkCount:   chunkCount,
	}
}

// NewSearchReply create
func NewSearchReply(ogname, destination string, hops uint32, results []*SearchResult) *SearchReply {
	return &SearchReply{
		Origin:      ogname,
		Destination: destination,
		HopLimit:    hops,
		Results:     results,
	}
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
