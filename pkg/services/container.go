package services

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/gerald-lbn/refrain/config"
	"github.com/gerald-lbn/refrain/pkg/music/lrclib"
)

// Container contains all services used by the application and provides an easy way to handle dependency
// injection including within tests.
type Container struct {
	// Config stores the application configuration.
	Config *config.Config

	// Database stores the connection to the databas
	Database *sql.DB

	// Watcher is the file watcher service.
	Watcher *WatcherService

	// Worker is the service responsible for background tasks.
	Worker *WorkerService

	LyricsProvider *lrclib.LRCLibProvider
}

// NewContainer creates and initializes a new Container.
func NewContainer() *Container {
	c := new(Container)
	c.initConfig()
	c.initDatabase()
	c.initLyricsProvider()
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
		if err := c.Worker.Stop(); err != nil {
			return fmt.Errorf("failed to stop worker: %w", err)
		}
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

func (c *Container) initDatabase() {
	var err error
	var connection string

	switch c.Config.App.Environment {
	case config.EnvTest:
		connection = c.Config.Database.TestConnection
	default:
		connection = c.Config.Database.Connection
	}

	c.Database, err = openDB(c.Config.Database.Driver, connection)
	if err != nil {
		panic(err)
	}
}

// initLyrics providers initializes the lyrics provider
func (c *Container) initLyricsProvider() {
	c.LyricsProvider = lrclib.NewLRCLibProvider()
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
		if err := c.Watcher.AddPaths(c.Config.Libraries.Paths); err != nil {
			slog.Warn("failed to add some library paths to watcher", "error", err)
		}
	}

	// Start watching
	c.Watcher.Start()
}

// initWorker initializes the background worker service.
func (c *Container) initWorker() {
	worker := NewWorkerService(c.Config.Redis.Addr, c.Config.Worker.Concurrency)

	c.Worker = worker

	// Start the worker service
	c.Worker.Start()
}

// openDB opens a database connection.
func openDB(driver, connection string) (*sql.DB, error) {
	if driver == "sqlite3" {
		d := strings.Split(connection, "/")
		if len(d) > 1 {
			dirpath := strings.Join(d[:len(d)-1], "/")

			if err := os.MkdirAll(dirpath, 0755); err != nil {
				return nil, err
			}
		}
	}

	return sql.Open(driver, connection)
}
