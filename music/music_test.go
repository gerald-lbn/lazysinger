package music_test

import (
	"os"
	"testing"

	"github.com/gerald-lbn/lazysinger/log"
	"github.com/gerald-lbn/lazysinger/music"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestMusic(t *testing.T) {
	log.SetLevel(log.FatalLevel)
	RegisterFailHandler(Fail)
	RunSpecs(t, "Music Suite")
}

var _ = Describe("Music", func() {
	Context("Metadata", func() {
		var (
			BAD_OMENS_IMPOSE               = "test_data/Impose.wav"
			BAD_OMENS_IMPOSE_SYNCED_LYRICS = "test_data/Impose.lrc"
			BAD_OMENS_IMPOSE_PLAIN_LYRICS  = "test_data/Impose.txt"
		)

		When("extracting metadata from a music file", func() {
			It("extracts metadata correctly", func() {
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

			It("returns an error for non-existing file", func() {
				_, err := music.ExtractMetadaFromMusicFile("non_existing_file.mp3")
				Expect(err).To(HaveOccurred())
			})

			It("returns an error for non-music file", func() {
				_, err := music.ExtractMetadaFromMusicFile("metadata.go")
				Expect(err).To(HaveOccurred())
			})
		})

		When("handling lyrics paths", func() {
			It("gets lyrics paths correctly", func() {
				musicFilePath := "/music/artist/album/song.mp3"
				lyricsPath := music.GetLyricsPathFromMusicFilePath(musicFilePath)

				Expect(lyricsPath.PlainLyrics).To(Equal("/music/artist/album/song.txt"))
				Expect(lyricsPath.SyncedLyrics).To(Equal("/music/artist/album/song.lrc"))
			})

			It("checks lyrics existence correctly", func() {
				lyricsPath := music.LyricsPath{
					PlainLyrics:  BAD_OMENS_IMPOSE_PLAIN_LYRICS,
					SyncedLyrics: BAD_OMENS_IMPOSE_SYNCED_LYRICS,
				}
				Expect(music.CheckFileExistance(lyricsPath.PlainLyrics)).To(BeTrue())
				Expect(music.CheckFileExistance(lyricsPath.SyncedLyrics)).To(BeTrue())
			})
		})
	})

	Context("LyricsProcessor", func() {
		var (
			txtFilePath string
			lrcFilePath string
		)

		BeforeEach(func() {
			txtFilePath = "test_lyrics.txt"
			lrcFilePath = "test_lyrics.lrc"
		})

		When("downloading lyrics", func() {
			AfterEach(func() {
				os.Remove(txtFilePath)
				os.Remove(lrcFilePath)
			})

			It("creates a txt file with plain lyrics", func() {
				lyrics := "These are the test lyrics."

				err := music.DownloadLyrics(txtFilePath, lyrics)
				Expect(err).NotTo(HaveOccurred())

				data, err := os.ReadFile(txtFilePath)
				Expect(err).NotTo(HaveOccurred())
				Expect(string(data)).To(Equal(lyrics))
			})

			It("creates a lrc file with synced lyrics", func() {
				lyrics := "[00:12.00]These are the test lyrics.\n[00:34.00]With timestamps."

				err := music.DownloadLyrics(lrcFilePath, lyrics)
				Expect(err).NotTo(HaveOccurred())

				data, err := os.ReadFile(lrcFilePath)
				Expect(err).NotTo(HaveOccurred())
				Expect(string(data)).To(Equal(lyrics))
			})
		})

		When("deleting lyrics", func() {
			BeforeEach(func() {
				err := os.WriteFile(txtFilePath, []byte("Temporary lyrics"), 0644)
				Expect(err).NotTo(HaveOccurred())
				err = os.WriteFile(lrcFilePath, []byte("[00:12.00]Test lyrics"), 0644)
				Expect(err).NotTo(HaveOccurred())
			})

			AfterEach(func() {
				os.Remove(txtFilePath)
				os.Remove(lrcFilePath)
			})

			It("deletes txt file successfully", func() {
				err := music.DeleteLyrics(txtFilePath)
				Expect(err).NotTo(HaveOccurred())
				_, err = os.Stat(txtFilePath)
				Expect(os.IsNotExist(err)).To(BeTrue())
			})

			It("deletes lrc file successfully", func() {
				err := music.DeleteLyrics(lrcFilePath)
				Expect(err).NotTo(HaveOccurred())
				_, err = os.Stat(lrcFilePath)
				Expect(os.IsNotExist(err)).To(BeTrue())
			})

			It("returns error for non-existent file", func() {
				err := music.DeleteLyrics("non_existent_file.txt")
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Context("LyricsProvider", func() {
		var provider music.LyricsProvider

		BeforeEach(func() {
			provider = *music.NewLyricsProvider()
		})

		When("searching for lyrics regardless of specifying duration", func() {
			When("lyrics are not found", func() {
				It("returns error", func() {
					duration := 150
					lyricsResult, err := provider.Get(music.GetParameters{
						TrackName:  "abc",
						ArtistName: "abc",
						AlbumName:  "abc",
						Duration:   &duration,
					})

					Expect(lyricsResult).To(BeNil())
					Expect(err).To(HaveOccurred())
				})

				It("returns error", func() {
					lyricsResult, err := provider.Get(music.GetParameters{
						TrackName:  "abc",
						ArtistName: "abc",
						AlbumName:  "abc",
						Duration:   nil,
					})

					Expect(lyricsResult).To(BeNil())
					Expect(err).To(HaveOccurred())
				})
			})

			When("lyrics are found", func() {
				It("returns lyrics", func() {
					duration := 476
					lyricsResult, err := provider.Get(music.GetParameters{
						TrackName:  "Everglow",
						ArtistName: "STARSET",
						AlbumName:  "Vessels 2.0",
						Duration:   &duration,
					})

					Expect(lyricsResult).NotTo(BeNil())
					Expect(err).NotTo(HaveOccurred())
				})

				It("returns lyrics", func() {
					lyricsResult, err := provider.Get(music.GetParameters{
						TrackName:  "Everglow",
						ArtistName: "STARSET",
						AlbumName:  "Vessels 2.0",
						Duration:   nil,
					})

					Expect(lyricsResult).NotTo(BeNil())
					Expect(err).NotTo(HaveOccurred())
				})
			})
		})
	})

	Context("Utils", func() {
		BeforeEach(func() {
			os.Create("test_data/song.mp3")
			os.Create("test_data/song.MP3")
			os.Create("test_data/song.flac")
			os.Create("test_data/song.FLAC")
			os.Create("test_data/document.txt")
			os.Create("test_data/image.jpg")
			os.Create("test_data/video.mp4")
			os.Create("test_data/archive.zip")
		})

		AfterEach(func() {
			os.Remove("test_data/song.mp3")
			os.Remove("test_data/song.MP3")
			os.Remove("test_data/song.flac")
			os.Remove("test_data/song.FLAC")
			os.Remove("test_data/document.txt")
			os.Remove("test_data/image.jpg")
			os.Remove("test_data/video.mp4")
			os.Remove("test_data/archive.zip")
		})

		When("checking music file types", func() {
			It("identifies MP3 files correctly", func() {
				Expect(music.IsMusicFile("test_data/song.mp3")).To(BeTrue())
				Expect(music.IsMusicFile("test_data/song.MP3")).To(BeTrue())
			})

			It("identifies FLAC files correctly", func() {
				Expect(music.IsMusicFile("test_data/song.flac")).To(BeTrue())
				Expect(music.IsMusicFile("test_data/song.FLAC")).To(BeTrue())
			})

			It("rejects non-music files", func() {
				Expect(music.IsMusicFile("test_data/document.txt")).To(BeFalse())
				Expect(music.IsMusicFile("test_data/image.jpg")).To(BeFalse())
				Expect(music.IsMusicFile("test_data/video.mp4")).To(BeFalse())
				Expect(music.IsMusicFile("test_data/archive.zip")).To(BeFalse())
			})

			It("handles non-existent files", func() {
				Expect(music.IsMusicFile("test_data/non_existent.mp3")).To(BeFalse())
			})

			It("handles directories", func() {
				Expect(music.IsMusicFile("test_data/test_data")).To(BeFalse())
			})
		})
	})
})
