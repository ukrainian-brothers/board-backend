package board

import (
	"context"
	"github.com/ukrainian-brothers/board-backend/domain/user"
)

type AddUser struct {
	UserRepo user.Repository
}

func NewAddUser(userRepo user.Repository) *AddUser {
	return &AddUser{UserRepo: userRepo}
}

func (s AddUser) Execute(ctx context.Context, user *user.User) error {
	err := s.UserRepo.Add(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

