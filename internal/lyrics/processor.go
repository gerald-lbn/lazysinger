package lyrics

import "os"

type LyricsProcessor struct{}

func NewLyricsProcessor() *LyricsProcessor {
	return &LyricsProcessor{}
}

func DownloadLyrics(pathToFile string, lyrics string) error {
	file, err := os.Create(pathToFile)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(lyrics)

	return err
}

func DeleteLyrics(pathToFile string) error {
	err := os.Remove(pathToFile)
	return err
}
