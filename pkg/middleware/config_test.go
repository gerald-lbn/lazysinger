package middleware_test

import (
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/gerald-lbn/refrain/config"
	"github.com/gerald-lbn/refrain/pkg/context"
	"github.com/gerald-lbn/refrain/pkg/middleware"
	"github.com/labstack/echo/v4"
)

var _ = Describe("Config Middleware", func() {
	var (
		e   *echo.Echo
		ctx echo.Context
		cfg *config.Config
	)

	BeforeEach(func() {
		e = echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		ctx = e.NewContext(req, httptest.NewRecorder())
		cfg = &config.Config{}
	})

	Describe("Config", func() {
		It("should set the config in the context", func() {
			handler := middleware.Config(cfg)(func(c echo.Context) error {
				return nil
			})

			err := handler(ctx)
			Expect(err).NotTo(HaveOccurred())

			got, ok := ctx.Get(context.ConfigKey).(*config.Config)
			Expect(ok).To(BeTrue())
			Expect(got).To(BeIdenticalTo(cfg))
		})

		It("should make config available to the next handler", func() {
			var retrievedConfig *config.Config
			var retrievalOk bool

			handler := middleware.Config(cfg)(func(c echo.Context) error {
				retrievedConfig, retrievalOk = c.Get(context.ConfigKey).(*config.Config)
				return nil
			})

			err := handler(ctx)
			Expect(err).NotTo(HaveOccurred())
			Expect(retrievalOk).To(BeTrue())
			Expect(retrievedConfig).To(BeIdenticalTo(cfg))
		})

		It("should call the next handler", func() {
			nextCalled := false
			handler := middleware.Config(cfg)(func(c echo.Context) error {
				nextCalled = true
				return nil
			})

			err := handler(ctx)
			Expect(err).NotTo(HaveOccurred())
			Expect(nextCalled).To(BeTrue())
		})

		It("should propagate errors from the next handler", func() {
			expectedErr := echo.NewHTTPError(http.StatusInternalServerError, "test error")
			handler := middleware.Config(cfg)(func(c echo.Context) error {
				return expectedErr
			})

			err := handler(ctx)
			Expect(err).To(Equal(expectedErr))
		})
	})
})
