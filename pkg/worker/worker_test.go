package worker_test

import (
	"context"

	"github.com/gerald-lbn/refrain/pkg/music"
	"github.com/gerald-lbn/refrain/pkg/worker"
	"github.com/hibiken/asynq"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func ptrString(s string) *string {
	return &s
}

var _ = Describe("NewDownloadLyricsTask", func() {
	var metadata music.Metadata

	BeforeEach(func() {
		metadata.Path = "/path/to/file.flac"
		metadata.Title = ptrString("Caramel")
		metadata.Artist = ptrString("Sleep Token")
		metadata.Album = ptrString("Even In Arcadia")
		metadata.Duration = 290.0
		metadata.HasPlainLyrics = true
		metadata.PlainLyricsPath = "/path/to/file.txt"
		metadata.HasSyncedLyrics = true
		metadata.SyncedLyricsPath = "/path/to/file.lrc"
	})

	When("creating a new task", func() {
		Context("with metadata", func() {
			It("should succeed", func() {
				task, err := worker.NewDownloadLyricsTask(metadata)

				Expect(err).ToNot(HaveOccurred())
				Expect(task).ToNot(BeNil())
			})
		})
	})
})

var _ = Describe("HandleDownloadLyricsTask", func() {
	var metadata music.Metadata
	var ctx context.Context

	BeforeEach(func() {
		ctx = context.Background()

		metadata.Path = "/path/to/file.flac"
		metadata.Title = ptrString("Caramel")
		metadata.Artist = ptrString("Sleep Token")
		metadata.Album = ptrString("Even In Arcadia")
		metadata.Duration = 290.0
		metadata.HasPlainLyrics = true
		metadata.PlainLyricsPath = "/path/to/file.txt"
		metadata.HasSyncedLyrics = true
		metadata.SyncedLyricsPath = "/path/to/file.lrc"
	})

	When("handling task", func() {
		Context("with a song which already has both lyrics", func() {
			It("should skip the track", func() {
				task, err := worker.NewDownloadLyricsTask(metadata)
				Expect(err).ToNot(HaveOccurred())
				Expect(task).ToNot(BeNil())

				err = worker.HandleDownloadLyricsTask(ctx, task)
				Expect(err).To(MatchError(asynq.RevokeTask))

			})
		})

		Context("with a song which has no lyrics", func() {
			BeforeEach(func() {
				metadata.HasPlainLyrics = false
				metadata.HasSyncedLyrics = false
			})

			Context("with only the title and the artist set", func() {
				BeforeEach(func() {
					metadata.Album = nil
				})

				It("should fetch lyrics and return no errors", func() {
					task, err := worker.NewDownloadLyricsTask(metadata)
					Expect(err).ToNot(HaveOccurred())
					Expect(task).ToNot(BeNil())

					err = worker.HandleDownloadLyricsTask(ctx, task)
					Expect(err).ToNot(HaveOccurred())
				})
			})

			Context("with no title, no artist and no album name", func() {
				BeforeEach(func() {
					metadata.Title = nil
					metadata.Artist = nil
					metadata.Album = nil
				})

				It("shouldn't fetch lyrics and return an error", func() {
					task, err := worker.NewDownloadLyricsTask(metadata)
					Expect(err).ToNot(HaveOccurred())
					Expect(task).ToNot(BeNil())

					err = worker.HandleDownloadLyricsTask(ctx, task)
					Expect(err).To(MatchError(asynq.SkipRetry))
				})
			})

			Context("with the title, the artist and the album name specified", func() {
				It("should fetch lyrics and return no errors", func() {
					task, err := worker.NewDownloadLyricsTask(metadata)
					Expect(err).ToNot(HaveOccurred())
					Expect(task).ToNot(BeNil())

					err = worker.HandleDownloadLyricsTask(ctx, task)
					Expect(err).ToNot(HaveOccurred())
				})
			})
		})
	})
})
