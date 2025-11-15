package tasks

import (
	"context"
	"time"

	"github.com/gerald-lbn/refrain/pkg/log"
	"github.com/gerald-lbn/refrain/pkg/services"
	"github.com/mikestefanello/backlite"
)

type DownloadLyricsTask struct {
	Path string
}

func (t DownloadLyricsTask) Config() backlite.QueueConfig {
	return backlite.QueueConfig{
		Name:        "music.sync_lyrics",
		MaxAttempts: 10,
		Backoff:     24 * time.Hour,
		Retention: &backlite.Retention{
			OnlyFailed: false,
			Data: &backlite.RetainData{
				OnlyFailed: false,
			},
		},
	}
}

func NewDownloadLyricsTaskQueue(c *services.Container) backlite.Queue {
	return backlite.NewQueue(func(ctx context.Context, dlt DownloadLyricsTask) error {
		log.Default().Info("download lyrics task received", "path", dlt.Path)
		return nil
	})
}
