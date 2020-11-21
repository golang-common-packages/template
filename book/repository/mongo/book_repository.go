package bookMongoRepository

import (
	"reflect"

	"github.com/golang-common-packages/storage"
	
	"github.com/golang-common-packages/template/domain"
)

type mongoBookRepository struct {
	Conn storage.INoSQLDocument
}

// New will create an object that represent the domain.BookRepository interface
func New(Conn storage.INoSQLDocument) domain.BookRepository {
	return &mongoBookRepository{Conn}
}

// Create ...
func (mb *mongoBookRepository) CreateMany(databaseName, collectionName string, books []domain.Book) (interface{}, error) {

	newBooks := make([]interface{}, len(books))
	for i, v := range books {
		newBooks[i] = v
	}

	result, err := mb.Conn.Create(databaseName, collectionName, newBooks)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Read ...
func (mb *mongoBookRepository) Read(databaseName, collectionName string, filter interface{}, limit int64, dataModel reflect.Type) (interface{}, error) {
	return mb.Conn.Read(databaseName, collectionName, filter, limit, dataModel)
}

// Update ...
func (mb *mongoBookRepository) Update(databaseName, collectionName string, filter, update interface{}) (interface{}, error) {
	return mb.Conn.Update(databaseName, collectionName, filter, update)
}

// Delete ...
func (mb *mongoBookRepository) Delete(databaseName, collectionName string, filter interface{}) (interface{}, error) {
	return mb.Conn.Delete(databaseName, collectionName, filter)
}
