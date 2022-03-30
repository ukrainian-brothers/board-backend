package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	logrus "github.com/sirupsen/logrus"
	application "github.com/ukrainian-brothers/board-backend/app"
	"github.com/ukrainian-brothers/board-backend/domain"
	"github.com/ukrainian-brothers/board-backend/domain/user"
	"net/http"
)

type UserAPI struct {
	log *logrus.Entry
	r   *mux.Router
	app application.Application
}

func NewUserAPI(r *mux.Router, log *logrus.Entry, app application.Application) *UserAPI {
	usrApi := UserAPI{r: r, app: app, log: log}
	r.HandleFunc("/api/user/register", usrApi.Register).Methods("POST")
	r.HandleFunc("/api/user/login", usrApi.login).Methods("POST")
	return &usrApi
}

type registerPayload struct {
	Login     string `json:"login"`
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
		log.WithError(err).Error("failed decoding payload")
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
		"login": payload.Login,
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

	err = u.app.Commands.AddUser.Execute(ctx, usr)
	if err != nil {
		log.WithError(err).Error("failed to execute AddUser command")
		WriteError(w, http.StatusInternalServerError, "")
		return
	}

	WriteJSON(w, 201, map[string]string{"status": "ok"})
	return
}
func (u UserAPI) login(w http.ResponseWriter, r *http.Request) {}
