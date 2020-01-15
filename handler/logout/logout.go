package logout

import (
	"log"
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

// Handler function will register all APIs path at login package
func (h *Handler) Handler(e *echo.Group) {
	e.GET("/logout", h.Logout(), h.Monitor.Middleware(), h.JWT.Middleware(h.Config.Token.Accesstoken.PublicKey), h.Cache.Middleware(h.Hash))
}

// Logout function will handle login request
func (h *Handler) Logout() echo.HandlerFunc {
	return func(c echo.Context) error {
		accesstoken := c.Request().Header.Get(echo.HeaderAuthorization)
		key := h.Hash.SHA512(accesstoken)

		val, err := h.Cache.Get(key)
		if err != nil {
			log.Printf("Can not delete accesstoken from redis in logout handler: %s", err.Error())
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
		if val != "" {
			if err := h.Cache.Delete(key); err != nil {
				return c.NoContent(http.StatusNotFound)
			}
		}

		return c.NoContent(http.StatusOK)
	}
}
