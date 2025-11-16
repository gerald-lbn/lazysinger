package handlers

import (
	"context"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/gerald-lbn/refrain/pkg/log"
	"github.com/gerald-lbn/refrain/pkg/repository"
	"github.com/gerald-lbn/refrain/pkg/services"
	"github.com/gerald-lbn/refrain/pkg/tasks"
	"github.com/gerald-lbn/refrain/pkg/utils/file"
)

// HandleCreate handles create events emitted by the file system watcher.
func HandleCreate(c *services.Container, ctx context.Context) services.FileEventHandler {
	return func(event fsnotify.Event, ctx context.Context) error {
		log.Default().Debug("create event detected",
			"operation", event.Op.String(),
			"path", event.Name)

		if isDir, err := file.IsDirectory(event.Name); err != nil {
			return err
		} else if isDir {
			return nil
		}

		if isAudio, err := file.IsAudioFile(event.Name); err != nil {
			return err
		} else if !isAudio {
			return nil
		}

		_, err := c.Tasks.Add(tasks.DownloadLyricsTask{
			Path: event.Name,
		}).Wait(5 * time.Second).Save()

		if err != nil {
			return err
		}

		return nil
	}
}

// HandleDelete handles delete events emitted by the file system watcher.
func HandleDelete(c *services.Container, ctx context.Context) services.FileEventHandler {
	return func(event fsnotify.Event, ctx context.Context) error {
		log.Default().Debug("delete event detected",
			"operation", event.Op.String(),
			"path", event.Name)

		// Since the file is deleted, there is no way to know what has been deleted except for the path
		repo := repository.New(c.Database)
		err := repo.DeleteTrack(ctx, event.Name)

		// TODO: Remove potential tasks in queue which have the deleted track as a dependency

		return err
	}
}

// HandleRename handles rename events emitted by the file system watcher.
func HandleRename(c *services.Container, ctx context.Context) services.FileEventHandler {
	return HandleDelete(c, ctx)
}

// HandleWrite
