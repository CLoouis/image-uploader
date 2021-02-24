package auth

import (
	"context"
	"net/http"
)

type (
	AuthenticationToken struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}

	AccessToken struct {
		AccessToken string `json:"access_token"`
	}

	AuthService interface {
		Authenticate(context.Context, string, string) (AuthenticationToken, error)
		Refresh(context.Context, *http.Cookie) (AccessToken, error)
	}
)

func (a AuthenticationToken) GetAccessToken() AccessToken {
	return AccessToken{a.AccessToken}
}
