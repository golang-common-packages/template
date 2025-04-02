package bookUsecase_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/golang-common-packages/template/book/usecase"
	"github.com/golang-common-packages/template/domain"
	"github.com/golang-common-packages/template/mocks/domain"
)

func TestInsertBooks(t *testing.T) {
	mockBookRepo := new(mocks.BookRepository)
	mockDBName := string(mock.AnythingOfType("string"))
	mockColName := string(mock.AnythingOfType("string"))
	
	validBooks := []domain.Book{
		{
			Title: "Valid Book",
			Author: "Valid Author",
		},
	}

	invalidBooks := []domain.Book{
		{
			Title: "", // Invalid - empty title
			Author: "Author",
		},
	}

	t.Run("success", func(t *testing.T) {
		mockBookRepo.On("CreateMany", mockDBName, mockColName, validBooks).
			Return(1, nil).
			Once()

		u := bookUsecase.New(mockBookRepo, mockDBName, mockColName)
		r, err := u.InsertBooks(&validBooks)

		assert.Equal(t, 1, r)
		assert.NoError(t, err)

		mockBookRepo.AssertExpectations(t)
	})

	t.Run("error-failed", func(t *testing.T) {
		mockBookRepo.On("CreateMany", mockDBName, mockColName, validBooks).
			Return(nil, errors.New("Unexpected Error")).
			Once()

		u := bookUsecase.New(mockBookRepo, mockDBName, mockColName)
		r, err := u.InsertBooks(&validBooks)

		assert.Empty(t, r)
		assert.Error(t, err)

		mockBookRepo.AssertExpectations(t)
	})

	t.Run("validation-failed", func(t *testing.T) {
		u := bookUsecase.New(mockBookRepo, mockDBName, mockColName)
		r, err := u.InsertBooks(&invalidBooks)

		assert.Empty(t, r)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "validation failed")

		mockBookRepo.AssertNotCalled(t, "CreateMany")
	})
}

func TestListBooks(t *testing.T) {
	mockBookRepo := new(mocks.BookRepository)
	mockDBName := string(mock.AnythingOfType("string"))
	mockColName := string(mock.AnythingOfType("string"))
	mockLimit := int64(10)
	mockDataModel := reflect.TypeOf(domain.Book{})

	t.Run("success", func(t *testing.T) {
		mockBookRepo.On("Read", mockDBName, mockColName, mock.AnythingOfType("primitive.D"), mock.AnythingOfType("int64"), mock.AnythingOfType("*reflect.rtype")).
			Return(1, nil).
			Once()

		u := bookUsecase.New(mockBookRepo, mockDBName, mockColName)
		r, err := u.ListBooks(mockLimit, mockDataModel)

		assert.Equal(t, 1, r)
		assert.NoError(t, err)

		mockBookRepo.AssertExpectations(t)
	})
}

func TestUpdateBook(t *testing.T) {
	mockBookRepo := new(mocks.BookRepository)
	mockDBName := string(mock.AnythingOfType("string"))
	mockColName := string(mock.AnythingOfType("string"))

	mockBook := domain.Book{
		ID:     "000000000000000000000000",
		Title:  "A",
		Author: "B",
	}

	t.Run("success", func(t *testing.T) {
		mockBookRepo.On("Update", mockDBName, mockColName, mock.AnythingOfType("primitive.D"), mock.AnythingOfType("primitive.D")).
			Return(1, nil).
			Once()

		u := bookUsecase.New(mockBookRepo, mockDBName, mockColName)
		r, err := u.UpdateBook(mockBook)

		assert.Equal(t, 1, r)
		assert.NoError(t, err)

		mockBookRepo.AssertExpectations(t)
	})
}

func TestDeleteBook(t *testing.T) {
	mockBookRepo := new(mocks.BookRepository)
	mockDBName := string(mock.AnythingOfType("string"))
	mockColName := string(mock.AnythingOfType("string"))
	mockBookID := "000000000000000000000000"

	t.Run("success", func(t *testing.T) {
		mockBookRepo.On("Delete", mockDBName, mockColName, mock.AnythingOfType("primitive.D")).
			Return(1, nil).
			Once()

		u := bookUsecase.New(mockBookRepo, mockDBName, mockColName)
		r, err := u.DeleteBook(mockBookID)

		assert.Equal(t, 1, r)
		assert.NoError(t, err)

		mockBookRepo.AssertExpectations(t)
	})
}
