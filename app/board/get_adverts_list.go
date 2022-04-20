package board

import (
	"context"
	"github.com/ukrainian-brothers/board-backend/domain/advert"
	. "github.com/ukrainian-brothers/board-backend/pkg/translation"
)

type GetAdvertsList struct {
	repo advert.Repository
}

func NewGetAdvertsList(advertRepo advert.Repository) GetAdvertsList {
	return GetAdvertsList{repo: advertRepo}
}

func (a GetAdvertsList) Execute(ctx context.Context, languages LanguageTags, limit int, offset int) ([]*advert.Advert, error) {
	return a.repo.GetList(ctx, languages, limit, offset)
}
