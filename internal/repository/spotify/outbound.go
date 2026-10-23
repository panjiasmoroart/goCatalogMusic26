package spotify

import (
	"time"

	"github.com/panjiasmoroart/goCatalogMusic26/internal/configs"
	"github.com/panjiasmoroart/goCatalogMusic26/pkg/httpclient"
)

type outbound struct {
	cfg         *configs.Config
	client      httpclient.HTTPClient
	AccessToken string
	TokenType   string
	ExpiredAt   time.Time
}

func NewSpotifyOutbound(cfg *configs.Config, client httpclient.HTTPClient) *outbound {
	return &outbound{
		cfg:    cfg,
		client: client,
	}
}
