package advert

import (
	"encoding/json"
	"errors"
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

var (
	ContactEmptyErr     = errors.New("contact is empty")
	NoUserProvidedErr   = errors.New("no user provided")
	MissingBasicInfoErr = errors.New("advert is missing basic info")
)

type AdvertLog struct {
	AdvertID uuid.UUID
	Trigger  AdvertLogTrigger
	Meta     json.RawMessage
}

type Advert struct {
	ID          uuid.UUID
	Details     domain.AdvertDetails
	User        *user.User
	CreatedAt   time.Time
	UpdatedAt   *time.Time
	DestroyedAt *time.Time
}

type AdvertOption func(advert *Advert) error

func WithContactDetails(contactDetails domain.ContactDetails) AdvertOption {
	return func(advert *Advert) error {
		if contactDetails.IsEmpty() {
			return ContactEmptyErr
		}
		advert.Details.ContactDetails = contactDetails
		return nil
	}
}

func NewAdvert(user *user.User, title string, description string, advertType domain.AdvertType, opts ...AdvertOption) (*Advert, error) {
	if user == nil {
		return nil, NoUserProvidedErr
	}

	advert := &Advert{ID: uuid.New()}

	for _, option := range opts {
		err := option(advert)
		if err != nil {
			return nil, err
		}
	}

	if title == "" || description == "" {
		return nil, MissingBasicInfoErr
	}
	advert.User = user
	advert.Details.Title = title
	advert.Details.Description = description
	advert.Details.Type = advertType

	if advert.Details.ContactDetails.IsEmpty() {
		advert.Details.ContactDetails = user.ContactDetails
	}

	return advert, nil
}
