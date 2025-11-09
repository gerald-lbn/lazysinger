package music_test

import (
	"github.com/gerald-lbn/refrain/music"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Music", func() {
	var voreAudioPath = "../test_data/Vore.flac"
	var vorePlainLyrics = "../test_data/Vore.txt"
	var voreSyncedLyrics = "../test_data/Vore.lrc"

	When("extracting music metadata", func() {
		Context("from a audio file", func() {
			metadata, err := music.ExtractMetadata(voreAudioPath)

			It("should return metadata and no error", func() {
				Expect(metadata).ToNot(BeNil())
				Expect(err).ToNot(HaveOccurred())
			})

			It("should return metadata with the correct informations", func() {
				Expect(*metadata.Title).To(Equal("Vore"))
				Expect(*metadata.Artist).To(Equal("Sleep Token"))
				Expect(*metadata.Album).To(Equal("Take Me Back To Eden"))
				Expect(metadata.Duration).To(BeNumerically("~", 4, 1))
				Expect(metadata.HasPlainLyrics).To(BeTrue())
				Expect(metadata.HasSyncedLyrics).To(BeTrue())
				Expect(metadata.PlainLyricsPath).To(Equal(vorePlainLyrics))
				Expect(metadata.SyncedLyricsPath).To(Equal(voreSyncedLyrics))
			})

			It("should return true if all metadata are set", func() {
				Expect(metadata.HasAllMetadata()).To(BeTrue())
			})

			It("should return true if both lyrics are stored locally", func() {
				Expect(metadata.HasBothLyricsStoredLocally()).To(BeTrue())
			})
		})

		Context("from a non audio file", func() {
			metadata, err := music.ExtractMetadata(vorePlainLyrics)

			It("should retuning an error and no metadata", func() {
				Expect(metadata).To(BeNil())
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("invalid file"))
			})
		})
	})

	When("generating path", func() {
		Context("for synced lyrics", func() {
			It("should generate a synced lyrics file path with .lrc extension", func() {
				audioPath := "/path/to/audio.mp3"
				expectedPath := "/path/to/audio.lrc"
				lyricsPath, err := music.GenerateSyncedLyricsFilePathFromAudioFilePath(audioPath)
				Expect(lyricsPath).To(Equal(expectedPath))
				Expect(err).ToNot(HaveOccurred())
			})

			It("should handle audio files with multiple dots in their name correctly", func() {
				audioPath := "/path/to/my.audio.file.flac"
				expectedPath := "/path/to/my.audio.file.lrc"
				lyricsPath, err := music.GenerateSyncedLyricsFilePathFromAudioFilePath(audioPath)
				Expect(lyricsPath).To(Equal(expectedPath))
				Expect(err).ToNot(HaveOccurred())
			})

			It("should return an error if the audio file path has no extension", func() {
				audioPath := "/path/to/audiofile"
				lyricsPath, err := music.GenerateSyncedLyricsFilePathFromAudioFilePath(audioPath)
				Expect(lyricsPath).To(BeEmpty())
				Expect(err).To(MatchError(music.ErrNoExtensionInPath))
			})
		})

		Context("for plain lyrics", func() {
			It("should generate a plain lyrics file path with .txt extension", func() {
				audioPath := "/path/to/audio.flac"
				expectedPath := "/path/to/audio.txt"
				lyricsPath, err := music.GeneratePlainLyricsFilePathFromAudioFilePath(audioPath)
				Expect(lyricsPath).To(Equal(expectedPath))
				Expect(err).ToNot(HaveOccurred())
			})

			It("should handle audio files with multiple dots in their name correctly", func() {
				audioPath := "/path/to/my.audio.file.mp3"
				expectedPath := "/path/to/my.audio.file.txt"
				lyricsPath, err := music.GeneratePlainLyricsFilePathFromAudioFilePath(audioPath)
				Expect(lyricsPath).To(Equal(expectedPath))
				Expect(err).ToNot(HaveOccurred())
			})

			It("should return an error if the audio file path has no extension", func() {
				audioPath := "/path/to/audiofile"
				lyricsPath, err := music.GeneratePlainLyricsFilePathFromAudioFilePath(audioPath)
				Expect(lyricsPath).To(BeEmpty())
				Expect(err).To(MatchError(music.ErrNoExtensionInPath))
			})
		})
	})
})
