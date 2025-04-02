package BookHttpMiddleware

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

// GoMiddleware represent the data-struct for middleware
type GoMiddleware struct {
	logger *log.Logger
}

// New initialize the middleware
func New() *GoMiddleware {
	return &GoMiddleware{
		logger: log.New("middleware"),
	}
}

// CORS will handle the CORS middleware
func (m *GoMiddleware) CORS(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set("Access-Control-Allow-Origin", "https://trusted-domain.com")
		c.Response().Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Response().Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Response().Header().Set("Access-Control-Max-Age", "86400")
		return next(c)
	}
}

// SecurityHeaders adds security related headers
func (m *GoMiddleware) SecurityHeaders(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set("X-Content-Type-Options", "nosniff")
		c.Response().Header().Set("X-Frame-Options", "DENY")
		c.Response().Header().Set("X-XSS-Protection", "1; mode=block")
		c.Response().Header().Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
		return next(c)
	}
}

// RequestID adds unique request ID to each request
func (m *GoMiddleware) RequestID(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		reqID := uuid.New().String()
		c.Response().Header().Set("X-Request-ID", reqID)
		c.SetRequest(c.Request().WithContext(context.WithValue(
			c.Request().Context(),
			"requestID",
			reqID,
		)))
		return next(c)
	}
}

// Logger logs incoming requests
func (m *GoMiddleware) Logger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		start := time.Now()
		err := next(c)
		stop := time.Now()

		m.logger.Infof(
			"method=%s uri=%s status=%d latency=%s request_id=%s",
			c.Request().Method,
			c.Request().URL.Path,
			c.Response().Status,
			stop.Sub(start).String(),
			c.Response().Header().Get("X-Request-ID"),
		)
		return err
	}
}
