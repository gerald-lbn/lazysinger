package db_test

import (
	"database/sql"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/gerald-lbn/refrain/pkg/utils/db"
)

var _ = Describe("Db", func() {
	Context("StringToNullString", func() {
		It("should convert an empty string to a null string", func() {
			Expect(db.StringToNullString("")).To(Equal(sql.NullString{}))
		})

		It("should convert a non-empty string to a non-null string", func() {
			Expect(db.StringToNullString("test")).To(Equal(sql.NullString{String: "test", Valid: true}))
		})
	})

	Context("Like", func() {
		It("should return a LIKE pattern with wildcards", func() {
			Expect(db.Like("abc")).To(Equal("%abc%"))
		})
	})
})
