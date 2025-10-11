package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/gerald-lbn/lazysinger/internal/lyrics"
	"github.com/gerald-lbn/lazysinger/internal/music"
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

		if results.Instrumental {
			log.Printf("Track at %s is flagged as instrumental by the provider", p.Filepath)
			return nil
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
			} else {
				log.Printf("No synced lyrics found from provider for %s", p.Filepath)
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
			} else {
				log.Printf("No plain lyrics found from provider for %s", p.Filepath)
			}
		}
	} else {
		log.Printf("Track at %s already has both lyrics stored locally", p.Filepath)
	}

	return nil
}

func HandleDeleteLyrics(ctx context.Context, t *asynq.Task) error {
	var p LyricsDeletePayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("Unable to json.Unmarshal task payload: %v: %w", err, asynq.SkipRetry)
	}

	metadata, err := music.ExtractMetadaFromMusicFile(p.Filepath)
	if err != nil {
		return fmt.Errorf("Unable to extract metadata: %v", err)
	}

	if metadata.HasPlainLyrics {
		log.Printf("Removing lyrics stored at %s", metadata.PathToPlainLyrics)
		if err := lyrics.DeleteLyrics(metadata.PathToPlainLyrics); err != nil {
			log.Printf("An error occured when deleting lyrics at %s: %v", metadata.PathToPlainLyrics, err)
			return err
		}
		log.Printf("Lyrics stored at %s removed", metadata.PathToPlainLyrics)
	}

	if metadata.HasSyncedLyrics {
		log.Printf("Removing lyrics stored at %s", metadata.PathToSyncedLyrics)
		if err := lyrics.DeleteLyrics(metadata.PathToSyncedLyrics); err != nil {
			log.Printf("An error occured when deleting lyrics at %s: %v", metadata.PathToSyncedLyrics, err)
			return err
		}
		log.Printf("Lyrics stored at %s removed", metadata.PathToSyncedLyrics)
	}

	return nil
}
