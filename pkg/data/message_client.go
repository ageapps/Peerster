package data

// Message to send
type Message struct {
	Text          string
	Destination   string
	FileName      string
	RequestHash   string
	IndexFilePath string
	ID            uint32
	Keywords      []string
	Budget        uint64
}

// IsDirectMessage check if is private message
func (msg *Message) IsDirectMessage() bool {
	return msg.Destination != ""
}

// FileToIndex check if is private message
func (msg *Message) FileToIndex() bool {
	return msg.IndexFilePath != ""
}

// HasRequest check if is private message
func (msg *Message) HasRequest() bool {
	return msg.RequestHash != ""
}

// IsSearchMessage check if is private message
func (msg *Message) IsSearchMessage() bool {
	return msg.Keywords != nil && len(msg.Keywords) > 0
}
