package domain

import (
	"context"
	"time"
)

// User model
type User struct {
	ID         string    `json:"id,omitempty"`
	Name       string    `json:"name,omitempty"`
	Age        int       `json:"age,omitempty"`
	Username   string    `json:"username,omitempty"`
	Password   *string   `json:"password,omitempty"`
	Email      string    `json:"email,omitempty"`
	IsActive   bool      `json:"isactive"`
	Created    time.Time `json:"created,omitempty"`
	Updated    time.Time `json:"updated,omitempty"`
	Expiration time.Time `json:"expiration,omitempty"`
}

// UserRepository represent the user's repository contract
type UserRepository interface {
	Fetch(ctx context.Context, lastID string, limit uint32) (users []*User, err error)
	GetByID(ctx context.Context, id string) (user *User, err error)
	Update(ctx context.Context, user *User) error
	Store(ctx context.Context, user *User) error
	Delete(ctx context.Context, id string) error
}
