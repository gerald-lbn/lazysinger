package music_test

import (
	"os"

	"github.com/gerald-lbn/lazysinger/internal/music"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("LyricsProcessor", func() {
	var (
		txtFilePath string
		lrcFilePath string
	)

	BeforeEach(func() {
		txtFilePath = "test_lyrics.txt"
		lrcFilePath = "test_lyrics.lrc"
	})

	Context("When downloading lyrics", func() {
		AfterEach(func() {
			// Clean up test files if they exist
			os.Remove(txtFilePath)
			os.Remove(lrcFilePath)
		})

		It("should create a txt file with the lyrics", func() {
			lyrics := "These are the test lyrics."

			err := music.DownloadLyrics(txtFilePath, lyrics)
			Expect(err).NotTo(HaveOccurred())

			data, err := os.ReadFile(txtFilePath)
			Expect(err).NotTo(HaveOccurred())
			Expect(string(data)).To(Equal(lyrics))
		})

		It("should create a lrc file with the lyrics", func() {
			lyrics := "[00:12.00]These are the test lyrics.\n[00:34.00]With timestamps."

			err := music.DownloadLyrics(lrcFilePath, lyrics)
			Expect(err).NotTo(HaveOccurred())

			data, err := os.ReadFile(lrcFilePath)
			Expect(err).NotTo(HaveOccurred())
			Expect(string(data)).To(Equal(lyrics))
		})
	})

	Context("When deleting lyrics", func() {
		BeforeEach(func() {
			// Create a test file to delete
			err := os.WriteFile(txtFilePath, []byte("Temporary lyrics"), 0644)
			Expect(err).NotTo(HaveOccurred())
			err = os.WriteFile(lrcFilePath, []byte("[00:12.00]These are the test lyrics.\n[00:34.00]With timestamps."), 0644)
			Expect(err).NotTo(HaveOccurred())
		})

		AfterEach(func() {
			// Clean up in case the test fails before deletion
			os.Remove(txtFilePath)
			os.Remove(lrcFilePath)
		})

		It("should delete the txt file", func() {
			err := music.DeleteLyrics(txtFilePath)
			Expect(err).NotTo(HaveOccurred())

			_, err = os.Stat(txtFilePath)
			Expect(os.IsNotExist(err)).To(BeTrue())
		})

		It("should delete the lrc file", func() {
			err := music.DeleteLyrics(lrcFilePath)
			Expect(err).NotTo(HaveOccurred())

			_, err = os.Stat(lrcFilePath)
			Expect(os.IsNotExist(err)).To(BeTrue())
		})

		It("should return an error when trying to delete a non-existent file", func() {
			err := music.DeleteLyrics("non_existent_file.txt")
			Expect(err).To(HaveOccurred())
		})
	})
})
