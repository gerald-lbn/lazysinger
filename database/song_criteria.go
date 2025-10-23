package database

type SongCriteria struct {
	ID              *uint   // Filter songs by their primary key ID.
	InPath          *string // Filter songs by their file path.
	Title           *string // Filter songs by their title.
	Artist          *string // Filter songs by their artist.
	Album           *string // Filter songs by their album.
	HasSyncedLyrics *bool   // Filter songs based on presence of synchronized lyrics (true/false).
	HasPlainLyrics  *bool   // Filter songs based on presence of plain text lyrics (true/false).
	IsInstrumental  *bool   // Filter songs based on whether they are instrumental (true/false).
}

// NewSongCriteria creates a new instance of SongCriteria with no filters set.
func NewSongCriteria() *SongCriteria {
	return &SongCriteria{}
}

// IsEmpty returns true if no filter criteria have been set.
func (sc *SongCriteria) IsEmpty() bool {
	return sc.ID == nil &&
		sc.InPath == nil &&
		sc.Title == nil &&
		sc.Artist == nil &&
		sc.Album == nil &&
		sc.HasSyncedLyrics == nil &&
		sc.HasPlainLyrics == nil &&
		sc.IsInstrumental == nil
}

// WithID sets the ID filter for the criteria.
func (sc *SongCriteria) WithID(id uint) *SongCriteria {
	sc.ID = &id
	return sc
}

// WithPath sets the path filter for the criteria.
func (sc *SongCriteria) WithPath(p string) *SongCriteria {
	sc.InPath = &p
	return sc
}

// WithTitle sets the title filter for the criteria.
func (sc *SongCriteria) WithTitle(title string) *SongCriteria {
	sc.Title = &title
	return sc
}

// WithArtist sets the artist filter for the criteria.
func (sc *SongCriteria) WithArtist(artist string) *SongCriteria {
	sc.Artist = &artist
	return sc
}

// WithAlbum sets the album filter for the criteria.
func (sc *SongCriteria) WithAlbum(album string) *SongCriteria {
	sc.Album = &album
	return sc
}

// WithHasSyncedLyrics sets the synchronized lyrics filter.
func (sc *SongCriteria) WithSyncedLyrics(has bool) *SongCriteria {
	sc.HasSyncedLyrics = &has
	return sc
}

// WithHasPlainLyrics sets the plain text lyrics filter.
func (sc *SongCriteria) WithPlainLyrics(has bool) *SongCriteria {
	sc.HasPlainLyrics = &has
	return sc
}

// WithIsInstrumental sets the instrumental filter.
func (sc *SongCriteria) WithInstrumental(is bool) *SongCriteria {
	sc.IsInstrumental = &is
	return sc
}
