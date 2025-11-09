package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/gerald-lbn/refrain/music"
	"github.com/gerald-lbn/refrain/music/lrclib"
	"github.com/hibiken/asynq"
)

const (
	TypeDownloadLyrics = "download:lyrics"
)

type DownloadLyricsPayload struct {
	music.Metadata
}

func NewDownloadLyricsTask(metadata music.Metadata) (*asynq.Task, error) {
	payload, err := json.Marshal(DownloadLyricsPayload{metadata})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(
		TypeDownloadLyrics,
		payload,
		asynq.Queue(DownloadLyricsQueue),
		asynq.TaskID(metadata.Path),
		asynq.MaxRetry(1),
	), nil
}

func HandleDownloadLyricsTask(ctx context.Context, t *asynq.Task) error {
	var p DownloadLyricsPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	// Skip task if track already has both lyrics
	if p.HasBothLyricsStoredLocally() {
		return asynq.RevokeTask
	}

	provider := lrclib.NewLRCLibProvider()
	var options lrclib.SearchLyricsOptions
	if p.HasAllMetadata() {
		options = lrclib.WithTrackArtistAndAlbumName(*p.Title, *p.Artist, *p.Album)
	} else if p.Artist == nil || p.Title == nil {
		log.Printf("skipping track: '%s'. Reason: Not enough metadata to search", p.Path)
		return asynq.SkipRetry
	} else {
		options = lrclib.WithTrackAndArtistName(*p.Title, *p.Artist)
	}

	_, err := provider.GetLyrics(ctx, options, int(p.Duration))
	if err != nil {
		log.Printf("skipping track: '%s'. Reason: %s", p.Path, err)
		return asynq.SkipRetry
	}

	return nil
}
