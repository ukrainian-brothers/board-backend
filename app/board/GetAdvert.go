package board

import (
	"context"
	"errors"
	"github.com/ukrainian-brothers/board-backend/domain/advert"
)

type GetAdvert struct {
	AdvertRepo interface{}
}

func NewGetAdvert(advertRepo interface{}) GetAdvert {
	return GetAdvert{AdvertRepo: advertRepo}
}

func (a GetAdvert) Execute(ctx context.Context) (advert.Advert, error) {
	return advert.Advert{}, errors.New("not implemented yet")
}