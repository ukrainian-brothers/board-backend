package board

import (
	"context"
	"github.com/google/uuid"
	"github.com/ukrainian-brothers/board-backend/domain/advert"
)

type GetAdvert struct {
	AdvertRepo advert.Repository
}

func NewGetAdvert(advertRepo advert.Repository) GetAdvert {
	return GetAdvert{AdvertRepo: advertRepo}
}

func (a GetAdvert) Execute(ctx context.Context, id uuid.UUID) (advert.Advert, error) {
	return a.AdvertRepo.Get(ctx, id)
}
