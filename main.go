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

	"github.com/golang-common-packages/caching"
	"github.com/golang-common-packages/cloud-storage"
	"github.com/golang-common-packages/database"
	"github.com/golang-common-packages/echo-jwt-middleware"
	"github.com/golang-common-packages/email"
	"github.com/golang-common-packages/hash"
	"github.com/golang-common-packages/log"
	"github.com/golang-common-packages/monitoring"
	"github.com/golang-common-packages/otp"

	"github.com/golang-common-packages/template/config"
	"github.com/golang-common-packages/template/model"

	"github.com/golang-common-packages/template/handler/document"
	"github.com/golang-common-packages/template/handler/healthcheck"
	"github.com/golang-common-packages/template/handler/login"
	"github.com/golang-common-packages/template/handler/logout"
	"github.com/golang-common-packages/template/handler/metrics"
	"github.com/golang-common-packages/template/handler/refreshtoken"
	"github.com/golang-common-packages/template/handler/user"
)

var (
	e         = echo.New()
	wg        sync.WaitGroup
	conf      = config.Load("backend-golang")
	logClient = log.New(false, log.FLUENT, &log.Fluent{
		Tag:    conf.Service.Fluent.Tag,
		Host:   conf.Service.Fluent.Host,
		Port:   conf.Service.Fluent.Port,
		Prefix: conf.Service.Fluent.Prefix,
	})
	env = &config.Environment{
		Config: &conf,
		Database: database.NewDatabase(database.MONGODB, &database.Database{MongoDB: database.MongoDB{
			User:     conf.Service.Database.MongoDB.User,
			Password: conf.Service.Database.MongoDB.Password,
			Hosts:    conf.Service.Database.MongoDB.Hosts,
			Options:  conf.Service.Database.MongoDB.Options,
			DB:       conf.Service.Database.MongoDB.DB,
		}}),
		//Cache: caching.New(caching.REDIS, &caching.Config{Redis: caching.Redis{
		//	Password: conf.Service.Database.Redis.Password,
		//	Host:     conf.Service.Database.Redis.Host,
		//	DB:       conf.Service.Database.Redis.DB,
		//}}),
		Cache: caching.New(caching.CUSTOM, &caching.Config{CustomCache: caching.CustomCache{
			CleaningInterval: 3600000000000,    // nanosecond
			CacheSize:        10 * 1024 * 1024, // byte
			SizeChecker:      true,
		}}),
		Storage: cloudStorage.NewFilestore(cloudStorage.DRIVE, nil),
		Email: email.NewMailClient(email.SENDGRID, &email.MailConfig{
			URL:       conf.Service.Email.Host,
			Port:      conf.Service.Email.Port,
			Username:  conf.Service.Email.Username,
			Password:  conf.Service.Email.Password,
			SecretKey: conf.Service.Email.Key,
		}),
		Monitor: monitoring.New(monitoring.PGO, conf.Server.Name, ""),
		JWT:     &jwtMiddleware.Client{},
		Hash:    &hash.Client{},
		OTP:     &otp.Client{},
	}
)

func main() {
	e.Use(middleware.RequestID())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(middleware.Gzip())

	// Apply log service to echo
	e.Use(middleware.LoggerWithConfig(
		middleware.LoggerConfig{ // add uuid header to log
			Format: `{"level":"info","time":"${time_rfc3339_nano}","id":"${id}","remote_ip":"${remote_ip}","host":"${host}",` +
				`"method":"${method}","uri":"${uri}","status":${status},"error":"${error}","latency":${latency},` +
				`"latency_human":"${latency_human}","bytes_in":${bytes_in},` +
				`"bytes_out":${bytes_out},"uuid":"${header:uuid}"}` + "\n",
			Output: logClient,
		},
	))

	// Setup API Group
	apiGroup := e.Group("/api/" + conf.Server.Version)

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
