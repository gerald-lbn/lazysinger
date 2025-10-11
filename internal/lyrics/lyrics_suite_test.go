package lyrics

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestLyrics(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Lyrics Suite")
}
