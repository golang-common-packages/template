package user

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"gopkg.in/go-playground/validator.v9"

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

// Handler function will register all path to echo.Echo
func (h *Handler) Handler(e *echo.Group) {
	e.GET("/user", h.list(), h.JWT.Middleware(h.Config.Token.Accesstoken.PublicKey), h.Cache.Middleware(h.Hash), h.Monitor.Middleware())
	e.POST("/user/register", h.register(), h.Monitor.Middleware())
	e.GET("/user/active/:otp", h.active(), h.Monitor.Middleware())
}

// localhost:3000/api/v1/user?username=nbxtruong
// localhost:3000/api/v1/user?limit=3
// localhost:3000/api/v1/user?limit=3&lastid=5cee0e7af554bfbe838882c2
func (h *Handler) list() echo.HandlerFunc {
	return func(c echo.Context) error {
		if c.QueryParam("username") != "" {
			user, err := h.Database.GetUser(c.QueryParam("username"))
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, err)
			}
			user.Password = nil // Remove password
			return c.JSON(http.StatusOK, user)
		}

		users, err := h.Database.GetUsers(c.QueryParam("lastid"), c.QueryParam("limit"))
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		return c.JSON(http.StatusOK, users)
	}
}

func (h *Handler) register() echo.HandlerFunc {
	return func(c echo.Context) error {
		validate := validator.New()
		request := model.User{}

		// Map request body to struct
		if err := c.Bind(&request); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		// Validate request body struct
		if err := validate.Struct(request); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		*request.Password = h.Hash.SHA512(*request.Password)
		request.IsActive = false

		if err := h.Database.SaveUser(request); err != nil {
			log.Printf("Can not store to database in user hanlder: %s", err.Error())
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		otp := h.OTP.Default(h.Config.Token.OTP.SecretKey).Now()
		log.Printf("####################################")
		log.Printf("[Debug Mode] The OTP code is: " + otp)
		log.Printf("####################################")

		// Save OTP token to redis with pattern: (hash521(otp), username)
		if err := h.Cache.Set(h.Hash.SHA512(otp), request.Username, time.Duration(h.Config.Token.OTP.Timeout)*time.Minute); err != nil {
			log.Printf("Can not store to redis in user hanlder: %s", err.Error())
			return c.NoContent(http.StatusInternalServerError)
		}

		// Send OTP via email
		if err := h.Email.Send(h.Config.Service.Email.From, request.Email, h.Config.Service.Email.Subject+h.Config.Server.Name, h.Config.Service.Email.Message+otp); err != nil {
			log.Printf("Can not sent mail in user handler: %s", err.Error())
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.NoContent(http.StatusCreated)
	}
}

func (h *Handler) active() echo.HandlerFunc {
	return func(c echo.Context) error {
		otp := c.Param("otp")

		if len(otp) == 0 {
			return echo.NewHTTPError(http.StatusBadRequest, errors.New("OTP does not exist in the request"))
		}

		// Get otp code from Redis, that is stored from register API
		username, err := h.Cache.Get(h.Hash.SHA512(otp))
		if err != nil { // Not found otp on Redis
			return echo.NewHTTPError(http.StatusBadRequest, errors.New("OTP does not exist"))
		}

		if err := h.Database.ActiveUser(username); err != nil {
			log.Printf("Can not store to database in user hanlder: %s", err.Error())
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		// Delete opt in Redis
		if err = h.Cache.Delete(h.Hash.SHA512(otp)); err != nil {
			log.Printf("Can not delete redis record in user hanlder: %s", err.Error())
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return c.NoContent(http.StatusOK)
	}
}
