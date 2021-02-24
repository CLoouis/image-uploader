package jwt

import (
	"errors"
	"fmt"
	"github.com/CLoouis/image-uploader/pkg/api/user"
	"github.com/dgrijalva/jwt-go"
	"strings"
	"time"
)

type (
	Service struct {
		key           []byte
		accessExpiry  time.Duration
		refreshExpiry time.Duration
		algo          jwt.SigningMethod
	}
)

func New(algo, secret string, accessTTL, refreshTTL int) (Service, error) {
	signingMethod := jwt.GetSigningMethod(algo)
	if signingMethod == nil {
		return Service{}, fmt.Errorf("invalid jwt signing method: %s", algo)
	}

	return Service{
		key:           []byte(secret),
		accessExpiry:  time.Duration(accessTTL) * time.Minute,
		refreshExpiry: time.Duration(refreshTTL) * time.Hour * 24,
		algo:          signingMethod,
	}, nil
}

func (s Service) GenerateAccessToken(userData user.User) (string, error) {
	accessToken, err := jwt.NewWithClaims(s.algo, jwt.MapClaims{
		"id":  userData.ID,
		"exp": time.Now().Add(s.accessExpiry).Unix(),
	}).SignedString(s.key)

	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func (s Service) GenerateRefreshToken(id string) (string, error) {
	refreshToken, err := jwt.NewWithClaims(s.algo, jwt.MapClaims{
		"id":  id,
		"exp": time.Now().Add(s.refreshExpiry).Unix(),
	}).SignedString(s.key)

	if err != nil {
		return "", err
	}

	return refreshToken, nil
}

func (s Service) ParseToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if s.algo != token.Method {
			return nil, errors.New("token error")
		}
		return s.key, nil
	})
}

func (s Service) ParseAuthorizationHeader(authHeader string) (*jwt.Token, error) {
	parts := strings.SplitN(authHeader, " ", 2)
	if !(len(parts) == 2 && parts[0] == "Bearer") {
		return nil, errors.New("token error")
	}

	return s.ParseToken(parts[1])
}
