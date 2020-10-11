package repository

import (
	"context"

	"github.com/golang-common-packages/template/domain"
)

// UserRepo ...
type UserRepo struct {
	handler DBHandler
}

// NewUserRepo ...
func NewUserRepo(handler DBHandler) UserRepo {
	return UserRepo{handler}
}

// UserRegister ...
func (ur UserRepo) UserRegister(ctx context.Context, user domain.User) error {
	err := ur.handler.Register(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

// UserLogin ...
func (ur UserRepo) UserLogin(ctx context.Context, username, password string) (*domain.User, error) {
	results, err := ur.handler.Login(ctx, username, password)
	if err != nil {
		return results, err
	}

	return results, nil
}

// UpdateUser ...
func (ur UserRepo) UpdateUser(ctx context.Context, user domain.User) error {
	if err := ur.handler.UpdateUser(ctx, user); err != nil {
		return err
	}

	return nil
}

// GetUserInfo ...
func (ur UserRepo) GetUserInfo(ctx context.Context, ID string) (*domain.User, error) {
	results, err := ur.handler.GetUserInfo(ctx, ID)
	if err != nil {
		return results, err
	}

	return results, nil
}

// FetchUsers ...
func (ur UserRepo) FetchUsers(ctx context.Context, lastID string, limit uint32) ([]*domain.User, error) {
	users, err := ur.handler.FetchUsers(ctx, lastID, limit)
	if err != nil {
		return users, err
	}

	return users, nil
}
