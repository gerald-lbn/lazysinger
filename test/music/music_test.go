package music_test

import (
	"github.com/gerald-lbn/refrain/music"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Music", func() {
	When("When extracting music metadata", func() {
		var voreAudioPath = "../test_data/05. Vore.flac"
		var vorePlainLyrics = "../test_data/05. Vore.txt"

		Context("from a audio file", func() {
			metadata, err := music.ExtractMetadata(voreAudioPath)

			It("should return metadata and no error", func() {
				Expect(metadata).ToNot(BeNil())
				Expect(err).ToNot(HaveOccurred())
			})

			It("should return metadata with the correct title, artist album and duration", func() {
				Expect(*metadata.Title).To(Equal("Vore"))
				Expect(*metadata.Artist).To(Equal("Sleep Token"))
				Expect(*metadata.Album).To(Equal("Take Me Back To Eden"))
				Expect(metadata.Duration).To(BeNumerically("~", 4, 1))
			})

			It("should return true if all metadata are set", func() {
				Expect(metadata.HasAllMetadata()).To(BeTrue())
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
})
