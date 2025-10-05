package lyrics_test

import (
	"testing"

	"github.com/gerald-lbn/lazysinger/internal/lyrics"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestLyricsProvider(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "LyricsProvider Suite")
}

var _ = Describe("LyricsProvider", func() {
	var (
		provider lyrics.LyricsProvider
	)

	BeforeEach(func() {
		provider = *lyrics.NewLyricsProvider()
	})

	Context("when lyrics are not found", func() {
		It("should return an error and no lyrics", func() {
			duration := 150
			lyricsResult, err := provider.Get(lyrics.GetParameters{
				TrackName:  "abc",
				ArtistName: "abc",
				AlbumName:  "abc",
				Duration:   &duration,
			})

			Expect(lyricsResult).To(Equal(lyrics.LyricsResponse{}))
			Expect(err).To(HaveOccurred())
		})
	})

	Context("when lyrics are found", func() {
		It("should return lyrics and no error", func() {
			duration := 476
			lyricsResult, err := provider.Get(lyrics.GetParameters{
				TrackName:  "Everglow",
				ArtistName: "STARSET",
				AlbumName:  "Vessels 2.0",
				Duration:   &duration,
			})

			Expect(lyricsResult).NotTo(Equal(lyrics.LyricsResponse{}))
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
