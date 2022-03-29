package user_test

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	cfg, err := common.NewConfigFromFile("../../config/configuration.test.local.json")
	assert.NoError(t, err)

	db, err := common.InitPostgres(&cfg.Postgres)
	require.NoError(t, err)

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
	cfg, err := common.NewConfigFromFile("../../config/configuration.test.local.json")
	assert.NoError(t, err)

	db, err := common.InitPostgres(&cfg.Postgres)
	require.NoError(t, err)

	repo := internalUser.NewPostgresUserRepository(db)

	type testCase struct {
		name        string
		user        *user.User
		pre         func(t *testing.T, user *user.User)
		cleanUp     func(t *testing.T, id string)
		expectedErr error
	}

	testCases := []testCase{
		{
			name: "EXISTING_USER",
			user: &user.User{
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
			},
			pre: func(t *testing.T, user *user.User) {
				usr := internalUser.UserDB{
					ID:          user.ID,
					Login:       user.Login,
					Password:    user.Password,
					FirstName:   user.Person.FirstName,
					Surname:     user.Person.Surname,
					Mail:        createString("foo@gmail.com"),
					PhoneNumber: createString("+482222222"),
				}

				err = db.Insert(&usr)
				assert.NoError(t, err)
			},
			cleanUp: func(t *testing.T, id string) {
				_, err := db.Exec("DELETE FROM users WHERE id=$1", id)
				assert.NoError(t, err)
			},
			expectedErr: nil,
		},
		{
			name: "NOT_EXISTING_USER",
			user: &user.User{
				ID: uuid.MustParse("5e8d9731-2437-4ba0-810e-468f983e1a0b"),
			},
			cleanUp:     func(t *testing.T, id string) {},
			pre:         func(t *testing.T, user *user.User) {},
			expectedErr: sql.ErrNoRows,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			tC.pre(t, tC.user)

			usr, err := repo.GetByID(context.Background(), tC.user.ID)
			assert.ErrorIs(t, err, tC.expectedErr)
			if tC.expectedErr == nil {
				user.Assert(t, tC.user, usr)
			}

			tC.cleanUp(t, tC.user.ID.String())
		})
	}
}

func TestGetByLogin(t *testing.T) { // TODO: Fix this test..
	cfg, err := common.NewConfigFromFile("../../config/configuration.test.local.json")
	assert.NoError(t, err)

	db, err := common.InitPostgres(&cfg.Postgres)
	require.NoError(t, err)

	repo := internalUser.NewPostgresUserRepository(db)

	type testCase struct {
		name        string
		pre         func(t *testing.T) (result *user.User)
		cleanUp     func(t *testing.T, id string)
		expectedErr error
	}

	testCases := []testCase{
		{
			name: "EXISTING_USER",
			pre: func(t *testing.T) (result *user.User) {
				var usr = user.User{
					ID:       uuid.MustParse("69129a87-cccb-49f0-98c8-fc9b7a5e04dc"),
					Login:    "the_login",
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

				err = db.Insert(&usr)
				assert.NoError(t, err)

				return &usr
			},
			cleanUp: func(t *testing.T, id string) {
				_, err := db.Exec("DELETE FROM users WHERE id=$1", id)
				assert.NoError(t, err)
			},
			expectedErr: nil,
		},
		{
			name: "NOT_EXISTING_USER",
			pre: func(t *testing.T) (result *user.User) {
				result = &user.User{
					ID: uuid.MustParse("5e8d9731-2437-4ba0-810e-468f983e1a0b"),
				}

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

			usr, err := repo.GetByLogin(context.Background(), testUser.Login)
			assert.ErrorIs(t, err, tC.expectedErr)
			if tC.expectedErr == nil {
				user.Assert(t, testUser, usr)
			}

			tC.cleanUp(t, testUser.ID.String())
		})
	}
}

func TestUserExists(t *testing.T) {
	cfg, err := common.NewConfigFromFile("../../config/configuration.test.local.json")
	assert.NoError(t, err)

	db, err := common.InitPostgres(&cfg.Postgres)
	require.NoError(t, err)

	repo := internalUser.NewPostgresUserRepository(db)

	type expected struct {
		err    error
		exists bool
	}

	type testCase struct {
		name     string
		user     *user.User
		pre      func(t *testing.T, usr *user.User)
		cleanUp  func(t *testing.T, usr *user.User)
		expected expected
	}

	testCases := []testCase{
		{
			name: "exists",
			user: &user.User{
				ID:       uuid.New(),
				Login:    "login",
				Password: createString("pass"),
				Person:   domain.Person{"abc", "dawdwa"},
				ContactDetails: domain.ContactDetails{
					Mail: createString("aaaa@wp.pl"),
				},
			},
			pre: func(t *testing.T, usr *user.User) {
				userDB := internalUser.UserDB{}
				userDB.LoadUser(usr)
				err := db.Insert(&userDB)
				assert.NoError(t, err)
			},
			cleanUp: func(t *testing.T, usr *user.User) {
				_, err := db.Exec("DELETE FROM users WHERE id=$1", usr.ID)
				assert.NoError(t, err)
			},
			expected: expected{
				err: nil,
				exists: true,
			},
		},
		{
			name: "not exists",
			user: &user.User{
				Login:    "login",
			},
			pre: func(t *testing.T, usr *user.User) {},
			cleanUp: func(t *testing.T, usr *user.User) {},
			expected: expected{
				err: nil,
				exists: false,
			},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			tC.pre(t, tC.user)
			exists, err := repo.Exists(context.Background(), tC.user.Login)
			assert.Equal(t, tC.expected.err, err)
			assert.Equal(t, tC.expected.exists, exists)
			tC.cleanUp(t, tC.user)

		})
	}
}
