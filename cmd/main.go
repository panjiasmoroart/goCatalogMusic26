package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/panjiasmoroart/goCatalogMusic26/internal/configs"
	membershipsHandler "github.com/panjiasmoroart/goCatalogMusic26/internal/handler/memberships"
	"github.com/panjiasmoroart/goCatalogMusic26/internal/models/memberships"
	membershipsRepo "github.com/panjiasmoroart/goCatalogMusic26/internal/repository/memberships"
	membershipsSvc "github.com/panjiasmoroart/goCatalogMusic26/internal/service/memberships"
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
	log.Println("config : ", cfg)

	db, err := internalsql.Connect(cfg.Database.DataSourceName)
	if err != nil {
		log.Fatalf("failed to connect to database, err: %+v", err)
	}
	db.AutoMigrate(&memberships.User{})

	r := gin.Default()

	membershipRepo := membershipsRepo.NewRepository(db)

	membershipSvc := membershipsSvc.NewService(cfg, membershipRepo)

	membershipHandler := membershipsHandler.NewHandler(r, membershipSvc)
	membershipHandler.RegisterRoute()

	r.Run(cfg.Service.Port)
}
