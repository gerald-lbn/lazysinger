package database_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gerald-lbn/lazysinger/config"
	"github.com/gerald-lbn/lazysinger/database"
	"github.com/gerald-lbn/lazysinger/log"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestDatabase(t *testing.T) {
	log.SetLevel(log.FatalLevel)
	RegisterFailHandler(Fail)
	RunSpecs(t, "Database Suite")
}

var _ = Describe("Database", func() {
	var tempDir string

	BeforeEach(func() {
		var err error
		tempDir, err = os.MkdirTemp("", "lazysinger-test-*")
		Expect(err).NotTo(HaveOccurred())

		config.ResetConfig()
		config.Server.General.DataPath = tempDir
		config.Server.General.DatabaseName = "test.db"
		config.Server.Logger.Level = "error"
	})

	AfterEach(func() {
		Expect(database.Close()).ToNot(HaveOccurred())
		os.RemoveAll(tempDir)
	})

	Context("GetInstance", func() {
		It("returns a valid database connection", func() {
			db := database.GetInstance()
			Expect(db).NotTo(BeNil())

			// Check if database file was created
			dbPath := filepath.Join(tempDir, "test.db")
			Expect(dbPath).To(BeARegularFile())
		})

		It("returns the same instance on multiple calls", func() {
			db1 := database.GetInstance()
			db2 := database.GetInstance()
			Expect(db1).To(Equal(db2))
		})
	})

	Context("Close", func() {
		It("successfully closes the database connection", func() {
			db := database.GetInstance()
			Expect(db).NotTo(BeNil())

			err := database.Close()
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
