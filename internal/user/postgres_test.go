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

func newStringPtr(s string) *string {
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
				Password: newStringPtr("awddwaawd"),
				Person: domain.Person{
					FirstName: "Adam",
					Surname:   "Małysz",
				},
				ContactDetails: domain.ContactDetails{
					Mail:        newStringPtr("adam@wp.pl"),
					PhoneNumber: newStringPtr("111222333"),
				},
			},
			cleanUp: func(t *testing.T, id string) {
				_, err := db.Exec("DELETE FROM users WHERE id=$1", id)
				assert.NoError(t, err)
			},
			expectedErr: nil,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			err := repo.Add(context.Background(), tC.user)
			assert.Equal(t, tC.expectedErr, err)
			tC.cleanUp(t, tC.user.ID.String())
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
			name: "existing user",
			user: &user.User{
				ID:       uuid.New(),
				Login:    "foo",
				Password: newStringPtr("foobar"),
				Person: domain.Person{
					FirstName: "Foo",
					Surname:   "Bar",
				},
				ContactDetails: domain.ContactDetails{
					Mail:        newStringPtr("foo@gmail.com"),
					PhoneNumber: newStringPtr("+482222222"),
				},
			},
			pre: func(t *testing.T, user *user.User) {
				usrDB := internalUser.UserDB{}
				usrDB.LoadUser(user)

				err = db.Insert(&usrDB)
				assert.NoError(t, err)
			},
			cleanUp: func(t *testing.T, id string) {
				_, err := db.Exec("DELETE FROM users WHERE id=$1", id)
				assert.NoError(t, err)
			},
			expectedErr: nil,
		},
		{
			name: "not existing user",
			user: &user.User{
				ID: uuid.New(),
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

func TestGetByLogin(t *testing.T) {
	cfg, err := common.NewConfigFromFile("../../config/configuration.test.local.json")
	assert.NoError(t, err)

	db, err := common.InitPostgres(&cfg.Postgres)
	require.NoError(t, err)

	repo := internalUser.NewPostgresUserRepository(db)

	type testCase struct {
		name        string
		user        *user.User
		pre         func(t *testing.T, usr *user.User)
		cleanUp     func(t *testing.T, usr *user.User)
		expectedErr error
	}

	testCases := []testCase{
		{
			name: "existing user",
			user: &user.User{
				ID:       uuid.New(),
				Login:    "the_login111231",
				Password: newStringPtr("foobar"),
				Person: domain.Person{
					FirstName: "Foo",
					Surname:   "Bar",
				},
				ContactDetails: domain.ContactDetails{
					Mail:        newStringPtr("foo@gmail.com"),
					PhoneNumber: newStringPtr("+482222222"),
				},
			},
			pre: func(t *testing.T, usr *user.User) {
				usrDB := internalUser.UserDB{}
				usrDB.LoadUser(usr)
				err = db.Insert(&usrDB)
				assert.NoError(t, err)
			},
			cleanUp: func(t *testing.T, usr *user.User) {
				_, err := db.Exec("DELETE FROM users WHERE id=$1", usr.ID)
				assert.NoError(t, err)
			},
			expectedErr: nil,
		},
		{
			name: "not existing user",
			user: &user.User{
				ID: uuid.New(),
				Login: "this_user_does_not_exists",
			},
			pre: func(t *testing.T, usr *user.User) {},
			cleanUp: func(t *testing.T, usr *user.User) {},
			expectedErr: sql.ErrNoRows,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			tC.pre(t, tC.user)

			usr, err := repo.GetByLogin(context.Background(), tC.user.Login)
			assert.ErrorIs(t, err, tC.expectedErr)
			if tC.expectedErr == nil {
				user.Assert(t, tC.user, usr)
			}

			tC.cleanUp(t, tC.user)
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
				Password: newStringPtr("pass"),
				Person:   domain.Person{FirstName: "abc", Surname: "dawdwa"},
				ContactDetails: domain.ContactDetails{
					Mail: newStringPtr("aaaa@wp.pl"),
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
				err:    nil,
				exists: true,
			},
		},
		{
			name: "not exists",
			user: &user.User{
				Login: "login",
			},
			pre:     func(t *testing.T, usr *user.User) {},
			cleanUp: func(t *testing.T, usr *user.User) {},
			expected: expected{
				err:    nil,
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
