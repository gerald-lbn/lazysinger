package queue

import (
	"encoding/json"
	"time"

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
		return nil, err
	}

	return asynq.NewTask(TypeLyricsDelete, payload), nil
}
