package bookUsecase

import (
	"fmt"
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/golang-common-packages/template/domain"
)

type bookUsecase struct {
	bookRepo domain.BookRepository
	dbName   string
	collName string
}

// New ...
func New(bookRepo domain.BookRepository, dbName, collectionName string) domain.BookUsecase {
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
	return bu.bookRepo.Read(bu.dbName, bu.collName, primitive.D{}, limit, dataModel)
}

func (bu *bookUsecase) UpdateBook(newData domain.Book) (interface{}, error) {
	
	idPrimitive, err := primitive.ObjectIDFromHex(fmt.Sprintf("%v", newData.ID))
	if err != nil {
		return nil, err
	}

	filter := primitive.D{
		primitive.E{Key: "_id", Value: idPrimitive},
	}

	book := primitive.D{
		primitive.E{Key: "$set", Value: primitive.D{
			primitive.E{Key: "title", Value: newData.Title},
			primitive.E{Key: "author", Value: newData.Author},
			primitive.E{Key: "updated", Value: time.Now()},
		}},
	}

	return bu.bookRepo.Update(bu.dbName, bu.collName, filter, book)
}

func (bu *bookUsecase) DeleteBook(bookID string) (interface{}, error) {

	idPrimitive, err := primitive.ObjectIDFromHex(bookID)
	if err != nil {
		return nil, err
	}

	filter := primitive.D{
		primitive.E{Key: "_id", Value: idPrimitive},
	}

	result, err := bu.bookRepo.Delete(bu.dbName, bu.collName, filter)
	if err != nil {
		return nil, err
	}

	return result, nil
}
