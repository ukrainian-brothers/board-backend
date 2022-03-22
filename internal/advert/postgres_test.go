package advert

import (
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/ukrainian-brothers/board-backend/domain"
	"github.com/ukrainian-brothers/board-backend/domain/advert"
	"github.com/ukrainian-brothers/board-backend/domain/user"
	"github.com/ukrainian-brothers/board-backend/internal/common"
	internalUser "github.com/ukrainian-brothers/board-backend/internal/user"
	"testing"
	"time"
)

func getContactDetails() domain.ContactDetails {
	return domain.ContactDetails{
		Mail:        newStringPtr("foo@gmail.com"),
		PhoneNumber: newStringPtr("+482222222"),
	}
}

func TestAdvertPostgresAdd(t *testing.T) {
	cfg, err := common.NewConfigFromFile("../../config/configuration.test.local.json")
	assert.NoError(t, err)

	db, err := common.InitPostgres(&cfg.Postgres)
	require.NoError(t, err)

	repo := NewPostgresAdvertRepository(db)
	db.AddTableWithName(internalUser.UserDB{}, "users").SetKeys(false, "id")

	type testCase struct {
		name        string
		advert      *advert.Advert
		pre         func(t *testing.T, adv *advert.Advert)
		cleanUp     func(t *testing.T, advertID string, userID string)
		expectedErr error
	}

	testCases := []testCase{
		{
			name: "Success",
			advert: &advert.Advert{
				ID: uuid.New(),
				Details: domain.AdvertDetails{
					Title:          "foo",
					Description:    "bar",
					Type:           domain.AdvertTypePlaceToStay,
					ContactDetails: getContactDetails(),
				},
				User: &user.User{
					ID:       uuid.MustParse("69129a87-cccb-49f0-98c8-fc9b7a5e04dc"),
					Login:    "the_login",
					Password: newStringPtr("foobar"),
					Person: domain.Person{
						FirstName: "Foo",
						Surname:   "Bar",
					},
					ContactDetails: getContactDetails(),
				},
				CreatedAt: time.Now(),
			},
			pre: func(t *testing.T, adv *advert.Advert) {
				usr := adv.User
				usrForDB := internalUser.UserDB{
					ID:          usr.ID,
					Login:       usr.Login,
					Password:    usr.Password,
					FirstName:   usr.Person.FirstName,
					Surname:     usr.Person.Surname,
					Mail:        usr.ContactDetails.Mail,
					PhoneNumber: usr.ContactDetails.PhoneNumber,
				}
				err = db.Insert(&usrForDB)
				assert.NoError(t, err)
			},
			cleanUp: func(t *testing.T, advertID string, userID string) {
				_, err := db.Exec("DELETE FROM adverts WHERE id=$1", advertID)
				assert.NoError(t, err)
				_, err = db.Exec("DELETE FROM users WHERE id=$1", userID)
				assert.NoError(t, err)
			},
			expectedErr: nil,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			test.pre(t, test.advert)
			err := repo.Add(context.Background(), test.advert)
			test.cleanUp(t, test.advert.ID.String(), test.advert.User.ID.String())
			assert.ErrorIs(t, test.expectedErr, err)
		})
	}
}
