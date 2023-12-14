package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type HealthcheckHandler struct{}

func NewHealthcheckHandler() Handler {
	return &HealthcheckHandler{}
}

func (h *HealthcheckHandler) Bind(e *echo.Group) {
	e.GET("/ping", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})
}
