package main

import (
	"context"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/golang-microservices/cloud-storage"
	"github.com/golang-microservices/echo-jwt-middleware"

	"github.com/golang-microservices/template/config"
	"github.com/golang-microservices/template/model"

	"github.com/golang-microservices/template/handler/document"
	"github.com/golang-microservices/template/handler/healthcheck"
	"github.com/golang-microservices/template/handler/login"
	"github.com/golang-microservices/template/handler/logout"
	"github.com/golang-microservices/template/handler/metrics"
	"github.com/golang-microservices/template/handler/refreshtoken"
	"github.com/golang-microservices/template/handler/user"

	"github.com/golang-microservices/template/common/service/cachestore"
	"github.com/golang-microservices/template/common/service/datastore"
	"github.com/golang-microservices/template/common/service/email"
	"github.com/golang-microservices/template/common/service/logger"
	"github.com/golang-microservices/template/common/service/monitor"

	"github.com/golang-microservices/template/common/util/apigroup"
	"github.com/golang-microservices/template/common/util/condition"
	"github.com/golang-microservices/template/common/util/hash"
	"github.com/golang-microservices/template/common/util/otp"
)

var (
	e         = echo.New()
	wg        sync.WaitGroup
	conf      = config.Load("backend-golang")
	logClient = logger.NewLoggerstore(logger.FLUENT, &conf.Server, &conf.Service)
	// messageQueue = stream.NewStreamClient(stream.RABBITMQ, &conf.Service)
)

func main() {
	e.Use(middleware.RequestID())
	//e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(middleware.Gzip())

	e.Use(middleware.LoggerWithConfig(
		middleware.LoggerConfig{ // add uuid header to log
			Format: `{"level":"info","time":"${time_rfc3339_nano}","id":"${id}","remote_ip":"${remote_ip}","host":"${host}",` +
				`"method":"${method}","uri":"${uri}","status":${status},"error":"${error}","latency":${latency},` +
				`"latency_human":"${latency_human}","bytes_in":${bytes_in},` +
				`"bytes_out":${bytes_out},"uuid":"${header:uuid}"}` + "\n",
			Output: logClient,
		},
	))

	// // Message Queue testing
	// c := make(chan string)
	// r := messageQueue.NewConsumerStore()
	// r.Consumer("topic", c)
	// p := messageQueue.NewProducerStore()
	// p.Producer("topic", "hello")

	// go func() {
	// 	for msg := range c {
	// 		log.Println("msg: ", msg)
	// 	}
	// }()

	// Setup API Group
	apiGroup := e.Group(apigroup.SetAPIGroup(conf.Server.Name, conf.Server.Version))

	// Setup environment variable
	env := &config.Environment{
		Config:    &conf,
		Database:  datastore.NewDatastore(datastore.MONGODB, &conf.Service),
		Cache:     cachestore.NewCachestore(cachestore.REDIS, &conf.Service),
		Storage:   cloudStorage.NewFilestore(cloudStorage.DRIVE, nil),
		Email:     email.NewMailClient(email.SENDGRID, &conf.Service),
		Monitor:   monitor.NewMonitorStore(monitor.PGO, &conf.Server, &conf.Service),
		JWT:       &jwtMiddleware.Client{},
		Condition: &condition.Client{},
		Hash:      &hash.Client{},
		OTP:       &otp.Client{},
	}

	// API routing
	healthcheck.New(env).Handler(apiGroup)
	login.New(env).Handler(apiGroup)
	logout.New(env).Handler(apiGroup)
	document.New(env).Handler(apiGroup)
	user.New(env).Handler(apiGroup)
	refreshtoken.New(env).Handler(apiGroup)

	if env.Config.Server.Monitoring {
		metrics.New(env).Handler(e.Group(""))
	}

	// Service run without listening channel
	// e.Logger.Fatal(e.Start(":" + strconv.Itoa(conf.Server.Port)))

	// Service run with listening channel
	startServer(e, conf.Server)
	wg.Wait()
}

func startServer(e *echo.Echo, server model.Server) {
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)

	go func() {
		// start server
		e.Logger.Info("Start listening")
		port := strconv.Itoa(server.Port)

		if conf.Server.Heroku {
			port = os.Getenv("PORT")
		}

		if conf.Server.HTTPS {
			if err := e.StartTLS(":"+port, "./key/cert.pem", "./key/key.pem"); err != nil {
				e.Logger.Infof("Https server has error: %v", err)
				close(sigint)
				return
			}
		}

		if err := e.Start(":" + port); err != nil {
			e.Logger.Infof("Http server has error: %v", err)
			close(sigint)
			return
		}
	}()

	// listen for terminate signal
	<-sigint
	e.Logger.Infof("Shutting down the service")
	var t time.Duration
	if server.ShutdownTimeout < 1 {
		t = 5 * time.Second
	} else {
		t = time.Duration(server.ShutdownTimeout) * time.Second
	}
	ctx, cancel := context.WithTimeout(context.Background(), t)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Errorf("Http server shutdown: %v", err)
	}
	e.Logger.Infof("Service gracefully stopped")
}
