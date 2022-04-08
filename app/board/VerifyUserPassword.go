package board

import (
	"context"
	"fmt"
	"github.com/ukrainian-brothers/board-backend/domain/user"
	"github.com/ukrainian-brothers/board-backend/pkg/password"
)

type VerifyUserPassword struct {
	repo user.Repository
}

func NewVerifyUserPassword(userRepo user.Repository) VerifyUserPassword {
	return VerifyUserPassword{repo: userRepo}
}

func (a VerifyUserPassword) Execute(ctx context.Context, login string, rawPassword string) (bool, error) {
	userDB, err := a.repo.GetByLogin(ctx, login)
	if err != nil {
		return false, fmt.Errorf("failed GetUserByLogin: %w", err)
	}

	valid, err := password.VerifyPassword(rawPassword, *userDB.Password)
	return valid, err
}
