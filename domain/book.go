package domain

import (
	"context"
	"time"
)

// Book model
type Book struct {
	ID         string    `json:"id,omitempty"`
	Title      string    `json:"title,omitempty"`
	Author     string    `json:"author,omitempty"`
	Created    time.Time `json:"created,omitempty"`
	Updated    time.Time `json:"updated,omitempty"`
	Expiration time.Time `json:"expiration,omitempty"`
}

// Pagination ...
type Pagination struct {
	LastID string `json:"lastID,omitempty"`
	Limit  uint32 `json:"limit,omitempty"`
}

// BookRepository represent the book's repository contract
type BookRepository interface {
	Fetch(ctx context.Context, lastID string, limit uint32) (books []*Book, err error)
	GetByID(ctx context.Context, id string) (book *Book, err error)
	Update(ctx context.Context, book *Book) error
	Store(ctx context.Context, book *Book) error
	Delete(ctx context.Context, id string) error
}
