package services

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sync"

	"github.com/fsnotify/fsnotify"
)

// FileEventHandler is a callback function for handling file system events.
type FileEventHandler func(event fsnotify.Event, ctx context.Context) error

// WatcherService manages file system watching for multiple directories.
type WatcherService struct {
	watcher  *fsnotify.Watcher
	handlers map[fsnotify.Op]FileEventHandler
	done     chan struct{}
	wg       sync.WaitGroup
	mu       sync.RWMutex
}

// NewWatcherService creates a new WatcherService.
func NewWatcherService() (*WatcherService, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("failed to create watcher: %w", err)
	}

	ws := &WatcherService{
		watcher:  watcher,
		handlers: make(map[fsnotify.Op]FileEventHandler),
		done:     make(chan struct{}),
	}

	return ws, nil
}

// AddPath adds a path recursively to watch.
func (ws *WatcherService) AddPath(path string) error {
	ws.mu.Lock()
	defer ws.mu.Unlock()

	return ws.addRecursive(path)
}

// AddPaths adds multiple paths to watch.
func (ws *WatcherService) AddPaths(paths []string) error {
	for _, path := range paths {
		if err := ws.AddPath(path); err != nil {
			slog.Error("failed to add path to watcher",
				"path", path,
				"error", err)

			continue
		}
		slog.Info("added path to watcher", "path", path)
	}
	return nil
}

// addPath adds a single path to the watcher.
func (ws *WatcherService) addPath(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("path does not exist: %s", path)
	}

	if err := ws.watcher.Add(path); err != nil {
		return fmt.Errorf("failed to watch path %s: %w", path, err)
	}

	return nil
}

// addRecursive adds a path and all its subdirectories to the watcher.
func (ws *WatcherService) addRecursive(root string) error {
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			slog.Warn("error walking path", "path", path, "error", err)
			return nil
		}

		if info.IsDir() {
			if err := ws.addPath(path); err != nil {
				slog.Warn("failed to add directory to watcher",
					"path", path,
					"error", err)
			} else {
				slog.Debug("watching directory", "path", path)
			}
		}

		return nil
	})
}

// RegisterCreateHandler registers a callback function for create events.
// Any previously registered handler for this event will be overwritten.
func (ws *WatcherService) RegisterCreateHandler(handler FileEventHandler) {
	ws.mu.Lock()
	defer ws.mu.Unlock()
	ws.handlers[fsnotify.Create] = handler
}

// RegisterWriteHandler registers a callback function for write events.
// Any previously registered handler for this event will be overwritten.
func (ws *WatcherService) RegisterWriteHandler(handler FileEventHandler) {
	ws.mu.Lock()
	defer ws.mu.Unlock()
	ws.handlers[fsnotify.Write] = handler
}

// RegisterRenameHandler registers a callback for rename events.
// Any previously registered handler for this event will be overwritten.
func (ws *WatcherService) RegisterRenameHandler(handler FileEventHandler) {
	ws.mu.Lock()
	defer ws.mu.Unlock()
	ws.handlers[fsnotify.Rename] = handler
}

// RegisterChmodHandler registers a callback for chmod events.
// Any previously registered handler for this event will be overwritten.
func (ws *WatcherService) RegisterChmodHandler(handler FileEventHandler) {
	ws.mu.Lock()
	defer ws.mu.Unlock()
	ws.handlers[fsnotify.Chmod] = handler
}

// RegisterDeleteHandler registers a callback for delete events.
// Any previously registered handler for this event will be overwritten.
func (ws *WatcherService) RegisterDeleteHandler(handler FileEventHandler) {
	ws.mu.Lock()
	defer ws.mu.Unlock()
	ws.handlers[fsnotify.Remove] = handler
}

// Start begins watching for file system events.
func (ws *WatcherService) Start(ctx context.Context) {
	ws.wg.Add(1)
	go ws.watch(ctx)
	slog.Info("file watcher started")
}

// watch is the main event loop for processing file system events.
func (ws *WatcherService) watch(ctx context.Context) {
	defer ws.wg.Done()

	for {
		select {
		case event, ok := <-ws.watcher.Events:
			if !ok {
				return
			}

			slog.Debug("file system event",
				"event", event.Op.String(),
				"path", event.Name)

			if event.Has(fsnotify.Create) {
				if info, err := os.Stat(event.Name); err == nil && info.IsDir() {
					if err := ws.AddPath(event.Name); err != nil {
						slog.Error("failed to watch new directory",
							"path", event.Name,
							"error", err)
					}
				}
			}

			ws.mu.RLock()
			for op, handler := range ws.handlers {
				if event.Has(op) {
					if err := handler(event, ctx); err != nil {
						slog.Error("handler error",
							"event", event.Op.String(),
							"path", event.Name,
							"error", err)
					}
				}
			}
			ws.mu.RUnlock()

		case err, ok := <-ws.watcher.Errors:
			if !ok {
				return
			}
			slog.Error("file watcher error", "error", err)

		case <-ws.done:
			return
		}
	}
}

// Stop stops the watcher and waits for the event loop to finish.
func (ws *WatcherService) Stop() error {
	close(ws.done)
	ws.wg.Wait()

	if err := ws.watcher.Close(); err != nil {
		return fmt.Errorf("failed to close watcher: %w", err)
	}

	slog.Info("file watcher stopped")
	return nil
}

// GetWatchedPaths returns a list of all currently watched paths.
func (ws *WatcherService) GetWatchedPaths() []string {
	ws.mu.RLock()
	defer ws.mu.RUnlock()
	return ws.watcher.WatchList()
}
