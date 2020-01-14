package monitor

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"contrib.go.opencensus.io/exporter/prometheus"
	"github.com/labstack/echo/v4"
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"

	"github.com/golang-microservices/template/model"
)

// PGOClient manage all slack action
type PGOClient struct {
	Handler *prometheus.Exporter
}

var (
	MLatencyMs             = stats.Float64("repl/latency", "Latency", "ms")
	total_request_accepted = stats.Int64("repl/total_request_accepted", "Total request", "ms")
	error_count            = stats.Int64("repl/error_count", "Error count", "ms")

	KeyMethod, _ = tag.NewKey("method")
	KeyStatus, _ = tag.NewKey("status")
	KeyError, _  = tag.NewKey("error")
	KeyPath, _   = tag.NewKey("path")

	LatencyView = &view.View{
		Name:        "api/latency",
		Measure:     MLatencyMs,
		Description: "Latency",
		Aggregation: view.Distribution(),
		TagKeys:     []tag.Key{KeyMethod, KeyPath, KeyStatus, KeyError}}

	TotallRequest = &view.View{
		Name:        "api/total_request_accepted",
		Measure:     total_request_accepted,
		Description: "Total request",
		Aggregation: view.Count(),
	}
)

/*
	@sessionPGOMapping: Mapping between spaceName and prometheus.Exporter for singleton pattern
*/
var (
	sessionPGOMapping = make(map[string]*PGOClient)
)

// NewPGOClient function return a exporter client based on singleton pattern
func NewPGOClient(config *model.Server, configService *model.Service) MonitorStore {
	hash := configService.Hash()
	currentSession := sessionPGOMapping[hash]
	if currentSession == nil {
		currentSession = nil

		// if err := view.Register(LatencyView, BytesInView, BytesOutView); err != nil {
		if err := view.Register(LatencyView, TotallRequest); err != nil {

			panic(err)
		}

		pe, err := prometheus.NewExporter(prometheus.Options{
			Namespace: configService.PGO.SpaceName,
		})
		if err != nil {
			panic(err)
		}

		view.RegisterExporter(pe)
		log.Println("Connected to PGO Server")

		currentSession = &PGOClient{pe}
		sessionPGOMapping[hash] = currentSession
	}

	return currentSession
}

// Middleware returns a middleware that collect request data for PGO
func (p *PGOClient) Middleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if p == nil {
				next(c)
				return nil
			}
			start := time.Now()
			err := next(c)
			stop := time.Now()
			l := float64(stop.Sub(start))
			if err != nil {
				ctx, err := tag.New(context.Background(), tag.Insert(KeyMethod, c.Request().Method), tag.Insert(KeyStatus, err.Error()), tag.Insert(KeyPath, c.Path()))
				if err != nil {
					return err
				}

				stats.Record(ctx, MLatencyMs.M(l))
			} else {
				ctx, err := tag.New(context.Background(), tag.Insert(KeyMethod, c.Request().Method), tag.Insert(KeyStatus, fmt.Sprintf("%d", c.Response().Status)), tag.Insert(KeyPath, c.Path()))
				if err != nil {
					return err
				}
				bytesIn, _ := strconv.ParseInt(c.Request().Header.Get(echo.HeaderContentLength), 10, 64)
				stats.Record(ctx, MLatencyMs.M(l), total_request_accepted.M(bytesIn), error_count.M(c.Response().Size))
			}

			return err
		}
	}
}
