package metrics

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/golang-microservices/template/common/service/monitor"
	"github.com/golang-microservices/template/config"
)

// Handler manage all request and dependency
type Handler struct {
	*config.Environment
}

var (
	notSupportHandlerFunction = func() echo.HandlerFunc {
		return func(c echo.Context) error {
			return echo.NewHTTPError(http.StatusInternalServerError, errors.New("not support"))
		}
	}
)

// New return a new Handler
func New(env *config.Environment) *Handler {
	return &Handler{env}
}

// Handler register all path to echo.Echo
func (h *Handler) Handler(e *echo.Group) {
	var handler echo.HandlerFunc
	switch h.Environment.Monitor.(type) {
	case *monitor.PGOClient:
		pgo, _ := h.Environment.Monitor.(*monitor.PGOClient)
		handler = echo.WrapHandler(pgo.Handler)
		break
	default:
		handler = notSupportHandlerFunction()
		break
	}
	e.GET("/metrics", handler, h.Monitor.Middleware())
}
