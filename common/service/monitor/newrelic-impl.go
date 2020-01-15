package monitor

import (
	"fmt"
	"log"

	"github.com/labstack/echo/v4"
	nr "github.com/newrelic/go-agent"

	"github.com/golang-common-packages/template/model"
)

// RelicClient manage all slack action
type RelicClient struct {
	Session    nr.Application
	Config     *model.Server
	LicenseKey string
}

/*
	@NEWRELIC_TXN: defines the context key used to save newrelic transaction
*/
const (
	NEWRELIC_TXN = "newrelic-txn"
)

/*
	@sessionMapping: Mapping between licenseKey and nr.Application for singleton pattern
*/
var (
	sessionMapping = make(map[string]*RelicClient)
)

// NewRelicClient function return a new relic client based on singleton pattern
func NewRelicClient(config *model.Server, licenseKey string) MonitorStore {
	currentSession := sessionMapping[licenseKey]
	if currentSession == nil {
		currentSession = nil
		configRelic := nr.NewConfig(config.Name, licenseKey)
		app, err := nr.NewApplication(configRelic)
		if err != nil {
			panic(fmt.Errorf("New relic: %s", err))
		}
		log.Println("Connected to Relic Server")

		currentSession = &RelicClient{app, config, licenseKey}
		sessionMapping[licenseKey] = currentSession
	}

	return currentSession
}

// Middleware returns a middleware that collect request data for NewRelic
func (n *RelicClient) Middleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if n.Session == nil {
				next(c)
				return nil
			}

			transactionName := fmt.Sprintf("%s [%s]", c.Path(), c.Request().Method)
			txn := n.Session.StartTransaction(transactionName, c.Response().Writer, c.Request())
			defer txn.End()

			c.Set(NEWRELIC_TXN, txn)
			err := next(c)
			if err != nil {
				txn.NoticeError(err)
			}

			return err
		}
	}
}
