package main

import (
	"log"

	"github.com/golang-common-packages/storage"
	"github.com/labstack/echo/v4"

	"github.com/golang-common-packages/template/book/delivery/http"
	"github.com/golang-common-packages/template/book/delivery/http/middleware"
	"github.com/golang-common-packages/template/book/repository/mongo"
	"github.com/golang-common-packages/template/book/usecase"
	"github.com/golang-common-packages/template/config"
)

var (
	cfg    config.IConfig
	dbConn storage.INoSQLDocument
)

func init() {

	cfg = config.New()
	if cfg.GetBool(`debug`) {
		log.Println("Service RUN on DEBUG mode")
	}

	dbConn = storage.New(storage.NOSQLDOCUMENT)(storage.MONGODB, &storage.Config{MongoDB: storage.MongoDB{
		User:     cfg.GetString("database.mongodb.user"),
		Password: cfg.GetString("database.mongodb.password"),
		Hosts:    cfg.GetStringSlice("database.mongodb.hosts"),
		Options:  cfg.GetStringSlice("database.mongodb.options"),
		DB:       cfg.GetString("database.mongodb.dbName"),
	}}).(storage.INoSQLDocument)
}

func main() {

	e := echo.New()
	middL := BookHttpMiddleware.New()
	e.Use(middL.CORS)

	bookRepo := bookMongoRepository.New(dbConn)
	bookUCase := bookUsecase.New(bookRepo, cfg.GetString("database.mongodb.dbName"), cfg.GetString("database.mongodb.collections.book"))
	bookHttpDelivery.New(e, bookUCase)

	e.Start(cfg.GetString("server.address"))
}
