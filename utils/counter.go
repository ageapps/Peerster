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
	counter.value++
	return counter.value
}

//SetValue function
func (counter *Counter) SetValue(value uint32) {
	counter.mux.Lock()
	defer counter.mux.Unlock()
	counter.value = value
}
//GetValue function
func (counter *Counter) GetValue() uint32 {
	counter.mux.Lock()
	defer counter.mux.Unlock()
	return counter.value
}
