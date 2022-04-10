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
	. "github.com/ukrainian-brothers/board-backend/pkg/translation"
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
					Title:          MultilingualString{Ukrainian: "x"},
					Description:    MultilingualString{Ukrainian: "x"},
					Type:           domain.AdvertTypePlaceToStay,
					ContactDetails: getContactDetails(),
				},
				User: &user.User{
					ID:       uuid.New(),
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
				usrDB := internalUser.UserDB{}
				usrDB.LoadUser(adv.User)
				err = db.Insert(&usrDB)
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

	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			tC.pre(t, tC.advert)
			err := repo.Add(context.Background(), tC.advert)
			tC.cleanUp(t, tC.advert.ID.String(), tC.advert.User.ID.String())
			assert.ErrorIs(t, tC.expectedErr, err)
		})
	}
}

func TestAdvertPostgresGet(t *testing.T) {
	cfg, err := common.NewConfigFromFile("../../config/configuration.test.local.json")
	assert.NoError(t, err)

	db, err := common.InitPostgres(&cfg.Postgres)
	require.NoError(t, err)

	repo := NewPostgresAdvertRepository(db)
	db.AddTableWithName(internalUser.UserDB{}, "users").SetKeys(false, "id")

	type input struct {
		userDB          internalUser.UserDB
		advertDB        advertDB
		advertDetailsDB []advertDetailsDB
	}
	type testCase struct {
		name        string
		input       input
		pre         func(t *testing.T, input input)
		cleanUp     func(t *testing.T, input input)
		expectedErr error
	}

	testCases := []testCase{
		{
			name: "success",
			input: input{
				userDB: internalUser.UserDB{
					ID:        uuid.MustParse("38e520dc-ac8c-44a6-be74-0c3bfb7a4576"),
					Login:     "login",
					Password:  newStringPtr("passwordC1$23"),
					FirstName: "Mac",
					Surname:   "Cheese",
					Mail:      newStringPtr("mail@wp.pl"),
				},
				advertDB: advertDB{
					ID:     uuid.MustParse("e8e1f982-992d-40c9-8389-1ca147c97ecd"),
					UserID: uuid.MustParse("38e520dc-ac8c-44a6-be74-0c3bfb7a4576"),
					Type:   domain.AdvertTypeTransport,
					ContactDetails: domain.ContactDetails{
						Mail: newStringPtr("mail@wp.pl"),
					},
				},
				advertDetailsDB: []advertDetailsDB{
					{
						ID:          uuid.New(),
						AdvertID:    uuid.MustParse("e8e1f982-992d-40c9-8389-1ca147c97ecd"),
						Language:    Ukrainian,
						Title:       "титул",
						Description: "опис",
					},
				},
			},
			pre: func(t *testing.T, input input) {
				err := db.Insert(&input.userDB, &input.advertDB)
				assert.NoError(t, err)

				for _, detailsDb := range input.advertDetailsDB {
					err := db.Insert(&detailsDb)
					assert.NoError(t, err)
				}
			},
			cleanUp: func(t *testing.T, input input) {
				_, err := db.Exec("DELETE FROM adverts WHERE id=$1", input.advertDB.ID)
				// advert_details should be removed due to fk policy
				assert.NoError(t, err)

				_, err = db.Exec("DELETE FROM users WHERE id=$1", input.userDB.ID)
				assert.NoError(t, err)
			},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			if tC.pre != nil {
				tC.pre(t, tC.input)
			}
			if tC.cleanUp != nil {
				defer tC.cleanUp(t, tC.input)
			}

			_, err := repo.Get(context.Background(), tC.input.advertDB.ID)
			assert.Equal(t, tC.expectedErr, err)

		})
	}
}

func TestAdvertPostgresGetList(t *testing.T) {
	cfg, err := common.NewConfigFromFile("../../config/configuration.test.local.json")
	assert.NoError(t, err)

	db, err := common.InitPostgres(&cfg.Postgres)
	require.NoError(t, err)

	repo := NewPostgresAdvertRepository(db)
	db.AddTableWithName(internalUser.UserDB{}, "users").SetKeys(false, "id")

	type dbInput struct {
		userDB          internalUser.UserDB
		advertDB        []advertDB
		advertDetailsDB []advertDetailsDB
	}

	type input struct {
		langs  []LanguageTag
		limit  int
		offset int
	}

	type expected struct {
		err        error
		advertsLen int
	}

	type testCase struct {
		name     string
		dbInput  dbInput
		input    input
		pre      func(t *testing.T, input dbInput)
		cleanUp  func(t *testing.T, input dbInput)
		expected expected
	}

	testCases := []testCase{
		{
			name: "success",
			dbInput: dbInput{
				userDB: internalUser.UserDB{
					ID:        uuid.MustParse("38e520dc-ac8c-44a6-be74-0c3bfb7a4576"),
					Login:     "login",
					Password:  newStringPtr("passwordC1$23"),
					FirstName: "Mac",
					Surname:   "Cheese",
					Mail:      newStringPtr("mail@wp.pl"),
				},
				advertDB: []advertDB{
					{
						ID:     uuid.MustParse("e8e1f982-992d-40c9-8389-1ca147c97ecd"),
						UserID: uuid.MustParse("38e520dc-ac8c-44a6-be74-0c3bfb7a4576"),
						Type:   domain.AdvertTypeTransport,
						ContactDetails: domain.ContactDetails{
							Mail: newStringPtr("mail@wp.pl"),
						},
					},
					{
						ID:     uuid.MustParse("d69ab3bb-ee0f-41f9-bd38-0865c89953c5"),
						UserID: uuid.MustParse("38e520dc-ac8c-44a6-be74-0c3bfb7a4576"),
						Type:   domain.AdvertTypeTransport,
						ContactDetails: domain.ContactDetails{
							Mail: newStringPtr("mail@wp.pl"),
						},
					},
				},
				advertDetailsDB: []advertDetailsDB{
					{
						ID:          uuid.New(),
						AdvertID:    uuid.MustParse("e8e1f982-992d-40c9-8389-1ca147c97ecd"),
						Language:    Ukrainian,
						Title:       "титул",
						Description: "опис",
					},
					{
						ID:          uuid.New(),
						AdvertID:    uuid.MustParse("e8e1f982-992d-40c9-8389-1ca147c97ecd"),
						Language:    English,
						Title:       "title",
						Description: "description",
					},
					{
						ID:          uuid.New(),
						AdvertID:    uuid.MustParse("e8e1f982-992d-40c9-8389-1ca147c97ecd"),
						Language:    "XX", // unknown language, shouldn't exist in result
						Title:       "title",
						Description: "description",
					},
					{
						ID:          uuid.New(),
						AdvertID:    uuid.MustParse("d69ab3bb-ee0f-41f9-bd38-0865c89953c5"),
						Language:    English,
						Title:       "title",
						Description: "description",
					},
				},
			},
			input: input{
				langs: []LanguageTag{English},
				limit: 10,
			},
			pre: func(t *testing.T, input dbInput) {
				err := db.Insert(&input.userDB)
				assert.NoError(t, err)

				for _, advDB := range input.advertDB {
					err := db.Insert(&advDB)
					assert.NoError(t, err)
				}

				for _, detailsDb := range input.advertDetailsDB {
					err := db.Insert(&detailsDb)
					assert.NoError(t, err)
				}
			},
			cleanUp: func(t *testing.T, input dbInput) {
				for _, advertDB := range input.advertDB {
					_, err := db.Exec("DELETE FROM adverts WHERE id=$1", advertDB.ID)
					// advert_details should be removed due to fk policy
					assert.NoError(t, err)
				}

				_, err = db.Exec("DELETE FROM users WHERE id=$1", input.userDB.ID)
				assert.NoError(t, err)
			},
			expected: expected{
				err:        nil,
				advertsLen: 2,
			},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			if tC.pre != nil {
				tC.pre(t, tC.dbInput)
			}
			if tC.cleanUp != nil {
				defer tC.cleanUp(t, tC.dbInput)
			}

			adverts, err := repo.GetList(context.Background(), tC.input.langs, tC.input.limit, tC.input.offset)
			assert.Equal(t, tC.expected.err, err)
			assert.Equal(t, tC.expected.advertsLen, len(adverts))

		})
	}
}
