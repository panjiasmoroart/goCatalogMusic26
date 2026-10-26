package spotify

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/panjiasmoroart/goCatalogMusic26/internal/configs"
	"github.com/panjiasmoroart/goCatalogMusic26/pkg/httpclient"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_outbound_GetRecommendation(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockHTTPClient := httpclient.NewMockHTTPClient(mockCtrl)

	type args struct {
		limit   int
		trackID string
	}
	tests := []struct {
		name    string
		args    args
		want    *SpotifyRecommendationResponse
		wantErr bool
		mockFn  func(args args)
	}{
		{
			name: "success",
			args: args{
				limit:   10,
				trackID: "trackID",
			},
			want: &SpotifyRecommendationResponse{
				Tracks: []SpotifyTrackObject{
					{
						Album: SpotifyAlbumObject{
							AlbumType:   "album",
							TotalTracks: 12,
							Images: []SpotifyAlbumImage{
								{
									URL: "https://i.scdn.co/image/ab67616d0000b273fdab4a163ab9f6db72c952ee",
								},
								{
									URL: "https://i.scdn.co/image/ab67616d00001e02fdab4a163ab9f6db72c952ee",
								},
								{
									URL: "https://i.scdn.co/image/ab67616d00004851fdab4a163ab9f6db72c952ee",
								},
							},
							Name: "A Night At The Opera",
						},
						Artists: []SpotifyArtistObject{
							{
								Href: "https://api.spotify.com/v1/artists/1dfeR4HaWDbWqFHLkxsg1d",
								Name: "Queen",
							},
						},
						Explicit: false,
						Href:     "https://api.spotify.com/v1/tracks/2JiDi0qAXsPwhPqA2qaKGt",
						ID:       "2JiDi0qAXsPwhPqA2qaKGt",
						Name:     "Bohemian Rhapsody",
					},
					{
						Album: SpotifyAlbumObject{
							AlbumType:   "album",
							TotalTracks: 27,
							Images: []SpotifyAlbumImage{
								{
									URL: "https://i.scdn.co/image/ab67616d0000b273545d78706aef46068cdd4056",
								},
								{
									URL: "https://i.scdn.co/image/ab67616d00001e02545d78706aef46068cdd4056",
								},
								{
									URL: "https://i.scdn.co/image/ab67616d00004851545d78706aef46068cdd4056",
								},
							},
							Name: "Queen Rock Montreal",
						},
						Artists: []SpotifyArtistObject{
							{
								Href: "https://api.spotify.com/v1/artists/1dfeR4HaWDbWqFHLkxsg1d",
								Name: "Queen",
							},
						},
						Explicit: false,
						Href:     "https://api.spotify.com/v1/tracks/0U49xOQCWh6zDWcu5r8Why",
						ID:       "0U49xOQCWh6zDWcu5r8Why",
						Name:     "Bohemian Rhapsody - Live At The Montreal Forum / November 1981",
					},
				},
			},
			wantErr: false,
			mockFn: func(args args) {
				params := url.Values{}
				params.Set("limit", strconv.Itoa(args.limit))
				params.Set("market", "ID")
				params.Set("seed_tracks", args.trackID)

				basePath := `https://api.spotify.com/v1/recommendations`
				urlPath := fmt.Sprintf("%s?%s", basePath, params.Encode())

				req, err := http.NewRequest(http.MethodGet, urlPath, nil)
				assert.NoError(t, err)

				req.Header.Set("Authorization", "Bearer accessToken")
				mockHTTPClient.EXPECT().Do(req).Return(&http.Response{
					StatusCode: 200,
					Body:       io.NopCloser(bytes.NewBufferString(recommendationResponse)),
				}, nil)
			},
		},
		{
			name: "failed",
			args: args{
				limit:   10,
				trackID: "trackID",
			},
			want:    nil,
			wantErr: true,
			mockFn: func(args args) {
				params := url.Values{}
				params.Set("limit", strconv.Itoa(args.limit))
				params.Set("market", "ID")
				params.Set("seed_tracks", args.trackID)

				basePath := `https://api.spotify.com/v1/recommendations`
				urlPath := fmt.Sprintf("%s?%s", basePath, params.Encode())

				req, err := http.NewRequest(http.MethodGet, urlPath, nil)
				assert.NoError(t, err)

				req.Header.Set("Authorization", "Bearer accessToken")
				mockHTTPClient.EXPECT().Do(req).Return(nil, assert.AnError)
			},
		},
		{
			name: "failed: 500",
			args: args{
				limit:   10,
				trackID: "trackID",
			},
			want:    nil,
			wantErr: true,
			mockFn: func(args args) {
				params := url.Values{}
				params.Set("limit", strconv.Itoa(args.limit))
				params.Set("market", "ID")
				params.Set("seed_tracks", args.trackID)

				basePath := `https://api.spotify.com/v1/recommendations`
				urlPath := fmt.Sprintf("%s?%s", basePath, params.Encode())

				req, err := http.NewRequest(http.MethodGet, urlPath, nil)
				assert.NoError(t, err)

				req.Header.Set("Authorization", "Bearer accessToken")
				mockHTTPClient.EXPECT().Do(req).Return(&http.Response{
					StatusCode: 500,
					Body:       io.NopCloser(bytes.NewBufferString(`Internal Server Error`)),
				}, nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn(tt.args)
			o := &outbound{
				cfg:         &configs.Config{},
				client:      mockHTTPClient,
				AccessToken: "accessToken",
				TokenType:   "Bearer",
				ExpiredAt:   time.Now().Add(1 * time.Hour),
			}
			got, err := o.GetRecommendation(context.Background(), tt.args.limit, tt.args.trackID)
			if (err != nil) != tt.wantErr {
				t.Errorf("outbound.GetRecommendation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("outbound.GetRecommendation() = %v, want %v", got, tt.want)
			}
		})
	}
}
