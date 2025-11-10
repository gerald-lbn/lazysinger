package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/fsnotify/fsnotify"
	"github.com/gerald-lbn/refrain/pkg/services"
)

func main() {
	c := services.NewContainer()
	defer c.Shutdown()

	slog.Info("application started",
		"name", c.Config.App.Name,
		"environment", c.Config.App.Environment,
		"libraries", c.Config.Libraries.Paths)

	c.Watcher.RegisterHandler(handleFileEvent)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	signal.Notify(quit, syscall.SIGHUP)
	signal.Notify(quit, syscall.SIGTERM)
	signal.Notify(quit, syscall.SIGABRT)
	<-quit
}

// handleFileEvent processes file system events
func handleFileEvent(event fsnotify.Event) error {
	slog.Info("file event detected",
		"operation", event.Op.String(),
		"path", event.Name)
	return nil
}
