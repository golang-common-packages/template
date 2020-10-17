package usercase

import (
	"reflect"

	"github.com/golang-common-packages/template/domain"
)

type bookUsecase struct {
	bookRepo domain.BookRepository
}

// NewBookUsecase ...
func NewBookUsecase(br domain.BookRepository) domain.BookUsecase {
	return &bookUsecase{
		bookRepo: br,
	}
}

func (bu *bookUsecase) InsertBook(books []domain.Book) (interface{}, error) {
	return bu.bookRepo.Create("databaseName", "collectionName", books)
}

func (bu *bookUsecase) ListBooks(limit int64, dataModel reflect.Type) (interface{}, error) {
	return bu.bookRepo.Read("databaseName", "collectionName", "filter", limit, dataModel)
}

func (bu *bookUsecase) UpdateBook(bookID string, update interface{}) (interface{}, error) {
	return bu.bookRepo.Update("databaseName", "collectionName", "filter", update)
}

func (bu *bookUsecase) DeleteBook(bookID string) (interface{}, error) {
	return bu.bookRepo.Delete("databaseName", "collectionName", "filter")
}
