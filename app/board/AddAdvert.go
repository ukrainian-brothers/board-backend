package board

import "context"

type AddAdvert struct {
	AdvertRepo interface{}
}

func NewAddAdvert(advertRepo interface{}) *AddAdvert {
	return &AddAdvert{AdvertRepo: advertRepo}
}

func (a AddAdvert) Execute(ctx context.Context) error {return nil}