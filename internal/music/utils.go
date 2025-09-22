package music

import (
	"path/filepath"
	"strings"
)

const (
	MP3_EXTENSION           = ".mp3"
	FLAC_EXTENSION          = ".flac"
	SYNCED_LYRICS_EXTENSION = ".lrc"
	PLAIN_LYRICS_EXTENSION  = ".txt"
)

func IsMusicFile(path string) bool {
	extName := strings.ToLower(filepath.Ext(path))
	switch extName {
	case MP3_EXTENSION, FLAC_EXTENSION:
		return true
	default:
		return false
	}
}
