package singleton

import (
	"reflect"
	"sync"

	"github.com/gerald-lbn/lazysinger/log"
)

var (
	instances = map[string]any{}
	pending   = map[string]chan struct{}{}
	lock      sync.RWMutex
)

func GetInstance[T any](constructor func() T) T {
	var v T
	name := reflect.TypeOf(v).String()

	// First check with read lock
	lock.RLock()
	if instance, ok := instances[name]; ok {
		defer lock.RUnlock()
		return instance.(T)
	}
	lock.RUnlock()

	// Now check if someone is already creating this type
	lock.Lock()

	// Check again with the write lock - someone might have created it
	if instance, ok := instances[name]; ok {
		lock.Unlock()
		return instance.(T)
	}

	// Check if creation is pending
	wait, isPending := pending[name]
	if !isPending {
		// We'll be the one creating it
		pending[name] = make(chan struct{})
		wait = pending[name]
	}
	lock.Unlock()

	// If someone else is creating it, wait for them
	if isPending {
		<-wait // Wait for creation to complete

		// Now it should be in the instances map
		lock.RLock()
		defer lock.RUnlock()
		return instances[name].(T)
	}

	// We're responsible for creating the instance
	newInstance := constructor()

	// Store it and signal other goroutines
	lock.Lock()
	instances[name] = newInstance
	close(wait)           // Signal that creation is complete
	delete(pending, name) // Clean up
	log.Trace().Str("name", name).Type("instance", newInstance).Msg("Created new singleton")
	lock.Unlock()

	return newInstance
}
