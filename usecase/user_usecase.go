package usecase

import (
	"context"
	"time"

	"github.com/golang-common-packages/template/domain"
)

// UserUsecase ...
type UserUsecase struct {
	userRepo       domain.UserRepository
	contextTimeout time.Duration
}

// NewUserUsecase will create new an articleUsecase object
func NewUserUsecase(userRepo domain.UserRepository, timeout time.Duration) UserUsecase {
	return UserUsecase{
		userRepo:       userRepo,
		contextTimeout: timeout,
	}
}

// FetchUsers ...
func (b *UserUsecase) FetchUsers(ctx context.Context, lastID string, limit uint32) (users []*domain.User, err error) {
	if limit == 0 {
		limit = 10
	}

	ctx, cancel := context.WithTimeout(ctx, b.contextTimeout)
	defer cancel()

	users, err = b.userRepo.Fetch(ctx, lastID, limit)
	if err != nil {
		return nil, err
	}

	return users, nil
}

// GetUserByID ...
func (b *UserUsecase) GetUserByID(ctx context.Context, id string) (user *domain.User, err error) {
	ctx, cancel := context.WithTimeout(ctx, b.contextTimeout)
	defer cancel()

	user, err = b.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// UpdateUser ...
func (b *UserUsecase) UpdateUser(ctx context.Context, user *domain.User) error {
	ctx, cancel := context.WithTimeout(ctx, b.contextTimeout)
	defer cancel()

	user.Updated = time.Now()
	return b.userRepo.Update(ctx, user)
}

// StoreUser ...
func (b *UserUsecase) StoreUser(ctx context.Context, user *domain.User) error {
	ctx, cancel := context.WithTimeout(ctx, b.contextTimeout)
	defer cancel()

	return b.userRepo.Store(ctx, user)
}

// DeleteUserByID ...
func (b *UserUsecase) DeleteUserByID(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, b.contextTimeout)
	defer cancel()

	existedUser, err := b.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil
	}

	if existedUser == (&domain.User{}) {
		return domain.ErrNotFound
	}

	return b.userRepo.Delete(ctx, id)
}
