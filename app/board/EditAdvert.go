package board

import "context"

type EditAdvert struct {
	AdvertRepo interface{}
}

func NewEditAdvert(advertRepo interface{}) EditAdvert {
	return EditAdvert{AdvertRepo: advertRepo}
}

func (a EditAdvert) Execute(ctx context.Context) error {return nil}