package spotify

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/rs/zerolog/log"
)

type SpotifySearchResponse struct {
	Tracks SpotifyTracks `json:"tracks"`
}

type SpotifyTracks struct {
	Href     string               `json:"href"`
	Limit    int                  `json:"limit"`
	Next     *string              `json:"next"` // menggunakan string pointer karena Nullable
	Offset   int                  `json:"offset"`
	Previous *string              `json:"previous"` // menggunakan string pointer karena Nullable
	Total    int                  `json:"total"`
	Items    []SpotifyTrackObject `json:"items"`
}

type SpotifyTrackObject struct {
	Album    SpotifyAlbumObject    `json:"album"`
	Artists  []SpotifyArtistObject `json:"artists"`
	Explicit bool                  `json:"explicit"`
	Href     string                `json:"href"`
	ID       string                `json:"id"`
	Name     string                `json:"name"`
}

type SpotifyAlbumObject struct {
	AlbumType   string              `json:"album_type"`
	TotalTracks int                 `json:"total_tracks"`
	Images      []SpotifyAlbumImage `json:"images"`
	Name        string              `json:"name"`
}

type SpotifyAlbumImage struct {
	URL string `json:"url"`
}

type SpotifyArtistObject struct {
	Href string `json:"href"`
	Name string `json:"name"`
}

func (o *outbound) Search(ctx context.Context, query string, limit, offset int) (*SpotifySearchResponse, error) {
	params := url.Values{}
	params.Set("q", query)
	params.Set("type", "track")
	params.Set("limit", strconv.Itoa(limit))
	params.Set("offset", strconv.Itoa(offset))

	basePath := `https://api.spotify.com/v1/search`
	urlPath := fmt.Sprintf("%s?%s", basePath, params.Encode())
	req, err := http.NewRequest(http.MethodGet, urlPath, nil)
	if err != nil {
		log.Error().Err(err).Msg("error create search request for spotify")
		return nil, err
	}

	accessToken, tokenType, err := o.GetTokenDetails()
	if err != nil {
		log.Error().Err(err).Msg("error get token details")
		return nil, err
	}

	bearerToken := fmt.Sprintf("%s %s", tokenType, accessToken)
	req.Header.Set("Authorization", bearerToken)

	resp, err := o.client.Do(req)
	if err != nil {
		log.Error().Err(err).Msg("error execute search request for spotify")
		return nil, err
	}
	defer resp.Body.Close()

	var response SpotifySearchResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		log.Error().Err(err).Msg("error unmarshal search response from spotify")
		return nil, err
	}
	return &response, nil
}
