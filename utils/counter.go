package utils

import (
	"sync"
)

//Counter struct
type Counter struct {
	value uint32
	mux   sync.Mutex
}

func NewCounter(initialValue uint32) *Counter {
	return &Counter{
		value: initialValue,
	}
}

//Increment function
func (counter *Counter) Increment() uint32 {
	counter.mux.Lock()
	defer counter.mux.Unlock()
	counter.value = +1
	return counter.value
}
