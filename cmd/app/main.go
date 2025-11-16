package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/gerald-lbn/refrain/pkg/handlers"
	"github.com/gerald-lbn/refrain/pkg/log"
	"github.com/gerald-lbn/refrain/pkg/router"
	"github.com/gerald-lbn/refrain/pkg/services"
	"github.com/gerald-lbn/refrain/pkg/tasks"
)

func main() {
	// Start a new container.
	c := services.NewContainer()
	defer func() {
		fatal("shutdown failed", c.Shutdown())
	}()

	// Build the router.
	if err := router.BuildRouter(c); err != nil {
		fatal("failed to build the router", err)
	}

	// Register all task queues.
	tasks.Register(c)

	// Start the task runner to execute queued tasks.
	ctx := context.Background()
	c.Tasks.Start(ctx)

	// Start the server.
	go func() {
		addr := fmt.Sprintf("%s:%d", c.Config.HTTP.Hostname, c.Config.HTTP.Port)
		if c.Config.HTTP.TLS.Enabled {
			if err := c.Web.ListenTLS(addr, c.Config.HTTP.TLS.Certificate, c.Config.HTTP.TLS.Key); err != nil {
				fatal("failed to start the server", err)
			}
		} else {
			if err := c.Web.Listen(addr); err != nil {
				fatal("failed to start the server", err)
			}
		}
	}()

	slog.Info("application started",
		"name", c.Config.App.Name,
		"environment", c.Config.App.Environment,
		"libraries", c.Config.Libraries.Paths,
		"tasks", c.Config.Tasks,
	)

	c.Watcher.RegisterCreateHandler(handlers.HandleCreate(c, ctx))
	c.Watcher.RegisterRenameHandler(handlers.HandleRename(c, ctx))
	c.Watcher.RegisterDeleteHandler(handlers.HandleDelete(c, ctx))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	signal.Notify(quit, syscall.SIGHUP)
	signal.Notify(quit, syscall.SIGTERM)
	signal.Notify(quit, syscall.SIGABRT)
	<-quit
}

// fatal logs an error and terminates the application, if the error is not nil.
func fatal(msg string, err error) {
	if err != nil {
		log.Default().Error(msg, "error", err)
		os.Exit(1)
	}
}
