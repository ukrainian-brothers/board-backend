package api

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	internal_user "github.com/ukrainian-brothers/board-backend/internal/user"
	"testing"
)

func TestRegistrationE2E(t *testing.T) {
	userRepo, db := getPostgresRepo()
	server, client := createUserAPI(userRepo)

	type expected struct {
		status      int
		errorStruct errorStruct
	}
	type testCase struct {
		name     string
		payload  registerPayload
		pre      func(t *testing.T, payload registerPayload)
		cleanUp  func(t *testing.T, payload registerPayload)
		expected expected
	}

	testCases := []testCase{
		{
			name:    "no payload",
			payload: registerPayload{},
			pre:     func(t *testing.T, payload registerPayload) {},
			cleanUp: func(t *testing.T, payload registerPayload) {},
			expected: expected{
				status: 422,
				errorStruct: errorStruct{
					Error:   "Unprocessable Entity",
					Details: "missing contact details",
				},
			},
		},
		{
			name:    "no personal data",
			payload: registerPayload{
				Login:     "the_new_user2115",
				Password:  "pass",
				Mail:      "mrosiak@wp.pl",
				Phone:     "+48 111 222 333",
			},
			pre:     func(t *testing.T, payload registerPayload) {},
			cleanUp: func(t *testing.T, payload registerPayload) {},
			expected: expected{
				status: 422,
				errorStruct: errorStruct{
					Error:   "Unprocessable Entity",
				},
			},
		},
		{
			name: "success",
			payload: registerPayload{
				Login:     "the_new_user2115",
				Password:  "pass",
				Firstname: "Mac",
				Surname:   "Cheese",
				Mail:      "mrosiak@wp.pl",
				Phone:     "+48 111 222 333",
			},
			pre: func(t *testing.T, payload registerPayload) {},
			cleanUp: func(t *testing.T, payload registerPayload) {
				_, err := db.Exec("DELETE FROM users WHERE login=$1", payload.Login)
				assert.NoError(t, err)
			},
			expected: expected{
				status: 201,
			},
		},
		{
			name: "already exists",
			payload: registerPayload{
				Login:     "the_new_user2115",
				Password:  "pass",
				Firstname: "Mac",
				Surname:   "Cheese",
				Mail:      "mrosiak@wp.pl",
				Phone:     "+48 111 222 333",
			},
			pre: func(t *testing.T, payload registerPayload) {
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
			cleanUp: func(t *testing.T, payload registerPayload) {
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
			resp := doRequest(t, client, "POST", fmt.Sprintf("%s/api/user/register", server.URL), tC.payload, &responseStruct)
			assert.Equal(t, tC.expected.status, resp.StatusCode)
			assert.Equal(t, tC.expected.errorStruct.Error, responseStruct.Error)
			assert.Equal(t, tC.expected.errorStruct.Details, responseStruct.Details)

			tC.cleanUp(t, tC.payload)
		})
	}
}

func TestRegistration(t *testing.T) {
	userRepo := getMockedRepo()
	server, client := createUserAPI(&userRepo)

	type expected struct {
		status      int
		errorStruct errorStruct
	}

	type testCase struct {
		name     string
		mock     func()
		payload  registerPayload
		expected expected
	}

	testCases := []testCase{
		{
			name: "UserExists query internal DB error",
			mock: func() {
				userRepo.On("Exists", mock.Anything, mock.Anything).Return(false, errors.New("err"))
			},
			payload: registerPayload{
				Firstname: "Mac",
				Surname: "Smith",
				Mail:  "the_mail",
				Phone: "+48 111 222 333",
			},
			expected: expected{
				status: 500,
				errorStruct: errorStruct{
					Error:   "Internal Server Error",
				},
			},
		},
		{
			name: "AddUser command internal DB error",
			mock: func() {
				userRepo.On("Exists", mock.Anything, mock.Anything).Return(false, nil)
				userRepo.On("Add", mock.Anything, mock.Anything).Return(errors.New("err"))
			},
			payload: registerPayload{
				Firstname: "Mac",
				Surname: "Smith",
				Mail:  "the_mail",
				Phone: "+48 111 222 333",
			},
			expected: expected{
				status: 500,
				errorStruct: errorStruct{
					Error:   "Internal Server Error",
				},
			},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			tC.mock()

			responseStruct := errorStruct{}
			resp := doRequest(t, client, "POST", fmt.Sprintf("%s/api/user/register", server.URL), tC.payload, &responseStruct)

			assert.Equal(t, tC.expected.status, resp.StatusCode)
			assert.Equal(t, tC.expected.errorStruct.Error, responseStruct.Error)
			assert.Equal(t, tC.expected.errorStruct.Details, responseStruct.Details)
		})
	}
}
