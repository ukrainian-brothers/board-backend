package board

import (
	"context"
	"github.com/ukrainian-brothers/board-backend/domain/user"
)

type AddUser struct {
	repo user.Repository
}

func NewAddUser(userRepo user.Repository) AddUser {
	return AddUser{repo: userRepo}
}

func (s AddUser) Execute(ctx context.Context, user *user.User) error {
	return s.repo.Add(ctx, user)
}

