package database

import (
	"time"

	"gorm.io/gorm"
)

// Song represents a music file in the library
type Song struct {
	gorm.Model
	Path            string `gorm:"uniqueIndex;not null"`
	Title           string `gorm:"not null"`
	Artist          string `gorm:"not null"`
	Album           string `gorm:"not null"`
	HasSyncedLyrics bool   `gorm:"check:has_synced_lyrics IN (0,1)"`
	HasPlainLyrics  bool   `gorm:"check:has_plain_lyrics IN (0,1)"`
	IsInstrumental  bool   `gorm:"check:is_instrumental IN (0,1)"`
	LastScanned     time.Time
	LastError       string
}

// Migration helper function to create/update database schema
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(&Song{})
}
