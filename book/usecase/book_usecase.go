package usercase

import (
	"fmt"
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/golang-common-packages/template/domain"
)

type bookUsecase struct {
	bookRepo domain.BookRepository
	dbName   string
	collName string
}

// NewBookUsecase ...
func NewBookUsecase(bookRepo domain.BookRepository, dbName, collectionName string) domain.BookUsecase {
	return &bookUsecase{
		bookRepo: bookRepo,
		dbName:   dbName,
		collName: collectionName,
	}
}

func (bu *bookUsecase) InsertBooks(books *[]domain.Book) (interface{}, error) {

	// Dereference the pointer and update the value
	currentTime := time.Now()
	for i := range *books {
		(*books)[i].Created = currentTime
		(*books)[i].Updated = currentTime
	}

	result, err := bu.bookRepo.CreateMany(bu.dbName, bu.collName, *books)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (bu *bookUsecase) ListBooks(limit int64, dataModel reflect.Type) (interface{}, error) {
	return bu.bookRepo.Read(bu.dbName, bu.collName, bson.D{}, limit, dataModel)
}

func (bu *bookUsecase) UpdateBook(newData domain.Book) (interface{}, error) {
	idPrimitive, err := primitive.ObjectIDFromHex(fmt.Sprintf("%v", newData.ID))
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": idPrimitive}

	book := bson.D{
		{"$set", bson.D{
			{"title", newData.Title},
			{"author", newData.Author},
			{"updated", time.Now()},
		}},
	}

	return bu.bookRepo.Update(bu.dbName, bu.collName, filter, book)
}

func (bu *bookUsecase) DeleteBook(bookID string) (interface{}, error) {

	idPrimitive, err := primitive.ObjectIDFromHex(bookID)
	if err != nil {
		return nil, err
	}

	result, err := bu.bookRepo.Delete(bu.dbName, bu.collName, bson.M{"_id": idPrimitive})
	if err != nil {
		return nil, err
	}

	return result, nil
}
