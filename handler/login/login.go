package login

import (
	"log"
	"net/http"
	"time"

	"github.com/golang-common-packages/template/config"
	"github.com/golang-common-packages/template/model"
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
	e.POST("/login", h.login(), h.Monitor.Middleware())
}

func (h *Handler) login() echo.HandlerFunc {
	return func(c echo.Context) error {
		// Map request body to struct
		requestBody := new(model.LoginInfo)
		if err := c.Bind(requestBody); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}
		// Query user from DB by username
		user, err := h.Database.GetUser(requestBody.Username)
		if err != nil {
			return echo.NewHTTPError(http.StatusNotFound, err)
		}
		if h.Hash.SHA512(requestBody.Password) == *user.Password && user.IsActive == true {
			accessToken, refreshToken, err := h.JWT.CreateNewTokens(h.Config.Token.Accesstoken.PrivateKey, h.Config.Token.Refreshtoken.PrivateKey, user.Email, "normal", h.Config.Token.Accesstoken.JWTTimeout, true)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, err)
			}

			// Save access token to redis with pattern: (sha512(token), sha512(username+token))
			err = h.Cache.Set(h.Hash.SHA512(accessToken), h.Hash.SHA512(user.Username+accessToken), time.Hour*time.Duration(h.Config.Token.Accesstoken.JWTTimeout))
			if err != nil {
				log.Printf("Can not save access token to redis in login handler: %s", err.Error())
				return echo.NewHTTPError(http.StatusInternalServerError)
			}

			return c.JSON(http.StatusOK, echo.Map{
				"accessToken":  accessToken,
				"refreshToken": refreshToken,
				"profile":      user,
			})
		}

		return c.NoContent(http.StatusBadGateway)
	}
}
