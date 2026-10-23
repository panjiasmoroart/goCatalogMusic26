package tracks

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/panjiasmoroart/goCatalogMusic26/internal/middleware"
	"github.com/panjiasmoroart/goCatalogMusic26/internal/models/spotify"
)

//go:generate mockgen -source=handler.go -destination=handler_mock_test.go -package=tracks
type service interface {
	Search(ctx context.Context, query string, pageSize, pageIndex int, userID uint) (*spotify.SearchResponse, error)
}

type Handler struct {
	*gin.Engine
	service service
}

func NewHandler(api *gin.Engine, service service) *Handler {
	return &Handler{
		api,
		service,
	}
}

func (h *Handler) RegisterRoute() {
	route := h.Group("/tracks")
	route.Use(middleware.AuthMiddleware())
	route.GET("/search", h.Search)
}
