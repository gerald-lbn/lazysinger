package tests

import (
	"testing"

	"github.com/gerald-lbn/lyrisync/internal/lyrics"
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
	trackName := "abc"
	artistName := "abc"
	albumName := "abc"
	duration := 150

	lyrics, err := suite.lyricsProvider.Get(trackName, artistName, albumName, duration)

	assert.Empty(suite.T(), lyrics)
	assert.Error(suite.T(), err)
}

func (suite *TestSuite) TestLyricsFoundDoesNotReturnError() {
	trackName := "Everglow"
	artistName := "STARSET"
	albumName := "Vessels 2.0"
	duration := 476

	lyrics, err := suite.lyricsProvider.Get(trackName, artistName, albumName, duration)

	assert.NotEmpty(suite.T(), lyrics)
	assert.NoError(suite.T(), err)
}

func TestLyricsTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
