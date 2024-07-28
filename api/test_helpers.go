package api

import (
	"bytes"
	"encoding/json"
	"github.com/go-gorp/gorp"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	application "github.com/ukrainian-brothers/board-backend/app"
	"github.com/ukrainian-brothers/board-backend/app/board"
	"github.com/ukrainian-brothers/board-backend/domain/advert"
	"github.com/ukrainian-brothers/board-backend/domain/user"
	internal_advert "github.com/ukrainian-brothers/board-backend/internal/advert"
	"github.com/ukrainian-brothers/board-backend/internal/common"
	internal_user "github.com/ukrainian-brothers/board-backend/internal/user"
	"github.com/ukrainian-brothers/board-backend/pkg/test_helpers"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newStringPtr(s string) *string {
	return &s
}

func getMockedRepo() (internal_user.RepositoryMock, internal_advert.RepositoryMock) {
	return internal_user.RepositoryMock{}, internal_advert.RepositoryMock{}
}

func getPostgresRepos(t *testing.T) (user.Repository, advert.Repository, *gorp.DbMap) {
	cfg := test_helpers.GetTestConfig(t)

	db, err := common.InitPostgres(&cfg.Postgres)
	if err != nil {
		log.WithError(err).Fatal("failed initializing postgres")
	}

	return internal_user.NewPostgresUserRepository(db), internal_advert.NewPostgresAdvertRepository(db), db
}

func createTestAPIs(t *testing.T, advertRepo advert.Repository, userRepo user.Repository) (*httptest.Server, http.Client, *sessions.CookieStore) {
	logger := log.NewEntry(log.New())

	app := application.Application{
		Commands: application.Commands{
			AddUser:   board.NewAddUser(userRepo),
			AddAdvert: board.NewAddAdvert(advertRepo),
		},
		Queries: application.Queries{
			UserExists:         board.NewUserExists(userRepo),
			GetUserByLogin:     board.NewGetUserByLogin(userRepo),
			VerifyUserPassword: board.NewVerifyUserPassword(userRepo),
			GetAdvertsList:     board.NewGetAdvertsList(advertRepo),
		},
	}

	cfg := test_helpers.GetTestConfig(t)

	sessionStore := sessions.NewCookieStore([]byte(cfg.Session.Secret))
	middleware := NewMiddlewareProvider(sessionStore, &app, cfg)

	router := mux.NewRouter()
	router.Use(middleware.BodyLimitMiddleware)
	router.Use(middleware.LoggingMiddleware(logger))
	NewUserAPI(router, logger, app, middleware, sessionStore, cfg)
	NewAdvertAPI(router, logger, app, middleware, sessionStore, cfg)

	server := httptest.NewServer(router)

	return server, http.Client{}, sessionStore
}

func responseToStruct(t *testing.T, resp *http.Response, response interface{}) {
	by, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	err = json.Unmarshal(by, response)
	assert.NoError(t, err)
}

func doRequest(t *testing.T, client http.Client, method string, url string, payload interface{}, response interface{}, cookies []*http.Cookie) *http.Response {
	by, err := json.Marshal(payload)
	assert.NoError(t, err)

	req, err := http.NewRequest(method, url, bytes.NewReader(by))
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}
	assert.NoError(t, err)

	resp, err := client.Do(req)
	assert.NoError(t, err)

	if response != nil {
		responseToStruct(t, resp, &response)
	}

	return resp
}
