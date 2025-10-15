package music_test

import (
	"github.com/gerald-lbn/lazysinger/internal/music"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Metadata", func() {
	var (
		BAD_OMENS_IMPOSE               = "test_data/Impose.wav"
		BAD_OMENS_IMPOSE_SYNCED_LYRICS = "test_data/Impose.lrc"
		BAD_OMENS_IMPOSE_PLAIN_LYRICS  = "test_data/Impose.txt"
	)

	Context("When extracting metadata from a music file", func() {
		It("should extract metadata correctly", func() {
			metadata, err := music.ExtractMetadaFromMusicFile(BAD_OMENS_IMPOSE)
			Expect(err).To(BeNil())
			Expect(metadata.FilePath).To(Equal(BAD_OMENS_IMPOSE))
			Expect(metadata.TrackName).To(Equal("Impose"))
			Expect(metadata.ArtistName).To(Equal("Bad Omens"))
			Expect(metadata.AlbumName).To(Equal("Impose"))
			Expect(metadata.HasPlainLyrics).To(BeTrue())
			Expect(metadata.PathToPlainLyrics).To(Equal(BAD_OMENS_IMPOSE_PLAIN_LYRICS))
			Expect(metadata.HasSyncedLyrics).To(BeTrue())
			Expect(metadata.PathToSyncedLyrics).To(Equal(BAD_OMENS_IMPOSE_SYNCED_LYRICS))
		})

		It("should return an error when extraction a non-existing file", func() {
			_, err := music.ExtractMetadaFromMusicFile("non_existing_file.mp3")
			Expect(err).ToNot(BeNil())
		})

		It("should return an error when extracting a non-music file", func() {
			_, err := music.ExtractMetadaFromMusicFile("metadata.go")
			Expect(err).ToNot(BeNil())
		})

		It("should get lyrics paths correctly", func() {
			musicFilePath := "/music/artist/album/song.mp3"
			lyricsPath := music.GetLyricsPathFromMusicFilePath(musicFilePath)

			Expect(lyricsPath.PlainLyrics).To(Equal("/music/artist/album/song.txt"))
			Expect(lyricsPath.SyncedLyrics).To(Equal("/music/artist/album/song.lrc"))
		})

		It("should check lyrics existence correctly", func() {
			lyricsPath := music.LyricsPath{
				PlainLyrics:  BAD_OMENS_IMPOSE_PLAIN_LYRICS,
				SyncedLyrics: BAD_OMENS_IMPOSE_SYNCED_LYRICS,
			}
			hasPlain := music.CheckLyricsExistance(lyricsPath.PlainLyrics)
			Expect(hasPlain).To(BeTrue())
			hasSynced := music.CheckLyricsExistance(lyricsPath.SyncedLyrics)
			Expect(hasSynced).To(BeTrue())
		})
	})
})
