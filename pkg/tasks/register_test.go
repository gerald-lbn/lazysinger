package tasks_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/gerald-lbn/refrain/pkg/services"
	"github.com/gerald-lbn/refrain/pkg/tasks"
)

var _ = Describe("Register", func() {
	var (
		container *services.Container
	)

	BeforeEach(func() {
		container = services.NewContainer()
	})

	It("should register all queues", func() {
		Expect(func() {
			tasks.Register(container)
		}).ToNot(Panic())
	})
})
