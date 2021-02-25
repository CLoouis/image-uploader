package transport

import (
	"context"
	"net/http"

	"github.com/CLoouis/image-uploader/pkg/api/image"
	"github.com/labstack/echo/v4"
)

type (
	HTTP struct {
		imageService image.ImageService
	}
)

func NewHTTP(service image.ImageService, r *echo.Group) {
	h := HTTP{imageService: service}

	r.POST("", h.saveImage)
	r.GET("/url", h.getUrl)
	r.GET("/:file_name", h.getImageByFileName)
}

func (h *HTTP) saveImage(c echo.Context) error {
	var imageData image.Image
	if err := c.Bind(&imageData); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	ctx := context.WithValue(c.Request().Context(), "id", c.Get("id"))
	err := h.imageService.SaveImageMetadata(ctx, imageData)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	return c.NoContent(http.StatusOK)
}

func (h *HTTP) getUrl(c echo.Context) error {
	result, err := h.imageService.HandleGetPresignUploadUrlRequest(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	return c.JSON(http.StatusOK, result)
}

func (h *HTTP) getImageByFileName(c echo.Context) error {
	ctx := context.WithValue(c.Request().Context(), "id", c.Get("id"))
	result, err := h.imageService.GetImageByFileName(ctx, c.Param("file_name"))
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	return c.JSON(http.StatusOK, result)
}
