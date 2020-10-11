package controller

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/golang-common-packages/template/domain"
	"github.com/golang-common-packages/template/usecase"
)

// BookController ...
type BookController struct {
	bookUsecase usecase.BookUsecase
}

// NewBookController ...
func NewBookController(bookInteractor usecase.BookUsecase) *BookController {
	return &BookController{bookInteractor}
}

// Add ...
func (controller *BookController) Add(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	var book domain.Book
	err := json.NewDecoder(req.Body).Decode(&book)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(res).Encode(ErrorResponse{Message: "Invalid Payload"})
		return
	}

	err2 := controller.bookUsecase.StoreBook(context.Background(), &book)
	if err2 != nil {
		res.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(res).Encode(ErrorResponse{Message: err2.Error()})
		return
	}

	res.WriteHeader(http.StatusOK)
}

// FindAll ...
func (controller *BookController) FindAll(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	var pagination domain.Pagination
	err := json.NewDecoder(req.Body).Decode(&pagination)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(res).Encode(ErrorResponse{Message: "Invalid Payload"})
		return
	}

	results, err2 := controller.bookUsecase.FetchBooks(context.Background(), pagination.LastID, pagination.Limit)
	if err2 != nil {
		res.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(res).Encode(ErrorResponse{Message: err2.Error()})
		return
	}
	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(results)
}
