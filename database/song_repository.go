package database

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Common repository errors
var (
	ErrSongNotFound     = errors.New("song not found")
	ErrSongsNotFound    = errors.New("songs not found")
	ErrLyricsNotFound   = errors.New("lyrics not found")
	ErrInvalidScanState = errors.New("invalid scan state")
)

// SongRepository defines the interface for song-related database operations.
type SongRepository interface {
	Repository[Song, *SongCriteria]
}

// NewSongRepository creates a new instance of SongRepository.
func NewSongRepository(ctx context.Context, db *gorm.DB) SongRepository {
	return &song_repository{ctx: ctx, db: db}
}

type song_repository struct {
	ctx context.Context
	db  *gorm.DB
}

// buildQuery applies the criteria filters to the query
func buildCriteria(ctx context.Context, db *gorm.DB, criteria *SongCriteria) *gorm.DB {
	if criteria == nil {
		return db
	}

	q := db.WithContext(ctx)

	// Use scopes to build the query conditionally
	return q.Scopes(
		withPath(criteria.InPath),
		withSyncedLyrics(criteria.HasSyncedLyrics),
		withPlainLyrics(criteria.HasPlainLyrics),
		withInstrumental(criteria.IsInstrumental),
		withLastScannedBefore(criteria.LastScannedBefore),
		withLastScannedAfter(criteria.LastScannedAfter),
	)
}

func withPath(path *string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(&Song{Path: *path})
	}
}

func withSyncedLyrics(hasSyncedLyrics *bool) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(&Song{HasSyncedLyrics: *hasSyncedLyrics})
	}
}

func withPlainLyrics(hasPlainLyrics *bool) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(&Song{HasPlainLyrics: *hasPlainLyrics})
	}
}

func withInstrumental(instrumental *bool) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(&Song{IsInstrumental: *instrumental})
	}
}

func withLastScannedBefore(t *time.Time) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("last_scanned < ?", *t)
	}
}

func withLastScannedAfter(t *time.Time) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("last_scanned > ?", *t)
	}
}

// Create persists a new song in the database.
func (sr *song_repository) Create(s Song) RepositoryOperationResult[Song] {
	var createdSong *Song
	result := sr.db.Model(&createdSong).WithContext(sr.ctx).Clauses(clause.Returning{}).Create(s)

	return RepositoryOperationResult[Song]{
		Data:  createdSong,
		Error: result.Error,
	}
}

// FindBy retrieves a single song matching the given criteria.
func (sr *song_repository) FindBy(criteria *SongCriteria) RepositoryOperationResult[Song] {
	var song *Song
	result := buildCriteria(sr.ctx, sr.db, criteria).First(&song)

	return RepositoryOperationResult[Song]{
		Data:  song,
		Error: result.Error,
	}
}

// FindManyBy retrieves multiple songs matching the given criteria.
func (sr *song_repository) FindManyBy(criteria *SongCriteria) RepositoryOperationResult[[]Song] {
	var songs *[]Song
	result := buildCriteria(sr.ctx, sr.db, criteria).Find(&songs)

	return RepositoryOperationResult[[]Song]{
		Data:  songs,
		Error: result.Error,
	}
}

// Update modifies an existing song in the
func (sr *song_repository) Update(s Song) RepositoryOperationResult[Song] {
	var updatedSong *Song
	result := sr.db.WithContext(sr.ctx).Model(&updatedSong).Clauses(clause.Returning{}).Updates(s)

	return RepositoryOperationResult[Song]{
		Data:  updatedSong,
		Error: result.Error,
	}
}

// Delete removes a song from the
func (sr *song_repository) Delete(s Song) RepositoryOperationResult[Song] {
	var removedSong *Song
	result := sr.db.WithContext(sr.ctx).Model(&removedSong).Clauses(clause.Returning{}).Delete(s)

	return RepositoryOperationResult[Song]{
		Data:  removedSong,
		Error: result.Error,
	}
}
