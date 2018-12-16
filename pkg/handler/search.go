package handler

import (
	"math"
	"strings"
	"sync"
	"time"

	"github.com/ageapps/Peerster/pkg/logger"

	"github.com/ageapps/Peerster/pkg/router"

	"github.com/ageapps/Peerster/pkg/data"
)

// MaxBudget achieved in search
const MaxBudget = 32

// DefaultBudget used
const DefaultBudget = uint64(2)

// DefaultMatchThreshold needed to Stop searching
const DefaultMatchThreshold = 2

// SearchHandler is a handler that will be in
// charge of requesting data from other peers
// FileName            string
// MetaHash            data.HashValue
// active              bool
// connection          *ConnectionHandler
// router              *router.Router
// mux                 sync.Mutex
// quitChannel         chan bool
// resetChannel        chan bool
//
type SearchHandler struct {
	Name                string
	budget              uint64
	originPeer          string
	matchThreshold      int
	fromClient          bool
	Keywords            []string
	destinationBudgets  map[string]int
	matchedDestinations map[string]*data.FileResult
	stopped             bool
	waiting             bool
	connection          *ConnectionHandler
	router              *router.Router
	mux                 sync.Mutex
	timer               *time.Timer
	quitChannel         chan bool
	ReplyChannel        chan *data.SearchReply
}

// NewSearchHandler function
func NewSearchHandler(name string, budget uint64, fromClient bool, originPeer string, Keywords []string, peerConection *ConnectionHandler, router *router.Router) *SearchHandler {
	defaultBudget := DefaultBudget
	if budget > 0 {
		defaultBudget = budget
	}
	matchThreshold := DefaultMatchThreshold
	handler := &SearchHandler{
		Name:                name,
		budget:              defaultBudget,
		originPeer:          originPeer,
		matchThreshold:      matchThreshold,
		fromClient:          fromClient,
		Keywords:            Keywords,
		destinationBudgets:  make(map[string]int),
		matchedDestinations: make(map[string]*data.FileResult),
		stopped:             false,
		waiting:             false,
		connection:          peerConection,
		router:              router,
		timer:               &time.Timer{},
		quitChannel:         make(chan bool),
		ReplyChannel:        make(chan *data.SearchReply),
	}
	return handler
}

// Start handler
func (handler *SearchHandler) Start(onFileReceviedHandler func(*data.FileResult), onStopHandler func()) {
	go func() {
		handler.resetTimer()
		handler.sendRequest()
		for {
			select {
			case reply := <-handler.ReplyChannel:
				for _, result := range reply.Results {
					fileResult := data.NewFileResult(result.FileName, reply.Origin, result.MetafileHash, result.ChunkMap)
					handler.matchedDestinations[reply.Origin] = fileResult
					onFileReceviedHandler(fileResult)
				}
				nrMatches := len(handler.matchedDestinations)

				if nrMatches >= handler.matchThreshold {
					logger.Logf("Found - %v matches, stopping", nrMatches)
					handler.Stop()
				}
			case <-handler.timer.C:
				logger.Logf("TIMEOUT - %v matches", len(handler.matchedDestinations))
				if handler.waiting {
					logger.Logf("Handler waiting for matches %v", handler.Keywords)
				} else {
					if handler.budget > MaxBudget {
						handler.waiting = true
						// create a timeout for waiting for answers
						timer2 := time.NewTimer(10 * time.Second)
						go func() {
							<-timer2.C
							handler.Stop()
						}()
					} else {
						handler.sendRequest()
						handler.resetTimer()
					}
				}
			case <-handler.quitChannel:
				logger.Log("Finishing search handler - " + handler.Name)
				if handler.timer.C != nil {
					handler.timer.Stop()
				}
				onStopHandler()
				return
			}
		}
	}()
}

func (handler *SearchHandler) updateMatchThreshHold() {
	if handler.router.GetTableSize() < handler.matchThreshold {
		handler.matchThreshold = handler.router.GetTableSize()
		return
	}
	handler.matchThreshold = DefaultMatchThreshold
}

func (handler *SearchHandler) loadDestinationBudgets() {
	nrPeers := handler.router.GetTableSize()
	if int(handler.budget) < nrPeers {
		nrPeers = int(handler.budget)
	}
	handler.destinationBudgets = make(map[string]int)
	if handler.originPeer != "" && !handler.fromClient {
		nrPeers--
		handler.destinationBudgets[handler.originPeer] = 0
	}
	logger.Logf("Loading destination budgets %v to %v peers", handler.budget, nrPeers)
	if nrPeers <= 0 {
		logger.Log("No peers to send search request")
		handler.Stop()
		return
	}
	peersChosen := nrPeers
	initialBudget := int(handler.budget)
	for peersChosen > 0 {
		logger.Logf("Assigning budget: %v/%v", float64(initialBudget), float64(peersChosen))
		assignedBudget := int(math.Ceil(float64(initialBudget) / float64(peersChosen)))
		destination := handler.router.GetRandomDestination(handler.destinationBudgets)
		logger.Logf("Destination: %v assigned budget: %v", destination, assignedBudget)
		handler.destinationBudgets[destination] = assignedBudget
		peersChosen--
		initialBudget -= assignedBudget
	}
}

func (handler *SearchHandler) resetTimer() {
	//logger.Log("Launching new timer")
	if handler.getTimer().C != nil {
		handler.getTimer().Stop()
	}
	handler.mux.Lock()
	handler.timer = time.NewTimer(1 * time.Second)
	handler.mux.Unlock()
}

func (handler *SearchHandler) getTimer() *time.Timer {
	handler.mux.Lock()
	defer handler.mux.Unlock()
	return handler.timer
}

func (handler *SearchHandler) sendRequest() {
	handler.loadDestinationBudgets()
	handler.updateMatchThreshHold()
	for destination, budget := range handler.destinationBudgets {
		if budget > 0 {
			msg := data.NewSearchRequest(handler.originPeer, uint64(budget), handler.Keywords)
			packet := &data.GossipPacket{SearchRequest: msg}
			if destinationAdress, ok := handler.router.GetDestination(destination); ok {
				logger.Logf("Sending SEARCH REQUEST Dest:%v", destination)
				handler.connection.SendPacketToPeer(destinationAdress.String(), packet)
			} else {
				logger.Logf("INVALID SEARCH REQUEST Dest:%v", destination)
			}
		}
	}
	if handler.budget <= MaxBudget {
		handler.budget = handler.budget * 2
	}
}

// MatchesResults function checks if the results match keywords in prcess
func (handler *SearchHandler) MatchesResults(results []*data.SearchResult) bool {
	matchFile := true
	for _, result := range results {
		fileName := result.FileName
		matchedKeyword := false
		for _, key := range handler.Keywords {
			if strings.Contains(fileName, key) {
				matchedKeyword = true
				break
			}
		}
		logger.Logf("Matched %v/%v - %v\n", fileName, matchedKeyword, handler.Keywords)
		matchFile = matchFile && matchedKeyword
	}
	return matchFile
}

// Stop search handler
func (handler *SearchHandler) Stop() {
	logger.Logf("Stopping search handler - " + handler.Name)
	if !handler.stopped {
		handler.stopped = true
		close(handler.quitChannel)
		return
	}
	logger.Log("Data Handler already stopped....")
}
