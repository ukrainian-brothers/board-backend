package main

import (
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	log "github.com/sirupsen/logrus"
	"github.com/ukrainian-brothers/board-backend/api"
	application "github.com/ukrainian-brothers/board-backend/app"
	"github.com/ukrainian-brothers/board-backend/app/board"
	"github.com/ukrainian-brothers/board-backend/internal/advert"
	"github.com/ukrainian-brothers/board-backend/internal/common"
	"github.com/ukrainian-brothers/board-backend/internal/user"
	"net/http"
	"time"
)

func main() {
	logger := log.NewEntry(log.New())

	cfg, err := common.NewConfigFromFile("config/configuration.local.json")
	if err != nil {
		log.WithError(err).Fatal("failed initializing config")
	}

	db, err := common.InitPostgres(&cfg.Postgres)
	if err != nil {
		log.WithError(err).Fatal("failed initializing postgres")
	}

	userRepo := user.NewPostgresUserRepository(db)
	advertRepo := advert.NewPostgresAdvertRepository(db)

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

	sessionStore := sessions.NewCookieStore([]byte(cfg.Session.Secret))
	middleware := api.NewMiddlewareProvider(sessionStore, &app, cfg)

	router := mux.NewRouter()
	router.Use(middleware.BodyLimitMiddleware)
	router.Use(middleware.LoggingMiddleware(logger))
	api.NewUserAPI(router, logger, app, middleware, sessionStore, cfg)
	api.NewAdvertAPI(router, logger, app, middleware, sessionStore, cfg)

	srv := &http.Server{
		Handler:      router,
		Addr:         "127.0.0.1:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
