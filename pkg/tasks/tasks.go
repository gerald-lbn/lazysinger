package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/gerald-lbn/refrain/pkg/music"
	"github.com/gerald-lbn/refrain/pkg/music/lrclib"
	"github.com/hibiken/asynq"
)

const (
	TypeDownloadLyrics = "download:lyrics"

	DownloadLyricsQueue = "download_lyrics"
)

type DownloadLyricsPayload struct {
	FilePath string
}

func NewDownloadLyricsTask(path string) (*asynq.Task, error) {
	payload, err := json.Marshal(DownloadLyricsPayload{FilePath: path})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(
		TypeDownloadLyrics,
		payload,
		asynq.TaskID(path),
		asynq.MaxRetry(1),
		asynq.Queue(DownloadLyricsQueue),
	), nil
}

func HandleDownloadLyricsTask(ctx context.Context, t *asynq.Task) error {
	var payload DownloadLyricsPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	track, err := music.ExtractMetadata(payload.FilePath)
	if err != nil {
		log.Printf("skipping track: '%s'. Reason: %s", payload.FilePath, err)
		return asynq.SkipRetry
	}

	// Skip task if track already has both lyrics
	if track.HasBothLyricsStoredLocally() {
		return asynq.RevokeTask
	}

	provider := lrclib.NewLRCLibProvider()
	var options lrclib.SearchLyricsOptions
	if track.HasAllMetadata() {
		options = lrclib.WithTrackArtistAndAlbumName(*track.Title, *track.Artist, *track.Album)
	} else if track.Artist == nil || track.Title == nil {
		log.Printf("skipping track: '%s'. Reason: Not enough metadata to search", track.Path)
		return asynq.SkipRetry
	} else {
		options = lrclib.WithTrackAndArtistName(*track.Title, *track.Artist)
	}

	lyrics, err := provider.GetLyrics(ctx, options, int(track.Duration))
	if err != nil {
		log.Printf("skipping track: '%s'. Reason: %s", track.Path, err)
		return asynq.SkipRetry
	}

	if lyrics == nil {
		log.Printf("skipping track: '%s'. Empty lyrics", track.Path)
		return nil
	}

	// Skip instrumental track
	if lyrics.Instrumental {
		return nil
	}

	if len(lyrics.PlainLyrics) > 0 && !track.HasPlainLyrics {
		err = os.WriteFile(track.PlainLyricsPath, []byte(lyrics.PlainLyrics), 0644)
		if err != nil {
			log.Printf("Failed to write to '%s'. Reason: %s", track.PlainLyricsPath, err)
		}
	}

	if len(lyrics.SyncedLyrics) > 0 && !track.HasSyncedLyrics {
		err = os.WriteFile(track.SyncedLyricsPath, []byte(lyrics.SyncedLyrics), 0644)
		if err != nil {
			log.Printf("Failed to write to '%s'. Reason: %s", track.SyncedLyricsPath, err)
		}
	}

	return nil
}
