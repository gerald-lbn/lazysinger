package database

import (
	"context"
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	ErrSongNotFound  = errors.New("song not found")
	ErrSongsNotFound = errors.New("songs not found")
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
	if criteria == nil || criteria.IsEmpty() {
		return db
	}

	q := db.WithContext(ctx)

	return q.Scopes(
		withID(criteria.ID),
		withPath(criteria.Path),
		withTitle(criteria.Title),
		withArtist(criteria.Artist),
		withAlbum(criteria.Album),
		withSyncedLyrics(criteria.HasSyncedLyrics),
		withPlainLyrics(criteria.HasPlainLyrics),
		withInstrumental(criteria.IsInstrumental),
	)
}

func withID(id *uint) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if id == nil {
			return db
		}
		return db.Where(&Song{ID: *id})
	}
}

func withTitle(title *string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if title == nil {
			return db
		}
		return db.Where(&Song{Title: *title})
	}
}

func withArtist(artist *string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if artist == nil {
			return db
		}
		return db.Where(&Song{Artist: *artist})
	}
}

func withAlbum(album *string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if album == nil {
			return db
		}
		return db.Where(&Song{Album: *album})
	}
}

func withPath(path *string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if path == nil {
			return db
		}
		return db.Where(&Song{Path: *path})
	}
}

func withSyncedLyrics(hasSyncedLyrics *bool) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if hasSyncedLyrics == nil {
			return db
		}
		return db.Where(&Song{HasSyncedLyrics: *hasSyncedLyrics})
	}
}

func withPlainLyrics(hasPlainLyrics *bool) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if hasPlainLyrics == nil {
			return db
		}
		return db.Where(&Song{HasPlainLyrics: *hasPlainLyrics})
	}
}

func withInstrumental(instrumental *bool) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if instrumental == nil {
			return db
		}
		return db.Where(&Song{IsInstrumental: *instrumental})
	}
}

// Create persists a new song in the database.
func (sr *song_repository) Create(s *Song) RepositoryOperationResult[*Song] {
	result := sr.db.WithContext(sr.ctx).Create(s)
	return RepositoryOperationResult[*Song]{
		Data:  s,
		Error: result.Error,
	}
}

// FindBy retrieves a single song matching the given criteria.
func (sr *song_repository) FindBy(criteria *SongCriteria) RepositoryOperationResult[*Song] {
	var song *Song
	result := buildCriteria(sr.ctx, sr.db, criteria).First(&song)

	if criteria.ID != nil && song.ID != *criteria.ID {
		return RepositoryOperationResult[*Song]{Data: nil, Error: ErrSongNotFound}
	}
	if criteria.Path != nil && song.Path != *criteria.Path {
		return RepositoryOperationResult[*Song]{Data: nil, Error: ErrSongNotFound}
	}
	if criteria.Title != nil && song.Title != *criteria.Title {
		return RepositoryOperationResult[*Song]{Data: nil, Error: ErrSongNotFound}
	}
	if criteria.Artist != nil && song.Artist != *criteria.Artist {
		return RepositoryOperationResult[*Song]{Data: nil, Error: ErrSongNotFound}
	}
	if criteria.Album != nil && song.Album != *criteria.Album {
		return RepositoryOperationResult[*Song]{Data: nil, Error: ErrSongNotFound}
	}
	if criteria.HasSyncedLyrics != nil && song.HasSyncedLyrics != *criteria.HasSyncedLyrics {
		return RepositoryOperationResult[*Song]{Data: nil, Error: ErrSongNotFound}
	}
	if criteria.HasPlainLyrics != nil && song.HasPlainLyrics != *criteria.HasPlainLyrics {
		return RepositoryOperationResult[*Song]{Data: nil, Error: ErrSongNotFound}
	}
	if criteria.IsInstrumental != nil && song.IsInstrumental != *criteria.IsInstrumental {
		return RepositoryOperationResult[*Song]{Data: nil, Error: ErrSongNotFound}
	}

	return RepositoryOperationResult[*Song]{
		Data:  song,
		Error: result.Error,
	}
}

// FindManyBy retrieves multiple songs matching the given criteria.
func (sr *song_repository) FindManyBy(criteria *SongCriteria) RepositoryOperationResult[*[]Song] {
	var songs *[]Song
	result := buildCriteria(sr.ctx, sr.db, criteria).Find(&songs)

	if songs == nil || len(*songs) == 0 {
		return RepositoryOperationResult[*[]Song]{
			Data:  nil,
			Error: ErrSongsNotFound,
		}
	}

	return RepositoryOperationResult[*[]Song]{
		Data:  songs,
		Error: result.Error,
	}
}

// Update modifies an existing song in the database
func (sr *song_repository) Update(s *Song) RepositoryOperationResult[*Song] {
	result := sr.db.WithContext(sr.ctx).Clauses(clause.Returning{}).Updates(&s)

	if result.Error != nil {
		return RepositoryOperationResult[*Song]{
			Data:  nil,
			Error: result.Error,
		}
	}
	if result.RowsAffected == 0 {
		return RepositoryOperationResult[*Song]{
			Data:  nil,
			Error: ErrSongNotFound,
		}
	}

	return RepositoryOperationResult[*Song]{
		Data:  s,
		Error: nil,
	}
}

// Delete removes a song from the database
func (sr *song_repository) Delete(s *Song) RepositoryOperationResult[*Song] {
	result := sr.db.WithContext(sr.ctx).Clauses(clause.Returning{}).Delete(&s)

	if result.Error != nil {
		return RepositoryOperationResult[*Song]{
			Data:  nil,
			Error: result.Error,
		}
	}
	if result.RowsAffected == 0 {
		return RepositoryOperationResult[*Song]{
			Data:  nil,
			Error: ErrSongNotFound,
		}
	}

	return RepositoryOperationResult[*Song]{
		Data:  s,
		Error: nil,
	}
}
