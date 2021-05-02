package bookMongoRepository_test

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/golang-common-packages/storage/mocks"
	"github.com/golang-common-packages/template/book/repository/mongo"
	"github.com/golang-common-packages/template/domain"
)

func TestCreateMany(t *testing.T) {
	mockStorage := new(mocks.INoSQLDocument)
	mockDBName := string(mock.AnythingOfType("string"))
	mockColName := string(mock.AnythingOfType("string"))
	mockBooks := []domain.Book{}

	t.Run("sucess", func(t *testing.T) {
		mockStorage.On("Create", mockDBName, mockColName, mock.AnythingOfType("[]interface {}")).
			Return(1, nil).
			Once()

		repo := bookMongoRepository.New(mockStorage)
		r, err := repo.CreateMany(mockDBName, mockColName, mockBooks)

		assert.NotEmpty(t, r)
		assert.NoError(t, err)

		mockStorage.AssertExpectations(t)
	})
}

func TestRead(t *testing.T) {
	mockStorage := new(mocks.INoSQLDocument)
	mockDBName := string(mock.AnythingOfType("string"))
	mockColName := string(mock.AnythingOfType("string"))
	mockLimit := int64(10)
	mockDataModel := reflect.TypeOf(domain.Book{})
	mockFilter := primitive.D{}

	t.Run("sucess", func(t *testing.T) {
		mockStorage.On("Read", mockDBName, mockColName, mock.AnythingOfType("primitive.D"), mock.AnythingOfType("int64"), mock.AnythingOfType("*reflect.rtype")).
			Return(1, nil).
			Once()

		repo := bookMongoRepository.New(mockStorage)
		r, err := repo.Read(mockDBName, mockColName, mockFilter, mockLimit, mockDataModel)

		assert.NotEmpty(t, r)
		assert.NoError(t, err)

		mockStorage.AssertExpectations(t)
	})
}

func TestUpdate(t *testing.T) {
	mockStorage := new(mocks.INoSQLDocument)
	mockDBName := string(mock.AnythingOfType("string"))
	mockColName := string(mock.AnythingOfType("string"))
	mockFilter := primitive.D{}
	mockUpdate := primitive.D{}

	t.Run("sucess", func(t *testing.T) {
		mockStorage.On("Update", mockDBName, mockColName, mock.AnythingOfType("primitive.D"), mock.AnythingOfType("primitive.D")).
			Return(1, nil).
			Once()

		repo := bookMongoRepository.New(mockStorage)
		r, err := repo.Update(mockDBName, mockColName, mockFilter, mockUpdate)

		assert.NotEmpty(t, r)
		assert.NoError(t, err)

		mockStorage.AssertExpectations(t)
	})
}

func TestDelete(t *testing.T) {
	mockStorage := new(mocks.INoSQLDocument)
	mockDBName := string(mock.AnythingOfType("string"))
	mockColName := string(mock.AnythingOfType("string"))
	mockFilter := primitive.D{}

	t.Run("sucess", func(t *testing.T) {
		mockStorage.On("Delete", mockDBName, mockColName, mock.AnythingOfType("primitive.D")).
			Return(1, nil).
			Once()

		repo := bookMongoRepository.New(mockStorage)
		r, err := repo.Delete(mockDBName, mockColName, mockFilter)

		assert.NotEmpty(t, r)
		assert.NoError(t, err)

		mockStorage.AssertExpectations(t)
	})
}
