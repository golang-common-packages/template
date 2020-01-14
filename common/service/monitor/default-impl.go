package monitor

import (
	"github.com/labstack/echo/v4"
)

// DefaultClient manage all slack action
type DefaultClient struct{}

// NewDefaultClient function return empty struct
func NewDefaultClient() MonitorStore {
	return &DefaultClient{}
}

// Middleware returns a empty middleware for monitoring (Default)
func (n *DefaultClient) Middleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			next(c)
			return nil
		}
	}
}
