package config

import (
	"os"
	"strconv"
	"testing"

	"github.com/gerald-lbn/lazysinger/log"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestConfig(t *testing.T) {
	log.SetLevel(log.FatalLevel)
	RegisterFailHandler(Fail)
	RunSpecs(t, "Config Suite")
}

var _ = Describe("Config", func() {
	BeforeEach(func() {
		// Clear environment variables
		os.Unsetenv(DATA_PATH)
		os.Unsetenv(DATABASE_NAME)
		os.Unsetenv(DATABASE_PURGE_SCHEDULE)
		os.Unsetenv(LOG_LEVEL)
		os.Unsetenv(MUSIC_LIBRARY)
		os.Unsetenv(SCANNING_SCHEDULE)
		os.Unsetenv(WORKER_CONCURRENCY)
		os.Unsetenv(REDIS_ADDR)

		// Reset config to clean state
		ResetConfig()
	})

	Context("LoadWithDefault", func() {
		When("environment variable is not set", func() {
			It("returns the default value", func() {
				result := LoadWithDefault("NONEXISTENT_VAR", "default_value")
				Expect(result).To(Equal("default_value"))
			})
		})

		When("environment variable is set", func() {
			It("returns the environment value", func() {
				os.Setenv(LOG_LEVEL, "debug")
				result := LoadWithDefault(LOG_LEVEL, DEFAULT_LOG_LEVEL)
				Expect(result).To(Equal("debug"))
			})

			It("returns empty string when env var is empty", func() {
				os.Setenv(LOG_LEVEL, "")
				result := LoadWithDefault(LOG_LEVEL, DEFAULT_LOG_LEVEL)
				Expect(result).To(Equal(""))
			})
		})
	})

	Context("Setup", func() {
		When("all environment variables are at defaults", func() {
			BeforeEach(func() {
				Expect(Setup()).To(Succeed())
			})

			It("sets default general options", func() {
				Expect(Server.General.DataPath).To(Equal(DEFAULT_DATA_PATH))
				Expect(Server.Database.Name).To(Equal(DEFAULT_DATABASE_NAME))
				Expect(Server.Database.PurgeSchedule).To(Equal(DEFAULT_DATABASE_PURGE_SCHEDULE))
			})

			It("sets default logger options", func() {
				Expect(Server.Logger.Level).To(Equal(DEFAULT_LOG_LEVEL))
			})

			It("sets default scanner options", func() {
				Expect(Server.Scanner.Directory).To(Equal(DEFAULT_MUSIC_LIBRARY))
				Expect(Server.Scanner.Schedule).To(Equal(DEFAULT_SCANNING_SCHEDULE))
			})

			It("sets default worker options", func() {
				Expect(strconv.Itoa(Server.Worker.Concurrency)).To(Equal(DEFAULT_WORKER_CONCURRENCY))
				Expect(Server.Worker.RedisAddr).To(Equal(DEFAULT_REDIS_ADDR))
			})
		})

		When("environment variables are set", func() {
			BeforeEach(func() {
				os.Setenv(DATA_PATH, "/custom/data")
				os.Setenv(DATABASE_NAME, "custom.db")
				os.Setenv(DATABASE_PURGE_SCHEDULE, "0 */1 * * *")
				os.Setenv(LOG_LEVEL, "debug")
				os.Setenv(MUSIC_LIBRARY, "/custom/music")
				os.Setenv(SCANNING_SCHEDULE, "*/30 * * * *")
				os.Setenv(WORKER_CONCURRENCY, "4")
				os.Setenv(REDIS_ADDR, "redis:6379")
				Expect(Setup()).To(Succeed())
			})

			It("uses environment values for general options", func() {
				Expect(Server.General.DataPath).To(Equal("/custom/data"))
				Expect(Server.Database.Name).To(Equal("custom.db"))            // Updated struct field
				Expect(Server.Database.PurgeSchedule).To(Equal("0 */1 * * *")) // Assert new purge schedule
			})

			It("uses environment values for logger options", func() {
				Expect(Server.Logger.Level).To(Equal("debug"))
			})

			It("uses environment values for scanner options", func() {
				Expect(Server.Scanner.Directory).To(Equal("/custom/music"))
				Expect(Server.Scanner.Schedule).To(Equal("*/30 * * * *"))
			})

			It("uses environment values for worker options", func() {
				Expect(Server.Worker.Concurrency).To(Equal(4))
				Expect(Server.Worker.RedisAddr).To(Equal("redis:6379"))
			})
		})

		When("worker concurrency is invalid", func() {
			It("returns error for non-numeric value", func() {
				os.Setenv(WORKER_CONCURRENCY, "invalid")
				Expect(Setup()).To(HaveOccurred())
			})

			It("returns error for zero value", func() {
				os.Setenv(WORKER_CONCURRENCY, "0")
				err := Setup()
				Expect(err).To(HaveOccurred())
			})

			It("returns error for negative value", func() {
				os.Setenv(WORKER_CONCURRENCY, "-1")
				err := Setup()
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Context("ResetConfig", func() {
		When("Some values are set", func() {
			// Set values in a BeforeEach to ensure they are reset before each test in this context
			BeforeEach(func() {
				Server.General.DataPath = "/custom/data"
				Server.Database.Name = "custom.db"            // Updated struct field
				Server.Database.PurgeSchedule = "0 */1 * * *" // New database option
				Server.Logger.Level = "debug"
				Server.Scanner.Directory = "/custom/music"
				Server.Scanner.Schedule = "*/30 * * * *"
				Server.Worker.Concurrency = 4
				Server.Worker.RedisAddr = "redis:6379"
			})

			It("resets all configuration values", func() {
				ResetConfig()

				Expect(Server.General.DataPath).To(BeEmpty())
				Expect(Server.Database.Name).To(BeEmpty())          // Updated assertion
				Expect(Server.Database.PurgeSchedule).To(BeEmpty()) // New assertion
				Expect(Server.Logger.Level).To(BeEmpty())
				Expect(Server.Scanner.Directory).To(BeEmpty())
				Expect(Server.Scanner.Schedule).To(BeEmpty())
				Expect(Server.Worker.Concurrency).To(BeZero())
				Expect(Server.Worker.RedisAddr).To(BeEmpty())
			})
		})
	})
})
