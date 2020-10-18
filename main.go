package main

import (
	"log"

	"github.com/golang-common-packages/storage"
	"github.com/labstack/echo/v4"

	_httpDeliver "github.com/golang-common-packages/template/book/delivery/http"
	_httpMiddleware "github.com/golang-common-packages/template/book/delivery/http/middleware"
	_bookRepo "github.com/golang-common-packages/template/book/repository/mongo"
	_bookUsecase "github.com/golang-common-packages/template/book/usecase"
	_config "github.com/golang-common-packages/template/config"
)

var (
	config _config.IConfig
	dbConn storage.INoSQLDocument
)

func init() {

	config = _config.NewViperConfig()
	if config.GetBool(`debug`) {
		log.Println("Service RUN on DEBUG mode")
	}

	dbConn = storage.New(storage.NOSQLDOCUMENT)(storage.MONGODB, &storage.Config{MongoDB: storage.MongoDB{
		User:     config.GetString("database.mongodb.user"),
		Password: config.GetString("database.mongodb.password"),
		Hosts:    config.GetStringSlice("database.mongodb.hosts"),
		Options:  config.GetStringSlice("database.mongodb.options"),
		DB:       config.GetString("database.mongodb.dbName"),
	}}).(storage.INoSQLDocument)
}

func main() {

	e := echo.New()
	middL := _httpMiddleware.InitMiddleware()
	e.Use(middL.CORS)

	bookRepo := _bookRepo.NewMongoBookRepository(dbConn)
	bookUsecase := _bookUsecase.NewBookUsecase(bookRepo, config.GetString("database.mongodb.dbName"), config.GetString("database.mongodb.bookCollName"))

	_httpDeliver.NewBookHandler(e, bookUsecase)

	e.Start(config.GetString("server.address"))
}
