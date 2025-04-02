package bookHttpDelivery_test

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/golang-common-packages/template/domain"
)

func NewBookHandler(bu domain.BookUsecase) *BookHandler {
	return &BookHandler{
		BUsercase: bu,
	}
}

type BookHandler struct {
	BUsercase domain.BookUsecase
}

type mockBookUsecase struct {
	mock.Mock
}

func (m *mockBookUsecase) InsertBooks(books *[]domain.Book) (interface{}, error) {
	args := m.Called(books)
	return args.Get(0), args.Error(1)
}

func (m *mockBookUsecase) ListBooks(limit int64, dataModel reflect.Type) (interface{}, error) {
	args := m.Called(limit, dataModel)
	return args.Get(0), args.Error(1)
}

func (m *mockBookUsecase) UpdateBook(update domain.Book) (interface{}, error) {
	args := m.Called(update)
	return args.Get(0), args.Error(1)
}

func (m *mockBookUsecase) DeleteBook(bookID string) (interface{}, error) {
	args := m.Called(bookID)
	return args.Get(0), args.Error(1)
}

func (h *BookHandler) Fetch(c echo.Context) error {
	limit := 10 // default
	if l := c.QueryParam("limit"); l != "" {
		limit, _ = strconv.Atoi(l)
	}

	list, err := h.BUsercase.ListBooks(int64(limit), reflect.TypeOf(domain.Book{}))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, list)
}

func (h *BookHandler) StoreMany(c echo.Context) error {
	var books []domain.Book
	if err := c.Bind(&books); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	result, err := h.BUsercase.InsertBooks(&books)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, result)
}

func TestBookHandler_Fetch(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	m := new(mockBookUsecase)
	h := NewBookHandler(m)

	// Mock data
	books := []domain.Book{
		{Title: "Test Book", Author: "Test Author"},
	}

	// Expectations
	m.On("ListBooks", int64(10), reflect.TypeOf(domain.Book{})).Return(books, nil)

	// Test
	if assert.NoError(t, h.Fetch(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "Test Book")
	}

	// Assert expectations
	m.AssertExpectations(t)
}

func TestBookHandler_StoreMany(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/books", strings.NewReader(`[{"title":"New Book","author":"Author"}]`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder() 
	c := e.NewContext(req, rec)

	m := new(mockBookUsecase)
	h := NewBookHandler(m)

	// Expectations
	m.On("InsertBooks", mock.AnythingOfType("*[]domain.Book")).Return(1, nil)

	// Test
	if assert.NoError(t, h.StoreMany(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Contains(t, rec.Body.String(), "1")
	}

	// Assert expectations
	m.AssertExpectations(t)
}
