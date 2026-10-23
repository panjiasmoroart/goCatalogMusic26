package tracks

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/panjiasmoroart/goCatalogMusic26/internal/models/spotify"
	"github.com/panjiasmoroart/goCatalogMusic26/pkg/jwt"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestHandler_Search(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockSvc := NewMockservice(mockCtrl)

	tests := []struct {
		name               string
		expectedStatusCode int
		expectedBody       spotify.SearchResponse
		wantErr            bool
		mockFn             func()
	}{
		{
			name:               "success",
			expectedStatusCode: 200,
			expectedBody: spotify.SearchResponse{
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
			mockFn: func() {
				mockSvc.EXPECT().Search(gomock.Any(), "bohemian rhapsody", 10, 1, uint(1)).Return(&spotify.SearchResponse{
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
				}, nil)
			},
		},
		{
			name:               "failed",
			expectedStatusCode: 400,
			expectedBody:       spotify.SearchResponse{},
			wantErr:            true,
			mockFn: func() {
				mockSvc.EXPECT().Search(gomock.Any(), "bohemian rhapsody", 10, 1, uint(1)).Return(nil, assert.AnError)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()
			api := gin.New()

			h := &Handler{
				Engine:  api,
				service: mockSvc,
			}
			h.RegisterRoute()
			w := httptest.NewRecorder()

			endpoint := `/tracks/search?query=bohemian+rhapsody&pageSize=10&pageIndex=1`

			req, err := http.NewRequest(http.MethodGet, endpoint, nil)
			assert.NoError(t, err)
			token, err := jwt.CreateToken(1, "username", "")
			assert.NoError(t, err)
			req.Header.Set("Authorization", token)

			h.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatusCode, w.Code)

			if !tt.wantErr {
				res := w.Result()
				defer res.Body.Close()

				response := spotify.SearchResponse{}
				err = json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)

				assert.Equal(t, tt.expectedBody, response)
			}
		})
	}
}
