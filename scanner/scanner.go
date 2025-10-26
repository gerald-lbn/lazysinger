package scanner

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"path/filepath"

	"github.com/gerald-lbn/lazysinger/config"
	"github.com/gerald-lbn/lazysinger/log"
	"github.com/gerald-lbn/lazysinger/music"
	"github.com/gerald-lbn/lazysinger/worker"
	"github.com/hibiken/asynq"
)

func ScanAll(ctx context.Context) error {
	asynqClient := asynq.NewClient(asynq.RedisClientOpt{
		Addr: config.Server.Worker.RedisAddr,
	})

	return filepath.WalkDir(config.Server.Scanner.Directory, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Error().Err(err).Str("path", path).Msg("Unable to walk down the file tree")
		}

		if !d.IsDir() && music.IsMusicFile(path) {
			task, err := worker.NewDownloadLyricsTaskHandler(path)
			if err != nil {
				log.Error().Err(err).Str("path", path).Msg("Unable to create download task")
				return err
			} else {
				log.Debug().Str("path", path).Msg("Download task created successfully")
			}

			taskID := fmt.Sprintf("downloading-lyrics:%s", path)
			taskInfo, err := asynqClient.Enqueue(task, asynq.TaskID(taskID))
			switch {
			case errors.Is(err, asynq.ErrTaskIDConflict), errors.Is(err, asynq.ErrDuplicateTask):
				{
					log.Warn().Err(err).Msg("This track already has a task in queue for lyrics download, skipping...")
				}
			case err != nil:
				{
					log.Error().Err(err).Str("path", path).Msg("Unable to enqueue download task")
					return err
				}
			default:
				{
					log.Info().Str("path", path).Str("task_id", taskInfo.ID).Msg("Task pushed to queue")
				}
			}
		}

		return nil
	})
}
