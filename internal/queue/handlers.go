package queue

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gerald-lbn/lazysinger/internal/log"
	"github.com/gerald-lbn/lazysinger/internal/lyrics"
	"github.com/gerald-lbn/lazysinger/internal/music"
	"github.com/hibiken/asynq"
)

func HandleDownloadLyrics(ctx context.Context, t *asynq.Task) error {
	var p LyricsDownloadPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		returnError := fmt.Errorf("Unable to json.Unmarshal task payload: %v: %w", err, asynq.SkipRetry)
		log.Error().Err(returnError).Interface("payload", p)
		return returnError
	}

	metadata, err := music.ExtractMetadaFromMusicFile(p.Filepath)
	if err != nil {
		return err
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
			log.Info().Str("track", p.Filepath).Msg("The track is flagged as an instrumental track by the provider")
			return nil
		}

		if !metadata.HasSyncedLyrics {
			log.Info().Str("track", p.Filepath).Msg("The track doesn't have synced lyrics stored locally")
			if results.SyncedLyrics != "" {
				log.Debug().Str("path", p.Filepath).Msg("Found synced lyrics, saving them locally.")
				err := lyrics.DownloadLyrics(metadata.PathToSyncedLyrics, results.SyncedLyrics)
				if err != nil {
					return err
				}
				log.Info().Str("path", p.Filepath).Str("synced-lyrics_path", metadata.PathToSyncedLyrics).Msg("Successfully saved synced lyrics")
			} else {
				log.Warn().Str("track", p.Filepath).Msg("No synced lyrics found from provider")
			}
		}
		if !metadata.HasPlainLyrics {
			log.Info().Str("track", p.Filepath).Msg("The track doesn't have plain lyrics stored locally")
			if results.PlainLyrics != "" {
				log.Debug().Str("path", p.Filepath).Msg("Found plain lyrics, saving them locally")
				err := lyrics.DownloadLyrics(metadata.PathToPlainLyrics, results.PlainLyrics)
				if err != nil {
					return err
				}
				log.Info().Str("track", p.Filepath).Str("plain_lyrics_path", metadata.PathToPlainLyrics).Msg("Successfully saved plain lyrics")
			} else {
				log.Warn().Str("track", p.Filepath).Msg("No plain lyrics found from the provider")
			}
		}
	} else {
		log.Info().Str("track", metadata.FilePath).Msg("The track already has both lyrics stored locally")
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
		log.Debug().Str("lyrics", metadata.PathToPlainLyrics).Msg("Removing lyrics")
		if err := lyrics.DeleteLyrics(metadata.PathToPlainLyrics); err != nil {
			log.Error().Err(err).Str("lyrics", metadata.PathToPlainLyrics).Msg("Unable to remove lyrics")
			return err
		}
		log.Info().Str("lyrics", metadata.PathToPlainLyrics).Msg("Lyrics removed")
	}

	if metadata.HasSyncedLyrics {
		log.Debug().Str("lyrics", metadata.PathToPlainLyrics).Msg("Removing lyrics")
		if err := lyrics.DeleteLyrics(metadata.PathToSyncedLyrics); err != nil {
			log.Error().Err(err).Str("lyrics", metadata.PathToSyncedLyrics).Msg("Unable to remove lyrics")
			return err
		}
		log.Info().Str("lyrics", metadata.PathToSyncedLyrics).Msg("Lyrics removed")
	}

	return nil
}
