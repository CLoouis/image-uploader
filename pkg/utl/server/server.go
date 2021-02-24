package server

import (
	conf "github.com/CLoouis/image-uploader"
	"github.com/labstack/echo/v4"
	"net/http"
)

func New() *echo.Echo {
	e := echo.New()

	e.GET("/", healthCheckServer)

	return e
}

func healthCheckServer(c echo.Context) error {
	return c.JSON(http.StatusOK, "OK")
}

func Start(e *echo.Echo, cfg *conf.Configuration) {
	serverConfig := &http.Server{Addr: cfg.Server.Port}
	e.Logger.Fatal(e.StartServer(serverConfig))
}
