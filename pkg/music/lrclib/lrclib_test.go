package lrclib_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/gerald-lbn/refrain/pkg/music/lrclib"
)

type mockRoundTripper struct {
	roundTrip func(req *http.Request) (*http.Response, error)
}

func (m *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.roundTrip(req)
}

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

		When("the API returns an HTTP error status", func() {
			It("should return an error indicating the API failure", func() {
				status := http.StatusNotFound
				body := "{\"message\":\"Failed to find specified track\",\"name\":\"TrackNotFound\",\"statusCode\":404}"
				mockClient := &http.Client{
					Transport: &mockRoundTripper{
						roundTrip: func(req *http.Request) (*http.Response, error) {
							return &http.Response{
								StatusCode: status,
								Body:       io.NopCloser(bytes.NewBufferString(body)),
								Header:     make(http.Header),
							}, nil
						},
					},
				}
				lrclibClient = lrclib.NewLRCLibProvider(lrclib.WithHttpClient(mockClient))
				results, err := lrclibClient.SearchLyrics(ctx, lrclib.WithQuery("anything"))

				Expect(results).To(BeNil())
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal(fmt.Sprintf("LRCLib API request failed with status: %d. Reason: %s", status, body)))
			})
		})
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

		When("the API returns an HTTP error status", func() {
			It("should return an error indicating the API failure", func() {
				status := http.StatusInternalServerError
				body := "Internal Server Error"
				mockClient := &http.Client{
					Transport: &mockRoundTripper{
						roundTrip: func(req *http.Request) (*http.Response, error) {
							return &http.Response{
								StatusCode: status,
								Body:       io.NopCloser(bytes.NewBufferString(body)),
								Header:     make(http.Header),
							}, nil
						},
					},
				}
				lrclibClient = lrclib.NewLRCLibProvider(lrclib.WithHttpClient(mockClient))

				res, err := lrclibClient.GetLyrics(
					ctx,
					lrclib.WithTrackAndArtistName("Impose", "Bad Omens"),
					263,
				)

				Expect(res).To(BeNil())
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal(fmt.Sprintf("LRCLib API request failed with status: %d. Reason: %s", status, body)))
			})
		})
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
				status := http.StatusInternalServerError
				body := "Internal Server Error"
				mockClient := &http.Client{
					Transport: &mockRoundTripper{
						roundTrip: func(req *http.Request) (*http.Response, error) {
							return &http.Response{
								StatusCode: status,
								Body:       io.NopCloser(bytes.NewBufferString(body)),
								Header:     make(http.Header),
							}, nil
						},
					},
				}
				lrclibClient = lrclib.NewLRCLibProvider(lrclib.WithHttpClient(mockClient))

				res, err := lrclibClient.GetLyricsByID(ctx, TakeMeBackToEdenLrcLibID)

				Expect(res).To(BeNil())
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal(fmt.Sprintf("LRCLib API request failed with status: %d. Reason: %s", status, body)))
			})
		})
	})
})
