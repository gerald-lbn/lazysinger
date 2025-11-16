package tasks

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/gerald-lbn/refrain/pkg/log"
	"github.com/gerald-lbn/refrain/pkg/music"
	"github.com/gerald-lbn/refrain/pkg/music/lrclib"
	"github.com/gerald-lbn/refrain/pkg/services"
	"github.com/gerald-lbn/refrain/pkg/utils/file"
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
		if isAudio, err := file.IsAudioFile(dlt.Path); err != nil {
			return err
		} else if !isAudio {
			return nil
		}

		track, err := music.ExtractMetadata(dlt.Path)
		if err != nil {
			return err
		}

		// Skip task if track already has both lyrics
		if track.HasBothLyricsStoredLocally() {
			return nil
		}

		var options lrclib.SearchLyricsOptions
		if track.HasAllMetadata() {
			options = lrclib.WithTrackArtistAndAlbumName(*track.Title, *track.Artist, *track.Album)
		} else if track.Artist == nil || track.Title == nil {
			log.Default().Warn("skipping track",
				slog.String("path", dlt.Path),
				slog.String("reason", "not enough metadata to search"),
			)
			return lrclib.ErrMissingTrackOrArtistName
		} else {
			options = lrclib.WithTrackAndArtistName(*track.Title, *track.Artist)
		}

		lyrics, err := c.LyricsProvider.GetLyrics(ctx, options, int(track.Duration))
		if err != nil {
			return err
		}

		// Skip instrumental track
		if lyrics.Instrumental {
			return nil
		}

		// Write plain lyrics
		if len(lyrics.PlainLyrics) > 0 && !track.HasPlainLyrics {
			err = os.WriteFile(track.PlainLyricsPath, []byte(lyrics.PlainLyrics), 0644)
			if err != nil {
				log.Default().Error("failed to write file",
					slog.String("path", track.PlainLyricsPath),
					slog.String("content", lyrics.PlainLyrics),
					slog.String("error", err.Error()),
				)

				return err
			}
		}

		// Write synced lyrics
		if len(lyrics.SyncedLyrics) > 0 && !track.HasSyncedLyrics {
			err = os.WriteFile(track.SyncedLyricsPath, []byte(lyrics.SyncedLyrics), 0644)
			if err != nil {
				log.Default().Error("failed to write file",
					slog.String("path", track.SyncedLyricsPath),
					slog.String("content", lyrics.SyncedLyrics),
					slog.String("error", err.Error()),
				)

				return err
			}
		}

		return nil
	})
}
