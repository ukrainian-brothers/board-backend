package user

import (
	"context"
	"github.com/google/uuid"
)

type Repository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)
	GetByLogin(ctx context.Context, login string) (*User, error)
	Add(ctx context.Context, user *User) error
	Delete(ctx context.Context, id uuid.UUID) error
}
