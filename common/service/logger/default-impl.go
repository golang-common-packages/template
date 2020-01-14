package logger

// DefaultClient manage all slack action
type DefaultClient struct{}

// NewDefaultClient function return empty struct
func NewDefaultClient() LoggerStore {
	return &DefaultClient{}
}

// Write functions return empty write function (Default)
func (d *DefaultClient) Write(p []byte) (n int, err error) {
	return len(p), nil
}

// Close functions return empty close function (Default)
func (d *DefaultClient) Close() {

}
