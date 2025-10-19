package database

import "time"

// SongCriteria defines the available search criteria for song queries.
type SongCriteria struct {
	InPath            *string    // Filter songs by their file path
	HasSyncedLyrics   *bool      // Filter songs based on presence of synchronized lyrics
	HasPlainLyrics    *bool      // Filter songs based on presence of plain text lyrics
	IsInstrumental    *bool      // Filter songs based on whether they are instrumental
	LastScannedBefore *time.Time // Filter songs scanned before this time
	LastScannedAfter  *time.Time // Filter songs scanned after this time
}

// NewSongCriteria creates a new instance of SongCriteria with no filters set.
func NewSongCriteria() *SongCriteria {
	return &SongCriteria{}
}

// WithPath sets the path filter for the criteria.
func (sc *SongCriteria) WithPath(p string) *SongCriteria {
	sc.InPath = &p
	return sc
}

// WithSyncedLyrics sets the synchronized lyrics filter.
func (sc *SongCriteria) WithSyncedLyrics(has bool) *SongCriteria {
	sc.HasSyncedLyrics = &has
	return sc
}

// WithPlainLyrics sets the plain lyrics filter.
func (sc *SongCriteria) WithPlainLyrics(has bool) *SongCriteria {
	sc.HasPlainLyrics = &has
	return sc
}

// WithInstrumental sets the instrumental filter.
func (sc *SongCriteria) WithInstrumental(is bool) *SongCriteria {
	sc.IsInstrumental = &is
	return sc
}

// WithLastScannedBefore sets the upper time boundary for when songs were last scanned.
func (sc *SongCriteria) WithLastScannedBefore(t time.Time) *SongCriteria {
	sc.LastScannedBefore = &t
	return sc
}

// WithLastScannedAfter sets the lower time boundary for when songs were last scanned.
func (sc *SongCriteria) WithLastScannedAfter(t time.Time) *SongCriteria {
	sc.LastScannedAfter = &t
	return sc
}
