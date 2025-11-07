package watcher

import (
	"context"
	"log"

	"github.com/fsnotify/fsnotify"
)

type fs_watcher struct {
	watcher *fsnotify.Watcher
	done    chan struct{}
}

type Watcher interface {
	Start(ctx context.Context, path string) error
	Stop() error
}

// NewWatcher
func NewWatcher() (Watcher, error) {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	return &fs_watcher{watcher: w, done: make(chan struct{})}, nil
}

func (w *fs_watcher) Start(ctx context.Context, path string) error {
	log.Printf("Starting fs watcher for directory: %s", path)

	defer w.watcher.Close()

	// Start listening for events.
	go func() {
		for {
			select {
			case event, ok := <-w.watcher.Events:
				if !ok {
					return
				}
				log.Println("event:", event)
				if event.Has(fsnotify.Write) {
					log.Println("modified file:", event.Name)
				}
			case err, ok := <-w.watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	err := w.watcher.Add(path)
	if err != nil {
		return err
	}

	select {
	case <-ctx.Done():
		return nil
	case <-w.done:
		return nil
	}
}

// Stop removes all watches
func (w *fs_watcher) Stop() error {
	log.Println("Stopping watcher...")
	close(w.done)
	return w.watcher.Close()
}
