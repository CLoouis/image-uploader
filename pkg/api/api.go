package api

import (
	"context"
	"time"

	config "github.com/CLoouis/image-uploader"
	authService "github.com/CLoouis/image-uploader/pkg/api/auth/service"
	authTransport "github.com/CLoouis/image-uploader/pkg/api/auth/transport"
	imageRepository "github.com/CLoouis/image-uploader/pkg/api/image/repository"
	imageService "github.com/CLoouis/image-uploader/pkg/api/image/service"
	imageTransport "github.com/CLoouis/image-uploader/pkg/api/image/transport"
	userRepository "github.com/CLoouis/image-uploader/pkg/api/user/repository/mongo"
	userService "github.com/CLoouis/image-uploader/pkg/api/user/service"
	userTransport "github.com/CLoouis/image-uploader/pkg/api/user/transport"
	"github.com/CLoouis/image-uploader/pkg/utl/jwt"
	"github.com/CLoouis/image-uploader/pkg/utl/middleware/auth"
	"github.com/CLoouis/image-uploader/pkg/utl/server"
	"github.com/CLoouis/image-uploader/pkg/utl/uploader/s3"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Start(cfg *config.Configuration) error {
	// init database connection
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.Database.URI))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	e := server.New()
	api := e.Group("/api")
	db := client.Database(cfg.Database.Name)
	timeout := time.Duration(cfg.Database.Timeout) * time.Second

	//init jwt
	jwtUtl, err := jwt.New(cfg.JWT.SigningAlgorithm, cfg.JWT.SecretKey, cfg.JWT.AccessTokenExpiry, cfg.JWT.RefreshTokenExpiry)
	if err != nil {
		panic(err)
	}

	//init auth middleware
	authMw := auth.Middleware(jwtUtl)
	api.Use(authMw)

	// init user component
	userCollection := db.Collection("user")
	userRepository := userRepository.NewUserRepositoryImpl(userCollection)
	userService := userService.NewUserServiceImpl(userRepository, timeout)
	userTransport.NewHTTP(userService, e.Group("/user"), authMw)

	// init auth component
	authService := authService.NewAuthService(userRepository, jwtUtl, timeout)
	authTransport.NewHTTP(authService, cfg.Server.CookieName, cfg.JWT.RefreshTokenExpiry, e)

	// init uploader
	uploader := s3.NewS3Uploader(cfg.CloudStorage.AccessKeyId, cfg.CloudStorage.S3Region, cfg.CloudStorage.SecretAccessKey, cfg.CloudStorage.S3Bucket, "")

	// init image component
	imageCollection := db.Collection("image")
	imageRepository := imageRepository.NewImageRepositoryImpl(imageCollection)
	imageService := imageService.NewImageServiceImpl(imageRepository, userRepository, timeout, uploader)
	imageTransport.NewHTTP(imageService, api.Group("/image"))

	server.Start(e, cfg)
	return nil
}
