package music

import "go.senan.xyz/taglib"

type Metadata struct {
	Title    *string
	Artist   *string
	Album    *string
	Duration float64
}

func (m *Metadata) HasAllMetadata() bool {
	return m.Title != nil && m.Artist != nil && m.Album != nil
}
func ExtractMetadata(p string) (*Metadata, error) {
	tags, err := taglib.ReadTags(p)
	if err != nil {
		return nil, err
	}

	properties, err := taglib.ReadProperties(p)
	if err != nil {
		return nil, err
	}

	var title string
	if len(tags[taglib.Title]) > 0 {
		title = tags[taglib.Title][0]
	}

	var artist string
	if len(tags[taglib.AlbumArtist]) > 0 {
		artist = tags[taglib.AlbumArtist][0]
	}

	var album string
	if len(tags[taglib.Album]) > 0 {
		album = tags[taglib.Album][0]
	}

	return &Metadata{
		Title:    &title,
		Album:    &album,
		Artist:   &artist,
		Duration: properties.Length.Seconds(),
	}, nil
}
