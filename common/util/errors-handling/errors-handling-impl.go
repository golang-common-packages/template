package errhandling

import (
	"errors"
)

// Client manage all errors handling function
type Client struct{}

// New returns an error with the given text.
func (c *Client) New(text string) error {
	return errors.New(text)
}

// DefaultErrorIfNil checks if the err is nil, if true returns the default message otherwise err.Error()
func (c *Client) DefaultErrorIfNil(err error, defaultMessage string) string {
	if err != nil {
		return err.Error()
	}
	return defaultMessage
}
