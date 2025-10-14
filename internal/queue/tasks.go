package queue

import (
	"encoding/json"
	"time"

	"github.com/gerald-lbn/lazysinger/internal/log"
	"github.com/hibiken/asynq"
)

// A list of task types
const (
	TypeLyricsDownload = "lyrics:download"
	TypeLyricsDelete   = "lyrics:delete"
)

// Task payload
type LyricsDownloadPayload struct {
	Filepath string
}
type LyricsDeletePayload = LyricsDownloadPayload

func NewDownloadLyricsTask(filepath string) (*asynq.Task, error) {
	payload, err := json.Marshal(LyricsDownloadPayload{
		Filepath: filepath,
	})
	if err != nil {
		log.Error().Err(err).Msg("An error occurred while marshaling the lyrics download task payload")
		return nil, err
	}

	// Let times for the lyrics (.lrc and .txt) to be present on the filesystem.
	// Prevent false-negative local lyrics detection
	return asynq.NewTask(TypeLyricsDownload, payload, asynq.ProcessIn(10*time.Second)), nil
}

func NewDeleteLyricsTask(filepath string) (*asynq.Task, error) {
	payload, err := json.Marshal(LyricsDeletePayload{
		Filepath: filepath,
	})
	if err != nil {
		log.Error().Err(err).Msg("An error occurred while marshaling the lyrics delete task payload")
		return nil, err
	}

	return asynq.NewTask(TypeLyricsDelete, payload), nil
}
