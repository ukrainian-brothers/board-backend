package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	logrus "github.com/sirupsen/logrus"
	application "github.com/ukrainian-brothers/board-backend/app"
	"github.com/ukrainian-brothers/board-backend/domain"
	"github.com/ukrainian-brothers/board-backend/domain/user"
	"github.com/ukrainian-brothers/board-backend/internal/common"
	"net/http"
)

type UserAPI struct {
	log          *logrus.Entry
	router       *mux.Router
	app          application.Application
	sessionStore *sessions.CookieStore
	cfg          *common.Config
}

func NewUserAPI(r *mux.Router, log *logrus.Entry, app application.Application, middleware *MiddlewareProvider, sessionStore *sessions.CookieStore, cfg *common.Config) *UserAPI {
	usrApi := UserAPI{router: r, app: app, log: log, sessionStore: sessionStore, cfg: cfg}
	r.HandleFunc("/api/user/register", usrApi.Register).Methods("POST")
	r.HandleFunc("/api/user/login", usrApi.Login).Methods("POST")
	return &usrApi
}

type registerPayload struct {
	Login     string `json:"Login"`
	Password  string `json:"password"`
	Firstname string `json:"firstname"`
	Surname   string `json:"surname"`
	Mail      string `json:"mail"`
	Phone     string `json:"phone"`
}

func (u UserAPI) Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := u.log

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	payload := registerPayload{}
	err := dec.Decode(&payload)
	if err != nil {
		log.WithError(err).Error("failed decoding register payload")
		WriteError(w, http.StatusUnprocessableEntity, "invalid payload")
		return
	}

	contactDetails, err := domain.NewContactDetails(payload.Mail, payload.Phone)
	if err != nil {
		log.WithError(err).Error("failed creating contact details")
		WriteError(w, http.StatusUnprocessableEntity, "missing contact details")
		return
	}

	log = log.WithFields(logrus.Fields{
		"Login": payload.Login,
		"mail":  payload.Mail,
	})

	usr, err := user.NewUser(payload.Firstname, payload.Surname, payload.Login, payload.Password, contactDetails)
	if err != nil {
		log.WithError(err).Error("failed creating User struct")
		WriteError(w, http.StatusUnprocessableEntity, "")
		return
	}

	userExists, err := u.app.Queries.UserExists.Execute(ctx, usr.Login)
	if err != nil {
		log.WithError(err).Error("failed to execute UserExists query")
		WriteError(w, http.StatusInternalServerError, "")
		return
	}

	if userExists {
		log.Info("user already exists")
		WriteError(w, http.StatusUnprocessableEntity, "user already exists")
		return
	}

	err = u.app.Commands.AddUser.Execute(ctx, *usr)
	if err != nil {
		log.WithError(err).Error("failed to execute AddUser command")
		WriteError(w, http.StatusInternalServerError, "")
		return
	}

	WriteJSON(w, 201, map[string]string{"status": "ok"})
}

type loginPayload struct {
	Login    string
	Password string
}

func (u UserAPI) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := u.log

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	payload := loginPayload{}
	err := dec.Decode(&payload)
	if err != nil {
		log.WithError(err).Error("failed decoding login payload")
		WriteError(w, http.StatusUnprocessableEntity, "invalid payload")
		return
	}

	exists, err := u.app.Queries.UserExists.Execute(ctx, payload.Login)
	if err != nil {
		log.WithError(err).Error("failed verifying user existence while logging in")
		WriteError(w, http.StatusInternalServerError, "")
		return
	}

	if !exists {
		log.Info("failed login, user does not exists")
		WriteError(w, http.StatusUnprocessableEntity, "user does not exists")
		return
	}

	valid, err := u.app.Queries.VerifyUserPassword.Execute(ctx, payload.Login, payload.Password)
	if err != nil {
		log.WithError(err).Error("failed verifying user password")
		WriteError(w, http.StatusInternalServerError, "")
		return
	}

	if !valid {
		log.Info("wrong credentials")
		WriteError(w, http.StatusForbidden, "wrong credentials")
		return
	}

	session, err := u.sessionStore.Get(r, u.cfg.Session.SessionKey)
	if err != nil {
		log.WithError(err).Error("Login failed getting session")
		WriteError(w, http.StatusInternalServerError, "")
		return
	}

	session.Values["user_login"] = payload.Login
	err = session.Save(r, w)
	if err != nil {
		log.WithError(err).Error("failed saving session")
	}
	WriteJSON(w, 200, map[string]string{"status": "ok"})
}
