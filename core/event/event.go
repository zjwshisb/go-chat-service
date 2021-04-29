package event

import (
	"sync"
)

type Handle func(i ...interface{})

type Event interface {
	register(name string, handle Handle)
	call(name string, i ...interface{})
}

type BaseEvent struct {
	events map[string] []Handle
	lock sync.RWMutex
}
func (e *BaseEvent) Register(name string, handle Handle) {
	e.lock.Lock()
	defer e.lock.Unlock()
	if e.events == nil {
		e.events = map[string][]Handle{}
	}
	callback, exist := e.events[name]
	if exist {
		callback = append(callback, handle)
	} else {
		callback = []Handle{
			handle,
		}
	}
	e.events[name] = callback
}
func (e *BaseEvent) Call(name string, i ...interface{}) {
	e.lock.RLock()
	defer e.lock.RUnlock()
	callback, exist := e.events[name]
	if exist {
		for _, f := range callback {
			f(i...)
		}
	}
}