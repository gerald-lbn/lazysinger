package music

import (
	"os"

	"github.com/gerald-lbn/lazysinger/log"
)

type LyricsProcessor struct{}

func NewLyricsProcessor() *LyricsProcessor {
	return &LyricsProcessor{}
}

func DownloadLyrics(pathToFile string, lyrics string) error {
	log.Debug().Str("file", pathToFile).Msg("Creating file")
	file, err := os.Create(pathToFile)
	if err != nil {
		log.Error().Err(err).Str("file", pathToFile).Msg("An error occured when creating file")
		return err
	}
	defer file.Close()

	log.Debug().Str("lyrics", lyrics).Str("file", pathToFile).Msg("Writing lyrics to file")
	_, err = file.WriteString(lyrics)
	if err != nil {
		log.Error().Err(err).Str("file", pathToFile).Str("lyrics", lyrics).Msg("An error occured when writing lyrics to file")
	}

	return err
}

func DeleteLyrics(pathToFile string) error {
	log.Debug().Str("file", pathToFile).Msg("Deleting lyrics")
	err := os.Remove(pathToFile)
	if err != nil {
		log.Error().Err(err).Msg("An occured while deleting lyrics")
	}
	return err
}
