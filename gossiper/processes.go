package gossiper

import (
	"strings"

	"github.com/ageapps/Peerster/pkg/data"
	"github.com/ageapps/Peerster/pkg/handler"
	"github.com/ageapps/Peerster/pkg/logger"
)

// ProcessType string
type ProcessType string

const (
	// PROCESS_MONGUER var
	PROCESS_MONGUER ProcessType = "PROCESS_MONGUER"
	// PROCESS_DATA var
	PROCESS_DATA ProcessType = "PROCESS_DATA"
	// PROCESS_SEARCH var
	PROCESS_SEARCH ProcessType = "PROCESS_SEARCH"
)

func (gossiper *Gossiper) registerProcess(process interface{}, ptype ProcessType) {
	go func() {

		name := ""
		switch ptype {
		case PROCESS_MONGUER:
			regProcess := process.(*handler.MongerHandler)
			name = regProcess.Name
			gossiper.mux.Lock()
			gossiper.monguerPocesses[name] = regProcess
			gossiper.mux.Unlock()

		case PROCESS_DATA:
			regProcess := process.(*handler.FileHandler)
			name = regProcess.Name
			gossiper.mux.Lock()
			gossiper.fileProcesses[name] = regProcess
			gossiper.mux.Unlock()

		case PROCESS_SEARCH:
			regProcess := process.(*handler.SearchHandler)
			name = regProcess.Name
			gossiper.mux.Lock()
			gossiper.searchProcesses[name] = regProcess
			gossiper.mux.Unlock()
		}
		logger.Logf("%v - Registering %v - %v", gossiper.Name, ptype, name)
	}()

}

func (gossiper *Gossiper) unregisterProcess(name string, ptype ProcessType) {
	go func() {
		found := true
		switch ptype {
		case PROCESS_MONGUER:
			_, found = gossiper.monguerPocesses[name]
			gossiper.mux.Lock()
			gossiper.monguerPocesses[name] = nil
			delete(gossiper.monguerPocesses, name)
			gossiper.mux.Unlock()
		case PROCESS_DATA:
			gossiper.mux.Lock()
			_, found = gossiper.fileProcesses[name]
			gossiper.fileProcesses[name] = nil
			delete(gossiper.fileProcesses, name)
			gossiper.mux.Unlock()
		case PROCESS_SEARCH:
			gossiper.mux.Lock()
			_, found = gossiper.searchProcesses[name]
			gossiper.searchProcesses[name] = nil
			delete(gossiper.searchProcesses, name)
			gossiper.mux.Unlock()
		}
		logger.Logf("%v - Unregistering %v - %v found:%v", gossiper.Name, ptype, name, found)
	}()
}

func (gossiper *Gossiper) duplicateProcess(name string, ptype ProcessType) bool {
	gossiper.mux.Lock()
	exists := false
	switch ptype {
	case PROCESS_MONGUER:
		_, exists = gossiper.monguerPocesses[name]
	case PROCESS_DATA:
		_, exists = gossiper.fileProcesses[name]

	case PROCESS_SEARCH:
		_, exists = gossiper.searchProcesses[name]
	}
	gossiper.mux.Unlock()
	return exists
}

func (gossiper *Gossiper) getMongerProcesses() map[string]*handler.MongerHandler {
	gossiper.mux.Lock()
	defer gossiper.mux.Unlock()
	return gossiper.monguerPocesses
}

func (gossiper *Gossiper) getDataProcesses() map[string]*handler.FileHandler {
	gossiper.mux.Lock()
	defer gossiper.mux.Unlock()
	return gossiper.fileProcesses
}

func (gossiper *Gossiper) getSeachProcesses() map[string]*handler.SearchHandler {
	gossiper.mux.Lock()
	defer gossiper.mux.Unlock()
	return gossiper.searchProcesses
}

func (gossiper *Gossiper) findMonguerProcess(originAddress string, routeMonguer bool) *handler.MongerHandler {
	processes := gossiper.getMongerProcesses()
	gossiper.mux.Lock()
	defer gossiper.mux.Unlock()

	for _, process := range processes {
		if process.GetMonguerPeer() == originAddress && process.IsRouteMonguer() == routeMonguer {
			return process
		}
	}
	return nil
}
func (gossiper *Gossiper) findDataProcess(origin string, hash string) *handler.FileHandler {
	processes := gossiper.getDataProcesses()
	for _, process := range processes {
		if strings.Contains(process.Name, origin) && process.GetExpectingHashStr() == hash {
			return process
		}
	}
	return nil
}
func (gossiper *Gossiper) findSearchProcess(results []*data.SearchResult) *handler.SearchHandler {

	processes := gossiper.getSeachProcesses()
	// map of filename and processes that have matching keywords
	for _, process := range processes {
		if process.MatchesResults(results) {
			return process
		}
	}
	return nil
}
