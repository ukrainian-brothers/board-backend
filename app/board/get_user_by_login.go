package board

import (
	"context"
	"github.com/ukrainian-brothers/board-backend/domain/user"
)

type GetUserByLogin struct {
	repo user.Repository
}

func NewGetUserByLogin(userRepo user.Repository) GetUserByLogin {
	return GetUserByLogin{repo: userRepo}
}

func (a GetUserByLogin) Execute(ctx context.Context, login string) (*user.User, error) {
	return a.repo.GetByLogin(ctx, login)
}
