package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
	application "github.com/ukrainian-brothers/board-backend/app"
	"github.com/ukrainian-brothers/board-backend/domain"
	"github.com/ukrainian-brothers/board-backend/domain/advert"
	"github.com/ukrainian-brothers/board-backend/internal/common"
	"net/http"
	"time"
)

type AdvertAPI struct {
	log          *logrus.Entry
	router       *mux.Router
	app          application.Application
	sessionStore *sessions.CookieStore
	cfg          *common.Config
}

func NewAdvertAPI(r *mux.Router, log *logrus.Entry, app application.Application, middleware *MiddlewareProvider, sessionStore *sessions.CookieStore, cfg *common.Config) *AdvertAPI {
	advertApi := AdvertAPI{router: r, app: app, log: log, sessionStore: sessionStore, cfg: cfg}
	r.HandleFunc("/api/adverts", middleware.AuthMiddleware(advertApi.AddAdvert, log)).Methods("POST")
	return &advertApi
}

type contactPayload struct {
	Mail        string
	PhoneNumber string
}
type newAdvertPayload struct {
	Title          string            `json:"title"`
	Description    string            `json:"description"`
	Type           domain.AdvertType `json:"type"`
	ContactDetails contactPayload    `json:"contact_details"`
}

func (a AdvertAPI) AddAdvert(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := a.log

	userLogin := ctx.Value("user_login")
	if userLogin == nil {
		log.Info("not authorized user tries to add advert")
		WriteError(w, http.StatusForbidden, "not authorized")
		return
	}

	log = log.WithField("user_login", userLogin.(string))

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	payload := newAdvertPayload{}
	err := dec.Decode(&payload)
	if err != nil {
		log.WithError(err).Error("failed decoding newAdvert payload")
		WriteError(w, http.StatusUnprocessableEntity, "invalid payload")
		return
	}

	usr, err := a.app.Queries.GetUserByLogin.Execute(ctx, userLogin.(string))
	if err != nil {
		if errors.Unwrap(err) == sql.ErrNoRows {
			log.Info("not authorized user tries to add advert")
			WriteError(w, http.StatusForbidden, "user does not exists anymore")
			return
		}
		log.WithError(err).Error("AddAdvert failed getting user by login")
		WriteError(w, http.StatusInternalServerError, "")
		return
	}

	advertContact, err := domain.NewContactDetails(payload.ContactDetails.Mail, payload.ContactDetails.PhoneNumber)
	if err != nil {
		log.WithError(err).Error("AddAdvert failed creating contact details")
		WriteError(w, http.StatusUnprocessableEntity, "invalid payload")
		return
	}

	adv, err := advert.NewAdvert(usr, payload.Title, payload.Description, payload.Type, advert.WithContactDetails(advertContact))
	if err != nil {
		log.WithError(err).Error("AddAdvert failed creating advert")
		WriteError(w, http.StatusUnprocessableEntity, "invalid advert details")
		return
	}

	err = a.app.Commands.AddAdvert.Execute(ctx, adv)
	if err != nil {
		log.WithError(err).Error("AddAdvert failed inserting advert")
		WriteError(w, http.StatusInternalServerError, "")
		return
	}

	response := advertResponse{}
	response.LoadAdvert(adv)
	WriteJSON(w, 201, response)
}

type advertResponse struct {
	ID             string            `json:"id"`
	Title          string            `json:"title"`
	Description    string            `json:"description"`
	Type           domain.AdvertType `json:"type"`
	ContactDetails contactPayload    `json:"contact_details"`
	CreatedAt      time.Time         `json:"created_at"`
	UpdatedAt      *time.Time        `json:"updated_at,omitempty"`
	DestroyedAt    *time.Time        `json:"destroyed_at,omitempty"`
}

func (a *advertResponse) LoadAdvert(adv *advert.Advert) {
	a.ID = adv.ID.String()
	a.Title = adv.Details.Title
	a.Description = adv.Details.Description
	a.Type = adv.Details.Type
	a.ContactDetails.Mail = *adv.Details.ContactDetails.Mail
	a.ContactDetails.PhoneNumber = *adv.Details.ContactDetails.PhoneNumber
	a.CreatedAt = adv.CreatedAt
	a.UpdatedAt = adv.UpdatedAt
	a.DestroyedAt = adv.DestroyedAt
}
