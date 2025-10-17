package music

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/gerald-lbn/lazysinger/log"
	"go.senan.xyz/taglib"
)

type Metadata struct {
	FilePath           string
	TrackName          string
	ArtistName         string
	AlbumName          string
	HasPlainLyrics     bool
	PathToPlainLyrics  string
	HasSyncedLyrics    bool
	PathToSyncedLyrics string
}

type LyricsPath struct {
	PlainLyrics  string
	SyncedLyrics string
}

func (m *Metadata) HasBothLyrics() bool {
	return m.HasPlainLyrics && m.HasSyncedLyrics
}

func CheckFileExistance(path string) bool {
	_, err := os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		log.Debug().Str("path", path).Msg("Lyrics file does not exist")
		return false
	}
	if err != nil {
		log.Error().Err(err).Str("path", path).Msg("An error occured when checking lyrics existance")
		return false
	}
	return err == nil
}

func GetLyricsPathFromMusicFilePath(path string) LyricsPath {
	log.Debug().Str("path", path).Msg("Crafting lyrics path from path")
	dir := filepath.Dir(path)
	fileNameWithoutExt := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	lyricsPaths := LyricsPath{
		PlainLyrics:  filepath.Join(dir, fileNameWithoutExt+PLAIN_LYRICS_EXTENSION),
		SyncedLyrics: filepath.Join(dir, fileNameWithoutExt+SYNCED_LYRICS_EXTENSION),
	}
	log.Debug().Interface("lyrics_path", lyricsPaths)

	return lyricsPaths
}

func ExtractMetadaFromMusicFile(path string) (Metadata, error) {
	log.Debug().Str("path", path).Msg("Extracting audio metadata from filepath")
	metadata, err := taglib.ReadTags(path)
	if err != nil {
		log.Error().Err(err).Str("path", path).Msg("An error occured when extracting metadata from file")
		return Metadata{}, err
	}

	lyricsPath := GetLyricsPathFromMusicFilePath(path)

	return Metadata{
		FilePath:           path,
		TrackName:          metadata[taglib.Title][0],
		ArtistName:         metadata[taglib.Artist][0],
		AlbumName:          metadata[taglib.Album][0],
		HasPlainLyrics:     CheckFileExistance(lyricsPath.PlainLyrics),
		PathToPlainLyrics:  lyricsPath.PlainLyrics,
		HasSyncedLyrics:    CheckFileExistance(lyricsPath.SyncedLyrics),
		PathToSyncedLyrics: lyricsPath.SyncedLyrics,
	}, nil
}
