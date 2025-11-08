package music

import (
	"errors"
	"path/filepath"
	"strings"

	"go.senan.xyz/taglib"
)

var ErrNoExtensionInPath = errors.New("No extension found in path")
var ErrEmptyFilePath = errors.New("Path is empty")

// Metadata contains the metadata properties of an audio file
type Metadata struct {
	// Title is the name of the audio
	Title *string
	// Artist is the album artist's name of the audio
	Artist *string
	// Album is the name of the album the audio belongs to
	Album *string
	// Duration is the length of the audio in seconds
	Duration float64
}

func (m *Metadata) HasAllMetadata() bool {
	return m.Title != nil && m.Artist != nil && m.Album != nil
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

	return &Metadata{
		Title:    &title,
		Album:    &album,
		Artist:   &artist,
		Duration: properties.Length.Seconds(),
	}, nil
}

// generateLyricsFilePathFromAudioFilePath is a helper function to create a lyrics file path
// with a specified extension from an audio file path.
func generateLyricsFilePathFromAudioFilePath(p, ext string) (string, error) {
	if p == "" || p == "." {
		return "", ErrEmptyFilePath
	}

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
	return generateLyricsFilePathFromAudioFilePath(p, "txt")
}

// GenerateSyncedLyricsFilePathFromAudioFilePath generates a .lrc (synced lyrics) file path
// from an audio file path.
func GenerateSyncedLyricsFilePathFromAudioFilePath(p string) (string, error) {
	return generateLyricsFilePathFromAudioFilePath(p, "lrc")
}
