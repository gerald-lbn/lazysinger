package music_test

import (
	"github.com/gerald-lbn/lazysinger/internal/music"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("LyricsProvider", func() {
	var (
		provider music.LyricsProvider
	)

	BeforeEach(func() {
		provider = *music.NewLyricsProvider()
	})

	Context("when lyrics are not found", func() {
		It("should return an error and no lyrics", func() {
			duration := 150
			lyricsResult, err := provider.Get(music.GetParameters{
				TrackName:  "abc",
				ArtistName: "abc",
				AlbumName:  "abc",
				Duration:   &duration,
			})

			Expect(lyricsResult).To(Equal(music.LyricsResponse{}))
			Expect(err).To(HaveOccurred())
		})

		It("should return an error even if duration is nil", func() {
			lyricsResult, err := provider.Get(music.GetParameters{
				TrackName:  "abc",
				ArtistName: "abc",
				AlbumName:  "abc",
				Duration:   nil,
			})

			Expect(lyricsResult).To(Equal(music.LyricsResponse{}))
			Expect(err).To(HaveOccurred())
		})
	})

	Context("when lyrics are found", func() {
		It("should return lyrics and no error", func() {
			duration := 476
			lyricsResult, err := provider.Get(music.GetParameters{
				TrackName:  "Everglow",
				ArtistName: "STARSET",
				AlbumName:  "Vessels 2.0",
				Duration:   &duration,
			})

			Expect(lyricsResult).NotTo(Equal(music.LyricsResponse{}))
			Expect(err).NotTo(HaveOccurred())
		})

		It("should return lyrics even if duration is nil", func() {
			lyricsResult, err := provider.Get(music.GetParameters{
				TrackName:  "Everglow",
				ArtistName: "STARSET",
				AlbumName:  "Vessels 2.0",
				Duration:   nil,
			})

			Expect(lyricsResult).NotTo(Equal(music.LyricsResponse{}))
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
