package api

import (
	"fmt"

	"github.com/labstack/echo/v4"
)

func PingPong() echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		fmt.Printf("faf")
		return c.String(200, "pong")
	}
}
