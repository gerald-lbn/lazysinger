package purger

import (
	"context"
	"errors"
	"fmt"

	"github.com/gerald-lbn/lazysinger/config"
	"github.com/gerald-lbn/lazysinger/database"
	"github.com/gerald-lbn/lazysinger/log"
	"github.com/gerald-lbn/lazysinger/worker"
	"github.com/hibiken/asynq"
)

func PurgeAll(ctx context.Context) error {
	asynqClient := asynq.NewClient(asynq.RedisClientOpt{
		Addr: config.Server.Worker.RedisAddr,
	})

	db := database.Connect()

	sr := database.NewSongRepository(ctx, db)
	findManyByResult := sr.FindManyBy(database.NewSongCriteria())
	if findManyByResult.Error != nil {
		log.Error().Err(findManyByResult.Error).Msg("An error occured while querying all songs")
		return findManyByResult.Error
	}

	for _, song := range *findManyByResult.Data {
		task, err := worker.NewDatabasePurge(song.Path)
		if err != nil {
			log.Error().Err(err).Str("path", song.Path).Msg("Unable to create task to purge the entry from the database")
			return err
		}

		taskID := fmt.Sprintf("purging:%s", song.Path)
		taskInfo, err := asynqClient.Enqueue(task, asynq.TaskID(taskID))
		switch {
		case errors.Is(err, asynq.ErrTaskIDConflict), errors.Is(err, asynq.ErrDuplicateTask):
			{
				log.Warn().Err(err).Msg("This track already has a task in queue for purging, skipping...")
			}
		case err != nil:
			{
				log.Error().Err(err).Str("path", song.Path).Msg("Unable to enqueue purging task")
				return err
			}
		default:
			{
				log.Info().Str("path", song.Path).Str("task_id", taskInfo.ID).Msg("Task pushed to queue")
			}
		}
	}

	return nil
}
