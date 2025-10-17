package database

import (
	"context"

	models "github.com/gerald-lbn/lazysinger/database/models"
	"github.com/gerald-lbn/lazysinger/log"
	"gorm.io/gorm"
)

type TrackRepository struct {
	ctx context.Context
	db  *gorm.DB
}

func NewTrackRepository(db *gorm.DB, ctx context.Context) *TrackRepository {
	return &TrackRepository{db: db, ctx: ctx}
}

func (r *TrackRepository) Create(track models.Track) RepositoryResult[models.Track] {
	log.Debug().Interface("track", track).Msg("Creating track in database")
	err := gorm.G[models.Track](r.db).Create(r.ctx, &track)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create track in database")
	} else {
		log.Debug().Interface("track", track).Msg("Created track in database")
	}
	return RepositoryResult[models.Track]{Result: track, Error: err}
}

func (r *TrackRepository) FindAll() RepositoryResult[[]models.Track] {
	log.Debug().Msg("Finding all tracks ordered by filepath")
	tracks, err := gorm.G[models.Track](r.db).Select("*").Order("file_path").Find(r.ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to find all tracks")
	} else {
		log.Debug().Int("count", len(tracks)).Msg("Found all tracks")
	}
	return RepositoryResult[[]models.Track]{Result: tracks, Error: err}
}

func (r *TrackRepository) FindByFilePath(filepath string) RepositoryResult[models.Track] {
	log.Debug().Str("filepath", filepath).Msg("Finding track by filepath")
	track, err := gorm.G[models.Track](r.db).Where("file_path = ?", filepath).First(r.ctx)
	if err != nil {
		log.Error().Err(err).Str("filepath", filepath).Msg("Failed to find track by filepath")
	} else {
		log.Debug().Interface("track", track).Msg("Found track by filepath")
	}
	return RepositoryResult[models.Track]{Result: track, Error: err}
}

func (r *TrackRepository) FindById(id uint) RepositoryResult[models.Track] {
	log.Debug().Uint("id", id).Msg("Finding track by ID")
	track, err := gorm.G[models.Track](r.db).Where("id = ?", id).First(r.ctx)
	if err != nil {
		log.Error().Err(err).Uint("id", id).Msg("Failed to find track by ID")
	} else {
		log.Debug().Interface("track", track).Msg("Found track by ID")
	}
	return RepositoryResult[models.Track]{Result: track, Error: err}
}

func (r *TrackRepository) Update(track models.Track) RepositoryResult[models.Track] {
	log.Info().Interface("track", track).Msg("Updating track in database")
	_, err := gorm.G[models.Track](r.db).Where("id = ?", track.ID).Updates(r.ctx, track)
	if err != nil {
		log.Error().Err(err).Interface("track", track).Msg("Failed to update track in database")
	} else {
		log.Info().Interface("track", track).Msg("Updated track in database")
	}
	return RepositoryResult[models.Track]{Result: track, Error: err}
}

func (r *TrackRepository) DeleteByFilePath(filepath string) RepositoryResult[models.Track] {
	log.Debug().Str("filepath", filepath).Msg("Deleting track by filepath")
	_, err := gorm.G[models.Track](r.db).Where("file_path = ?", filepath).Delete(r.ctx)
	if err != nil {
		log.Error().Err(err).Str("filepath", filepath).Msg("Failed to delete track by filepath")
	} else {
		log.Debug().Str("filepath", filepath).Msg("Deleted track by filepath")
	}
	return RepositoryResult[models.Track]{Error: err}
}
