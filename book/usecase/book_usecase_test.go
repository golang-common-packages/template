package bookUsecase_test

import (
	"reflect"
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/golang-common-packages/template/constant"
	"github.com/golang-common-packages/template/domain"
	"github.com/golang-common-packages/template/book/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockBookRepository struct {
	mock.Mock
}

func (m *MockBookRepository) CreateMany(databaseName, collectionName string, books []domain.Book) (interface{}, error) {
	args := m.Called(databaseName, collectionName, books)
	return args.Get(0), args.Error(1)
}

func (m *MockBookRepository) Read(databaseName, collectionName string, filter interface{}, limit int64, dataModel reflect.Type) (interface{}, error) {
	args := m.Called(databaseName, collectionName, filter, limit, dataModel)
	return args.Get(0), args.Error(1)
}

func (m *MockBookRepository) Update(databaseName, collectionName string, filter, update interface{}) (interface{}, error) {
	args := m.Called(databaseName, collectionName, filter, update)
	return args.Get(0), args.Error(1)
}

func (m *MockBookRepository) Delete(databaseName, collectionName string, filter interface{}) (interface{}, error) {
	args := m.Called(databaseName, collectionName, filter)
	return args.Get(0), args.Error(1)
}

func TestInsertBooks_Success(t *testing.T) {
	mockRepo := new(MockBookRepository)
	usecase := bookUsecase.New(mockRepo, "test_db", "books")

	books := []domain.Book{
		{Title: "Test Book", Author: "Author"},
	}

	mockRepo.On("CreateMany", "test_db", "books", books).Return(1, nil)

	result, err := usecase.InsertBooks(&books)

	assert.NoError(t, err)
	assert.Equal(t, 1, result)
	mockRepo.AssertExpectations(t)
}

func TestListBooks_Success(t *testing.T) {
	mockRepo := new(MockBookRepository)
	usecase := bookUsecase.New(mockRepo, "test_db", "books")

	expectedBooks := []domain.Book{
		{Title: "Test Book", Author: "Author"},
	}

	mockRepo.On("Read", "test_db", "books", primitive.D{}, int64(10), reflect.TypeOf(domain.Book{})).Return(expectedBooks, nil)

	result, err := usecase.ListBooks(10, reflect.TypeOf(domain.Book{}))

	assert.NoError(t, err)
	assert.Equal(t, expectedBooks, result)
	mockRepo.AssertExpectations(t)
}

func TestInsertBooks_InvalidData(t *testing.T) {
	mockRepo := new(MockBookRepository)
	usecase := bookUsecase.New(mockRepo, "test_db", "books")

	books := []domain.Book{
		{Title: "", Author: ""}, // Invalid data
	}

	// Expect validation to fail before calling repository
	_, err := usecase.InsertBooks(&books)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "validation failed")
	mockRepo.AssertNotCalled(t, "CreateMany")
}

func TestListBooks_InvalidLimit(t *testing.T) {
	mockRepo := new(MockBookRepository)
	usecase := bookUsecase.New(mockRepo, "test_db", "books")

	// Test negative limit
	_, err := usecase.ListBooks(-1, reflect.TypeOf(domain.Book{}))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid limit")

	// Test zero limit
	_, err = usecase.ListBooks(0, reflect.TypeOf(domain.Book{}))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid limit")

	// Verify repository was not called
	mockRepo.AssertNotCalled(t, "Read")
}

func TestUpdateBook_NotFound(t *testing.T) {
	mockRepo := new(MockBookRepository)
	usecase := bookUsecase.New(mockRepo, "test_db", "books")

	book := domain.Book{
		ID:    primitive.NewObjectID(),
		Title: "Updated Title",
		Author: "Updated Author",
	}

	// Mock repository to return ErrNotFound
	mockRepo.On("Update", "test_db", "books", mock.Anything, mock.Anything).Return(0, constant.ErrNotFound)

	_, err := usecase.UpdateBook(book)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
	mockRepo.AssertExpectations(t)
}
