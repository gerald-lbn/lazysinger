package config_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/gerald-lbn/refrain/config"
)

var _ = Describe("Config", func() {
	Context("GetConfig", func() {
		It("should return a config without error", func() {
			_, err := config.GetConfig()
			Expect(err).NotTo(HaveOccurred())
		})

		It("should return the correct environment when switched", func() {
			var env config.Environment = "abc"

			config.SwitchEnvironment(env)
			cfg, err := config.GetConfig()

			Expect(err).NotTo(HaveOccurred())
			Expect(cfg.App.Environment).To(Equal(env))
		})
	})
})
