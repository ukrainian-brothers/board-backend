package board

import (
	"context"
	"github.com/ukrainian-brothers/board-backend/domain/advert"
	"github.com/ukrainian-brothers/board-backend/domain/user"
)

type AddAdvert struct {
	AdvertRepo advert.Repository
}

func NewAddAdvert(advertRepo advert.Repository) *AddAdvert {
	return &AddAdvert{AdvertRepo: advertRepo}
}

func (a AddAdvert) Execute(ctx context.Context, usr *user.User) error {
	return nil
}