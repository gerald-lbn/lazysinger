package services

import (
	"fmt"
	"log/slog"

	"github.com/gerald-lbn/refrain/config"
	"github.com/hibiken/asynq"
)

// Container contains all services used by the application and provides an easy way to handle dependency
// injection including within tests.
type Container struct {
	// Config stores the application configuration.
	Config *config.Config

	// Watcher is the file watcher service.
	Watcher *WatcherService

	// WorkerServer stores a worker server
	Worker *WorkerService
}

// NewContainer creates and initializes a new Container.
func NewContainer() *Container {
	c := new(Container)
	c.initConfig()
	c.initWorker()
	c.initWatcher()
	return c
}

// Shutdown gracefully shuts the Container down
func (c *Container) Shutdown() error {
	if c.Watcher != nil {
		if err := c.Watcher.Stop(); err != nil {
			return fmt.Errorf("failed to stop watcher: %w", err)
		}
	}

	if c.Worker != nil {
		c.Worker.server.Shutdown()
	}

	return nil
}

// initConfig initializes configuration.
func (c *Container) initConfig() {
	cfg, err := config.GetConfig()
	if err != nil {
		panic(fmt.Sprintf("failed to load config: %v", err))
	}
	c.Config = &cfg

	// Configure logging.
	switch cfg.App.Environment {
	case config.EnvProduction:
		slog.SetLogLoggerLevel(slog.LevelInfo)
	default:
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}
}

// initWatcher initializes the file watcher service.
func (c *Container) initWatcher() {
	watcher, err := NewWatcherService()
	if err != nil {
		panic(fmt.Sprintf("failed to create watcher service: %v", err))
	}

	c.Watcher = watcher

	// Add configured library paths
	if len(c.Config.Libraries.Paths) > 0 {
		if err := c.Watcher.AddPaths(c.Config.Libraries.Paths, true); err != nil {
			slog.Warn("failed to add some library paths to watcher", "error", err)
		}
	}

	// Start watching
	c.Watcher.Start()
}

func (c *Container) initWorker() {
	worker, err := NewWorkerService()
	if err != nil {
		panic(fmt.Sprintf("failed to create worker service: %v", err))
	}

	c.Worker = worker

	worker.server = asynq.NewServer(
		asynq.RedisClientOpt{
			Addr: c.Config.Redis.Addr,
		},
		asynq.Config{
			Concurrency: c.Config.Worker.Concurrency,
			Queues: map[string]int{
				DownloadLyricsQueueName: 10,
			},
		})

	// Start server
	c.Worker.Start()

}
