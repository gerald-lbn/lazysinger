package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/gerald-lbn/lyrisync/internal/lyrics"
	"github.com/gerald-lbn/lyrisync/internal/music"
	"github.com/hibiken/asynq"
)

func HandleDownloadLyrics(ctx context.Context, t *asynq.Task) error {
	var p LyricsDownloadPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("Unable to json.Unmarshal task payload: %v: %w", err, asynq.SkipRetry)
	}

	metadata, err := music.ExtractMetadaFromMusicFile(p.Filepath)
	if err != nil {
		return fmt.Errorf("Unable to extract metadata: %v", err)
	}

	if !metadata.HasBothLyrics() {
		lyricsProvider := lyrics.NewLyricsProvider()
		results, err := lyricsProvider.Get(lyrics.GetParameters{
			TrackName:  metadata.TrackName,
			ArtistName: metadata.ArtistName,
			AlbumName:  metadata.AlbumName,
		})
		if err != nil {
			return err
		}

		if !metadata.HasSyncedLyrics {
			log.Printf("Track at %s doesn't have synced lyrics stored locally", p.Filepath)
			if results.SyncedLyrics != "" {
				log.Printf("Found synced lyrics for %s, saving them locally", p.Filepath)
				err := lyrics.DownloadLyrics(metadata.PathToSyncedLyrics, results.SyncedLyrics)
				if err != nil {
					log.Printf("An error occured while saving synced lyrics to %s: %v", metadata.PathToSyncedLyrics, err)
					return err
				}
				log.Printf("Successfully saved synced lyrics at %s", metadata.PathToSyncedLyrics)
			}
		}
		if !metadata.HasPlainLyrics {
			log.Printf("Track at %s doesn't have plain lyrics stored locally", p.Filepath)
			if results.PlainLyrics != "" {
				log.Printf("Found plain lyrics for %s, saving them locally", p.Filepath)
				err := lyrics.DownloadLyrics(metadata.PathToPlainLyrics, results.PlainLyrics)
				if err != nil {
					log.Printf("An error occured while saving plain lyrics to %s: %v", metadata.PathToPlainLyrics, err)
					return err
				}
				log.Printf("Successfully saved plain lyrics at %s", metadata.PathToPlainLyrics)
			}
		}
	} else {
		log.Printf("Track at %s already has both lyrics stored locally", p.Filepath)
	}

	return nil
}
