package api

import "github.com/labstack/echo/v4"

func PingPong() echo.HandlerFunc {
	return func (c echo.Context) (err error)  {
		return c.String(200, "pong")
	}
}