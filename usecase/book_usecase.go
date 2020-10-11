package usecase

import (
	"context"
	"time"

	"github.com/golang-common-packages/template/domain"
)

// BookUsecase ...
type BookUsecase struct {
	bookRepo       domain.BookRepository
	contextTimeout time.Duration
}

// NewBookUsecase will create new an articleUsecase object
func NewBookUsecase(bookRepo domain.BookRepository, timeout time.Duration) BookUsecase {
	return BookUsecase{
		bookRepo:       bookRepo,
		contextTimeout: timeout,
	}
}

// FetchBooks ...
func (b *BookUsecase) FetchBooks(ctx context.Context, lastID string, limit uint32) (books []*domain.Book, err error) {
	if limit == 0 {
		limit = 10
	}

	ctx, cancel := context.WithTimeout(ctx, b.contextTimeout)
	defer cancel()

	books, err = b.bookRepo.Fetch(ctx, lastID, limit)
	if err != nil {
		return nil, err
	}

	return books, nil
}

// GetBookByID ...
func (b *BookUsecase) GetBookByID(ctx context.Context, id string) (book *domain.Book, err error) {
	ctx, cancel := context.WithTimeout(ctx, b.contextTimeout)
	defer cancel()

	book, err = b.bookRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return book, nil
}

// UpdateBook ...
func (b *BookUsecase) UpdateBook(ctx context.Context, book *domain.Book) error {
	ctx, cancel := context.WithTimeout(ctx, b.contextTimeout)
	defer cancel()

	book.Updated = time.Now()
	return b.bookRepo.Update(ctx, book)
}

// StoreBook ...
func (b *BookUsecase) StoreBook(ctx context.Context, book *domain.Book) error {
	ctx, cancel := context.WithTimeout(ctx, b.contextTimeout)
	defer cancel()

	return b.bookRepo.Store(ctx, book)
}

// DeleteBookByID ...
func (b *BookUsecase) DeleteBookByID(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, b.contextTimeout)
	defer cancel()

	existedBook, err := b.bookRepo.GetByID(ctx, id)
	if err != nil {
		return nil
	}

	if existedBook == (&domain.Book{}) {
		return domain.ErrNotFound
	}

	return b.bookRepo.Delete(ctx, id)
}
