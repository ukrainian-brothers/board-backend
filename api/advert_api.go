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
	. "github.com/ukrainian-brothers/board-backend/pkg/translation"
	"net/http"
	"strconv"
	"strings"
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
	r.HandleFunc("/api/adverts", advertApi.AdvertsList).Methods("GET")
	return &advertApi
}

type contactPayload struct {
	Mail        string `json:"mail"`
	PhoneNumber string `json:"phone"`
}
type newAdvertPayload struct {
	Title          MultilingualString `json:"title"`
	Description    MultilingualString `json:"description"`
	Type           domain.AdvertType  `json:"type"`
	ContactDetails contactPayload     `json:"contact_details"`
}

func (p *newAdvertPayload) RemoveUnsupportedLanguages() {
	p.Title.RemoveUnsupported()
	p.Description.RemoveUnsupported()
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
	payload.RemoveUnsupportedLanguages()

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
	ID             string             `json:"id"`
	Title          MultilingualString `json:"title"`
	Description    MultilingualString `json:"description"`
	Type           domain.AdvertType  `json:"type"`
	ContactDetails contactPayload     `json:"contact_details"`
	CreatedAt      time.Time          `json:"created_at"`
	UpdatedAt      *time.Time         `json:"updated_at,omitempty"`
	DestroyedAt    *time.Time         `json:"destroyed_at,omitempty"`
}

func (a *advertResponse) LoadAdvert(adv *advert.Advert) {
	a.ID = adv.ID.String()
	a.Title = adv.Details.Title
	a.Description = adv.Details.Description
	a.Type = adv.Details.Type
	a.CreatedAt = adv.CreatedAt
	a.UpdatedAt = adv.UpdatedAt
	a.DestroyedAt = adv.DestroyedAt
	if adv.Details.ContactDetails.Mail != nil {
		a.ContactDetails.Mail = *adv.Details.ContactDetails.Mail
	}

	if adv.Details.ContactDetails.PhoneNumber != nil {
		a.ContactDetails.PhoneNumber = *adv.Details.ContactDetails.PhoneNumber
	}
}

const MaxAdvertsInResponse = 50

func (a AdvertAPI) AdvertsList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	limit, err := strconv.Atoi(r.FormValue("limit"))
	if err != nil {
		limit = MaxAdvertsInResponse
	}

	offset, err := strconv.Atoi(r.FormValue("offset"))
	if err != nil {
		offset = 0
	}

	// Will load languages from url param &langs=ua,pl,en into slice
	langs := LanguageTags{}.FromStrings(strings.Split(r.FormValue("langs"), ","))

	log := a.log.WithFields(logrus.Fields{
		"limit":  limit,
		"offset": offset,
	})

	adverts, err := a.app.Queries.GetAdvertsList.Execute(ctx, langs, limit, offset) // TODO: pass real langage tags
	if err != nil {
		log.WithError(err).Error("AdvertsList failed while fetching list of adverts")
		WriteError(w, http.StatusInternalServerError, "")
		return
	}

	if len(adverts) == 0 {
		WriteJSON(w, 200, []advertResponse{})
		return
	}

	var response []advertResponse
	for _, adv := range adverts {
		advResponse := advertResponse{}
		advResponse.LoadAdvert(adv)
		response = append(response, advResponse)
	}

	WriteJSON(w, 200, response)
}
