package api

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/ukrainian-brothers/board-backend/domain"
	"net/http"
	"testing"
)

func TestAddAdvertE2E(t *testing.T) {
	userRepo, advertRepo, db := getPostgresRepos(t)
	server, client, sessionStore := createTestAPIs(t, advertRepo, userRepo)

	type expected struct {
		status      int
		errorStruct errorStruct
	}
	type testCase struct {
		name     string
		loggedIn bool
		payload  newAdvertPayload
		expected expected
		cleanUp  func(t *testing.T, userID uuid.UUID)
	}

	testCases := []testCase{
		{
			name: "not authorised",
			expected: expected{
				status: http.StatusForbidden,
				errorStruct: errorStruct{
					Error:   "Forbidden",
					Details: "not authorized",
				},
			},
			loggedIn: false,
			cleanUp:  func(t *testing.T, userID uuid.UUID) {},
		},
		{
			name: "invalid payload",
			expected: expected{
				status: http.StatusUnprocessableEntity,
				errorStruct: errorStruct{
					Error:   "Unprocessable Entity",
					Details: "invalid payload",
				},
			},
			loggedIn: true,
			cleanUp:  func(t *testing.T, userID uuid.UUID) {},
		},
		{
			name: "invalid title",
			expected: expected{
				status: http.StatusUnprocessableEntity,
				errorStruct: errorStruct{
					Error:   "Unprocessable Entity",
					Details: "invalid advert details",
				},
			},
			payload: newAdvertPayload{
				ContactDetails: contactPayload{
					Mail: "mmmmmm@wp.pl",
				},
			},
			loggedIn: true,
			cleanUp:  func(t *testing.T, userID uuid.UUID) {},
		},
		{
			name:     "success",
			loggedIn: true,
			payload: newAdvertPayload{
				Title:       "title",
				Description: "description",
				Type:        domain.AdvertTypeTransport,
				ContactDetails: contactPayload{
					Mail: "mmmmmm@wp.pl",
				},
			},
			expected: expected{
				status: http.StatusCreated,
			},
			cleanUp: func(t *testing.T, userID uuid.UUID) {
				_, err := db.Exec("DELETE FROM adverts WHERE user_id=$1", userID.String())
				assert.NoError(t, err)
			},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			var cookies []*http.Cookie
			if tC.loggedIn {
				testUser := createTestUser(t, "test_login1213930", userRepo)
				cookies = createTestSession(t, testUser, sessionStore)
				defer removeTestUser(t, testUser.ID, userRepo)
				defer tC.cleanUp(t, testUser.ID)
			}

			errResponse := errorStruct{}
			resp := doRequest(t, client, "POST", fmt.Sprintf("%s/api/adverts", server.URL), tC.payload, &errResponse, cookies)
			assert.Equal(t, tC.expected.status, resp.StatusCode)
			assert.Equal(t, tC.expected.errorStruct.Error, errResponse.Error)
			assert.Equal(t, tC.expected.errorStruct.Details, errResponse.Details)

		})
	}
}

func TestAddAdvert(t *testing.T) {
	// TODO: In another branch add tests for handling errors from DB
}
