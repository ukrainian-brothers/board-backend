package advert

import (
	"context"
	"errors"
	"github.com/google/uuid"
	. "github.com/ukrainian-brothers/board-backend/pkg/translation"
)

var (
	AdvertAlreadyExists = errors.New("advert already exists in repository")
	AdvertNotFound      = errors.New("advert not found in repository")
)

type Repository interface {
	Get(ctx context.Context, id uuid.UUID) (Advert, error)
	GetList(ctx context.Context, langs LanguageTags, limit int, offset int) ([]*Advert, error)
	Add(ctx context.Context, advert *Advert) error
	Delete(ctx context.Context, id uuid.UUID) error
}
