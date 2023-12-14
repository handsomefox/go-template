package handlers

import (
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

type ExampleHandler struct {
	log *zerolog.Logger
}

func NewExampleHandler(logger *zerolog.Logger) Handler {
	return &ExampleHandler{
		log: logger,
	}
}

func (h *ExampleHandler) Bind(e *echo.Group) {
	e.GET("/:id", h.GetUserByID)
}

type GetUserByIDRequest struct {
	ID uuid.UUID `json:"id" param:"id"`
}

func (r *GetUserByIDRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ID, validation.Required, validation.NilOrNotEmpty, is.UUID),
	)
}

func (h *ExampleHandler) GetUserByID(c echo.Context) error {
	var r GetUserByIDRequest
	if err := BindAndValidate(c, &r); err != nil {
		h.log.Err(err).Msg("failed to bind request body")
		return err
	}

	h.log.Debug().Any("request", r).Msg("got request with body")
	return c.JSON(http.StatusOK, map[string]any{
		"id":    r.ID,
		"name":  "dude",
		"email": "dude@gmail.com",
	})
}
