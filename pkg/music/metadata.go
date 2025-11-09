package music

import (
	"errors"
	"path/filepath"
	"strings"

	"github.com/gerald-lbn/refrain/pkg/utils/file"
	"go.senan.xyz/taglib"
)

const (
	SYNCED_LYRICS_EXTENSION = "lrc"
	PLAIN_LYRICS_EXTENSION  = "txt"
)

var (
	ErrNoExtensionInPath = errors.New("No extension found in path")
)

// Metadata contains the metadata properties of an audio file
type Metadata struct {
	// Path is the absolute filepath the audio
	Path string
	// Title is the name of the audio
	Title *string
	// Artist is the album artist's name of the audio
	Artist *string
	// Album is the name of the album the audio belongs to
	Album *string
	// Duration is the length of the audio in seconds
	Duration float64
	// HasPlainLyrics indicates whether the audio has plain lyrics stored locally
	HasPlainLyrics bool
	// PlainLyricsPath points to the plain lyrics stored locally
	PlainLyricsPath string
	// HasSyncedLyrics indicates whether the audio has synced lyrics stored locally
	HasSyncedLyrics bool
	// SyncedLyricsPath points to the synced lyrics stored locally
	SyncedLyricsPath string
}

func (m *Metadata) HasAllMetadata() bool {
	return m.Title != nil && m.Artist != nil && m.Album != nil
}

func (m *Metadata) HasBothLyricsStoredLocally() bool {
	return m.HasPlainLyrics && m.HasSyncedLyrics
}

// ExtractMetadata extracts metadata from an audio file specified by its path.
func ExtractMetadata(p string) (*Metadata, error) {
	tags, err := taglib.ReadTags(p)
	if err != nil {
		return nil, err
	}

	properties, err := taglib.ReadProperties(p)
	if err != nil {
		return nil, err
	}

	var title string
	if len(tags[taglib.Title]) > 0 {
		title = tags[taglib.Title][0]
	}

	var artist string
	if len(tags[taglib.AlbumArtist]) > 0 {
		artist = tags[taglib.AlbumArtist][0]
	}

	var album string
	if len(tags[taglib.Album]) > 0 {
		album = tags[taglib.Album][0]
	}

	plainLyricsPath, err := GeneratePlainLyricsFilePathFromAudioFilePath(p)
	if err != nil {
		return nil, err
	}

	hasPlainLyrics := file.Exists(plainLyricsPath)

	syncedLyricsPath, err := GenerateSyncedLyricsFilePathFromAudioFilePath(p)
	if err != nil {
		return nil, err
	}

	hasSyncedLyrics := file.Exists(syncedLyricsPath)

	return &Metadata{
		Path:             p,
		Title:            &title,
		Album:            &album,
		Artist:           &artist,
		Duration:         properties.Length.Seconds(),
		HasPlainLyrics:   hasPlainLyrics,
		PlainLyricsPath:  plainLyricsPath,
		HasSyncedLyrics:  hasSyncedLyrics,
		SyncedLyricsPath: syncedLyricsPath,
	}, nil
}

// generateLyricsFilePathFromAudioFilePath is a helper function to create a lyrics file path
// with a specified extension from an audio file path.
func generateLyricsFilePathFromAudioFilePath(p, ext string) (string, error) {
	audioExt := filepath.Ext(p)
	if audioExt == "" {
		return "", ErrNoExtensionInPath
	}

	base := strings.TrimSuffix(p, audioExt)
	return base + "." + ext, nil
}

// GeneratePlainLyricsFilePathFromAudioFilePath generates a .txt lyrics file path
// from an audio file path.
func GeneratePlainLyricsFilePathFromAudioFilePath(p string) (string, error) {
	return generateLyricsFilePathFromAudioFilePath(p, PLAIN_LYRICS_EXTENSION)
}

// GenerateSyncedLyricsFilePathFromAudioFilePath generates a .lrc (synced lyrics) file path
// from an audio file path.
func GenerateSyncedLyricsFilePathFromAudioFilePath(p string) (string, error) {
	return generateLyricsFilePathFromAudioFilePath(p, SYNCED_LYRICS_EXTENSION)
}
