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
	db.AddTableWithName(advertDB{}, "adverts").SetKeys(false, "id")
	db.AddTableWithName(advertDetailsDB{}, "adverts_details").SetKeys(false, "id")

	return &PostgresAdvertRepository{
		db: db,
	}
}

func newStringPtr(s string) *string {
	return &s
}

type advertDB struct {
	ID             uuid.UUID             `db:"id"`
	UserID         uuid.UUID             `db:"user_id"`
	Type           domain.AdvertType     `db:"type"`
	Views          int                   `db:"views"`
	ContactDetails domain.ContactDetails `db:"contact_details,json"`
	CreatedAt      time.Time             `db:"created_at"`
	UpdatedAt      *time.Time            `db:"updated_at"`
	DestroyedAt    *time.Time            `db:"destroyed_at"`
}

type advertDetailsDB struct {
	ID          uuid.UUID   `db:"id"`
	AdvertID    uuid.UUID   `db:"advert_id"`
	Language    LanguageTag `db:"language"` // ISO 639-1
	Title       string      `db:"title"`
	Description string      `db:"description"`
}

func (repo PostgresAdvertRepository) Get(ctx context.Context, id uuid.UUID) (advert.Advert, error) {
	sqlExec := repo.db.WithContext(ctx)

	type advertAndUser struct {
		advertDB
		ID          uuid.UUID `db:"id"`
		Login       string    `db:"login"`
		Password    *string   `db:"password"`
		FirstName   string    `db:"name"`
		Surname     string    `db:"surname"`
		Mail        *string   `db:"mail"`
		PhoneNumber *string   `db:"phone_number"`
	}

	adv := advertAndUser{}

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

	advertDetailsList, err := sqlExec.Select(advertDetailsDB{}, "SELECT * FROM adverts_details WHERE advert_id=$1", adv.ID)
	if err != nil {
		return advert.Advert{}, fmt.Errorf("failed getting advert translations: %w", err)
	}

	multilingualTitle := make(MultilingualString)
	multilingualDescription := make(MultilingualString)
	for _, val := range advertDetailsList {
		advertDetails := val.(advertDetailsDB)
		multilingualTitle[advertDetails.Language] = advertDetails.Title
		multilingualDescription[advertDetails.Language] = advertDetails.Description
	}

	multilingualTitle.RemoveUnsupported()
	multilingualDescription.RemoveUnsupported()

	return advert.Advert{
		ID: adv.ID,
		Details: domain.AdvertDetails{
			Title:          multilingualTitle,
			Description:    multilingualDescription,
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

	advertDb := advertDB{
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

		advertDetailsDB := advertDetailsDB{
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

func (repo PostgresAdvertRepository) Delete(ctx context.Context, id uuid.UUID) error {
	panic("implement me")
}
