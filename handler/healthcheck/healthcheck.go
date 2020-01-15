package healthcheck

import (
	"net/http"

	"github.com/golang-common-packages/template/config"
	"github.com/labstack/echo/v4"
)

// Handler manage all request and dependency
type Handler struct {
	*config.Environment
}

// New return a new Handler
func New(env *config.Environment) *Handler {
	return &Handler{env}
}

// Handler register all path to echo.Echo
func (h *Handler) Handler(e *echo.Group) {
	e.GET("/healthcheck", h.healthcheck(), h.Monitor.Middleware())
}

func (h *Handler) healthcheck() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	}
}
