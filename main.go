package main

import (
	"log"

	"github.com/labstack/echo/v4"

	"github.com/golang-common-packages/storage"
	"github.com/golang-common-packages/template/book/delivery/http/middleware"
	_bookRepo "github.com/golang-common-packages/template/book/repository/mongo"
	_cfg "github.com/golang-common-packages/template/config"
	_bookUsecase "github.com/golang-common-packages/template/book/usecase"
	_httpDeliver "github.com/golang-common-packages/template/book/delivery/http"
)

var (
	config _cfg.Config
	dbConn storage.INoSQLDocument
)

func init() {

	dbConn = storage.New(storage.NOSQLDOCUMENT)(storage.MONGODB, &storage.Config{MongoDB: storage.MongoDB{
		User:     config.GetString(`database.user`),
		Password: config.GetString(`database.password`),
		Hosts:    []string{config.GetString(`database.host`)},
		Options:  []string{config.GetString(`database.host`)},
		DB:       config.GetString(`database.database`),
	}}).(storage.INoSQLDocument)

	config = _cfg.NewViperConfig()
	if config.GetBool(`debug`) {
		log.Println("Service RUN on DEBUG mode")
	}
}

func main() {

	e := echo.New()
	middL := middleware.InitMiddleware()
	e.Use(middL.CORS)

	bookRepo := _bookRepo.NewMongoBookRepository(dbConn)
	bookUsecase := _bookUsecase.NewBookUsecase(bookRepo)

	_httpDeliver.NewBookHandler(e, bookUsecase)

	e.Start(config.GetString("server.address"))
}
