package domain

import (
	"reflect"
	"time"
)

// Book ...
type Book struct {
	ID      interface{} `json:"id,omitempty" bson:"_id,omitempty"`
	Title   string      `json:"title" bson:"title" validate:"required"`
	Author  string      `json:"author" bson:"author" validate:"required"`
	Updated time.Time   `json:"updated,omitempty" bson:"updated,omitempty"`
	Created time.Time   `json:"created,omitempty" bson:"created,omitempty"`
}

// BookRepository ...
type BookRepository interface {
	CreateMany(databaseName, collectionName string, books []Book) (interface{}, error)
	Read(databaseName, collectionName string, filter interface{}, limit int64, dataModel reflect.Type) (interface{}, error)
	Update(databaseName, collectionName string, filter, update interface{}) (interface{}, error)
	Delete(databaseName, collectionName string, filter interface{}) (interface{}, error)
}

// BookUsecase ..
type BookUsecase interface {
	InsertBooks(books *[]Book) (interface{}, error)
	ListBooks(limit int64, dataModel reflect.Type) (interface{}, error)
	UpdateBook(update Book) (interface{}, error)
	DeleteBook(bookID string) (interface{}, error)
}
