package tracks

import (
	"context"
	reflect "reflect"
	"testing"

	"github.com/panjiasmoroart/goCatalogMusic26/internal/models/spotify"
	spotifyRepo "github.com/panjiasmoroart/goCatalogMusic26/internal/repository/spotify"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_service_Search(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockSpotifyOutbound := NewMockspotifyOutbound(mockCtrl)

	next := "https://api.spotify.com/v1/search?offset=10&limit=10&query=bohemian%20rhapsody&type=track&market=ID&locale=en-US,en;q%3D0.8"
	// isLikedTrue := true
	// isLikedFalse := false
	type args struct {
		query     string
		pageSize  int
		pageIndex int
	}
	tests := []struct {
		name    string
		args    args
		want    *spotify.SearchResponse
		wantErr bool
		mockFn  func(args args)
	}{
		{
			name: "success",
			args: args{
				query:     "bohemian rhapsody",
				pageSize:  10,
				pageIndex: 1,
			},
			want: &spotify.SearchResponse{
				Limit:  10,
				Offset: 0,
				Items: []spotify.SpotifyTrackObject{
					{
						AlbumType:        "album",
						AlbumTotalTracks: 12,
						AlbumImagesURL:   []string{"https://i.scdn.co/image/ab67616d0000b273fdab4a163ab9f6db72c952ee", "https://i.scdn.co/image/ab67616d00001e02fdab4a163ab9f6db72c952ee", "https://i.scdn.co/image/ab67616d00004851fdab4a163ab9f6db72c952ee"},
						AlbumName:        "A Night At The Opera",
						ArtistsName:      []string{"Queen"},
						Explicit:         false,
						ID:               "2JiDi0qAXsPwhPqA2qaKGt",
						Name:             "Bohemian Rhapsody",
					},
					{
						AlbumType:        "album",
						AlbumTotalTracks: 27,
						AlbumImagesURL:   []string{"https://i.scdn.co/image/ab67616d0000b273545d78706aef46068cdd4056", "https://i.scdn.co/image/ab67616d00001e02545d78706aef46068cdd4056", "https://i.scdn.co/image/ab67616d00004851545d78706aef46068cdd4056"},
						AlbumName:        "Queen Rock Montreal",
						ArtistsName:      []string{"Queen"},
						Explicit:         false,
						ID:               "0U49xOQCWh6zDWcu5r8Why",
						Name:             "Bohemian Rhapsody - Live At The Montreal Forum / November 1981",
					},
				},
				Total: 18,
			},
			wantErr: false,
			mockFn: func(args args) {
				mockSpotifyOutbound.EXPECT().Search(gomock.Any(), args.query, 10, 0).Return(&spotifyRepo.SpotifySearchResponse{
					Tracks: spotifyRepo.SpotifyTracks{
						Href:     "https://api.spotify.com/v1/search?offset=0&limit=10&query=bohemian%20rhapsody&type=track&market=ID&locale=en-US,en;q%3D0.8",
						Limit:    10,
						Next:     &next,
						Offset:   0,
						Previous: nil,
						Total:    18,
						Items: []spotifyRepo.SpotifyTrackObject{
							{
								Album: spotifyRepo.SpotifyAlbumObject{
									AlbumType:   "album",
									TotalTracks: 12,
									Images: []spotifyRepo.SpotifyAlbumImage{
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
								Artists: []spotifyRepo.SpotifyArtistObject{
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
								Album: spotifyRepo.SpotifyAlbumObject{
									AlbumType:   "album",
									TotalTracks: 27,
									Images: []spotifyRepo.SpotifyAlbumImage{
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
								Artists: []spotifyRepo.SpotifyArtistObject{
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
				}, nil)
			},
		},
		{
			name: "failed",
			args: args{
				query:     "bohemian rhapsody",
				pageSize:  10,
				pageIndex: 1,
			},
			want:    nil,
			wantErr: true,
			mockFn: func(args args) {
				mockSpotifyOutbound.EXPECT().Search(gomock.Any(), args.query, 10, 0).Return(nil, assert.AnError)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn(tt.args)
			s := &service{
				spotifyOutbound: mockSpotifyOutbound,
			}
			got, err := s.Search(context.Background(), tt.args.query, tt.args.pageSize, tt.args.pageIndex, 1)
			if (err != nil) != tt.wantErr {
				t.Errorf("service.Search() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("service.Search() = %v, want %v", got, tt.want)
			}
		})
	}
}
