package api

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/ukrainian-brothers/board-backend/domain/user"
	internal_user "github.com/ukrainian-brothers/board-backend/internal/user"
	"net/http"
	"testing"
)

func TestRegistrationE2E(t *testing.T) {
	userRepo, advertRepo, db := getPostgresRepos(t)
	server, client, _ := createTestAPIs(t, advertRepo, userRepo)

	type expected struct {
		status      int
		errorStruct errorStruct
	}
	type testCase struct {
		name     string
		payload  *registerPayload
		pre      func(t *testing.T, payload *registerPayload)
		cleanUp  func(t *testing.T, payload *registerPayload)
		expected expected
	}

	testCases := []testCase{
		{
			name:    "no payload",
			payload: &registerPayload{},
			pre:     func(t *testing.T, payload *registerPayload) {},
			cleanUp: func(t *testing.T, payload *registerPayload) {},
			expected: expected{
				status: 422,
				errorStruct: errorStruct{
					Error:   "Unprocessable Entity",
					Details: "missing contact details",
				},
			},
		},
		{
			name: "no personal data",
			payload: &registerPayload{
				Login:    "the_new_user2115",
				Password: "pass",
				Mail:     "mrosiak@wp.pl",
				Phone:    "+48 111 222 333",
			},
			pre:     func(t *testing.T, payload *registerPayload) {},
			cleanUp: func(t *testing.T, payload *registerPayload) {},
			expected: expected{
				status: 422,
				errorStruct: errorStruct{
					Error: "Unprocessable Entity",
				},
			},
		},
		{
			name: "success",
			payload: &registerPayload{
				Login:     "the_new_user2115",
				Password:  "pass",
				Firstname: "Mac",
				Surname:   "Cheese",
				Mail:      "mrosiak@wp.pl",
				Phone:     "+48 111 222 333",
			},
			pre: func(t *testing.T, payload *registerPayload) {},
			cleanUp: func(t *testing.T, payload *registerPayload) {
				_, err := db.Exec("DELETE FROM users WHERE login=$1", payload.Login)
				assert.NoError(t, err)
			},
			expected: expected{
				status: 201,
			},
		},
		{
			name: "already exists",
			payload: &registerPayload{
				Login:     "the_new_user2115",
				Password:  "pass",
				Firstname: "Mac",
				Surname:   "Cheese",
				Mail:      "mrosiak@wp.pl",
				Phone:     "+48 111 222 333",
			},
			pre: func(t *testing.T, payload *registerPayload) {
				usrDB := internal_user.UserDB{
					ID:          uuid.New(),
					Login:       payload.Login,
					Password:    newStringPtr(payload.Password),
					FirstName:   payload.Firstname,
					Surname:     payload.Surname,
					Mail:        newStringPtr(payload.Mail),
					PhoneNumber: newStringPtr(payload.Phone),
				}
				err := db.Insert(&usrDB)
				assert.Nil(t, err)
			},
			cleanUp: func(t *testing.T, payload *registerPayload) {
				_, err := db.Exec("DELETE FROM users WHERE login=$1", payload.Login)
				assert.NoError(t, err)
			},
			expected: expected{
				status: 422,
				errorStruct: errorStruct{
					Error:   "Unprocessable Entity",
					Details: "user already exists",
				},
			},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			tC.pre(t, tC.payload)

			responseStruct := errorStruct{}
			resp := doRequest(t, client, "POST", fmt.Sprintf("%s/api/user/register", server.URL), tC.payload, &responseStruct, nil)
			assert.Equal(t, tC.expected.status, resp.StatusCode)
			assert.Equal(t, tC.expected.errorStruct.Error, responseStruct.Error)
			assert.Equal(t, tC.expected.errorStruct.Details, responseStruct.Details)

			tC.cleanUp(t, tC.payload)
		})
	}
}

func TestRegistration(t *testing.T) {
	userRepo, advertRepo := getMockedRepo()
	server, client, _ := createTestAPIs(t, &advertRepo, &userRepo)

	type expected struct {
		status      int
		errorStruct errorStruct
	}

	type testCase struct {
		name     string
		mock     func()
		payload  *registerPayload
		expected expected
	}

	testCases := []testCase{
		{
			name: "UserExists query internal DB error",
			mock: func() {
				userRepo.On("Exists", mock.Anything, mock.Anything).Return(false, errors.New("err"))
			},
			payload: &registerPayload{
				Firstname: "Mac",
				Surname:   "Smith",
				Mail:      "the_mail",
				Phone:     "+48 111 222 333",
			},
			expected: expected{
				status: 500,
				errorStruct: errorStruct{
					Error: "Internal Server Error",
				},
			},
		},
		{
			name: "AddUser command internal DB error",
			mock: func() {
				userRepo.On("Exists", mock.Anything, mock.Anything).Return(false, nil)
				userRepo.On("Add", mock.Anything, mock.Anything).Return(errors.New("err"))
			},
			payload: &registerPayload{
				Firstname: "Mac",
				Surname:   "Smith",
				Mail:      "the_mail",
				Phone:     "+48 111 222 333",
			},
			expected: expected{
				status: 500,
				errorStruct: errorStruct{
					Error: "Internal Server Error",
				},
			},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			tC.mock()

			responseStruct := errorStruct{}
			resp := doRequest(t, client, "POST", fmt.Sprintf("%s/api/user/register", server.URL), tC.payload, &responseStruct, nil)

			assert.Equal(t, tC.expected.status, resp.StatusCode)
			assert.Equal(t, tC.expected.errorStruct.Error, responseStruct.Error)
			assert.Equal(t, tC.expected.errorStruct.Details, responseStruct.Details)
		})
	}
}

func TestDecodingError(t *testing.T) {
	userRepo, advertRepo := getMockedRepo()
	server, client, _ := createTestAPIs(t, &advertRepo, &userRepo)
	responseStruct := errorStruct{}

	endpoints := []string{"/api/user/register", "/api/user/login"}
	for _, endpoint := range endpoints {
		resp := doRequest(t, client, "POST", server.URL+endpoint, byte(1), &responseStruct, nil)
		assert.Equal(t, 422, resp.StatusCode)
		assert.Equal(t, "Unprocessable Entity", responseStruct.Error)
		assert.Equal(t, "invalid payload", responseStruct.Details)
	}

}

func TestLoginE2E(t *testing.T) {
	userRepo, advertRepo, db := getPostgresRepos(t)
	server, client, sessionStore := createTestAPIs(t, advertRepo, userRepo)

	type expected struct {
		sessionExists bool
		status        int
		errorStruct   errorStruct
	}

	type testCase struct {
		name     string
		payload  loginPayload
		pre      func(t *testing.T, payload loginPayload)
		cleanUp  func(t *testing.T, payload loginPayload)
		expected expected
	}

	testCases := []testCase{
		{
			name: "success",
			payload: loginPayload{
				Login:    "the_test_user",
				Password: "secret_pass",
			},
			pre: func(t *testing.T, payload loginPayload) {
				usrDB := internal_user.UserDB{
					ID:        uuid.New(),
					Login:     payload.Login,
					Password:  newStringPtr("$argon2id$v=19$m=65536,t=3,p=2$2AhHwyZVY7yNE8PJjOOIrg$ZvmY83U2SXVdzKl3CKM7z8U1R8CzFj3HO5J2p4LDBXo"),
					FirstName: "Mac",
					Surname:   "Cheese",
					Mail:      newStringPtr("macncheese@wp.pl"),
				}
				err := db.Insert(&usrDB)
				assert.NoError(t, err)
			},
			cleanUp: func(t *testing.T, payload loginPayload) {
				_, err := db.Exec("DELETE FROM users WHERE login=$1", payload.Login)
				assert.NoError(t, err)
			},
			expected: expected{
				sessionExists: true,
				status:        200,
			},
		},
		{
			name: "invalid credentials",
			payload: loginPayload{
				Login:    "the_test_user",
				Password: "wrong_password",
			},
			pre: func(t *testing.T, payload loginPayload) {
				usrDB := internal_user.UserDB{
					ID:        uuid.New(),
					Login:     payload.Login,
					Password:  newStringPtr("$argon2id$v=19$m=65536,t=3,p=2$2AhHwyZVY7yNE8PJjOOIrg$ZvmY83U2SXVdzKl3CKM7z8U1R8CzFj3HO5J2p4LDBXo"),
					FirstName: "Mac",
					Surname:   "Cheese",
					Mail:      newStringPtr("macncheese@wp.pl"),
				}
				err := db.Insert(&usrDB)
				assert.NoError(t, err)
			},
			cleanUp: func(t *testing.T, payload loginPayload) {
				_, err := db.Exec("DELETE FROM users WHERE login=$1", payload.Login)
				assert.NoError(t, err)
			},
			expected: expected{
				sessionExists: false,
				status:        403,
				errorStruct: errorStruct{
					Error:   "Forbidden",
					Details: "wrong credentials",
				},
			},
		},
		{
			name: "user doesnt exists",
			payload: loginPayload{
				Login:    "the_test_user",
				Password: "password",
			},
			pre:     func(t *testing.T, payload loginPayload) {},
			cleanUp: func(t *testing.T, payload loginPayload) {},
			expected: expected{
				sessionExists: false,
				status:        422,
				errorStruct: errorStruct{
					Error:   "Unprocessable Entity",
					Details: "user does not exists",
				},
			},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			tC.pre(t, tC.payload)

			responseStruct := errorStruct{}
			resp := doRequest(t, client, "POST", fmt.Sprintf("%s/api/user/login", server.URL), tC.payload, &responseStruct, []*http.Cookie{})

			assert.Equal(t, tC.expected.status, resp.StatusCode)
			assert.Equal(t, tC.expected.errorStruct.Details, responseStruct.Details)
			assert.Equal(t, tC.expected.errorStruct.Error, responseStruct.Error)

			if tC.expected.sessionExists {
				resp.Request.Header["Cookie"] = resp.Header["Set-Cookie"]

				session, err := sessionStore.Get(resp.Request, getTestConfig(t).Session.SessionKey)
				assert.NoError(t, err)

				assert.Equal(t, tC.payload.Login, session.Values["user_login"])
			}

			tC.cleanUp(t, tC.payload)
		})
	}
}

func TestLogin(t *testing.T) {
	userRepo, advertRepo := getMockedRepo()
	server, client, _ := createTestAPIs(t, &advertRepo, &userRepo)

	type expected struct {
		status      int
		errorStruct errorStruct
	}

	type testCase struct {
		name     string
		payload  loginPayload
		mock     func(payload loginPayload)
		expected expected
	}

	testCases := []testCase{
		{
			name: "failed at UserExists query",
			payload: loginPayload{
				Login:    "login",
				Password: "password",
			},
			mock: func(payload loginPayload) {
				userRepo.On("Exists", mock.Anything, mock.Anything).Return(false, errors.New("x"))
			},
			expected: expected{
				status: 500,
				errorStruct: errorStruct{
					Error: "Internal Server Error",
				},
			},
		},
		{
			name: "failed at VerifyUserPassword query",
			payload: loginPayload{
				Login:    "login",
				Password: "password",
			},
			mock: func(payload loginPayload) {
				userRepo.On("Exists", mock.Anything, mock.Anything).Return(true, nil)
				userRepo.On("GetByLogin", mock.Anything, mock.Anything).Return(&user.User{}, errors.New("x"))
			},
			expected: expected{
				status: 500,
				errorStruct: errorStruct{
					Error: "Internal Server Error",
				},
			},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			tC.mock(tC.payload)
			responseStruct := errorStruct{}
			resp := doRequest(t, client, "POST", fmt.Sprintf("%s/api/user/login", server.URL), tC.payload, &responseStruct, []*http.Cookie{})

			assert.Equal(t, tC.expected.status, resp.StatusCode)
			assert.Equal(t, tC.expected.errorStruct.Details, responseStruct.Details)
			assert.Equal(t, tC.expected.errorStruct.Error, responseStruct.Error)
		})
	}
}
