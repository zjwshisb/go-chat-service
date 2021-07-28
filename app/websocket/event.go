package websocket

import "sync"

type event interface {
	Register(name int, handle Handle)
	Call(name int, i ...interface{})
}

type baseEvent struct {
	events map[int] []Handle
	lock sync.RWMutex
}
func (e *baseEvent) Register(name int, handle Handle) {
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
func (e *baseEvent) Call(name int, i ...interface{}) {
	e.lock.RLock()
	defer e.lock.RUnlock()
	callback, exist := e.events[name]
	if exist {
		for _, f := range callback {
			go f(i...)
		}
	}
}
