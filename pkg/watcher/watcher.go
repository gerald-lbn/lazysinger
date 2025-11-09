package watcher

import (
	"context"
	"io/fs"
	"log"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/gerald-lbn/refrain/pkg/watcher/handlers"
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

func (w *fs_watcher) Start(ctx context.Context, dir string) error {
	log.Printf("Starting fs watcher for directory: %s", dir)

	defer w.watcher.Close()

	// Start listening for events.
	go func() {
		for {
			select {
			case event, ok := <-w.watcher.Events:
				if !ok {
					return
				}

				log.Printf("Event: %s, File: %s\n", event.Op.String(), event.Name)

			case err, ok := <-w.watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)

			case <-w.done:
				return

			case <-ctx.Done():
				return
			}
		}
	}()

	err := w.watcher.Add(dir)
	if err != nil {
		return err
	}

	err = filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			w.watcher.Add(path)
		} else {
			handlers.OnInitialScan(path)
		}

		return nil
	})

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
