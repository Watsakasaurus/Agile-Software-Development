package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// Ping returns a ping response
// - suitable for any check not wanting to hit db backend
func Ping(c echo.Context) error {
	return c.JSON(http.StatusOK, echo.Map{"message": "pong"})
}
