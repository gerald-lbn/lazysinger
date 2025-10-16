package database

import "time"

type Track struct {
	ID              uint      `gorm:"<-:create;primaryKey;autoIncrement"`
	FilePath        string    `gorm:"uniqueIndex;not null"`
	Name            string    `gorm:"not null"`
	Artist          string    `gorm:"not null"`
	Album           string    `gorm:"not null"`
	HasPlainLyrics  bool      `gorm:"default:false"`
	HasSyncedLyrics bool      `gorm:"default:false"`
	CreatedAt       time.Time `gorm:"autoCreateTime"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime"`
}
