package tests

import (
	"testing"

	"github.com/gerald-lbn/lazysinger/internal/lyrics"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	suite.Suite
	lyricsProvider lyrics.LyricsProvider
}

func (suite *TestSuite) SetupTest() {
	suite.lyricsProvider = *lyrics.NewLyricsProvider()
}

func (suite *TestSuite) TestLyricsNotFoundReturnsError() {
	duration := 150
	lyrics, err := suite.lyricsProvider.Get(lyrics.GetParameters{
		TrackName:  "abc",
		ArtistName: "abc",
		AlbumName:  "abc",
		Duration:   &duration,
	})

	assert.Empty(suite.T(), lyrics)
	assert.Error(suite.T(), err)
}

func (suite *TestSuite) TestLyricsFoundDoesNotReturnError() {
	duration := 476
	lyrics, err := suite.lyricsProvider.Get(lyrics.GetParameters{
		TrackName:  "Everglow",
		ArtistName: "STARSET",
		AlbumName:  "Vessels 2.0",
		Duration:   &duration,
	})

	assert.NotEmpty(suite.T(), lyrics)
	assert.NoError(suite.T(), err)
}

func TestLyricsTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
