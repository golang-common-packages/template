package bookHttpDelivery

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	RequestID string `json:"request_id,omitempty"`
}

type AppError struct {
	Err        error
	StatusCode int
	Message    string
}

func (e *AppError) Error() string {
	return e.Message
}

func ErrorHandler(err error, c echo.Context) {
	reqID := c.Response().Header().Get("X-Request-ID")
	logger := log.New("error-handler")

	// Default error response
	statusCode := http.StatusInternalServerError
	message := "Internal Server Error"

	// Handle different error types
	switch e := err.(type) {
	case *echo.HTTPError:
		statusCode = e.Code
		message = e.Message.(string)
	case *AppError:
		statusCode = e.StatusCode
		message = e.Message
		if e.Err != nil {
			logger.Errorf("RequestID: %s, Error: %v", reqID, e.Err)
		}
	default:
		// Handle other errors
		if strings.Contains(err.Error(), "not found") {
			statusCode = http.StatusNotFound
			message = "Resource not found"
		} else if strings.Contains(err.Error(), "validation") {
			statusCode = http.StatusBadRequest
			message = "Validation error"
		}
		logger.Errorf("RequestID: %s, Error: %v", reqID, err)
	}

	// Send error response
	c.JSON(statusCode, ErrorResponse{
		Code:      statusCode,
		Message:   message,
		RequestID: reqID,
	})
}
