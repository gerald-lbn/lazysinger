package queue

import (
	"encoding/json"

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
		return nil, err
	}

	return asynq.NewTask(TypeLyricsDownload, payload), nil
}
