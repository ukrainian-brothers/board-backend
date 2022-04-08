package board

import (
	"context"
	"github.com/ukrainian-brothers/board-backend/domain/user"
	"github.com/ukrainian-brothers/board-backend/pkg/password"
)

type VerifyUserPassword struct {
	repo user.Repository
}

func NewVerifyUserPassword(userRepo user.Repository) VerifyUserPassword {
	return VerifyUserPassword{repo: userRepo}
}

func (a VerifyUserPassword) Execute(ctx context.Context, user user.User) (bool, error) {
	userDB, err := a.repo.GetByLogin(ctx, user.Login)
	if err != nil {
		return false, err
	}

	valid, err := password.VerifyPassword(*user.Password, *userDB.Password)
	return valid, err
}
