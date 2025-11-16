package tasks

import (
	"context"
	"time"

	"github.com/gerald-lbn/refrain/pkg/music"
	"github.com/gerald-lbn/refrain/pkg/repository"
	"github.com/gerald-lbn/refrain/pkg/services"
	dbUtils "github.com/gerald-lbn/refrain/pkg/utils/db"
	"github.com/gerald-lbn/refrain/pkg/utils/file"
	"github.com/mikestefanello/backlite"
)

type PersistTrackInfoTask struct {
	Path string
}

func (t PersistTrackInfoTask) Config() backlite.QueueConfig {
	return backlite.QueueConfig{
		Name:        "music.persist_info",
		MaxAttempts: 10,
		Backoff:     1 * time.Minute,
		Retention: &backlite.Retention{
			OnlyFailed: false,
			Data: &backlite.RetainData{
				OnlyFailed: false,
			},
		},
	}
}

func NewPersistTrackInfoQueue(c *services.Container) backlite.Queue {
	return backlite.NewQueue(func(ctx context.Context, ptit PersistTrackInfoTask) error {
		if exists := file.Exists(ptit.Path); !exists {
			return nil
		}

		if isAudio, err := music.IsAudio(ptit.Path); err != nil {
			return err
		} else if !isAudio {
			return nil
		}

		track, err := music.ExtractMetadata(ptit.Path)
		if err != nil {
			return err
		}

		repo := repository.New(c.Database)

		// Update track info if it already exists
		if _, err := repo.GetTrackByPath(ctx, track.Path); err == nil {
			return repo.UpdateTrack(ctx, repository.UpdateTrackParams{
				Path:            track.Path,
				Title:           dbUtils.StringToNullString(*track.Title),
				Artist:          dbUtils.StringToNullString(*track.Artist),
				Album:           dbUtils.StringToNullString(*track.Album),
				Duration:        track.Duration,
				HasPlainLyrics:  track.HasPlainLyrics,
				HasSyncedLyrics: track.HasSyncedLyrics,
			})
		}

		// Create track if it doesn't exist
		err = repo.CreateTrack(ctx, repository.CreateTrackParams{
			Path:            track.Path,
			Title:           dbUtils.StringToNullString(*track.Title),
			Artist:          dbUtils.StringToNullString(*track.Artist),
			Album:           dbUtils.StringToNullString(*track.Album),
			Duration:        track.Duration,
			HasPlainLyrics:  track.HasPlainLyrics,
			HasSyncedLyrics: track.HasSyncedLyrics,
		})

		return err
	})
}
