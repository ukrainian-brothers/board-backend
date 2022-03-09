package advert_repo

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/ukrainian-brothers/board-backend/domain/advert"
)

var (
	AdvertAlreadyExists = errors.New("advert already exists in repository")
	AdvertNotFound = errors.New("advert not found in repository")
)

type Repository interface {
	Get(ctx context.Context, id uuid.UUID) (advert.Advert, error)
	Add(ctx context.Context, advert advert.Advert) error
	Delete(ctx context.Context, id uuid.UUID) error
}