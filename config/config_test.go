package config

import (
	"os"
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
		os.Unsetenv(LOG_LEVEL)
		os.Unsetenv(MUSIC_LIBRARY)
		os.Unsetenv(SCHEDULE)
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
})
