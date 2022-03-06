package advert_repo

import (
	"errors"
	"github.com/google/uuid"
	"github.com/ukrainian-brothers/board-backend/domain/advert"
)

var (
	AdvertAlreadyExists = errors.New("advert already exists in repository")
	AdvertNotFound = errors.New("advert not found in repository")
)

type Repository interface {
	Get(id uuid.UUID) (advert.Advert, error)
	Add(advert advert.Advert) error
	Delete(id uuid.UUID) error
}