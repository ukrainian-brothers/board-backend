package board

import (
	"context"
	"github.com/ukrainian-brothers/board-backend/domain/user"
)

type UserExists struct {
	repo user.Repository
}

func NewUserExists(repo user.Repository) UserExists {
	return UserExists{repo: repo}
}

func (a UserExists) Execute(ctx context.Context, login string) (bool, error) {
	return a.repo.Exists(ctx, login)
}
