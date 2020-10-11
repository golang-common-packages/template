package repository

import (
	"context"

	"github.com/golang-common-packages/template/domain"
)

// BookRepo ...
type BookRepo struct {
	handler DBHandler
}

// NewBookRepo ...
func NewBookRepo(handler DBHandler) BookRepo {
	return BookRepo{handler}
}

// SaveBook ...
func (ur UserRepo) SaveBook(ctx context.Context, book domain.Book) error {
	err := ur.handler.SaveBook(ctx, book)
	if err != nil {
		return err
	}

	return nil
}

// UpdateBook ...
func (ur UserRepo) UpdateBook(ctx context.Context, book domain.Book) error {
	if err := ur.handler.UpdateBook(ctx, book); err != nil {
		return err
	}

	return nil
}

// FetchBooks ...
func (ur UserRepo) FetchBooks(ctx context.Context, lastID string, limit uint32) ([]*domain.Book, error) {
	books, err := ur.handler.FetchBooks(ctx, lastID, limit)
	if err != nil {
		return nil, err
	}

	return books, nil
}
