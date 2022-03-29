package main

import (
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/ukrainian-brothers/board-backend/api"
	application "github.com/ukrainian-brothers/board-backend/app"
	"github.com/ukrainian-brothers/board-backend/app/board"
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

	app := application.Application{
		Commands: application.Commands{
			AddUser: board.NewAddUser(userRepo),
		},
		Queries:  application.Queries{
			UserExists: board.NewUserExists(userRepo),
		},
	}

	router := mux.NewRouter()
	router.Use(api.BodyLimitMiddleware)
	router.Use(api.LoggingMiddleware(logger))
	api.NewUserAPI(router, logger, app)

	srv := &http.Server{
		Handler: router,
		Addr:    "127.0.0.1:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
