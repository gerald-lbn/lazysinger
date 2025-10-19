package database

import (
	"context"
	"errors"

	"github.com/gerald-lbn/lazysinger/database"
	"gorm.io/gorm"
)

var (
	ErrSongNotFound     = errors.New("song not found")
	ErrSongsNotFound    = errors.New("songs not found")
	ErrLyricsNotFound   = errors.New("lyrics not found")
	ErrInvalidScanState = errors.New("invalid scan state")
)

func NewSongRepository(ctx context.Context, db *gorm.DB) Repository[database.Song] {
	return &song_repository{ctx: ctx, db: db}
}

type song_repository struct {
	ctx context.Context
	db  *gorm.DB
}

func (sr *song_repository) Create(s database.Song) RepositoryOperationResult[database.Song] {
	return RepositoryOperationResult[database.Song]{
		Data:  nil,
		Error: nil,
	}
}

func (sr *song_repository) FindBy(key any) RepositoryOperationResult[database.Song] {
	return RepositoryOperationResult[database.Song]{
		Data:  nil,
		Error: nil,
	}
}

func (sr *song_repository) FindManyBy(key any) RepositoryOperationResult[database.Song] {
	return RepositoryOperationResult[database.Song]{
		Data:  nil,
		Error: nil,
	}
}

func (sr *song_repository) Update(s database.Song) RepositoryOperationResult[database.Song] {
	return RepositoryOperationResult[database.Song]{
		Data:  nil,
		Error: nil,
	}
}

func (sr *song_repository) Delete(s database.Song) RepositoryOperationResult[database.Song] {
	return RepositoryOperationResult[database.Song]{
		Data:  nil,
		Error: nil,
	}
}
