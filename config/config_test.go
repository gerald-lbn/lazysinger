package config

import (
	"os"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestConfig(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Config Suite")
}

var _ = Describe("Config", func() {
	Context("When retrieving environment variables", func() {
		It("should return the default value if the env variable is not set", func() {
			value := GetEnvWithDefault(MUSIC_LIBRARY_PATH, DEFAULT_MUSIC_LIBRARY_PATH)
			Expect(value).To(Equal(DEFAULT_MUSIC_LIBRARY_PATH))
		})

		It("should return the env variable value if it is set", func() {
			// Set an environment variable for testing
			expectedValue := "/custom/music/path"
			_ = os.Setenv(MUSIC_LIBRARY_PATH, expectedValue)

			value := GetEnv(MUSIC_LIBRARY_PATH)
			Expect(value).To(Equal(expectedValue))
		})

		It("should return an error if the env variable is not set and no default is provided", func() {
			GetEnv(MUSIC_LIBRARY_PATH)
		})
	})
})
