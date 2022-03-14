package user_test

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/ukrainian-brothers/board-backend/domain"
	"github.com/ukrainian-brothers/board-backend/domain/user"
	"github.com/ukrainian-brothers/board-backend/internal/common"
	internalUser "github.com/ukrainian-brothers/board-backend/internal/user"
	"testing"
)

func createString(s string) *string {
	return &s
}

func TestUserPostgresAdd(t *testing.T) {
	cfg, err := common.NewConfigFromFile("../../config/configuration.test.local.json") // TODO: Describe how to achive this config in README.MD
	assert.NoError(t, err)
	db := common.InitPostgres(&cfg.Postgres)
	repo := internalUser.NewPostgresUserRepository(db)

	type testCase struct {
		name        string
		user        *user.User
		cleanUp     func(t *testing.T, id string)
		expectedErr error
	}

	testCases := []testCase{
		{
			name: "Success",
			user: &user.User{
				ID:       uuid.MustParse("6858fe22-2c04-4a13-bc75-eafeeb3cf767"),
				Login:    "Adam",
				Password: createString("awddwaawd"),
				Person: domain.Person{
					FirstName: "Adam",
					Surname:   "Małysz",
				},
				ContactDetails: domain.ContactDetails{
					Mail:        createString("adam@wp.pl"),
					PhoneNumber: createString("111222333"),
				},
			},
			cleanUp: func(t *testing.T, id string) {
				_, err := db.Exec("DELETE FROM users WHERE id=$1", id)
				assert.NoError(t, err)
			},
			expectedErr: nil,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			err := repo.Add(context.Background(), test.user)
			assert.Equal(t, test.expectedErr, err)
			test.cleanUp(t, test.user.ID.String())
		})
	}
}

func TestGetById(t *testing.T) {
	cfg, err := common.NewConfigFromFile("../../config/configuration.test.local.json") // TODO: Describe how to achive this config in README.MD
	assert.NoError(t, err)
	db := common.InitPostgres(&cfg.Postgres)
	repo := internalUser.NewPostgresUserRepository(db)

	type testCase struct {
		name        string
		pre         func(t *testing.T) (result user.User)
		cleanUp     func(t *testing.T, id string)
		expectedErr error
	}

	testCases := []testCase{
		{
			name: "EXISTING_USER",
			pre: func(t *testing.T) (result user.User) {
				var user = user.User{
					ID:       uuid.MustParse("69129a87-cccb-49f0-98c8-fc9b7a5e04dc"),
					Login:    "foo",
					Password: createString("foobar"),
					Person: domain.Person{
						FirstName: "Foo",
						Surname:   "Bar",
					},
					ContactDetails: domain.ContactDetails{
						Mail:        createString("foo@gmail.com"),
						PhoneNumber: createString("+482222222"),
					},
				}

				err = db.Insert(&user)
				assert.NoError(t, err)

				return user
			},
			cleanUp: func(t *testing.T, id string) {
				_, err := db.Exec("DELETE FROM users WHERE id=$1", id)
				assert.NoError(t, err)
			},
			expectedErr: nil,
		},
		{
			name: "NOT_EXISTING_USER",
			pre: func(t *testing.T) (result user.User) {
				result.ID = uuid.MustParse("5e8d9731-2437-4ba0-810e-468f983e1a0b")

				return result
			},
			cleanUp: func(t *testing.T, id string) {
			},
			expectedErr: sql.ErrNoRows,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			t.Parallel()

			testUser := tC.pre(t)

			_, err := repo.GetByID(context.Background(), testUser.ID)
			assert.ErrorIs(t, err, tC.expectedErr)

			tC.cleanUp(t, testUser.ID.String())
		})
	}
}
