package advert

import (
	"context"
	"fmt"
	"github.com/go-gorp/gorp"
	"github.com/google/uuid"
	"github.com/ukrainian-brothers/board-backend/domain"
	"github.com/ukrainian-brothers/board-backend/domain/advert"
	"github.com/ukrainian-brothers/board-backend/domain/user"
	. "github.com/ukrainian-brothers/board-backend/pkg/translation"
	"time"
)

type PostgresAdvertRepository struct {
	db *gorp.DbMap
}

func NewPostgresAdvertRepository(db *gorp.DbMap) *PostgresAdvertRepository {
	db.AddTableWithName(AdvertDB{}, "adverts").SetKeys(false, "id")
	db.AddTableWithName(AdvertDetailsDB{}, "adverts_details").SetKeys(false, "id")

	return &PostgresAdvertRepository{
		db: db,
	}
}

func newStringPtr(s string) *string {
	return &s
}

type AdvertDB struct {
	ID             uuid.UUID             `db:"id"`
	UserID         uuid.UUID             `db:"user_id"`
	Type           domain.AdvertType     `db:"type"`
	Views          int                   `db:"views"`
	ContactDetails domain.ContactDetails `db:"contact_details,json"`
	CreatedAt      time.Time             `db:"created_at"`
	UpdatedAt      *time.Time            `db:"updated_at"`
	DestroyedAt    *time.Time            `db:"destroyed_at"`
}

type AdvertDetailsDB struct {
	ID          uuid.UUID   `db:"id"`
	AdvertID    uuid.UUID   `db:"advert_id"`
	Language    LanguageTag `db:"language"` // ISO 639-1
	Title       string      `db:"title"`
	Description string      `db:"description"`
}

type advertTranslations struct {
	Title       MultilingualString
	Description MultilingualString
}

func (tr *advertTranslations) Filter(langs []LanguageTag) {
	newTitle := make(MultilingualString)
	newDesc := make(MultilingualString)
	for _, expectedLang := range langs {
		t, ok := tr.Title[expectedLang]
		if ok {
			newTitle[expectedLang] = t
		}

		d, ok := tr.Description[expectedLang]
		if ok {
			newDesc[expectedLang] = d
		}
	}
	tr.Title = newTitle
	tr.Description = newDesc
}

func (repo PostgresAdvertRepository) getAdvertTranslations(ctx context.Context, advertID uuid.UUID) (advertTranslations, error) {
	sqlExec := repo.db.WithContext(ctx)
	var advDetailsDB []AdvertDetailsDB
	_, err := sqlExec.Select(&advDetailsDB, "SELECT * FROM adverts_details WHERE advert_id=$1", advertID.String())
	if err != nil {
		return advertTranslations{}, fmt.Errorf("failed getting advert translations: %w", err)
	}

	translation := advertTranslations{
		Title:       make(MultilingualString),
		Description: make(MultilingualString),
	}

	for _, val := range advDetailsDB {
		translation.Title[val.Language] = val.Title
		translation.Description[val.Language] = val.Description
	}

	translation.Title.RemoveUnsupported()
	translation.Description.RemoveUnsupported()
	return translation, nil
}

func (repo PostgresAdvertRepository) Get(ctx context.Context, id uuid.UUID) (advert.Advert, error) {
	sqlExec := repo.db.WithContext(ctx)

	type advertAndUserDB struct {
		AdvertDB
		ID          uuid.UUID `db:"id"`
		Login       string    `db:"login"`
		Password    *string   `db:"password"`
		FirstName   string    `db:"name"`
		Surname     string    `db:"surname"`
		Mail        *string   `db:"mail"`
		PhoneNumber *string   `db:"phone_number"`
	}

	adv := advertAndUserDB{}

	err := sqlExec.SelectOne(&adv, "SELECT * FROM adverts JOIN users ON (adverts.user_id = users.id) WHERE adverts.id=$1;", id.String())
	if err != nil {
		return advert.Advert{}, fmt.Errorf("getting advert failed while selecting from db %w", err)
	}

	usr, err := user.NewUser(adv.FirstName, adv.Surname, adv.Login, *adv.Password, domain.ContactDetails{
		Mail:        adv.Mail,
		PhoneNumber: adv.PhoneNumber,
	})
	if err != nil {
		return advert.Advert{}, fmt.Errorf("getting advert failed while performing NewUser() %w", err)
	}

	translation, err := repo.getAdvertTranslations(ctx, adv.ID)
	if err != nil {
		return advert.Advert{}, err
	}

	return advert.Advert{
		ID: adv.ID,
		Details: domain.AdvertDetails{
			Title:          translation.Title,
			Description:    translation.Description,
			Type:           adv.Type,
			Views:          adv.Views,
			ContactDetails: adv.ContactDetails,
		},
		User:        usr,
		CreatedAt:   adv.CreatedAt,
		UpdatedAt:   adv.UpdatedAt,
		DestroyedAt: adv.DestroyedAt,
	}, nil
}

func (repo PostgresAdvertRepository) Add(ctx context.Context, advert *advert.Advert) error {
	trans, err := repo.db.Begin()
	if err != nil {
		return fmt.Errorf("failed creating transaction for adding advert to repo: %w", err)
	}
	sqlExecutor := trans.WithContext(ctx)

	advertDb := AdvertDB{
		ID:             advert.ID,
		UserID:         advert.User.ID,
		Type:           advert.Details.Type,
		ContactDetails: advert.Details.ContactDetails,
		CreatedAt:      advert.CreatedAt,
		DestroyedAt:    advert.DestroyedAt,
		UpdatedAt:      advert.UpdatedAt,
	}

	err = sqlExecutor.Insert(&advertDb)
	if err != nil {
		return fmt.Errorf("adding advert failed while performing sql %w", err)
	}

	for lang, title := range advert.Details.Title {
		description, ok := advert.Details.Description[lang]
		if !ok {
			continue
		}

		advertDetailsDB := AdvertDetailsDB{
			ID:          uuid.New(),
			AdvertID:    advert.ID,
			Language:    lang,
			Title:       title,
			Description: description,
		}
		err = sqlExecutor.Insert(&advertDetailsDB)
		if err != nil {
			return fmt.Errorf("failed inserting advertDetails: %w", err)
		}
	}

	return trans.Commit()
}

func (repo PostgresAdvertRepository) GetList(ctx context.Context, langs LanguageTags, limit int, offset int) ([]*advert.Advert, error) {
	sqlExec := repo.db.WithContext(ctx)

	var advertsDB []AdvertDB
	_, err := sqlExec.Select(&advertsDB, "SELECT * FROM adverts LIMIT $1 OFFSET $2", limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed selecting many adverts with translations: %w", err)
	}

	var adverts []*advert.Advert
	for _, advDB := range advertsDB {
		translation, err := repo.getAdvertTranslations(ctx, advDB.ID)
		if err != nil {
			return nil, err
		}

		// don't filter if there are no langs selected
		if !langs.Empty() {
			translation.Filter(langs)
		}

		if translation.Title.Empty() || translation.Description.Empty() {
			continue
		}

		adverts = append(adverts, &advert.Advert{
			ID: advDB.ID,
			Details: domain.AdvertDetails{
				Title:       translation.Title,
				Description: translation.Description,
				Type:        advDB.Type,
				Views:       advDB.Views,
				ContactDetails: domain.ContactDetails{
					Mail:        advDB.ContactDetails.Mail,
					PhoneNumber: advDB.ContactDetails.PhoneNumber,
				},
			},
			User:        &user.User{ID: advDB.UserID},
			CreatedAt:   advDB.CreatedAt,
			UpdatedAt:   advDB.UpdatedAt,
			DestroyedAt: advDB.DestroyedAt,
		})
	}

	return adverts, nil
}

func (repo PostgresAdvertRepository) Delete(ctx context.Context, id uuid.UUID) error {
	panic("implement me")
}
