package context_test

import (
	"context"
	"errors"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	cntx "github.com/gerald-lbn/refrain/pkg/context"
	"github.com/labstack/echo/v4"
)

var _ = Describe("Context", func() {
	Describe("IsCanceledError", func() {
		It("should return false when context is not canceled", func() {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			Expect(cntx.IsCanceledError(ctx.Err())).To(BeFalse())
		})

		It("should return true when context is canceled", func() {
			ctx, cancel := context.WithCancel(context.Background())
			cancel()
			Expect(cntx.IsCanceledError(ctx.Err())).To(BeTrue())
		})

		It("should return false when context is timed out", func() {
			ctx, cancel := context.WithTimeout(context.Background(), time.Microsecond*5)
			<-ctx.Done()
			defer cancel()
			Expect(cntx.IsCanceledError(ctx.Err())).To(BeFalse())
		})

		It("should return false for non-context errors", func() {
			err := errors.New("test error")
			Expect(cntx.IsCanceledError(err)).To(BeFalse())
		})
	})

	Describe("Cache", func() {
		var (
			echoCtx   echo.Context
			key       string
			value     string
			callCount int
			callback  func(echo.Context) string
		)

		BeforeEach(func() {
			echoCtx = echo.New().NewContext(nil, nil)
			key = "testing"
			value = "hello"
			callCount = 0
			callback = func(ctx echo.Context) string {
				callCount++
				return value
			}
		})

		It("should not have a value before caching", func() {
			Expect(echoCtx.Get(key)).To(BeNil())
		})

		It("should generate and cache a value on first call", func() {
			got := cntx.Cache(echoCtx, key, callback)
			Expect(got).To(Equal(value))
			Expect(callCount).To(Equal(1))
		})

		It("should return cached value on subsequent calls", func() {
			got := cntx.Cache(echoCtx, key, callback)
			Expect(got).To(Equal(value))
			Expect(callCount).To(Equal(1))

			got = cntx.Cache(echoCtx, key, callback)
			Expect(got).To(Equal(value))
			Expect(callCount).To(Equal(1))
		})

		It("should call the callback only once for the same key", func() {
			cntx.Cache(echoCtx, key, callback)
			cntx.Cache(echoCtx, key, callback)
			cntx.Cache(echoCtx, key, callback)

			Expect(callCount).To(Equal(1))
		})

		It("should cache different values for different keys", func() {
			key1 := "key1"
			key2 := "key2"
			value1 := "value1"
			value2 := "value2"

			got1 := cntx.Cache(echoCtx, key1, func(ctx echo.Context) string {
				return value1
			})
			got2 := cntx.Cache(echoCtx, key2, func(ctx echo.Context) string {
				return value2
			})

			Expect(got1).To(Equal(value1))
			Expect(got2).To(Equal(value2))
		})
	})
})
