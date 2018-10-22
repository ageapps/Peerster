package logger

import (
	"fmt"

	"github.com/ageapps/Peerster/data"
)

var l = Logger{"", "", false}

type Logger struct {
	address string
	name    string
	debug   bool
}

func LogRumor(msg data.RumorMessage, from string) {
	fmt.Printf("RUMOR origin %v from %v ID %v contents %v\n", msg.Origin, from, msg.ID, msg.Text)
}

func LogStatus(msg data.PeerStatus, from string) {
	fmt.Printf("STATUS from %v peer %v nextID %v\n", from, msg.Identifier, msg.NextID)
}

func LogSimple(msg data.SimpleMessage) {
	fmt.Printf("SIMPLE MESSAGE origin %v from %v contents %v\n", msg.OriginalName, msg.RelayPeerAddr, msg.Contents)
}
func LogPeers(peers string) {
	fmt.Println("PEERS " + peers)
}
func LogInSync(peer string) {
	fmt.Println("IN SYNC WITH " + peer)
}

func LogClient(msg data.Message) {
	fmt.Printf("CLIENT MESSAGE %v\n", msg.Text)
}
func LogCoin(address string) {
	fmt.Printf("FLIPPED COIN sending rumor to %v\n", address)
}
func LogMonguer(address string) {
	fmt.Printf("MONGERING with %v \n", address)
}

func CreateLogger(name, address string, debug bool) {
	l.name = name
	l.address = address
	l.debug = debug
	if l.debug {
		fmt.Printf("*******  Logger Created  ******\nNAME: %v\nADRESS: %v\n", l.name, l.address)
		fmt.Printf("*******  **************  ******\n")
	}
}

func Log(text string) {
	if l.debug {
		fmt.Println("[" + l.name + "/" + l.address + "]: " + text)
	}
}
