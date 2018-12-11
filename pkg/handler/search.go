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

// MatchThreshold needed to stop searching
const MatchThreshold = 2

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
	origin              string
	Keywords            []string
	destinationBudgets  map[string]int
	matchedDestinations map[string]*data.FileResult
	stopped             bool
	connection          *ConnectionHandler
	router              *router.Router
	mux                 sync.Mutex
	timer               *time.Timer
	quitChannel         chan bool
	ReplyChannel        chan *data.SearchReply
}

// NewSearchHandler function
func NewSearchHandler(name string, budget uint64, origin string, Keywords []string, peerConection *ConnectionHandler, router *router.Router) *SearchHandler {
	defaultBudget := DefaultBudget
	if budget > 0 {
		defaultBudget = budget
	}
	handler := &SearchHandler{
		Name:                name,
		budget:              defaultBudget,
		origin:              origin,
		Keywords:            Keywords,
		destinationBudgets:  make(map[string]int),
		matchedDestinations: make(map[string]*data.FileResult),
		stopped:             false,
		connection:          peerConection,
		router:              router,
		timer:               &time.Timer{},
		quitChannel:         make(chan bool),
		ReplyChannel:        make(chan *data.SearchReply),
	}
	handler.loadDestinationBudgets()
	return handler
}

// Start handler
func (handler *SearchHandler) Start(onStopHandler func(map[string]*data.FileResult)) {
	go func() {
		go handler.resetTimer()
		handler.sendRequest()
		for {
			select {
			case reply := <-handler.ReplyChannel:
				for _, result := range reply.Results {
					handler.matchedDestinations[result.FileName] = data.NewFileResult(result.FileName, reply.Origin, result.MetafileHash)
				}
				nrMatches := len(handler.matchedDestinations)
				if nrMatches >= MatchThreshold {
					handler.Stop()
				}
			case <-handler.timer.C:
				logger.Logf("TIMEOUT - %v matches", len(handler.matchedDestinations))
				handler.sendRequest()
				handler.resetTimer()
			case <-handler.quitChannel:
				logger.Log("Finishing search handler - " + handler.Name)
				if handler.timer.C != nil {
					handler.timer.Stop()
				}
				onStopHandler(handler.matchedDestinations)
				return
			}
		}
	}()
}

func (handler *SearchHandler) loadDestinationBudgets() {
	nrPeers := handler.router.GetTableSize()
	if int(handler.budget) < nrPeers {
		nrPeers = int(handler.budget)
	}
	nrPeers = nrPeers - len(handler.destinationBudgets)
	logger.Logf("Loading destination budgets to %v peers", nrPeers)
	peersChosen := 0
	initialBudget := int(handler.budget)

	for peersChosen < nrPeers {
		assignedBudget := int(math.Ceil(float64(initialBudget) / float64(nrPeers)))
		destination := handler.router.GetRandomDestination(handler.destinationBudgets)
		logger.Logf("Destination: %v assigned budget: %v", destination, assignedBudget)
		handler.destinationBudgets[destination] = assignedBudget
		peersChosen++
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
	for destination, budget := range handler.destinationBudgets {
		if budget > 0 {
			msg := data.NewSearchRequest(handler.origin, uint64(budget), handler.Keywords)
			packet := &data.GossipPacket{SearchRequest: msg}
			if destinationAdress, ok := handler.router.GetDestination(destination); ok {
				logger.Logf("Sending SEARCH REQUEST Dest:%v", destination)
				handler.connection.SendPacketToPeer(destinationAdress.String(), packet)
			} else {
				logger.Logf("INVALID SEARCH REQUEST Dest:%v", destination)
			}
		}
	}
	if handler.budget < MaxBudget {
		handler.budget = handler.budget * 2
		handler.loadDestinationBudgets()
	} else {
		// stop searching when max budget achieved
		handler.Stop()
	}
}

// MatchesResults funtion checks if the results match keywords in prcess
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
		matchFile = matchFile && matchedKeyword
	}
	return matchFile
}

// Stop search handler
func (handler *SearchHandler) Stop() {
	logger.Logf("Stopping search handler - " + handler.Name)
	handler.stopped = true
	close(handler.quitChannel)
}
