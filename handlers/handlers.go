package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Handler interface {
	Bind(e *echo.Group)
}

type Validatable interface {
	Validate() error
}

func BindAndValidate[T Validatable](c echo.Context, dst T) error {
	if err := c.Bind(dst); err != nil {
		return err
	}
	if err := dst.Validate(); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return nil
}
