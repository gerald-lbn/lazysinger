package music

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/gerald-lbn/lazysinger/log"
)

const (
	MP3_EXTENSION           = ".mp3"
	FLAC_EXTENSION          = ".flac"
	SYNCED_LYRICS_EXTENSION = ".lrc"
	PLAIN_LYRICS_EXTENSION  = ".txt"
)

func IsMusicFile(path string) bool {
	log.Debug().Str("path", path).Msg("Checking if filepath points to a music file")
	fi, err := os.Stat(path)
	if err != nil {
		return false
	}

	if fi.Mode().IsRegular() {
		extName := strings.ToLower(filepath.Ext(path))
		switch extName {
		case MP3_EXTENSION, FLAC_EXTENSION:
			return true
		default:
			return false
		}
	}

	return false
}
