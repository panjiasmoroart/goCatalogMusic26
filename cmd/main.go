package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/panjiasmoroart/goCatalogMusic26/internal/configs"
	membershipsHandler "github.com/panjiasmoroart/goCatalogMusic26/internal/handler/memberships"
	tracksHandler "github.com/panjiasmoroart/goCatalogMusic26/internal/handler/tracks"
	"github.com/panjiasmoroart/goCatalogMusic26/internal/models/memberships"
	"github.com/panjiasmoroart/goCatalogMusic26/internal/models/trackactivities"
	membershipsRepo "github.com/panjiasmoroart/goCatalogMusic26/internal/repository/memberships"
	"github.com/panjiasmoroart/goCatalogMusic26/internal/repository/spotify"
	trackactivitiesRepo "github.com/panjiasmoroart/goCatalogMusic26/internal/repository/trackactivities"
	membershipsSvc "github.com/panjiasmoroart/goCatalogMusic26/internal/service/memberships"
	"github.com/panjiasmoroart/goCatalogMusic26/internal/service/tracks"
	"github.com/panjiasmoroart/goCatalogMusic26/pkg/httpclient"
	"github.com/panjiasmoroart/goCatalogMusic26/pkg/internalsql"
)

func main() {
	var (
		cfg *configs.Config
	)

	err := configs.Init(
		configs.WithConfigFolder([]string{
			"./configs/",
			"./internal/configs/", // for local configs file path
		}),
		configs.WithConfigFile("config"),
		configs.WithConfigType("yaml"),
	)
	if err != nil {
		log.Fatalf("failed to initialize configs: %v", err)
	}
	cfg = configs.Get()
	log.Println("config -->>> : ", cfg)

	db, err := internalsql.Connect(cfg.Database.DataSourceName)
	if err != nil {
		log.Fatalf("failed to connect to database, err: %+v", err)
	}
	db.AutoMigrate(&memberships.User{})
	db.AutoMigrate(&trackactivities.TrackActivity{})

	r := gin.Default()

	// httpClient := httpclient.NewClient(&httpclient.Client{})
	httpClient := httpclient.NewClient(&http.Client{})

	spotifyOutbound := spotify.NewSpotifyOutbound(cfg, httpClient)

	membershipRepo := membershipsRepo.NewRepository(db)
	trackAvtivitiesRepo := trackactivitiesRepo.NewRepository(db)

	membershipSvc := membershipsSvc.NewService(cfg, membershipRepo)
	tracksSvc := tracks.NewService(spotifyOutbound, trackAvtivitiesRepo)

	membershipHandler := membershipsHandler.NewHandler(r, membershipSvc)
	membershipHandler.RegisterRoute()

	tracksHandler := tracksHandler.NewHandler(r, tracksSvc)
	tracksHandler.RegisterRoute()

	r.Run(cfg.Service.Port)
}
