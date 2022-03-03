package advert

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/ukrainian-brothers/board-backend/domain"
	"github.com/ukrainian-brothers/board-backend/domain/user"
	"time"
)

type AdvertLogTrigger string

const (
	AdvertCreatedEvent AdvertLogTrigger = "created"
	AdvertUpdatedEvent                  = "updated"
	AdvertDeletedEvent                  = "deleted"
)

type AdvertLog struct {
	AdvertID uuid.UUID
	Trigger  AdvertLogTrigger
	Meta     json.RawMessage
}

type Advert struct {
	ID          uuid.UUID
	Advert      domain.Advert
	Contact     domain.ContactDetails
	Creator     user.User
	CreatedAt   time.Time
	UpdatedAt   *time.Time
	DestroyedAt *time.Time
}
