package handlers

import (
	"errors"
	"log"

	"github.com/gerald-lbn/refrain/music"
	"github.com/gerald-lbn/refrain/utils/file"
	"github.com/gerald-lbn/refrain/worker"
	"github.com/hibiken/asynq"
)

func OnInitialScan(p string) error {
	ok, err := file.IsAudioFile(p)
	if err != nil || !ok {
		return err
	}

	// Checks if track has both lyrics stored already
	metadata, err := music.ExtractMetadata(p)
	if err != nil {
		return err
	}
	if !metadata.HasAllMetadata() {
		return nil
	}

	task, err := worker.NewDownloadLyricsTask(*metadata)
	if err != nil {
		if errors.Is(err, asynq.ErrTaskIDConflict) {
			log.Printf("'%s' is already in queue for processing, skipping", p)
			return nil
		}

		return err
	}

	client := worker.NewClient()
	defer client.Close()

	info, err := client.Enqueue(task)
	if err != nil {
		log.Print("An error occured while enqueuing task", err)
		return err
	}

	log.Printf("Task #%s pushed to '%s' to process %s", info.ID, info.Queue, p)

	return nil
}
