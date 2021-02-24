package service

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/CLoouis/image-uploader/pkg/api/auth"
	"github.com/CLoouis/image-uploader/pkg/api/user"
	"github.com/CLoouis/image-uploader/pkg/utl/hash"
	"github.com/CLoouis/image-uploader/pkg/utl/jwt"

	jwt2 "github.com/dgrijalva/jwt-go"
)

type (
	AuthServiceImpl struct {
		userRepo       user.UserRepository
		tokenGenerator jwt.Service
		timeout        time.Duration
	}
)

func NewAuthService(userRepo user.UserRepository, tokenGenerator jwt.Service, timeout time.Duration) auth.AuthService {
	return AuthServiceImpl{userRepo, tokenGenerator, timeout}
}

func (a AuthServiceImpl) Authenticate(c context.Context, email, password string) (auth.AuthenticationToken, error) {
	ctx, cancel := context.WithTimeout(c, a.timeout)
	defer cancel()

	userData, err := a.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return auth.AuthenticationToken{}, err
	}

	valid := hash.HashMathcesPassword(userData.Password, password)
	if !valid {
		return auth.AuthenticationToken{}, errors.New("wrong password")
	}

	accessToken, err := a.tokenGenerator.GenerateAccessToken(userData)
	if err != nil {
		return auth.AuthenticationToken{}, err
	}

	refreshToken, err := a.tokenGenerator.GenerateRefreshToken(userData.ID.Hex())
	if err != nil {
		return auth.AuthenticationToken{}, err
	}

	return auth.AuthenticationToken{accessToken, refreshToken}, nil
}

func (a AuthServiceImpl) Refresh(c context.Context, refreshCookie *http.Cookie) (auth.AccessToken, error) {
	ctx, cancel := context.WithTimeout(c, a.timeout)
	defer cancel()

	token, err := a.tokenGenerator.ParseToken(refreshCookie.Value)
	if err != nil {
		return auth.AccessToken{}, err
	}

	claims := token.Claims.(jwt2.MapClaims)
	userID := claims["id"]

	userData, err := a.userRepo.FindById(ctx, userID.(string))
	if err != nil {
		return auth.AccessToken{}, errors.New("invalid token")
	}

	accessToken, err := a.tokenGenerator.GenerateAccessToken(userData)
	if err != nil {
		return auth.AccessToken{}, err
	}

	return auth.AccessToken{accessToken}, nil
}
