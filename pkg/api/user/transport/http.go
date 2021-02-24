package transport

import (
	"context"
	"net/http"

	"github.com/CLoouis/image-uploader/pkg/api/user"
	"github.com/labstack/echo/v4"
)

type (
	HTTP struct {
		userService user.UserService
	}
)

func NewHTTP(service user.UserService, r *echo.Group, authMiddleware echo.MiddlewareFunc) {
	h := HTTP{userService: service}

	r.POST("", h.createUser)
	r.GET("/me", h.getUserInfo, authMiddleware)
}

func (h *HTTP) createUser(c echo.Context) error {
	var userDataFromRequest *user.User
	if err := c.Bind(&userDataFromRequest); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	ctx := context.WithValue(c.Request().Context(), "id", c.Get("id"))
	createdUser, err := h.userService.Create(ctx, *userDataFromRequest)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	return c.JSON(http.StatusOK, createdUser)
}

func (h *HTTP) getUserInfo(c echo.Context) error {
	ctx := context.WithValue(c.Request().Context(), "id", c.Get("id"))
	userInfo, err := h.userService.Me(ctx)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	return c.JSON(http.StatusOK, userInfo)
}
