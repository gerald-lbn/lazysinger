package lyrics

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const (
	lrcLibBaseURl = "https://lrclib.net/api"
)

type LyricsProvider struct {
}

type LyricsResponse struct {
	ID           int     `json:"id"`
	TrackName    string  `json:"trackName"`
	ArtistName   string  `json:"artistName"`
	AlbumName    string  `json:"albumName"`
	Duration     float32 `json:"duration"`
	Instrumental bool    `json:"instrumental"`
	PlainLyrics  string  `json:"plainLyrics"`
	SyncedLyrics string  `json:"syncedLyrics"`
}

type BadLyricsResponse struct {
	Code    int    `json:"code"`
	Name    string `json:"name"`
	Message string `json:"message"`
}

type RequestChallenge struct {
	Prefix string `json:"prefix"`
	Target string `json:"target"`
}

func NewLyricsProvider() *LyricsProvider {
	return &LyricsProvider{}
}

func (lp *LyricsProvider) Get(trackName string, artistName string, albumName string, duration int) (LyricsResponse, error) {
	client := http.Client{}

	trackName = strings.ReplaceAll(trackName, " ", "+")
	artistName = strings.ReplaceAll(artistName, " ", "+")
	albumName = strings.ReplaceAll(albumName, " ", "+")

	url := fmt.Sprintf("%s/get?artist_name=%s&track_name=%s&album_name=%s&duration=%d", lrcLibBaseURl, artistName, trackName, albumName, duration)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return LyricsResponse{}, err
	}

	res, errReq := client.Do(req)
	if errReq != nil {
		return LyricsResponse{}, errReq
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, errRead := io.ReadAll(res.Body)
	if errRead != nil {
		return LyricsResponse{}, errRead
	}

	// Handle bad response
	if res.StatusCode != 200 {
		response := BadLyricsResponse{}
		errJson := json.Unmarshal(body, &response)
		if errJson != nil {
			return LyricsResponse{}, errJson
		}

		return LyricsResponse{}, fmt.Errorf("[%s] %s (Status code: %d)", response.Name, response.Message, response.Code)
	}

	lyrics := LyricsResponse{}
	errJson := json.Unmarshal(body, &lyrics)
	if errJson != nil {
		return LyricsResponse{}, errJson
	}

	return lyrics, nil
}
