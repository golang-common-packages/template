package domain

import (
	"reflect"
	"time"
)

// Book ...
type Book struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title" validate:"required"`
	Author    string    `json:"author" validate:"required"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}

// BookRepository ...
type BookRepository interface {
	Create(databaseName, collectionName string, books []Book) (interface{}, error)
	Read(databaseName, collectionName string, filter interface{}, limit int64, dataModel reflect.Type) (interface{}, error)
	Update(databaseName, collectionName string, filter, update interface{}) (interface{}, error)
	Delete(databaseName, collectionName string, filter interface{}) (interface{}, error)
}

// BookUsecase ..
type BookUsecase interface {
	InsertBook(books []Book) (interface{}, error)
	ListBooks(limit int64, dataModel reflect.Type) (interface{}, error)
	UpdateBook(bookID string, update interface{}) (interface{}, error)
	DeleteBook(bookID string) (interface{}, error)
}
