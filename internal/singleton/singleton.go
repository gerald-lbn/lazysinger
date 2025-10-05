package singleton

import (
	"reflect"
	"sync"
)

var (
	instances = map[string]any{}
	pending   = map[string]chan struct{}{}
	lock      sync.RWMutex
)

func GetInstance[T any](constructor func() T) T {
	var t T
	name := reflect.TypeOf(t).String()

	// First, try to read the instance with a read lock
	lock.RLock()
	if instance, exists := instances[name]; exists {
		defer lock.RUnlock()
		return instance.(T)
	}
	lock.RUnlock()

	lock.Lock()

	// Double-check if the instance was created while acquiring the write lock
	if instance, exists := instances[name]; exists {
		lock.Unlock()
		return instance.(T)
	}

	// Check if creation is already pending
	wait, isPending := pending[name]
	if !isPending {
		pending[name] = make(chan struct{})
		wait = pending[name]
	}
	lock.Unlock()

	// If creation is pending, wait for it to complete
	if isPending {
		<-wait

		lock.RLock()
		defer lock.RUnlock()
		return instances[name].(T)
	}

	// Create the instance
	instance := constructor()

	// Store the instance and notify any waiters
	lock.Lock()
	instances[name] = instance
	close(wait)
	delete(pending, name)
	lock.Unlock()

	return instance
}
