package util

import "sync"

type Hook struct {
	Hooks map[string] []func(i ...interface{})
	Lock sync.RWMutex
}

func (hook *Hook) RegisterHook(name string, callback func(i ...interface{}))  {
	hook.Lock.Lock()
	defer hook.Lock.Unlock()
	if hook.Hooks == nil {
		hook.Hooks = make(map[string][]func(i ...interface{}))
	}
	if callbacks, ok := hook.Hooks[name]; ok {
		newCallbacks := append(callbacks, callback)
		hook.Hooks[name] = newCallbacks
	} else {
		callbacks := []func(i ...interface{}){
			callback,
		}
		hook.Hooks[name] = callbacks
	}
}

func (hook *Hook) TriggerHook(name string, payload ...interface{})  {
	hook.Lock.RLock()
	defer hook.Lock.RUnlock()
	callbacks, ok := hook.Hooks[name]
	if ok {
		for _, callback := range callbacks {
			go callback(payload...)
		}
	}
}

