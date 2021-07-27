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
	events map[int] []Handle
	lock sync.RWMutex
}
func (e *BaseEvent) Register(name int, handle Handle) {
	e.lock.Lock()
	defer e.lock.Unlock()
	if e.events == nil {
		e.events = map[int][]Handle{}
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
func (e *BaseEvent) Call(name int, i ...interface{}) {
	e.lock.RLock()
	defer e.lock.RUnlock()
	callback, exist := e.events[name]
	if exist {
		for _, f := range callback {
			go f(i...)
		}
	}
}