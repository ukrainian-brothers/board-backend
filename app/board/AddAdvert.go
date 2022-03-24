package board

import (
	"context"
	"github.com/ukrainian-brothers/board-backend/domain/advert"
)

type AddAdvert struct {
	AdvertRepo advert.Repository
}

func NewAddAdvert(advertRepo advert.Repository) *AddAdvert {
	return &AddAdvert{AdvertRepo: advertRepo}
}

func (a AddAdvert) Execute(ctx context.Context, advert *advert.Advert) error {
	err := a.AdvertRepo.Add(ctx, advert)
	if err != nil {
		return err
	}
	return nil
}