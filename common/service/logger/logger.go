package logger

import "github.com/golang-microservices/template/model"

// Loggerstore store function in logging package
type LoggerStore interface {
	Write(p []byte) (n int, err error)
	Close()
}

const (
	FLUENT = iota
	DEFAULT
)

// NewLoggerstore function for Factory Pattern
func NewLoggerstore(loggerstoreType int, config *model.Server, configService *model.Service) LoggerStore {
	if !config.Logging {
		loggerstoreType = DEFAULT
	}

	switch loggerstoreType {
	case FLUENT:
		return NewFluentClient(configService)
	default:
		return NewDefaultClient()
	}

	return nil
}
