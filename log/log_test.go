package log

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Log", func() {
	Context("When setting log level from string", func() {
		It("should set the correct log level", func() {
			SetLevelString("panic")
			Expect(GetLevel()).To(Equal(PanicLevel))

			SetLevelString("fatal")
			Expect(GetLevel()).To(Equal(FatalLevel))

			SetLevelString("error")
			Expect(GetLevel()).To(Equal(ErrorLevel))

			SetLevelString("warn")
			Expect(GetLevel()).To(Equal(WarnLevel))

			SetLevelString("info")
			Expect(GetLevel()).To(Equal(InfoLevel))

			SetLevelString("debug")
			Expect(GetLevel()).To(Equal(DebugLevel))

			SetLevelString("trace")
			Expect(GetLevel()).To(Equal(TraceLevel))

			SetLevelString("unknown")
			Expect(GetLevel()).To(Equal(InfoLevel))
		})
	})

	Context("When creating a new logger", func() {
		It("should return a non-nil logger", func() {
			Expect(logger).NotTo(BeNil())
		})
	})
})
