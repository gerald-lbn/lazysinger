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

type GetParameters struct {
	TrackName  string
	ArtistName string
	AlbumName  string
	Duration   *int
}

func (lp *LyricsProvider) Get(parameters GetParameters) (LyricsResponse, error) {
	client := http.Client{}

	trackName := strings.ReplaceAll(parameters.TrackName, " ", "+")
	artistName := strings.ReplaceAll(parameters.ArtistName, " ", "+")
	albumName := strings.ReplaceAll(parameters.AlbumName, " ", "+")

	url := fmt.Sprintf("%s/get?artist_name=%s&track_name=%s", lrcLibBaseURl, artistName, trackName)
	if albumName != "" {
		url += fmt.Sprintf("&album_name=%s", albumName)
	}
	if parameters.Duration != nil {
		url += fmt.Sprintf("&duration=%d", *parameters.Duration)
	}
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
