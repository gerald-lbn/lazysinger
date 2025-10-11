package music

import (
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("IsMusicFile", func() {
	Context("when the file is a music file", func() {
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

		It("should return true for .mp3 files", func() {
			Expect(IsMusicFile("test_data/song.mp3")).To(BeTrue())
			Expect(IsMusicFile("test_data/song.MP3")).To(BeTrue())
		})
		It("should return true for .flac files", func() {
			Expect(IsMusicFile("test_data/song.flac")).To(BeTrue())
			Expect(IsMusicFile("test_data/song.FLAC")).To(BeTrue())
		})
		It("should return false for other file types", func() {
			Expect(IsMusicFile("test_data/document.txt")).To(BeFalse())
			Expect(IsMusicFile("test_data/image.jpg")).To(BeFalse())
			Expect(IsMusicFile("test_data/video.mp4")).To(BeFalse())
			Expect(IsMusicFile("test_data/archive.zip")).To(BeFalse())
		})
		It("should return false for non-existent files", func() {
			Expect(IsMusicFile("test_data/non_existent.mp3")).To(BeFalse())
		})
		It("should return false for directories", func() {
			Expect(IsMusicFile("test_data/test_data")).To(BeFalse())
		})
	})
})
