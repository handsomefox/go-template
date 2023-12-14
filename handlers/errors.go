package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

var (
	ErrNotAllowed   = echo.NewHTTPError(http.StatusForbidden, "user does not have required permissions")
	ErrInvalidToken = echo.NewHTTPError(http.StatusUnauthorized, "token is invalid")
	ErrExpiredToken = echo.NewHTTPError(http.StatusUnauthorized, "token is expired")
)
