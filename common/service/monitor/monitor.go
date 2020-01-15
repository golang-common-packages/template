package monitor

import (
	"github.com/labstack/echo/v4"

	"github.com/golang-common-packages/template/model"
)

// MonitorStore function in newrelic package
type MonitorStore interface {
	Middleware() echo.MiddlewareFunc
}

const (
	NEWRELIC = iota
	PGO
	DEFAULT
)

// NewMonitorStore function for Factory Pattern
func NewMonitorStore(monitorType int, config *model.Server, configService *model.Service) MonitorStore {
	if !config.Monitoring {
		monitorType = DEFAULT
	}

	switch monitorType {
	case NEWRELIC:
		return NewRelicClient(config, configService.NewRelic.LicenseKey)
	case PGO:
		return NewPGOClient(config, configService)
	default:
		return NewDefaultClient()
	}

	return nil
}
