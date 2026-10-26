package tracks

import (
	"context"
	"reflect"
	"testing"

	"github.com/panjiasmoroart/goCatalogMusic26/internal/models/spotify"
	"github.com/panjiasmoroart/goCatalogMusic26/internal/models/trackactivities"
	spotifyRepo "github.com/panjiasmoroart/goCatalogMusic26/internal/repository/spotify"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_service_GetRecommendation(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockSpotifyOutbound := NewMockspotifyOutbound(mockCtrl)
	mockTrackActivitiesRepo := NewMocktrackActivitiesRepository(mockCtrl)

	isLikedTrue := true
	isLikedFalse := false

	type args struct {
		userID  uint
		limit   int
		trackID string
	}
	tests := []struct {
		name    string
		args    args
		want    *spotify.RecommendationResponse
		wantErr bool
		mockFn  func(args args)
	}{
		{
			name: "success",
			args: args{
				userID:  1,
				limit:   10,
				trackID: "trackID",
			},
			want: &spotify.RecommendationResponse{
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
						IsLiked:          &isLikedTrue,
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
						IsLiked:          &isLikedFalse,
					},
				},
			},
			wantErr: false,
			mockFn: func(args args) {
				mockSpotifyOutbound.EXPECT().GetRecommendation(gomock.Any(), 10, "trackID").Return(&spotifyRepo.SpotifyRecommendationResponse{
					Tracks: []spotifyRepo.SpotifyTrackObject{
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
				}, nil)

				mockTrackActivitiesRepo.EXPECT().GetBulkSpotifyIDs(gomock.Any(), uint(1), []string{"2JiDi0qAXsPwhPqA2qaKGt", "0U49xOQCWh6zDWcu5r8Why"}).
					Return(map[string]trackactivities.TrackActivity{
						"2JiDi0qAXsPwhPqA2qaKGt": {
							IsLiked: &isLikedTrue,
						},
						"0U49xOQCWh6zDWcu5r8Why": {
							IsLiked: &isLikedFalse,
						},
					}, nil)
			},
		},
		{
			name: "failed: when get bulk spotify id",
			args: args{
				userID:  1,
				limit:   10,
				trackID: "trackID",
			},
			want:    nil,
			wantErr: true,
			mockFn: func(args args) {
				mockSpotifyOutbound.EXPECT().GetRecommendation(gomock.Any(), 10, "trackID").Return(&spotifyRepo.SpotifyRecommendationResponse{
					Tracks: []spotifyRepo.SpotifyTrackObject{
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
				}, nil)

				mockTrackActivitiesRepo.EXPECT().GetBulkSpotifyIDs(gomock.Any(), uint(1), []string{"2JiDi0qAXsPwhPqA2qaKGt", "0U49xOQCWh6zDWcu5r8Why"}).
					Return(nil, assert.AnError)
			},
		},
		{
			name: "failed: when get recommendation from spotify outbound",
			args: args{
				userID:  1,
				limit:   10,
				trackID: "trackID",
			},
			// want:    nil,
			// wantErr: true,
			want: &spotify.RecommendationResponse{
				Items: []spotify.SpotifyTrackObject{},
			},
			wantErr: false,
			mockFn: func(args args) {
				mockSpotifyOutbound.EXPECT().GetRecommendation(gomock.Any(), 10, "trackID").Return(nil, assert.AnError)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn(tt.args)
			s := &service{
				spotifyOutbound:     mockSpotifyOutbound,
				trackActivitiesRepo: mockTrackActivitiesRepo,
			}
			got, err := s.GetRecommendation(context.Background(), tt.args.userID, tt.args.limit, tt.args.trackID)
			if (err != nil) != tt.wantErr {
				t.Errorf("service.GetRecommendation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("service.GetRecommendation() = %v, want %v", got, tt.want)
			}
		})
	}
}
