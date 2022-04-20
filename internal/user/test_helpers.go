package user

import (
	"context"
	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/stretchr/testify/assert"
	"github.com/ukrainian-brothers/board-backend/domain"
	"github.com/ukrainian-brothers/board-backend/domain/user"
	"github.com/ukrainian-brothers/board-backend/pkg/test_helpers"
	"net/http"
	"net/http/httptest"
	"testing"
)

func GetValidContactDetails() domain.ContactDetails {
	return domain.ContactDetails{
		Mail:        test_helpers.NewStringPtr("adam@wp.pl"),
		PhoneNumber: test_helpers.NewStringPtr("111222333"),
	}
}

func CreateTestUser(t *testing.T, login string, userRepo user.Repository) *user.User {
	usr := &user.User{
		ID:       uuid.New(),
		Login:    login,
		Password: test_helpers.NewStringPtr("2Fzs0V!@~4m;'13.!#"),
		Person: domain.Person{
			FirstName: "Mac",
			Surname:   "Cheese",
		},
		ContactDetails: GetValidContactDetails(),
	}
	err := userRepo.Add(context.Background(), usr)
	assert.NoError(t, err)
	return usr
}

func RemoveTestUser(t *testing.T, id uuid.UUID, userRepo user.Repository) {
	err := userRepo.Delete(context.Background(), id)
	assert.NoError(t, err)
}

func CreateTestSession(t *testing.T, usr *user.User, store sessions.Store) []*http.Cookie {
	cfg := test_helpers.GetTestConfig(t)
	r := &http.Request{}
	session, err := store.Get(r, cfg.Session.SessionKey)
	assert.NoError(t, err)

	session.Values["user_login"] = usr.Login
	w := httptest.NewRecorder()
	err = session.Save(r, w)
	assert.NoError(t, err)

	return w.Result().Cookies()
}

func GenerateUserDB(id uuid.UUID) UserDB {
	usr := UserDB{
		ID:        id,
		Login:     test_helpers.RandomString(10),
		FirstName: test_helpers.RandomString(10),
		Surname:   test_helpers.RandomString(10),
		Mail:      test_helpers.RandomMail(),
	}
	return usr
}
