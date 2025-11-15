package middleware_test

import (
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/gerald-lbn/refrain/pkg/middleware"
	"github.com/gerald-lbn/refrain/pkg/session"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
)

var _ = Describe("Session Middleware", func() {
	var (
		e     *echo.Echo
		ctx   echo.Context
		store sessions.Store
	)

	BeforeEach(func() {
		e = echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		ctx = e.NewContext(req, httptest.NewRecorder())
		store = sessions.NewCookieStore([]byte("secret"))
	})

	Describe("Session", func() {
		It("should set the session store in the context", func() {
			handler := middleware.Session(store)(func(c echo.Context) error {
				_, err := session.Get(c, "test")
				Expect(err).ToNot(HaveOccurred())
				return nil
			})

			err := handler(ctx)
			Expect(err).ToNot(HaveOccurred())
		})

		It("should allow retrieving sessions after middleware execution", func() {
			var retrievedSession *sessions.Session
			var retrievalErr error

			handler := middleware.Session(store)(func(c echo.Context) error {
				retrievedSession, retrievalErr = session.Get(c, "test")
				return nil
			})

			err := handler(ctx)
			Expect(err).ToNot(HaveOccurred())
			Expect(retrievalErr).ToNot(HaveOccurred())
			Expect(retrievedSession).ToNot(BeNil())
		})
	})
})
