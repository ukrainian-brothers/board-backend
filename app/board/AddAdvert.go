package board

import (
	"context"
	"github.com/ukrainian-brothers/board-backend/domain/user"
	"github.com/ukrainian-brothers/board-backend/internal/advert_repo"
)

type AddAdvert struct {
	AdvertRepo advert_repo.Repository
}

func NewAddAdvert(advertRepo advert_repo.Repository) *AddAdvert {
	return &AddAdvert{AdvertRepo: advertRepo}
}

func (a AddAdvert) Execute(ctx context.Context, usr *user.User) error {
	return nil
}