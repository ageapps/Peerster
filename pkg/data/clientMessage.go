package data

// Message to send
type Message struct {
	Text          string
	Destination   string
	FileName      string
	RequestHash   string
	IndexFilePath string
	ID            uint32
}

// IsPrivate check if is private message
func (msg *Message) IsPrivate() bool {
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
