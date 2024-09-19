package middleware

import (
	"backend/infrastructure/rest/services"

	"github.com/labstack/echo/v4"
)

func SessionAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		accessToken := c.Request().Header.Get("Authorization")
		if !services.VerifyToken(accessToken) {
			return c.String(401, "Bad token")
		}
		return next(c)
	}
}
