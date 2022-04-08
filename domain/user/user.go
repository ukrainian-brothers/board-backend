package user

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/ukrainian-brothers/board-backend/domain"
	"time"
)

type User struct {
	ID             uuid.UUID
	Login          string
	Password       *string
	Person         domain.Person
	ContactDetails domain.ContactDetails
}

var (
	MissingPersonalDataErr = errors.New("missing personal data")
	MissingContactDataErr  = errors.New("missing contact data")
)

func NewUser(firstName string, sureName string, login string, password string, contactDetails domain.ContactDetails) (*User, error) {
	// TODO: regex for login, password and write tests

	if firstName == "" || sureName == "" {
		return nil, MissingPersonalDataErr
	}

	if contactDetails.IsEmpty() {
		return nil, MissingContactDataErr
	}

	usr := &User{
		ID:       uuid.New(),
		Login:    login,
		Password: &password,
		Person: domain.Person{
			FirstName: firstName,
			Surname:   sureName,
		},
		ContactDetails: contactDetails,
	}

	return usr, nil
}

type Social struct {
	UserID       uuid.UUID       `json:"user_id"`
	Social       string          `json:"social"`
	SocialId     string          `json:"social_id"`
	AccessToken  string          `json:"access_token"`
	RefreshToken string          `json:"refresh_token"`
	UserData     json.RawMessage `json:"user_data"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
	DestroyedAt  time.Time       `json:"destroyed_at"`
}
