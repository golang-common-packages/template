package jwt

import (
	"github.com/labstack/echo/v4"

	"github.com/golang-microservices/template/model"
)

// Storage store function in jwt package
type Storage interface {
	Middleware(c model.Accesstoken) echo.MiddlewareFunc
	RefreshTokentMiddleware(c model.Refreshtoken) echo.MiddlewareFunc
	CreateNewTokens(c model.Token, data string, tokenType string, isAdmin bool) (accessToken, refreshToken string, err error)
	GenerateAccessToken(conf model.Accesstoken, data string, tokenType string, isAdmin bool) (string, error)
	GenerateRefreshToken(conf model.Refreshtoken, data string, tokenType string, isAdmin bool) (string, error)
	Validate(scope string, c echo.Context) error
}
