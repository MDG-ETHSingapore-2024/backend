package middleware

import (
	"backend/infrastructure/rest/services"

	"github.com/labstack/echo/v4"
)

func SessionAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		walletId := c.Request().Header.Get("X-Wallet-Address")

		if !services.VerifyWalletID(walletId) {
			return c.String(401, "Invalid wallet ID")
		}

		// If valid, proceed to the next handler
		return next(c)
	}
}
