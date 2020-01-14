package otp

import (
	"github.com/xlzd/gotp"
)

// OTPClient manage all OTP action
type Client struct{}

// NewOTP function will generate an otp code
func (c *Client) NewOTP(key string) string {
	return gotp.NewDefaultTOTP(key).Now()
}
