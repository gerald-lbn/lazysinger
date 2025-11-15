package session_test

import (
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/gerald-lbn/refrain/pkg/session"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
)

var _ = Describe("Session", func() {
	var (
		e   *echo.Echo
		ctx echo.Context
	)

	BeforeEach(func() {
		e = echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		ctx = e.NewContext(req, httptest.NewRecorder())
	})

	Describe("Get", func() {
		It("should return ErrStoreNotFound when store is not set", func() {
			_, err := session.Get(ctx, "test")
			Expect(err).To(MatchError(session.ErrStoreNotFound))
		})

		It("should return a session when store is set", func() {
			store := sessions.NewCookieStore([]byte("secret"))
			session.Store(ctx, store)

			sess, err := session.Get(ctx, "test")
			Expect(err).ToNot(HaveOccurred())
			Expect(sess).ToNot(BeNil())
		})
	})

	Describe("Store", func() {
		It("should set the session store in the context", func() {
			store := sessions.NewCookieStore([]byte("secret"))
			session.Store(ctx, store)

			// Verify store was set by attempting to retrieve a session
			sess, err := session.Get(ctx, "test")
			Expect(err).ToNot(HaveOccurred())
			Expect(sess).ToNot(BeNil())
		})
	})
})
