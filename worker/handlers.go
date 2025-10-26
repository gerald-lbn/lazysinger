package worker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gerald-lbn/lazysinger/database"
	"github.com/gerald-lbn/lazysinger/log"
	"github.com/gerald-lbn/lazysinger/music"
	"github.com/hibiken/asynq"
)

func HandleDownloadLyricsTask(ctx context.Context, t *asynq.Task) error {
	var p LyricsDownloadPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		returnError := fmt.Errorf("Unable to json.Unmarshal task payload: %v: %w", err, asynq.SkipRetry)
		log.Error().Err(err).Interface("payload", p)
		return returnError
	}

	// extract metadata
	metadata, err := music.ExtractMetadaFromMusicFile(p.Filepath)
	if err != nil {
		log.Error().Err(err).Send()
		return err
	}

	// Skip lyrics fetching if both are present locally
	lyricsPaths := music.GetLyricsPathFromMusicFilePath(metadata.FilePath)
	if metadata.HasBothLyrics() {
		log.Info().Str("path", metadata.FilePath).Msg("Both lyrics type are already stored locally, skipping...")
		return nil
	}

	// Fetch lyrics
	lyricsProvider := music.NewLyricsProvider()
	results, err := lyricsProvider.Get(music.GetParameters{
		TrackName:  metadata.TrackName,
		ArtistName: metadata.ArtistName,
		AlbumName:  metadata.AlbumName,
	})
	if err != nil {
		log.Error().Str("path", metadata.FilePath).Err(err).Send()
		return err
	}

	// Download lyrics
	if len(results.PlainLyrics) > 0 {
		err := music.DownloadLyrics(lyricsPaths.PlainLyrics, results.PlainLyrics)
		if err != nil {
			log.Error().Str("plain_lyrics_path", lyricsPaths.PlainLyrics).Str("lyrics", results.PlainLyrics).Err(err).Msg("Unable to write")
			return err
		} else {
			log.Debug().Str("plain_lyrics_path", lyricsPaths.PlainLyrics).Msg("Plain lyrics downloaded successfully")
		}
	} else {
		log.Debug().Str("path", p.Filepath).Msg("No lyrics found")
	}
	if len(results.SyncedLyrics) > 0 {
		err := music.DownloadLyrics(lyricsPaths.SyncedLyrics, results.SyncedLyrics)
		if err != nil {
			log.Error().Str("synced_lyrics_path", lyricsPaths.SyncedLyrics).Str("lyrics", results.SyncedLyrics).Err(err).Msg("Unable to write")
			return err
		} else {
			log.Debug().Str("synced_lyrics_path", lyricsPaths.SyncedLyrics).Msg("Synced lyrics downloaded successfully")
		}
	} else {
		log.Debug().Str("path", p.Filepath).Msg("No lyrics found")
	}

	return nil
}

func HandleDatabasePurgeTask(ctx context.Context, t *asynq.Task) error {
	ERROR_MSG := "An error occured while cleaning up database."

	var p LyricsDownloadPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		returnError := fmt.Errorf("Unable to json.Unmarshal task payload: %v: %w", err, asynq.SkipRetry)
		log.Error().Err(err).Interface("payload", p).Msg(ERROR_MSG)
		return returnError
	}

	db := database.Connect()
	sr := database.NewSongRepository(ctx, db)

	// Make sure the song exists in the database
	findByResult := sr.FindBy(&database.SongCriteria{Path: &p.Filepath})
	if findByResult.Error != nil {
		log.Error().Err(findByResult.Error).Msg(ERROR_MSG)
		return findByResult.Error
	}

	// Delete the song
	deleteResult := sr.Delete(findByResult.Data)
	if deleteResult.Error != nil {
		log.Error().Err(deleteResult.Error).Str("path", p.Filepath).Msg(ERROR_MSG)
	}

	log.Info().Str("path", p.Filepath).Msg("Entry purged")

	if err := database.Close(db); err != nil {
		log.Error().Err(err).Msg(ERROR_MSG)
		return err
	}

	return nil
}
