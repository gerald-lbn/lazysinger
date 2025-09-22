package music

import (
	"os"
	"path/filepath"
	"strings"

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
	plainLyrics  string
	syncedLyrics string
}

func (m *Metadata) HasBothLyrics() bool {
	return m.HasPlainLyrics && m.HasSyncedLyrics
}

func CheckLyricsExistance(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func GetLyricsPathFromMusicFilePath(path string) LyricsPath {
	dir := filepath.Dir(path)
	fileNameWithoutExt := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	lyricsPaths := LyricsPath{
		plainLyrics:  filepath.Join(dir, fileNameWithoutExt+PLAIN_LYRICS_EXTENSION),
		syncedLyrics: filepath.Join(dir, fileNameWithoutExt+SYNCED_LYRICS_EXTENSION),
	}

	return lyricsPaths
}

func ExtractMetadaFromMusicFile(path string) (Metadata, error) {
	metadata, err := taglib.ReadTags(path)
	if err != nil {
		return Metadata{}, err
	}

	lyricsPath := GetLyricsPathFromMusicFilePath(path)

	return Metadata{
		FilePath:           path,
		TrackName:          metadata[taglib.Title][0],
		ArtistName:         metadata[taglib.Artist][0],
		AlbumName:          metadata[taglib.Album][0],
		HasPlainLyrics:     CheckLyricsExistance(lyricsPath.plainLyrics),
		PathToPlainLyrics:  lyricsPath.plainLyrics,
		HasSyncedLyrics:    CheckLyricsExistance(lyricsPath.syncedLyrics),
		PathToSyncedLyrics: lyricsPath.syncedLyrics,
	}, nil
}
