package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	application "github.com/ukrainian-brothers/board-backend/app"
	"github.com/ukrainian-brothers/board-backend/domain"
	"github.com/ukrainian-brothers/board-backend/domain/user"
	"net/http"
)

type UserAPI struct {
	log *log.Entry
	r   *mux.Router
	app application.Application
}

func NewUserAPI(r *mux.Router, log *log.Entry, app application.Application) *UserAPI {
	usrApi := UserAPI{r: r, app: app, log: log}
	r.HandleFunc("/api/user/register", usrApi.register).Methods("POST")
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

func (u UserAPI) register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	payload := registerPayload{}
	err := dec.Decode(&payload)

	contactDetails, err := domain.NewContactDetails(payload.Mail, payload.Phone)
	if err != nil {
		u.log.WithError(err).Error("failed creating contact details")
		WriteError(w, InvalidPayload)
		return
	}

	usr, err := user.NewUser(payload.Firstname, payload.Surname, payload.Login, payload.Password, contactDetails)
	if err != nil {
		u.log.WithError(err).Error("failed creating User struct")
		WriteError(w, InvalidPayload)
		return
	}

	err = u.app.Commands.AddUser.Execute(ctx, usr)
	if err != nil {
		u.log.WithError(err).Error("failed to execute AddUser command")
		WriteError(w, InternalError)
		return
	}

	//WriteJSON(w, 200, struct {
	//	Success bool
	//}{Success: true}
	WriteJSON(w, 200, map[string]string{"status": "ok"})
	return
}
func (u UserAPI) login(w http.ResponseWriter, r *http.Request) {}
