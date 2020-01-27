package controllers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

// Ping returns a ping response
// - suitable for any check not wanting to hit db backend
func Ping(c echo.Context) error {
	return c.JSON(http.StatusOK, echo.Map{"message": "pong"})
}

func getOptionalString(c echo.Context, param string) *string {
	val := c.QueryParam(param)
	if val != "" {
		return &val
	}
	return nil
}

func getOptionalInt(c echo.Context, param string) *int {
	val := c.QueryParam(param)
	if val != "" {
		num, err := strconv.Atoi(val)
		if err != nil {
			log.Errorf("Failed to convert optional string to int: %s", err)
		}
		return &num
	}
	return nil
}

func getOptionalFloat64(c echo.Context, param string) *float64 {
	val := c.QueryParam(param)
	if val != "" {
		num, err := strconv.ParseFloat(val, 64)
		if err != nil {
			log.Errorf("Failed to convert optional string to float64: %s", err)
		}
		return &num
	}
	return nil
}
