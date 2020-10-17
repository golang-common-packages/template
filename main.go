package main

import (
	"log"

	"github.com/golang-common-packages/storage"
	"github.com/labstack/echo/v4"

	_httpDeliver "github.com/golang-common-packages/template/book/delivery/http"
	_httpMiddleware "github.com/golang-common-packages/template/book/delivery/http/middleware"
	_bookRepo "github.com/golang-common-packages/template/book/repository/mongo"
	_bookUsecase "github.com/golang-common-packages/template/book/usecase"
	_cfg "github.com/golang-common-packages/template/config"
)

var (
	config _cfg.IConfig
	dbConn storage.INoSQLDocument
)

func init() {

	config = _cfg.NewViperConfig()
	if config.GetBool(`debug`) {
		log.Println("Service RUN on DEBUG mode")
	}

	dbConn = storage.New(storage.NOSQLDOCUMENT)(storage.MONGODB, &storage.Config{MongoDB: storage.MongoDB{
		User:     config.GetString("database.user"),
		Password: config.GetString("database.password"),
		Hosts:    config.GetStringSlice("database.hosts"),
		Options:  config.GetStringSlice("database.options"),
		DB:       config.GetString("database.db"),
	}}).(storage.INoSQLDocument)
}

func main() {

	e := echo.New()
	middL := _httpMiddleware.InitMiddleware()
	e.Use(middL.CORS)

	bookRepo := _bookRepo.NewMongoBookRepository(dbConn)
	bookUsecase := _bookUsecase.NewBookUsecase(bookRepo)

	_httpDeliver.NewBookHandler(e, bookUsecase)

	e.Start(config.GetString("server.address"))
}
