package refreshtoken

import (
	"net/http"
	"reflect"

	"github.com/labstack/echo/v4"

	"github.com/golang-common-packages/template/config"
	"github.com/golang-common-packages/template/model"
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
	e.GET("/refresh", h.refreshToken(), h.JWT.RefreshTokentMiddleware(h.Config.Token.Refreshtoken.PublicKey))
}

func (h *Handler) refreshToken() echo.HandlerFunc {
	return func(c echo.Context) error {
		// Parse query string for query user by username
		request := model.GetUserByQueryString{}
		if err := c.Bind(&request); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		result, err := h.Database.GetByField(h.Config.Service.Database.MongoDB.DB, h.Config.Service.Database.Collection.User, "username", request.Username, reflect.TypeOf(model.User{}))
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		user, ok := result.(*model.User)
		if !ok {
			return echo.NewHTTPError(http.StatusInternalServerError, "Can not bind the result to model")
		}

		if user.Username != request.Username {
			return echo.NewHTTPError(http.StatusUnauthorized)
		}

		accessToken, refreshToken, err := h.JWT.CreateNewTokens(h.Config.Token.Accesstoken.PrivateKey, h.Config.Token.Refreshtoken.PrivateKey, user.Email, "normal", h.Config.Token.Accesstoken.JWTTimeout, true)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, echo.Map{
			"accessToken":  accessToken,
			"refreshToken": refreshToken,
		})
	}
}
