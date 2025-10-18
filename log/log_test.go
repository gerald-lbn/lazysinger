package log

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestLog(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Log Suite")
}

var _ = Describe("Log", func() {
	Context("Level Management", func() {
		When("getting log level", func() {
			It("returns the current global level", func() {
				SetLevel(DebugLevel)
				Expect(GetLevel()).To(Equal(DebugLevel))
			})
		})

		When("setting level from string", func() {
			It("sets panic level correctly", func() {
				SetLevelString("panic")
				Expect(GetLevel()).To(Equal(PanicLevel))
			})

			It("sets fatal level correctly", func() {
				SetLevelString("fatal")
				Expect(GetLevel()).To(Equal(FatalLevel))
			})

			It("sets error level correctly", func() {
				SetLevelString("error")
				Expect(GetLevel()).To(Equal(ErrorLevel))
			})

			It("sets warn level correctly", func() {
				SetLevelString("warn")
				Expect(GetLevel()).To(Equal(WarnLevel))
			})

			It("sets info level correctly", func() {
				SetLevelString("info")
				Expect(GetLevel()).To(Equal(InfoLevel))
			})

			It("sets debug level correctly", func() {
				SetLevelString("debug")
				Expect(GetLevel()).To(Equal(DebugLevel))
			})

			It("sets trace level correctly", func() {
				SetLevelString("trace")
				Expect(GetLevel()).To(Equal(TraceLevel))
			})
		})

		When("setting invalid level string", func() {
			It("defaults to info level", func() {
				SetLevelString("invalid")
				Expect(GetLevel()).To(Equal(InfoLevel))
			})

			It("defaults to info level for empty string", func() {
				SetLevelString("")
				Expect(GetLevel()).To(Equal(InfoLevel))
			})
		})
	})
})
