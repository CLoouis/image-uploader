package transport

import (
	"net/http"
	"time"

	"github.com/CLoouis/image-uploader/pkg/api/auth"
	"github.com/labstack/echo/v4"
)

type (
	HTTP struct {
		svc                auth.AuthService
		cookieName         string
		refreshTokenExpiry int
	}

	credentials struct {
		Email    string `json:"email" validate:"required"`
		Password string `json:"password" validate:"required"`
	}
)

func NewHTTP(svc auth.AuthService, cookieName string, refreshTokenExpiry int, e *echo.Echo) {
	h := HTTP{svc, cookieName, refreshTokenExpiry}

	e.POST("/login", h.login)
	e.POST("/refresh", h.refresh)
	e.POST("/logout", h.logout)
}

func (h HTTP) createNewHttpOnlyCookie(auth auth.AuthenticationToken) *http.Cookie {
	cookie := new(http.Cookie)
	cookie.Name = h.cookieName
	cookie.Value = auth.RefreshToken
	cookie.HttpOnly = true
	cookie.Path = "/"
	cookie.Secure = true
	cookie.SameSite = 4
	cookie.Expires = time.Now().Add(24 * time.Hour * time.Duration(h.refreshTokenExpiry))
	return cookie
}

func (h HTTP) deleteHttpOnlyCookie() *http.Cookie {
	cookie := new(http.Cookie)
	cookie.Name = h.cookieName
	cookie.Value = ""
	cookie.HttpOnly = true
	cookie.Path = "/"
	cookie.Secure = true
	cookie.SameSite = 4
	cookie.Expires = time.Now().Add(-1 * time.Hour)
	cookie.MaxAge = -1
	return cookie
}

func (h HTTP) login(c echo.Context) error {
	ctx := c.Request().Context()
	cred := new(credentials)
	err := c.Bind(cred)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	authResult, err := h.svc.Authenticate(ctx, cred.Email, cred.Password)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	c.SetCookie(h.createNewHttpOnlyCookie(authResult))

	return c.JSON(http.StatusOK, authResult.GetAccessToken())
}

func (h HTTP) refresh(c echo.Context) error {
	ctx := c.Request().Context()
	refreshCookie, err := c.Cookie(h.cookieName)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	accessToken, err := h.svc.Refresh(ctx, refreshCookie)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, accessToken)
}

func (h HTTP) logout(c echo.Context) error {
	c.SetCookie(h.deleteHttpOnlyCookie())

	return c.NoContent(http.StatusOK)
}
