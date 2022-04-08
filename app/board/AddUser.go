package board

import (
	"context"
	"github.com/ukrainian-brothers/board-backend/domain/user"
	"github.com/ukrainian-brothers/board-backend/pkg/password"
)

type AddUser struct {
	repo user.Repository
}

func NewAddUser(userRepo user.Repository) AddUser {
	return AddUser{repo: userRepo}
}

func (s AddUser) Execute(ctx context.Context, user user.User) error {
	hashedPassword, err := password.HashPassword(*user.Password, password.GetHashingParams())
	if err != nil {
		return err
	}
	user.Password = &hashedPassword

	return s.repo.Add(ctx, &user)
}