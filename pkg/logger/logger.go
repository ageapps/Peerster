package logger

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/ageapps/Peerster/pkg/data"
)

// Logger struct
type Logger struct {
	address string
	name    string
	debug   bool
	log     *log.Logger
}

var instance = Logger{
	address: "",
	name:    "",
	debug:   false,
	log:     log.New(os.Stdout, "LOG: ", log.Ltime),
}

// LogRumor func
func LogRumor(msg data.RumorMessage, from string) {
	fmt.Printf("RUMOR origin %v from %v ID %v contents %v\n", msg.Origin, from, msg.ID, msg.Text)
}

// LogStatus func
func LogStatus(msg data.StatusPacket, from string) {
	logStr := ""
	for _, status := range msg.Want {
		logStr += fmt.Sprintf("peer %v nextID %v ", status.Identifier, status.NextID)
	}
	fmt.Printf("STATUS from %v %v\n", from, logStr)
}

// LogSimple func
func LogSimple(msg data.SimpleMessage) {
	fmt.Printf("SIMPLE MESSAGE origin %v from %v contents %v\n", msg.OriginalName, msg.RelayPeerAddr, msg.Contents)
}

// LogPeers func
func LogPeers(peers string) {
	fmt.Println("PEERS " + peers)
}

// LogInSync func
func LogInSync(peer string) {
	fmt.Println("IN SYNC WITH " + peer)
}

// LogClient func
func LogClient(msg data.Message) {
	fmt.Printf("CLIENT MESSAGE %v\n", msg.Text)
}

// LogCoin func
func LogCoin(address string) {
	fmt.Printf("FLIPPED COIN sending rumor to %v\n", address)
}

// LogMonguer func
func LogMonguer(address string) {
	fmt.Printf("MONGERING with %v \n", address)
}

// LogDSDV func
func LogDSDV(origin, address string) {
	fmt.Printf("DSDV %v %v\n", origin, address)
}

// LogPrivate func
func LogPrivate(msg data.PrivateMessage) {
	fmt.Printf("PRIVATE origin %v hop-limit %v contents %v \n", msg.Origin, msg.HopLimit, msg.Text)
}

// LogMetafile func
func LogMetafile(filename, peer string) {
	fmt.Printf("DOWNLOADING metafile of %v from %v \n", filename, peer)
}

// LogChunk func
func LogChunk(filename, peer string, chunk int) {
	fmt.Printf("DOWNLOADING %v chunk %v from %v \n", filename, chunk, peer)
}

// LogReconstructed func
func LogReconstructed(filename string) {
	fmt.Printf("RECONSTRUCTED file %v \n", filename)
}

// LogFound func
func LogFound(filename, origin, metafile string, chunks []uint64) {
	indexes := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(chunks)), ","), "[]")
	fmt.Printf("FPUND match %v at %v metafile=%v chunks=%v\n", filename, origin, metafile, indexes)
}

// CreateLogger func
func CreateLogger(name, address string, debug bool) {
	instance.name = name
	instance.address = address
	instance.debug = debug

	instance.log = log.New(os.Stdout, "["+name+"/"+address+"]: ", log.Ltime)
	if debug {
		fmt.Printf("*******  Logger Created  ******\nNAME: %v\nADRESS: %v\n", name, address)
		fmt.Printf("*******  **************  ******\n")
	}
}

// Log func
func Log(text string) {
	if instance.debug {
		instance.log.Println(text)
	}
}

// Logf func
func Logf(format string, v ...interface{}) {
	if instance.debug {
		instance.log.Printf(format, v...)
	}
}
