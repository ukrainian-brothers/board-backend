package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-gorp/gorp"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	application "github.com/ukrainian-brothers/board-backend/app"
	"github.com/ukrainian-brothers/board-backend/app/board"
	"github.com/ukrainian-brothers/board-backend/domain"
	"github.com/ukrainian-brothers/board-backend/domain/advert"
	"github.com/ukrainian-brothers/board-backend/domain/user"
	internal_advert "github.com/ukrainian-brothers/board-backend/internal/advert"
	"github.com/ukrainian-brothers/board-backend/internal/common"
	internal_user "github.com/ukrainian-brothers/board-backend/internal/user"
	"github.com/ukrainian-brothers/board-backend/pkg/translation"
	"io"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func newStringPtr(s string) *string {
	return &s
}

func getMockedRepo() (internal_user.RepositoryMock, internal_advert.RepositoryMock) {
	return internal_user.RepositoryMock{}, internal_advert.RepositoryMock{}
}

func getPostgresRepos(t *testing.T) (user.Repository, advert.Repository, *gorp.DbMap) {
	cfg := getTestConfig(t)

	db, err := common.InitPostgres(&cfg.Postgres)
	if err != nil {
		log.WithError(err).Fatal("failed initializing postgres")
	}

	return internal_user.NewPostgresUserRepository(db), internal_advert.NewPostgresAdvertRepository(db), db
}

func getValidContactDetails() domain.ContactDetails {
	return domain.ContactDetails{
		Mail:        newStringPtr("macncheese@wp.pl"),
		PhoneNumber: newStringPtr("+48 111 222 333"),
	}
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

	cfg := getTestConfig(t)

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

func createTestUser(t *testing.T, login string, userRepo user.Repository) *user.User {
	usr := &user.User{
		ID:       uuid.New(),
		Login:    login,
		Password: newStringPtr("2Fzs0V!@~4m;'13.!#"),
		Person: domain.Person{
			FirstName: "Mac",
			Surname:   "Cheese",
		},
		ContactDetails: getValidContactDetails(),
	}
	err := userRepo.Add(context.Background(), usr)
	assert.NoError(t, err)
	return usr
}

func removeTestUser(t *testing.T, id uuid.UUID, userRepo user.Repository) {
	err := userRepo.Delete(context.Background(), id)
	assert.NoError(t, err)
}

func createTestSession(t *testing.T, usr *user.User, store sessions.Store) []*http.Cookie {
	cfg := getTestConfig(t)
	r := &http.Request{}
	session, err := store.Get(r, cfg.Session.SessionKey)
	assert.NoError(t, err)

	session.Values["user_login"] = usr.Login
	w := httptest.NewRecorder()
	err = session.Save(r, w)
	assert.NoError(t, err)

	return w.Result().Cookies()
}

func getTestConfig(t *testing.T) *common.Config {
	cfg, err := common.NewConfigFromFile("../config/configuration.test.local.json")
	assert.NoError(t, err)
	return cfg
}

func randomString(length int) string {
	by := make([]byte, length)
	for i := 0; i < length; i++ {
		by[i] = byte(65 + rand.Intn(25))
	}
	return string(by)
}

func randomNumberRange(min, max int) int {
	return rand.Intn(max-min) + min
}

func randomMail() *string {
	return newStringPtr(fmt.Sprintf("%s@wp.pl", randomString(8)))
}

func generateUserDB(id uuid.UUID) internal_user.UserDB {
	usr := internal_user.UserDB{
		ID:        id,
		Login:     randomString(10),
		FirstName: randomString(10),
		Surname:   randomString(10),
		Mail:      randomMail(),
	}
	return usr
}

func generateAdvertDB(advertID uuid.UUID, userID uuid.UUID) internal_advert.AdvertDB {
	return internal_advert.AdvertDB{
		ID:     advertID,
		UserID: userID,
		Type:   domain.AdvertTypeTransport,
		Views:  randomNumberRange(0, 150),
		ContactDetails: domain.ContactDetails{
			Mail: randomMail(),
		},
		CreatedAt: time.Now(),
	}
}
func generateAdvertDetailsDB(advertID uuid.UUID, language translation.LanguageTag) internal_advert.AdvertDetailsDB {
	return internal_advert.AdvertDetailsDB{
		ID:          uuid.New(),
		AdvertID:    advertID,
		Language:    language,
		Title:       randomString(10),
		Description: randomString(50),
	}
}
