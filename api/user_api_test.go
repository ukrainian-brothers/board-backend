package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-gorp/gorp"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	application "github.com/ukrainian-brothers/board-backend/app"
	"github.com/ukrainian-brothers/board-backend/app/board"
	"github.com/ukrainian-brothers/board-backend/internal/common"
	"github.com/ukrainian-brothers/board-backend/internal/user"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// var host = "127.0.0.1:8001"
//func createUserAPI() http.Client {
//	logger := log.NewEntry(log.New())
//
//	cfg, err := common.NewConfigFromFile("../config/configuration.test.local.json")
//	if err != nil {
//		log.WithError(err).Fatal("failed initializing config")
//	}
//
//	db, err := common.InitPostgres(&cfg.Postgres)
//	if err != nil {
//		log.WithError(err).Fatal("failed initializing postgres")
//	}
//
//	userRepo := user.NewPostgresUserRepository(db)
//
//	app := application.Application{
//		Commands: application.Commands{
//			AddUser: board.NewAddUser(userRepo),
//		},
//		Queries:  application.Queries{
//			UserExists: board.NewUserExists(userRepo),
//		},
//	}
//
//	router := mux.NewRouter()
//	router.Use(BodyLimitMiddleware)
//	router.Use(LoggingMiddleware(logger))
//	NewUserAPI(router, logger, app)
//
//	srv := &http.Server{
//		Handler: router,
//		Addr:    host,
//		WriteTimeout: 4 * time.Second,
//		ReadTimeout:  4 * time.Second,
//	}
//
//	go log.Fatal(srv.ListenAndServe())
//	return http.Client{
//		Timeout: time.Second*3,
//	}
//}

func newStringPtr(s string) *string {
	return &s
}

func createUserAPI() (*httptest.Server, http.Client, *gorp.DbMap) {
	logger := log.NewEntry(log.New())

	cfg, err := common.NewConfigFromFile("../config/configuration.test.local.json")
	if err != nil {
		log.WithError(err).Fatal("failed initializing config")
	}

	db, err := common.InitPostgres(&cfg.Postgres)
	if err != nil {
		log.WithError(err).Fatal("failed initializing postgres")
	}

	userRepo := user.NewPostgresUserRepository(db)

	app := application.Application{
		Commands: application.Commands{
			AddUser: board.NewAddUser(userRepo),
		},
		Queries: application.Queries{
			UserExists: board.NewUserExists(userRepo),
		},
	}

	router := mux.NewRouter()
	router.Use(BodyLimitMiddleware)
	router.Use(LoggingMiddleware(logger))
	NewUserAPI(router, logger, app)

	server := httptest.NewServer(router)

	return server, http.Client{
		Timeout: time.Second * 3,
	}, db
}

func TestRegistration(t *testing.T) {
	server, client, db := createUserAPI()

	type expected struct {
		status int
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
					Error: "Unprocessable Entity",
					Details: "missing contact details",
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
				usrDB := user.UserDB{
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
			by, err := json.Marshal(tC.payload)
			assert.NoError(t, err)

			req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/user/register", server.URL), bytes.NewReader(by))
			assert.NoError(t, err)

			resp, err := client.Do(req)
			assert.NoError(t, err)
			assert.Equal(t, tC.expected.status, resp.StatusCode)

			by, err = io.ReadAll(resp.Body)
			assert.NoError(t, err)

			responseStruct := errorStruct{}
			err = json.Unmarshal(by, &responseStruct)
			assert.NoError(t, err)

			assert.Equal(t, tC.expected.errorStruct.Error, responseStruct.Error)
			assert.Equal(t, tC.expected.errorStruct.Details, responseStruct.Details)


			tC.cleanUp(t, tC.payload)
		})
	}
}
