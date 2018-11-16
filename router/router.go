package router

import (
	"reflect"
	"sync"

	"github.com/ageapps/Peerster/logger"
	"github.com/ageapps/Peerster/utils"
)

// Router struct
type Router struct {
	table RoutingTable
	mux   sync.Mutex
}

// RoutingTable table map
type RoutingTable map[string]*utils.PeerAddress

// NewRouter func
func NewRouter() *Router {
	return &Router{
		table: make(map[string]*utils.PeerAddress),
	}
}

func (router *Router) GetTable() *RoutingTable {
	return &router.table
}

func (router *Router) SetEntry(origin, address string) bool {
	router.mux.Lock()
	defer router.mux.Unlock()
	isNew := false
	oldValue, ok := router.table[origin]

	if !ok || reflect.ValueOf(oldValue).IsNil() || oldValue.String() != address {
		isNew = true
		newEntry := utils.PeerAddress{}
		err := newEntry.Set(address)
		if err != nil {
			logger.Log("Error updating router entry")
			return false
		}
		router.table[origin] = &newEntry
		logger.LogDSDV(origin, address)
	}
	return isNew
}

func (router *Router) EntryExists(origin string) (isNew bool) {
	router.mux.Lock()
	defer router.mux.Unlock()
	_, ok := router.table[origin]
	return ok
}

func (router *Router) GetDestination(origin string) (entry *utils.PeerAddress, found bool) {
	router.mux.Lock()
	defer router.mux.Unlock()

	value, ok := router.table[origin]
	if !ok || reflect.ValueOf(value).IsNil() {
		return nil, false
	}
	return value, true
}
