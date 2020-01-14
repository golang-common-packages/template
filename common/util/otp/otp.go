package otp

// Storage store function in otp package
type Storage interface {
	NewOTP(key string) string
}
