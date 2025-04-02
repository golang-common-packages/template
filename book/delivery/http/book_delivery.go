package bookHttpDelivery

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

// New will initialize the book resources endpoint
func New(e *echo.Echo, bu domain.BookUsecase) {
	handler := &BookHandler{
		BUsercase: bu,
	}
	e.GET("/books", handler.Fetch)
	e.POST("/books", handler.StoreMany)
	e.PUT("/book", handler.Update)
	e.DELETE("/book/:id", handler.Delete)
}

// Fetch will fetch the book based on given params
func (b *BookHandler) Fetch(c echo.Context) error {
	limitS := c.QueryParam("limit")
	limit, _ := strconv.Atoi(limitS)

	listBook, err := b.BUsercase.ListBooks(int64(limit), reflect.TypeOf(domain.Book{}))
	if err != nil {
	return &AppError{
			Err: err,
			StatusCode: getStatusCode(err),
			Message: "Failed to fetch books",
		}
	}

	return c.JSON(http.StatusOK, listBook)
}

func (b *BookHandler) StoreMany(c echo.Context) error {
	books := new([]domain.Book)
	if err := c.Bind(&books); err != nil {
		return &AppError{
			Err: err,
			StatusCode: http.StatusUnprocessableEntity,
			Message: "Invalid request body",
		}
	}

	if ok, err := isRequestValidSlice(books); !ok {
		return &AppError{
			Err: err,
			StatusCode: http.StatusBadRequest,
			Message: "Validation failed",
		}
	}

	result, err := b.BUsercase.InsertBooks(books)
	if err != nil {
		return &AppError{
			Err: err,
			StatusCode: getStatusCode(err),
			Message: "Failed to insert books",
		}
	}

	return c.JSON(http.StatusCreated, result)
}

func (b *BookHandler) Update(c echo.Context) error {
	book := new(domain.Book)
	if err := c.Bind(&book); err != nil {
		return &AppError{
			Err: err,
			StatusCode: http.StatusUnprocessableEntity,
			Message: "Invalid request body",
		}
	}

	if _, err := b.BUsercase.UpdateBook(*book); err != nil {
		return &AppError{
			Err: err,
			StatusCode: getStatusCode(err),
			Message: "Failed to update book",
		}
	}

	return c.NoContent(http.StatusOK)
}

func (b *BookHandler) Delete(c echo.Context) error {
	if _, err := b.BUsercase.DeleteBook(c.Param("id")); err != nil {
		return &AppError{
			Err: err,
			StatusCode: getStatusCode(err),
			Message: "Failed to delete book",
		}
	}

	return c.NoContent(http.StatusOK)
}

func isRequestValid(book *domain.Book) (bool, error) {
	validate := validator.New()
	err := validate.Struct(book)
	if err != nil {
		return false, err
	}
	return true, nil
}

func isRequestValidSlice(books *[]domain.Book) (bool, error) {
	for _, book := range *books {
		return isRequestValid(&book)
	}

	return true, nil
}
