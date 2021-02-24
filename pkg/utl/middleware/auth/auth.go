package auth

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

type TokenParser interface {
	ParseAuthorizationHeader(string) (*jwt.Token, error)
}

func Middleware(tokenParser TokenParser) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token, err := tokenParser.ParseAuthorizationHeader(c.Request().Header.Get("Authorization"))
			if err != nil || !token.Valid {
				return c.JSON(http.StatusUnauthorized, "unauthorized")
			}

			claims := token.Claims.(jwt.MapClaims)

			id := claims["id"].(string)
			exp := time.Unix(int64(claims["exp"].(float64)), 0)

			if exp.Before(time.Now()) {
				return c.JSON(http.StatusUnauthorized, "token expired")
			}

			c.Set("id", id)
			return next(c)
		}
	}
}
