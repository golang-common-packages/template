package bookHttpDelivery_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/golang-common-packages/template/book/delivery/http"
	"github.com/golang-common-packages/template/domain"
	"github.com/golang-common-packages/template/mocks/domain"
)

func TestFetch(t *testing.T) {
	mockBookUCase := new(mocks.BookUsecase)
	currentTime := time.Now()
	mockLimit := int64(10)

	mockBook := domain.Book{
		ID:      "000000000000000000000000",
		Title:   "A",
		Author:  "B",
		Updated: currentTime,
		Created: currentTime,
	}
	mockListBook := make([]domain.Book, 0)
	mockListBook = append(mockListBook, mockBook)

	t.Run("sucess", func(t *testing.T) {
		mockBookUCase.On("ListBooks", mockLimit, mock.AnythingOfType("*reflect.rtype")).
			Return(mockListBook, nil)

		e := echo.New()
		req, err := http.NewRequest(echo.GET, "/books?limit="+fmt.Sprint(mockLimit), strings.NewReader(""))
		require.NoError(t, err)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		handler := bookHttpDelivery.BookHandler{
			BUsercase: mockBookUCase,
		}
		err = handler.Fetch(c)
		require.NoError(t, err)

		assert.Equal(t, "10", c.QueryParam("limit"))
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.NotEmpty(t, rec.Body.Len())

		mockBookUCase.AssertExpectations(t)
	})
}

func TestStoreMany(t *testing.T) {
	mockBookUCase := new(mocks.BookUsecase)
	mockBooks := []domain.Book{{
		Title:  "A",
		Author: "B",
	}}

	t.Run("sucess", func(t *testing.T) {
		mockBookUCase.On("InsertBooks", &mockBooks).
			Return(1, nil)

		mockJSON, err := json.Marshal(mockBooks)
		assert.NoError(t, err)

		e := echo.New()
		req, err := http.NewRequest(echo.POST, "/books", strings.NewReader(string(mockJSON)))
		require.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/books")
		handler := bookHttpDelivery.BookHandler{
			BUsercase: mockBookUCase,
		}
		err = handler.StoreMany(c)
		require.NoError(t, err)

		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Equal(t, "1\n", rec.Body.String())

		mockBookUCase.AssertExpectations(t)
	})
}

func TestUpdate(t *testing.T) {
	mockBookUCase := new(mocks.BookUsecase)
	currentTime := time.Now()

	mockBooks := domain.Book{
		ID:      "000000000000000000000000",
		Title:   "A",
		Author:  "B",
		Updated: currentTime,
		Created: currentTime,
	}

	t.Run("sucess", func(t *testing.T) {
		mockBookUCase.On("UpdateBook", mock.AnythingOfType("domain.Book")).
			Return(nil, nil)

		mockJSON, err := json.Marshal(mockBooks)
		assert.NoError(t, err)

		e := echo.New()
		req, err := http.NewRequest(echo.PUT, "/book", strings.NewReader(string(mockJSON)))
		require.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/book")
		handler := bookHttpDelivery.BookHandler{
			BUsercase: mockBookUCase,
		}
		err = handler.Update(c)
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, rec.Code)

		mockBookUCase.AssertExpectations(t)
	})
}

func TestDelete(t *testing.T) {
	mockBookUCase := new(mocks.BookUsecase)
	mockBookID := "000000000000000000000000"

	t.Run("sucess", func(t *testing.T) {
		mockBookUCase.On("DeleteBook", mockBookID).
			Return(nil, nil)

		e := echo.New()
		req, err := http.NewRequest(echo.DELETE, "/book/"+mockBookID, strings.NewReader(""))
		require.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/book/:id")
		c.SetParamNames("id")
		c.SetParamValues(mockBookID)
		handler := bookHttpDelivery.BookHandler{
			BUsercase: mockBookUCase,
		}
		err = handler.Delete(c)
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, rec.Code)

		mockBookUCase.AssertExpectations(t)
	})
}
