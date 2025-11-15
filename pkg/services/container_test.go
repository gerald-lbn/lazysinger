package services_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/gerald-lbn/refrain/pkg/services"
)

var _ = Describe("Container", func() {
	When("creating a new container", func() {
		It("should initialize all services", func() {
			c := services.NewContainer()

			Expect(c.Config).ToNot(BeNil())
			Expect(c.Watcher).ToNot(BeNil())
			Expect(c.LyricsProvider).ToNot(BeNil())
			Expect(c.Database).ToNot(BeNil())
			Expect(c.Tasks).ToNot(BeNil())
			Expect(c.Web).ToNot(BeNil())
		})
	})
})
