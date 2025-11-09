package lrclib_test

import (
	"context"
	"net/http"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/gerald-lbn/refrain/pkg/music/lrclib"
)

func TestProvider(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "LrcLib Suite")
}

var _ = Describe("LrcLib", func() {
	var (
		lrclibClient *lrclib.LRCLibProvider
		client       *http.Client
		ctx          context.Context
	)

	BeforeEach(func() {
		ctx = context.Background()
	})

	Context("NewLRCLibProvider", func() {
		It("should return a new provider", func() {
			lrclibClient = lrclib.NewLRCLibProvider()

			Expect(lrclibClient).ToNot(BeNil())
			Expect(lrclibClient.BaseURL).To(Equal(lrclib.BASE_URL))
			Expect(lrclibClient.HttpClient).ToNot(BeNil())
		})

		When("options are provided", func() {
			It("should apply the provided options", func() {
				lrclibClient = lrclib.NewLRCLibProvider(
					lrclib.WithHttpClient(client),
				)

				Expect(lrclibClient).ToNot(BeNil())
				Expect(lrclibClient.BaseURL).To(Equal(lrclib.BASE_URL))
				Expect(lrclibClient.HttpClient).To(Equal(client))
			})
		})
	})

	Context("SearchLyrics", func() {
		BeforeEach(func() {
			lrclibClient = lrclib.NewLRCLibProvider()
		})

		When("provided with a valid query", func() {
			It("should return lyrics without error", func() {
				results, err := lrclibClient.SearchLyrics(ctx, lrclib.WithQuery("Tate McRae - It's ok I'm ok"))

				Expect(err).ToNot(HaveOccurred())
				Expect(results).ToNot(BeEmpty())
			})
		})

		When("provided with valid track and artist names", func() {
			It("should return lyrics without error", func() {
				results, err := lrclibClient.SearchLyrics(ctx, lrclib.WithTrackAndArtistName("Impose", "Bad Omens"))

				Expect(err).ToNot(HaveOccurred())
				Expect(results).ToNot(BeEmpty())
			})
		})

		When("provided with track, artist, and album names", func() {
			It("should return lyrics without error", func() {
				results, err := lrclibClient.SearchLyrics(
					ctx,
					lrclib.WithTrackArtistAndAlbumName(
						"die for",
						"Maggie Lindemann",
						"HEADSPLIT",
					),
				)

				Expect(err).ToNot(HaveOccurred())
				Expect(results).ToNot(BeEmpty())
			})
		})

		When("insufficient search parameters are provided", func() {
			It("should return ErrInsufficientSearchParameters", func() {
				lyrics, err := lrclibClient.SearchLyrics(ctx, lrclib.SearchLyricsOptions{})

				Expect(lyrics).To(BeNil())
				Expect(err).To(MatchError(lrclib.ErrInsufficientSearchParameters))
			})
		})

		// When("the API returns an HTTP error status", func() {
		// 	It("should return an error indicating the API failure", func() {})
		// })

		// When("the API returns a malformed response", func() {
		// 	It("should return an unmarshaling error", func() {})
		// })
	})

	Context("GetLyrics", func() {
		BeforeEach(func() {
			lrclibClient = lrclib.NewLRCLibProvider()
		})

		When("provided with track and artist name", func() {
			It("should return lyrics without error", func() {
				res, err := lrclibClient.GetLyrics(
					ctx,
					lrclib.WithTrackAndArtistName("Impose", "Bad Omens"),
					263,
				)

				Expect(err).ToNot(HaveOccurred())
				Expect(res).ToNot(BeNil())
			})
		})

		When("provided with track, artist and album name", func() {
			It("should return lyrics without error", func() {
				res, err := lrclibClient.GetLyrics(
					ctx,
					lrclib.WithTrackArtistAndAlbumName("Impose", "Bad Omens", "Impose"),
					263,
				)

				Expect(err).ToNot(HaveOccurred())
				Expect(res).ToNot(BeNil())
			})
		})

		When("missing track or artist name", func() {
			It("should return ErrMissingTrackOrArtistName", func() {
				lyrics, err := lrclibClient.GetLyrics(
					ctx,
					lrclib.SearchLyricsOptions{},
					263,
				)

				Expect(lyrics).To(BeNil())
				Expect(err).To(MatchError(lrclib.ErrMissingTrackOrArtistName))
			})
		})

		When("invalid duration is provided", func() {
			It("should return ErrInvalidDuration", func() {
				lyrics, err := lrclibClient.GetLyrics(
					ctx,
					lrclib.WithTrackAndArtistName("Impose", "Bad Omens"),
					0,
				)

				Expect(lyrics).To(BeNil())
				Expect(err).To(MatchError(lrclib.ErrInvalidDuration))
			})
		})

		// When("the API returns an HTTP error status", func() {
		// 	It("should return an error indicating the API failure", func() {
		// 		Expect(res).To(BeNil())
		// 		Expect(err).To(HaveOccurred())
		// 		Expect(err.Error()).To(ContainSubstring("LRCLib API error (500): Internal Server Error"))
		// 	})
		// })

		// When("the API returns a malformed response", func() {
		// 	It("should return an unmarshaling error", func() {})
		// })
	})

	Context("GetLyricsByID", func() {
		TakeMeBackToEdenLrcLibID := "2288586"
		EmptyLrcLibID := ""

		BeforeEach(func() {
			lrclibClient = lrclib.NewLRCLibProvider()
		})

		When("provided with a valid ID", func() {
			It("should return lyrics without error", func() {
				res, err := lrclibClient.GetLyricsByID(ctx, TakeMeBackToEdenLrcLibID)

				Expect(err).ToNot(HaveOccurred())
				Expect(res).ToNot(BeNil())
			})
		})

		When("provided with an empty ID", func() {
			It("should return ErrMissingID", func() {
				res, err := lrclibClient.GetLyricsByID(ctx, EmptyLrcLibID)

				Expect(res).To(BeNil())
				Expect(err).To(MatchError(lrclib.ErrMissingID))
			})
		})

		When("the API returns an HTTP error status", func() {
			It("should return an error indicating the API failure", func() {
				// Logic for API HTTP error
			})
		})

		When("the API returns a malformed response", func() {
			It("should return an unmarshaling error", func() {
				// Logic for malformed response
			})
		})
	})
})
