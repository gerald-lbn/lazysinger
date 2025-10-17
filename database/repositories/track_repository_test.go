package database

import (
	"context"

	models "github.com/gerald-lbn/lazysinger/database/models"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var _ = Describe("TrackRepository", func() {
	var (
		db         *gorm.DB
		repository *TrackRepository
		err        error
	)

	BeforeEach(func() {
		db, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		Expect(err).NotTo(HaveOccurred())

		err = db.AutoMigrate(&models.Track{})
		Expect(err).NotTo(HaveOccurred())

		repository = NewTrackRepository(db, context.Background())
	})

	AfterEach(func() {
		sqlDB, err := db.DB()
		Expect(err).NotTo(HaveOccurred())
		err = sqlDB.Close()
		Expect(err).NotTo(HaveOccurred())
	})

	Context("when creating a track", func() {
		It("should create a new track successfully", func() {
			track := models.Track{
				FilePath:       "/path/to/song.mp3",
				Name:           "Test Song",
				Artist:         "Test Artist",
				Album:          "Test Album",
				HasPlainLyrics: true,
			}

			result := repository.Create(track)
			Expect(result.Error).NotTo(HaveOccurred())
			Expect(result.Result.ID).NotTo(BeZero())
			Expect(result.Result.FilePath).To(Equal(track.FilePath))
			Expect(result.Result.Name).To(Equal(track.Name))
			Expect(result.Result.Artist).To(Equal(track.Artist))
			Expect(result.Result.Album).To(Equal(track.Album))
			Expect(result.Result.HasPlainLyrics).To(Equal(track.HasPlainLyrics))
		})

		It("should fail when creating a track with duplicate filepath", func() {
			track := models.Track{
				FilePath: "/path/to/song.mp3",
				Name:     "Test Song",
			}

			result := repository.Create(track)
			Expect(result.Error).NotTo(HaveOccurred())

			result = repository.Create(track)
			Expect(result.Error).To(HaveOccurred())
		})
	})

	Context("when finding a track", func() {
		var testTrack models.Track

		BeforeEach(func() {
			testTrack = models.Track{
				FilePath: "/path/to/find/song.mp3",
				Name:     "Find Test Song",
				Artist:   "Find Test Artist",
				Album:    "Find Test Album",
			}
			result := repository.Create(testTrack)
			Expect(result.Error).NotTo(HaveOccurred())
		})

		It("should find an existing track by filepath", func() {
			result := repository.FindByFilePath(testTrack.FilePath)
			Expect(result.Error).NotTo(HaveOccurred())
			Expect(result.Result.FilePath).To(Equal(testTrack.FilePath))
			Expect(result.Result.Name).To(Equal(testTrack.Name))
		})

		It("should find all existing tracks", func() {
			result := repository.FindAll()
			Expect(result.Error).NotTo(HaveOccurred())
			Expect(len(result.Result)).To(BeNumerically("==", 1))
		})

		It("should return error when track not found", func() {
			result := repository.FindByFilePath("/nonexistent/path.mp3")
			Expect(result.Error).To(HaveOccurred())
		})
	})

	Context("when updating a track", func() {
		var testTrack models.Track

		BeforeEach(func() {
			testTrack = models.Track{
				FilePath: "/path/to/update/song.mp3",
				Name:     "Original Name",
				Artist:   "Original Artist",
			}
			result := repository.Create(testTrack)
			Expect(result.Error).NotTo(HaveOccurred())
			testTrack = result.Result
		})

		It("should update track details successfully", func() {
			testTrack.Name = "Updated Name"
			testTrack.HasSyncedLyrics = true

			result := repository.Update(testTrack)
			Expect(result.Error).NotTo(HaveOccurred())

			findResult := repository.FindByFilePath(testTrack.FilePath)
			Expect(findResult.Error).NotTo(HaveOccurred())
			Expect(findResult.Result.Name).To(Equal("Updated Name"))
			Expect(findResult.Result.HasSyncedLyrics).To(BeTrue())
			Expect(findResult.Result.Artist).To(Equal("Original Artist"))
		})
	})

	Context("when deleting a track", func() {
		var testTrack models.Track

		BeforeEach(func() {
			testTrack = models.Track{
				FilePath: "/path/to/delete/song.mp3",
				Name:     "Delete Test Song",
			}
			result := repository.Create(testTrack)
			Expect(result.Error).NotTo(HaveOccurred())
		})

		It("should delete an existing track", func() {
			result := repository.DeleteByFilePath(testTrack.FilePath)
			Expect(result.Error).NotTo(HaveOccurred())

			findResult := repository.FindByFilePath(testTrack.FilePath)
			Expect(findResult.Error).To(HaveOccurred())
		})

		It("should handle deletion of non-existent track", func() {
			result := repository.DeleteByFilePath("/nonexistent/path.mp3")
			Expect(result.Error).NotTo(HaveOccurred())
		})
	})
})
