package music_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestMusic(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Music Suite")
}
