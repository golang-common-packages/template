package http

import (
	"net/http"
	"reflect"
	"strconv"

	"github.com/labstack/echo/v4"
	"gopkg.in/go-playground/validator.v9"

	"github.com/golang-common-packages/template/domain"
)

// ResponseError represent the reseponse error struct
type ResponseError struct {
	Message string `json:"message"`
}

// BookHandler  represent the httphandler for book
type BookHandler struct {
	BUsercase domain.BookUsecase
}

// NewBookHandler will initialize the book resources endpoint
func NewBookHandler(e *echo.Echo, bu domain.BookUsecase) {
	handler := &BookHandler{
		BUsercase: bu,
	}
	e.GET("/books", handler.Fetch)
	e.POST("/books", handler.Store)
	e.GET("/books/:id", handler.Update)
	e.DELETE("/books/:id", handler.Delete)
}

// Fetch will fetch the book based on given params
func (b *BookHandler) Fetch(c echo.Context) error {
	limitS := c.QueryParam("limit")
	limit, _ := strconv.Atoi(limitS)

	listBook, err := b.BUsercase.ListBooks(int64(limit), reflect.TypeOf(domain.Book{}))
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, listBook)
}

// Store will store the book by given request body
func (b *BookHandler) Store(c echo.Context) (err error) {
	var book domain.Book
	err = c.Bind(&book)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	var ok bool
	if ok, err = isRequestValid(&book); !ok {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	var books = []domain.Book{book}
	_, err = b.BUsercase.InsertBook(books)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, book)
}

// Update will update book by given param
func (b *BookHandler) Update(c echo.Context) error {
	var book domain.Book
	err := c.Bind(&book)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	_, err = b.BUsercase.UpdateBook(c.Param("id"), book)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	return c.NoContent(http.StatusOK)
}

// Delete will delete book by given param
func (b *BookHandler) Delete(c echo.Context) error {
	_, err := b.BUsercase.DeleteBook(c.Param("id"))
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	return c.NoContent(http.StatusOK)
}

func isRequestValid(m *domain.Book) (bool, error) {
	validate := validator.New()
	err := validate.Struct(m)
	if err != nil {
		return false, err
	}
	return true, nil
}
