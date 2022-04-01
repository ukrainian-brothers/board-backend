package api

import (
	"bytes"
	"encoding/json"
	"github.com/go-gorp/gorp"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	application "github.com/ukrainian-brothers/board-backend/app"
	"github.com/ukrainian-brothers/board-backend/app/board"
	"github.com/ukrainian-brothers/board-backend/domain/user"
	"github.com/ukrainian-brothers/board-backend/internal/common"
	internal_user "github.com/ukrainian-brothers/board-backend/internal/user"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func newStringPtr(s string) *string {
	return &s
}

func getMockedRepo() internal_user.RepositoryMock {
	return internal_user.RepositoryMock{}
}

func getPostgresRepo() (user.Repository, *gorp.DbMap) {
	cfg, err := common.NewConfigFromFile("../config/configuration.test.local.json")
	if err != nil {
		log.WithError(err).Fatal("failed initializing config")
	}

	db, err := common.InitPostgres(&cfg.Postgres)
	if err != nil {
		log.WithError(err).Fatal("failed initializing postgres")
	}

	return internal_user.NewPostgresUserRepository(db), db
}

func createUserAPI(userRepo user.Repository) (*httptest.Server, http.Client) {
	logger := log.NewEntry(log.New())

	app := application.Application{
		Commands: application.Commands{
			AddUser: board.NewAddUser(userRepo),
		},
		Queries: application.Queries{
			UserExists:         board.NewUserExists(userRepo),
			VerifyUserPassword: board.NewVerifyUserPassword(userRepo),
		},
	}

	router := mux.NewRouter()
	router.Use(BodyLimitMiddleware)
	router.Use(LoggingMiddleware(logger))
	NewUserAPI(router, logger, app)

	server := httptest.NewServer(router)

	return server, http.Client{
		Timeout: time.Second * 3,
	}
}

func doRequest(t *testing.T, client http.Client, method string, url string, payload interface{}, response interface{}) *http.Response {
	by, err := json.Marshal(payload)
	assert.NoError(t, err)

	req, err := http.NewRequest(method, url, bytes.NewReader(by))
	assert.NoError(t, err)

	resp, err := client.Do(req)
	assert.NoError(t, err)

	if response != nil {
		by, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)

		err = json.Unmarshal(by, &response)
		assert.NoError(t, err)
	}

	return resp
}
