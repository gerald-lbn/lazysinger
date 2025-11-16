package db

import "database/sql"

// StringToNullString converts a string to a sql.NullString.
func StringToNullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{}
	}
	return sql.NullString{String: s, Valid: true}
}

// Like returns a LIKE pattern for the given value.
func Like(value string) string {
	return "%" + value + "%"
}
