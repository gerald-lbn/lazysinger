package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/fsnotify/fsnotify"
	"github.com/gerald-lbn/refrain/pkg/services"
	"github.com/gerald-lbn/refrain/pkg/tasks"
	"github.com/gerald-lbn/refrain/pkg/utils/file"
)

func main() {
	c := services.NewContainer()
	defer c.Shutdown()

	slog.Info("application started",
		"name", c.Config.App.Name,
		"environment", c.Config.App.Environment,
		"libraries", c.Config.Libraries.Paths,
		"redis", c.Config.Redis.Addr,
		"worker", c.Config.Worker.Concurrency,
	)

	c.Watcher.RegisterCreateHandler(func(event fsnotify.Event) error {
		slog.Info("create event detected",
			"operation", event.Op.String(),
			"path", event.Name)

		ok, err := file.IsAudioFile(event.Name)
		if err != nil {
			return err
		}

		if !ok {
			return nil
		}

		task, err := tasks.NewDownloadLyricsTask(event.Name)
		if err != nil {
			return err
		}

		info, err := c.Worker.EnqueueTask(task)
		if err != nil {
			return err
		}

		slog.Info("enqueued download lyrics task",
			"task_id", info.ID,
			"file_path", event.Name,
		)

		return nil
	})

	c.Watcher.RegisterWriteHandler(func(event fsnotify.Event) error {
		slog.Info("write event detected",
			"operation", event.Op.String(),
			"path", event.Name)
		return nil
	})

	c.Watcher.RegisterDeleteHandler(func(event fsnotify.Event) error {
		slog.Info("delete event detected",
			"operation", event.Op.String(),
			"path", event.Name,
		)
		return nil
	})

	c.Watcher.RegisterRenameHandler(func(event fsnotify.Event) error {
		slog.Info("rename event detected",
			"operation", event.Op.String(),
			"path", event.Name)
		return nil
	})

	c.Worker.RegisterHandler(tasks.TypeDownloadLyrics, tasks.HandleDownloadLyricsTask)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	signal.Notify(quit, syscall.SIGHUP)
	signal.Notify(quit, syscall.SIGTERM)
	signal.Notify(quit, syscall.SIGABRT)
	<-quit
}
