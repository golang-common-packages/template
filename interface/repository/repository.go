package repository

import (
	"context"

	"github.com/golang-common-packages/template/domain"
)

// DBHandler ...
type DBHandler interface {
	Register(ctx context.Context, user domain.User) error
	Login(ctx context.Context, username, password string) (user *domain.User, err error)
	UpdateUser(ctx context.Context, user domain.User) error
	GetUserInfo(ctx context.Context, ID string) (user *domain.User, err error)
	FetchUsers(ctx context.Context, lastID string, limit uint32) (users []*domain.User, err error)

	SaveBook(ctx context.Context, book domain.Book) error
	UpdateBook(ctx context.Context, book domain.Book) error
	FetchBooks(ctx context.Context, lastID string, limit uint32) (users []*domain.Book, err error)
}
