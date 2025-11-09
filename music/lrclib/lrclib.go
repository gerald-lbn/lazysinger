package lrclib

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

var (
	ErrInsufficientSearchParameters = errors.New("insufficient search parameters provided")
	ErrMissingTrackOrArtistName     = errors.New("track name and artist name are required")
	ErrInvalidDuration              = errors.New("duration must be a positive integer")
	ErrMissingID                    = errors.New("lyrics ID is required")
)

const (
	BASE_URL = "https://lrclib.net/api"

	ALBUM_NAME_PARAM  = "album_name"
	ARTIST_NAME_PARAM = "artist_name"
	DURATION_PARAM    = "duration"
	TRACK_NAME_PARAM  = "track_name"
)

type Lyrics struct {
	ID           int     `json:"id"`
	Name         string  `json:"name"`
	TrackName    string  `json:"trackName"`
	ArtistName   string  `json:"artistName"`
	AlbumName    string  `json:"albumName"`
	Duration     float64 `json:"duration"`
	Instrumental bool    `json:"instrumental"`
	PlainLyrics  string  `json:"plainLyrics"`
	SyncedLyrics string  `json:"syncedLyrics"`
}

type InvalidResponseError struct {
	Code    int    `json:"code"`
	Name    string `json:"name"`
	Message string `json:"message"`
}

type LRCLibProvider struct {
	BaseURL    string
	HttpClient *http.Client
}

type Option func(*LRCLibProvider)

func WithHttpClient(client *http.Client) Option {
	return func(p *LRCLibProvider) {
		p.HttpClient = client
	}
}

func NewLRCLibProvider(opts ...Option) *LRCLibProvider {
	p := &LRCLibProvider{
		BaseURL:    BASE_URL,
		HttpClient: &http.Client{},
	}

	for _, opt := range opts {
		opt(p)
	}

	return p
}

func craftLRCLibProviderURL(endpoint string, params map[string]string) string {
	if len(params) == 0 {
		return endpoint
	}

	queryParams := url.Values{}
	for key, value := range params {
		queryParams.Add(key, value)
	}

	url := endpoint + "?" + queryParams.Encode()
	log.Printf("URL: %s", url)
	return url
}

type SearchLyricsOptions struct {
	query      *string
	trackName  *string
	artistName *string
	albumName  *string
}

// WithQuery creates SearchLyricsOptions with a general query string.
func WithQuery(query string) SearchLyricsOptions {
	return SearchLyricsOptions{query: &query}
}

func WithTrackAndArtistName(trackName, artistName string) SearchLyricsOptions {
	return SearchLyricsOptions{trackName: &trackName, artistName: &artistName}
}

func WithTrackArtistAndAlbumName(trackName, artistName, albumName string) SearchLyricsOptions {
	return SearchLyricsOptions{trackName: &trackName, artistName: &artistName, albumName: &albumName}
}

func (p *LRCLibProvider) SearchLyrics(ctx context.Context, opts SearchLyricsOptions) ([]Lyrics, error) {
	endpoint := p.BaseURL + "/search"

	params := make(map[string]string)
	if opts.query != nil && *opts.query != "" {
		params["q"] = *opts.query
	} else if opts.trackName != nil && opts.artistName != nil && *opts.trackName != "" && *opts.artistName != "" {
		params[TRACK_NAME_PARAM] = *opts.trackName
		params[ARTIST_NAME_PARAM] = *opts.artistName
		if opts.albumName != nil {
			params[ALBUM_NAME_PARAM] = *opts.albumName
		}
	} else {
		return nil, ErrInsufficientSearchParameters
	}

	url := craftLRCLibProviderURL(endpoint, params)
	resp, err := p.HttpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	resBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		errorResp := &InvalidResponseError{}
		if err := json.Unmarshal(resBody, errorResp); err == nil && errorResp.Message != "" {
			return nil, fmt.Errorf("LRCLib API error (%d): %s", resp.StatusCode, errorResp.Message)
		}
		return nil, fmt.Errorf("LRCLib API search request failed with status: %d", resp.StatusCode)
	}

	var lyricsList []Lyrics
	if err := json.Unmarshal(resBody, &lyricsList); err != nil {
		return nil, fmt.Errorf("failed to unmarshal lyrics search response: %w", err)
	}

	return lyricsList, nil
}

func (p *LRCLibProvider) GetLyrics(ctx context.Context, opts SearchLyricsOptions, duration int) (*Lyrics, error) {
	var endpoint string
	endpoint = p.BaseURL + "/get"

	params := make(map[string]string)
	if opts.trackName != nil && opts.artistName != nil && *opts.trackName != "" && *opts.artistName != "" {
		params[TRACK_NAME_PARAM] = *opts.trackName
		params[ARTIST_NAME_PARAM] = *opts.artistName
		if opts.albumName != nil {
			params[ALBUM_NAME_PARAM] = *opts.albumName
		}
		if duration > 0 {
			params["duration"] = fmt.Sprintf("%d", duration)
		} else {
			return nil, ErrInvalidDuration
		}
	} else {
		return nil, ErrMissingTrackOrArtistName
	}

	url := craftLRCLibProviderURL(endpoint, params)
	resp, err := p.HttpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	resBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Check for invalid response
	if resp.StatusCode != http.StatusOK {
		errorResp := &InvalidResponseError{}
		if err := json.Unmarshal(resBody, errorResp); err == nil && errorResp.Message != "" {
			return nil, fmt.Errorf("LRCLib API error (%d): %s", resp.StatusCode, errorResp.Message)
		}
		return nil, fmt.Errorf("LRCLib API search request failed with status: %d", resp.StatusCode)
	}

	// Parse valid lyrics response
	lyrics := &Lyrics{}
	if err := json.Unmarshal(resBody, lyrics); err != nil {
		return nil, err
	}

	return lyrics, nil

}

func (p *LRCLibProvider) GetLyricsByID(ctx context.Context, id string) (*Lyrics, error) {
	if id == "" {
		return nil, ErrMissingID
	}

	endpoint := p.BaseURL + "/get/" + id
	url := craftLRCLibProviderURL(endpoint, map[string]string{})
	resp, err := p.HttpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	resBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Check for invalid response
	if resp.StatusCode != http.StatusOK {
		errorResp := &InvalidResponseError{}
		if err := json.Unmarshal(resBody, errorResp); err == nil && errorResp.Message != "" {
			return nil, fmt.Errorf("LRCLib API error (%d): %s", resp.StatusCode, errorResp.Message)
		}
		return nil, fmt.Errorf("LRCLib API search request failed with status: %d", resp.StatusCode)
	}

	// Parse valid lyrics response
	lyrics := &Lyrics{}
	if err := json.Unmarshal(resBody, lyrics); err != nil {
		return nil, err
	}

	return lyrics, nil
}
