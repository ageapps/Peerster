package router

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/ageapps/Peerster/pkg/logger"
	"github.com/ageapps/Peerster/pkg/utils"
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

// GetTable returns routing table
func (router *Router) GetTable() *RoutingTable {
	return &router.table
}

// SetEntry sets new entry
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
		router.addEntry(origin, &newEntry)
	}
	return isNew
}

func (router *Router) entryExists(origin string) (isNew bool) {
	_, ok := router.table[origin]
	return ok
}

// AddIfNotExists adds entry if there's none for the origin address
func (router *Router) AddIfNotExists(origin, address string) {
	router.mux.Lock()
	defer router.mux.Unlock()
	if !router.entryExists(origin) {
		newEntry := utils.PeerAddress{}
		err := newEntry.Set(address)
		if err != nil {
			logger.Log("Error updating router entry")
		} else {
			router.addEntry(origin, &newEntry)
		}
	}
}

func (router *Router) addEntry(origin string, entry *utils.PeerAddress) {
	logger.Log(fmt.Sprintf("Route entry appended - Origin:%v", origin))
	router.table[origin] = entry
	logger.LogDSDV(origin, entry.String())
}

// GetDestination returns de addess gibben an identifier
func (router *Router) GetDestination(origin string) (entry *utils.PeerAddress, found bool) {
	router.mux.Lock()
	defer router.mux.Unlock()

	value, ok := router.table[origin]
	if !ok || reflect.ValueOf(value).IsNil() {
		return nil, false
	}
	return value, true
}
