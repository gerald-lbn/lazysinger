package database_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/gerald-lbn/lazysinger/config"
	"github.com/gerald-lbn/lazysinger/database"
	"github.com/gerald-lbn/lazysinger/log"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gorm.io/gorm"
)

func TestDatabase(t *testing.T) {
	log.SetLevel(log.FatalLevel)
	RegisterFailHandler(Fail)
	RunSpecs(t, "Database Suite")
}

var _ = Describe("Database", func() {
	var tempDir string
	var db *gorm.DB
	var ctx context.Context
	var sr database.SongRepository

	BeforeEach(func() {
		var err error
		tempDir, err = os.MkdirTemp("", "lazysinger-test-*")
		Expect(err).NotTo(HaveOccurred())

		config.ResetConfig()
		config.Server.General.DataPath = tempDir
		config.Server.Database.Name = "test.db"
		config.Server.Logger.Level = "error"

		db = database.Connect()
	})

	AfterEach(func() {
		if db != nil {
			Expect(database.Close(db)).ToNot(HaveOccurred())
		}
		os.RemoveAll(tempDir)
	})

	Context("Connect", func() {
		It("returns a valid database connection", func() {
			Expect(db).NotTo(BeNil())

			dbPath := filepath.Join(tempDir, "test.db")
			Expect(dbPath).To(BeARegularFile())
		})
	})

	Context("Close", func() {
		When("closing a database connection", func() {
			It("should successfully closes the database connection", func() {
				localDB := database.Connect()
				err := database.Close(localDB)
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})

	Context("Song", func() {
		Context("Criteria", func() {
			When("Creating a criteria", func() {
				It("should return false if no criteria is set", func() {
					sr := database.NewSongCriteria()

					Expect(sr).ToNot(BeNil())
					Expect(sr.IsEmpty()).To(BeTrue())
				})

				It("should return true if at least one criteria is set", func() {
					sr := database.NewSongCriteria().WithID(1)

					Expect(sr).ToNot(BeNil())
					Expect(sr.IsEmpty()).ToNot(BeTrue())
				})
			})
		})

		Context("Repository", func() {
			assertSongEquals := func(expected *database.Song, actual *database.Song) {
				Expect(actual).ToNot(BeNil())

				Expect(actual.ID).To(Equal(expected.ID))
				Expect(actual.Path).To(Equal(expected.Path))
				Expect(actual.Title).To(Equal(expected.Title))
				Expect(actual.Artist).To(Equal(expected.Artist))
				Expect(actual.Album).To(Equal(expected.Album))
				Expect(actual.HasSyncedLyrics).To(Equal(expected.HasSyncedLyrics))
				Expect(actual.HasPlainLyrics).To(Equal(expected.HasPlainLyrics))
				Expect(actual.IsInstrumental).To(Equal(expected.IsInstrumental))
			}

			BeforeEach(func() {
				ctx = context.Background()
				sr = database.NewSongRepository(ctx, db)
			})

			When("Inserting songs", func() {
				songA := &database.Song{
					Path:            "/path/to/songA.flac",
					Title:           "Title",
					Artist:          "Artist",
					Album:           "Album",
					HasSyncedLyrics: false,
					HasPlainLyrics:  false,
					IsInstrumental:  false,
				}

				songB := &database.Song{
					Path:            "/path/to/songB.flac",
					Title:           "Title 2",
					Artist:          "Artist 2",
					Album:           "Album 2",
					HasSyncedLyrics: false,
					HasPlainLyrics:  false,
					IsInstrumental:  false,
				}

				It("should return the newly inserted song with the same data and a generated ID, and no errors", func() {
					result := sr.Create(songA)

					Expect(result.Error).ToNot(HaveOccurred())
					Expect(result.Data).ToNot(BeNil())

					Expect(songA.ID).To(Equal(result.Data.ID))
					Expect(result.Data.Path).To(Equal(songA.Path))
					Expect(result.Data.Title).To(Equal(songA.Title))
					Expect(result.Data.Artist).To(Equal(songA.Artist))
					Expect(result.Data.Album).To(Equal(songA.Album))
					Expect(result.Data.HasSyncedLyrics).To(Equal(songA.HasSyncedLyrics))
					Expect(result.Data.HasPlainLyrics).To(Equal(songA.HasPlainLyrics))
					Expect(result.Data.IsInstrumental).To(Equal(songA.IsInstrumental))
					Expect(result.Data.LastScanned).To(Equal(songA.LastScanned))
				})

				It("should insert multiple unique songs without errors", func() {
					result1 := sr.Create(songA)
					Expect(result1.Error).ToNot(HaveOccurred())

					result2 := sr.Create(songB)
					Expect(result2.Error).ToNot(HaveOccurred())
				})

				It("should complain about if two songs have the same path", func() {
					result1 := sr.Create(songA)
					Expect(result1.Error).ToNot(HaveOccurred())

					result2 := sr.Create(songA)
					Expect(result2.Error).To(HaveOccurred())
				})
			})

			When("Finding a single song", func() {
				var insertedSong *database.Song
				var now time.Time

				BeforeEach(func() {
					now = time.Now().Truncate(time.Second) // Truncate for consistent time comparisons
					insertedSong = &database.Song{
						Path:            "/find/test/song.flac",
						Title:           "Findable Title",
						Artist:          "Findable Artist",
						Album:           "Findable Album",
						HasSyncedLyrics: true,
						HasPlainLyrics:  false,
						IsInstrumental:  true,
						LastScanned:     &now,
					}

					createResult := sr.Create(insertedSong)
					Expect(createResult.Error).ToNot(HaveOccurred())
					Expect(createResult.Data).ToNot(BeNil())
				})

				It("should return the song when finding by its ID", func() {
					criteria := database.NewSongCriteria().WithID(insertedSong.ID)
					result := sr.FindBy(criteria)

					Expect(result.Error).ToNot(HaveOccurred())
					Expect(result.Data).ToNot(BeNil())
					assertSongEquals(insertedSong, result.Data)
				})

				It("should return the song when finding by its path", func() {
					criteria := database.NewSongCriteria().WithPath(insertedSong.Path)
					result := sr.FindBy(criteria)

					Expect(result.Error).ToNot(HaveOccurred())
					Expect(result.Data).ToNot(BeNil())
					assertSongEquals(insertedSong, result.Data)
				})

				It("should return the song when finding by its title", func() {
					criteria := database.NewSongCriteria().WithTitle(insertedSong.Title)
					result := sr.FindBy(criteria)

					Expect(result.Error).ToNot(HaveOccurred())
					Expect(result.Data).ToNot(BeNil())
					assertSongEquals(insertedSong, result.Data)
				})

				It("should return the song when finding by its artist", func() {
					criteria := database.NewSongCriteria().WithArtist(insertedSong.Artist)
					result := sr.FindBy(criteria)

					Expect(result.Error).ToNot(HaveOccurred())
					Expect(result.Data).ToNot(BeNil())
					assertSongEquals(insertedSong, result.Data)
				})

				It("should return the song when finding by its album", func() {
					criteria := database.NewSongCriteria().WithAlbum(insertedSong.Album)
					result := sr.FindBy(criteria)

					Expect(result.Error).ToNot(HaveOccurred())
					Expect(result.Data).ToNot(BeNil())
					assertSongEquals(insertedSong, result.Data)
				})

				It("should return the song when found by multiple criteria (e.g., title and artist)", func() {
					criteria := database.NewSongCriteria().
						WithTitle(insertedSong.Title).
						WithArtist(insertedSong.Artist)
					result := sr.FindBy(criteria)

					Expect(result.Error).ToNot(HaveOccurred())
					Expect(result.Data).ToNot(BeNil())
					assertSongEquals(insertedSong, result.Data)
				})

				It("should return the song when found by boolean HasSyncedLyrics: true", func() {
					criteria := database.NewSongCriteria().WithSyncedLyrics(true)
					result := sr.FindBy(criteria)

					Expect(result.Error).ToNot(HaveOccurred())
					Expect(result.Data).ToNot(BeNil())
					assertSongEquals(insertedSong, result.Data)
				})

				It("should not return the song when found by boolean HasSyncedLyrics: false (incorrect match)", func() {
					criteria := database.NewSongCriteria().WithSyncedLyrics(false)
					result := sr.FindBy(criteria)

					Expect(result.Error).To(MatchError(database.ErrSongNotFound))
					Expect(result.Data).To(BeNil())
				})

				It("should return the song when found by boolean IsInstrumental: true", func() {
					criteria := database.NewSongCriteria().WithInstrumental(true)
					result := sr.FindBy(criteria)

					Expect(result.Error).ToNot(HaveOccurred())
					Expect(result.Data).ToNot(BeNil())
					assertSongEquals(insertedSong, result.Data)
				})
			})

			When("Finding multiple songs", func() {
				var songA, songB, songC *database.Song

				BeforeEach(func() {
					songA = &database.Song{
						Path:            "/library/artist1/album1/song_a.mp3",
						Title:           "Song A Title",
						Artist:          "Artist One",
						Album:           "Album X",
						HasSyncedLyrics: true,
						HasPlainLyrics:  false,
						IsInstrumental:  false,
					}
					songB = &database.Song{
						Path:            "/library/artist1/album1/song_b.mp3",
						Title:           "Song B Title",
						Artist:          "Artist One",
						Album:           "Album Y",
						HasSyncedLyrics: false,
						HasPlainLyrics:  true,
						IsInstrumental:  false,
					}
					songC = &database.Song{
						Path:            "/library/artist2/album2/song_c.mp3",
						Title:           "Song C Title",
						Artist:          "Artist Two",
						Album:           "Album Z",
						HasSyncedLyrics: true,
						HasPlainLyrics:  true,
						IsInstrumental:  true,
					}

					result := sr.Create(songA)
					Expect(result.Error).ToNot(HaveOccurred())
					Expect(result.Data).ToNot(BeNil())

					result = sr.Create(songB)
					Expect(result.Error).ToNot(HaveOccurred())
					Expect(result.Data).ToNot(BeNil())

					result = sr.Create(songC)
					Expect(result.Error).ToNot(HaveOccurred())
					Expect(result.Data).ToNot(BeNil())
				})

				It("should return all songs if no criteria are provided", func() {
					criteria := database.NewSongCriteria()
					result := sr.FindManyBy(criteria)

					Expect(result.Error).ToNot(HaveOccurred())
					Expect(result.Data).ToNot(BeNil())
					Expect(len(*result.Data)).To(Equal(3))
				})

				It("should return songs matching a single criteria (Artist)", func() {
					criteria := database.NewSongCriteria().WithArtist("Artist One")
					result := sr.FindManyBy(criteria)

					Expect(result.Error).ToNot(HaveOccurred())
					Expect(result.Data).ToNot(BeNil())
					Expect(len(*result.Data)).To(Equal(2))
				})

				It("should return songs matching multiple criteria (HasSyncedLyrics and Artist)", func() {
					criteria := database.NewSongCriteria().WithSyncedLyrics(true).WithArtist("Artist One")
					result := sr.FindManyBy(criteria)

					Expect(result.Error).ToNot(HaveOccurred())
					Expect(result.Data).ToNot(BeNil())
					Expect(len(*result.Data)).To(Equal(1))
				})

				It("should return songs matching boolean HasPlainLyrics: true", func() {
					criteria := database.NewSongCriteria().WithPlainLyrics(true)
					result := sr.FindManyBy(criteria)

					Expect(result.Error).ToNot(HaveOccurred())
					Expect(result.Data).ToNot(BeNil())
					Expect(len(*result.Data)).To(Equal(2))
				})

				It("should return songs matching IsInstrumental: true", func() {
					criteria := database.NewSongCriteria().WithInstrumental(true)
					result := sr.FindManyBy(criteria)

					Expect(result.Error).ToNot(HaveOccurred())
					Expect(result.Data).ToNot(BeNil())
					Expect(len(*result.Data)).To(Equal(1))
				})

				It("should return ErrSongsNotFound if no songs match the criteria", func() {
					criteria := database.NewSongCriteria().WithArtist("Non Existent Artist")
					result := sr.FindManyBy(criteria)

					Expect(result.Error).To(MatchError(database.ErrSongsNotFound))
					Expect(result.Data).To(BeNil())
				})

				It("should return ErrSongsNotFound when an invalid criteria is provided (e.g., non-existent ID range)", func() {
					criteria := database.NewSongCriteria().WithID(9999999)
					result := sr.FindManyBy(criteria)

					Expect(result.Error).To(MatchError(database.ErrSongsNotFound))
					Expect(result.Data).To(BeNil())
				})
			})

			When("Updating a song", func() {
				var originalSong *database.Song

				BeforeEach(func() {
					originalSong = &database.Song{
						Path:            "/update/test/original.mp3",
						Title:           "Original Title",
						Artist:          "Original Artist",
						Album:           "Original Album",
						HasSyncedLyrics: false,
						HasPlainLyrics:  false,
						IsInstrumental:  true,
						LastScanned:     nil,
					}
					createResult := sr.Create(originalSong)
					Expect(createResult.Error).ToNot(HaveOccurred())
					Expect(createResult.Data).ToNot(BeNil())
				})

				It("should successfully update a single field of an existing song", func() {
					updatedTitle := "New Updated Title"
					songToUpdate := *originalSong
					songToUpdate.Title = updatedTitle

					result := sr.Update(&songToUpdate)

					Expect(result.Error).ToNot(HaveOccurred())
					Expect(result.Data).ToNot(BeNil())
					Expect(result.Data.ID).To(Equal(originalSong.ID))
					Expect(result.Data.Title).To(Equal(updatedTitle))

					fetchedResult := sr.FindBy(database.NewSongCriteria().WithID(originalSong.ID))
					Expect(fetchedResult.Error).ToNot(HaveOccurred())
					Expect(fetchedResult.Data).ToNot(BeNil())
					assertSongEquals(result.Data, fetchedResult.Data)
				})

				It("should successfully update multiple fields of an existing song", func() {
					updatedPath := "/update/test/updated.mp3"
					updatedArtist := "New Artist Name"
					updatedLyrics := true

					songToUpdate := *originalSong
					songToUpdate.Path = updatedPath
					songToUpdate.Artist = updatedArtist
					songToUpdate.HasPlainLyrics = updatedLyrics

					result := sr.Update(&songToUpdate)

					Expect(result.Error).ToNot(HaveOccurred())
					Expect(result.Data).ToNot(BeNil())
					Expect(result.Data.ID).To(Equal(originalSong.ID))
					Expect(result.Data.Path).To(Equal(updatedPath))
					Expect(result.Data.Artist).To(Equal(updatedArtist))
					Expect(result.Data.HasPlainLyrics).To(Equal(updatedLyrics))

					fetchedResult := sr.FindBy(database.NewSongCriteria().WithID(originalSong.ID))
					Expect(fetchedResult.Error).ToNot(HaveOccurred())
					Expect(fetchedResult.Data).ToNot(BeNil())
					assertSongEquals(result.Data, fetchedResult.Data)
				})

				It("should return ErrSongNotFound when trying to update a non-existent song", func() {
					nonExistentSong := &database.Song{
						Path:  "/non/existent/path.mp3",
						Title: "Dummy",
					}
					nonExistentSong.ID = 999999

					result := sr.Update(nonExistentSong)

					Expect(result.Error).To(MatchError(database.ErrSongNotFound))
					Expect(result.Data).To(BeNil())
				})

				It("should return an error when trying to update a unique field to a conflicting value", func() {
					conflictingSong := &database.Song{
						Path:            "/update/test/conflicting.mp3",
						Title:           "Conflicting Song",
						Artist:          "Another Artist",
						Album:           "Another Album",
						HasSyncedLyrics: false,
						HasPlainLyrics:  false,
						IsInstrumental:  false,
					}
					createResult := sr.Create(conflictingSong)
					Expect(createResult.Error).ToNot(HaveOccurred())

					songToUpdate := *originalSong
					songToUpdate.Path = conflictingSong.Path

					result := sr.Update(&songToUpdate)

					Expect(result.Error).To(HaveOccurred())
					Expect(result.Data).To(BeNil())
				})
			})

			When("Deleting a song", func() {
				var songToDelete *database.Song

				BeforeEach(func() {
					songToDelete = &database.Song{
						Path:            "/delete/test/song.mp3",
						Title:           "Delete Me",
						Artist:          "Delete Artist",
						Album:           "Delete Album",
						HasSyncedLyrics: false,
						HasPlainLyrics:  false,
						IsInstrumental:  false,
					}
					createResult := sr.Create(songToDelete)
					Expect(createResult.Error).ToNot(HaveOccurred())
					Expect(createResult.Data).ToNot(BeNil())
				})

				It("should successfully delete an existing song", func() {
					result := sr.Delete(songToDelete)

					Expect(result.Error).ToNot(HaveOccurred())
					Expect(result.Data).ToNot(BeNil())
					Expect(result.Data.ID).To(Equal(songToDelete.ID))

					criteria := database.NewSongCriteria().WithID(songToDelete.ID)
					fetchedResult := sr.FindBy(criteria)
					Expect(fetchedResult.Error).To(MatchError(database.ErrSongNotFound))
					Expect(fetchedResult.Data).To(BeNil())
				})

				It("should return ErrSongNotFound when trying to delete a non-existent song", func() {
					nonExistentSong := &database.Song{
						Path: "/non/existent/to/delete.mp3",
					}
					nonExistentSong.ID = 999999

					result := sr.Delete(nonExistentSong)

					Expect(result.Error).To(MatchError(database.ErrSongNotFound))
					Expect(result.Data).To(BeNil())
				})

				It("should return ErrSongNotFound when trying to delete an already deleted song (second attempt)", func() {
					firstDeleteResult := sr.Delete(songToDelete)
					Expect(firstDeleteResult.Error).ToNot(HaveOccurred())
					Expect(firstDeleteResult.Data).ToNot(BeNil())

					secondDeleteResult := sr.Delete(songToDelete)
					Expect(secondDeleteResult.Error).To(MatchError(database.ErrSongNotFound))
					Expect(secondDeleteResult.Data).To(BeNil())
				})
			})
		})
	})
})
