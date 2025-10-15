package queue

import (
	"encoding/json"
	"time"

	"github.com/gerald-lbn/lazysinger/log"
	"github.com/hibiken/asynq"
)

// A list of task types
const (
	TypeLyricsDownload = "lyrics:download"
)

// Task payload
type LyricsDownloadPayload struct {
	Filepath string
}

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
