package api

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/ukrainian-brothers/board-backend/domain"
	"github.com/ukrainian-brothers/board-backend/internal"
	"github.com/ukrainian-brothers/board-backend/internal/advert"
	"github.com/ukrainian-brothers/board-backend/internal/user"
	"github.com/ukrainian-brothers/board-backend/pkg/test_helpers"
	. "github.com/ukrainian-brothers/board-backend/pkg/translation"
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
					Mail: *test_helpers.RandomMail(),
				},
			},
			loggedIn: true,
		},
		{
			name:     "success",
			loggedIn: true,
			payload: newAdvertPayload{
				Title:       MultilingualString{English: "x"},
				Description: MultilingualString{English: "x"},
				Type:        domain.AdvertTypeTransport,
				ContactDetails: contactPayload{
					Mail: *test_helpers.RandomMail(),
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
				testUser := user.CreateTestUser(t, "test_login1213930", userRepo)
				cookies = user.CreateTestSession(t, testUser, sessionStore)
				defer user.RemoveTestUser(t, testUser.ID, userRepo)
				if tC.cleanUp != nil {
					defer tC.cleanUp(t, testUser.ID)
				}
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

func TestAdvertsListE2E(t *testing.T) {
	userRepo, advertRepo, db := getPostgresRepos(t)
	server, client, _ := createTestAPIs(t, advertRepo, userRepo)

	type input struct {
		limit  int
		offset int
		langs  LanguageTags
	}

	type inputDB struct {
		userDB          user.UserDB
		advertsDB       []advert.AdvertDB
		advertDetailsDB []advert.AdvertDetailsDB
	}

	type expected struct {
		status                int
		advertResponsesLength int
	}

	type testCase struct {
		name     string
		input    input
		expected expected
		inputDB  inputDB
		pre      func(t *testing.T, inputDB inputDB)
		cleanUp  func(t *testing.T, inputDB inputDB)
	}
	uuid_ := internal.HumanFriendlyUUID
	testCases := []testCase{
		{
			name: "success",
			input: input{
				langs: LanguageTags{},
			},
			expected: expected{
				status:                200,
				advertResponsesLength: 2,
			},
			inputDB: inputDB{
				/*
					User: first_user
						- Advert: first_advert
							- AdvertDetails: Ukrainian
							- AdvertDetails: English
							- AdvertDetails: "xx" // unknown language
						- Advert: second_advert
							- AdvertDetails: Ukrainian
							- AdvertDetails: English
				*/
				userDB: user.GenerateUserDB(uuid_("first_user")),
				advertsDB: []advert.AdvertDB{
					advert.GenerateTestAdvertDB(uuid_("first_advert"), uuid_("first_user")),
					advert.GenerateTestAdvertDB(uuid_("second_advert"), uuid_("first_user")),
				},
				advertDetailsDB: []advert.AdvertDetailsDB{
					advert.GenerateTestAdvertDetailsDB(uuid_("first_advert"), Ukrainian),
					advert.GenerateTestAdvertDetailsDB(uuid_("first_advert"), English),
					advert.GenerateTestAdvertDetailsDB(uuid_("first_advert"), "xx"), // unknown language

					advert.GenerateTestAdvertDetailsDB(uuid_("second_advert"), Ukrainian),
					advert.GenerateTestAdvertDetailsDB(uuid_("second_advert"), English),
				},
			},
			pre: func(t *testing.T, inputDB inputDB) {
				assert.NoError(t, db.Insert(&inputDB.userDB))

				for _, advertDB := range inputDB.advertsDB {
					assert.NoError(t, db.Insert(&advertDB))
				}

				for _, advertDetailsDB := range inputDB.advertDetailsDB {
					assert.NoError(t, db.Insert(&advertDetailsDB))
				}
			},
			cleanUp: func(t *testing.T, inputDB inputDB) {
				for _, advertDB := range inputDB.advertsDB {
					_, err := db.Exec("DELETE FROM adverts WHERE id=$1", advertDB.ID)
					assert.NoError(t, err)
				}

				// advert_details will be removed with adverts because of foreign key constraint

				_, err := db.Exec("DELETE FROM users WHERE login=$1", inputDB.userDB.Login)
				assert.NoError(t, err)

			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			if tC.pre != nil {
				tC.pre(t, tC.inputDB)
			}
			if tC.cleanUp != nil {
				defer tC.cleanUp(t, tC.inputDB)
			}

			var advertResponses []advertResponse
			resp := doRequest(t, client, "GET", fmt.Sprintf("%s/api/adverts", server.URL), nil, &advertResponses, []*http.Cookie{})
			assert.Equal(t, tC.expected.status, resp.StatusCode)
			assert.Equal(t, tC.expected.advertResponsesLength, len(advertResponses))
		})
	}
	// TODO: More test cases for testing limit, offset,
}
