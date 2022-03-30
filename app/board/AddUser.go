package board

import (
	"context"
	"github.com/ukrainian-brothers/board-backend/domain/user"
)

type AddUser struct {
	repo user.UserRepository
}

func NewAddUser(userRepo user.UserRepository) AddUser {
	return AddUser{repo: userRepo}
}

func (s AddUser) Execute(ctx context.Context, user *user.User) error {
	return s.repo.Add(ctx, user)
}

