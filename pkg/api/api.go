package api

import (
	"context"
	"time"

	config "github.com/CLoouis/image-uploader"
	authService "github.com/CLoouis/image-uploader/pkg/api/auth/service"
	authTransport "github.com/CLoouis/image-uploader/pkg/api/auth/transport"
	userRepository "github.com/CLoouis/image-uploader/pkg/api/user/repository/mongo"
	userService "github.com/CLoouis/image-uploader/pkg/api/user/service"
	userTransport "github.com/CLoouis/image-uploader/pkg/api/user/transport"
	"github.com/CLoouis/image-uploader/pkg/utl/jwt"
	"github.com/CLoouis/image-uploader/pkg/utl/middleware/auth"
	"github.com/CLoouis/image-uploader/pkg/utl/server"
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
	db := client.Database(cfg.Database.Name)

	//init jwt
	jwtUtl, err := jwt.New(cfg.JWT.SigningAlgorithm, cfg.JWT.SecretKey, cfg.JWT.AccessTokenExpiry, cfg.JWT.RefreshTokenExpiry)
	if err != nil {
		panic(err)
	}

	//init auth middleware
	authMw := auth.Middleware(jwtUtl)

	// init user component
	userCollection := db.Collection("user")
	userRepository := userRepository.NewUserRepositoryImpl(userCollection)
	userService := userService.NewUserServiceImpl(userRepository, time.Duration(cfg.Database.Timeout))
	userTransport.NewHTTP(userService, e.Group("/user"), authMw)

	// init auth component
	authService := authService.NewAuthService(userRepository, jwtUtl, time.Duration(cfg.Database.Timeout))
	authTransport.NewHTTP(authService, cfg.Server.CookieName, cfg.JWT.RefreshTokenExpiry, e)

	server.Start(e, cfg)
	return nil
}
