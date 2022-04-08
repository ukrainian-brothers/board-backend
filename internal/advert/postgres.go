package advert

import (
	"context"
	"fmt"
	"github.com/go-gorp/gorp"
	"github.com/google/uuid"
	"github.com/ukrainian-brothers/board-backend/domain"
	"github.com/ukrainian-brothers/board-backend/domain/advert"
	"github.com/ukrainian-brothers/board-backend/domain/user"
	"time"
)

type PostgresAdvertRepository struct {
	db *gorp.DbMap
}

func NewPostgresAdvertRepository(db *gorp.DbMap) *PostgresAdvertRepository {
	db.AddTableWithName(advertDB{}, "adverts").SetKeys(false, "id")

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
	Title          string                `db:"title"`
	Description    string                `db:"description"`
	Type           domain.AdvertType     `db:"type"`
	Views          int                   `db:"views"`
	ContactDetails domain.ContactDetails `db:"contact_details,json"`
	CreatedAt      time.Time             `db:"created_at"`
	UpdatedAt      *time.Time            `db:"updated_at"`
	DestroyedAt    *time.Time            `db:"destroyed_at"`
}

func (repo PostgresAdvertRepository) Get(ctx context.Context, id uuid.UUID) (advert.Advert, error) {
	repo.db.WithContext(ctx)

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

	err := repo.db.SelectOne(&adv, "SELECT * FROM adverts JOIN users ON (adverts.user_id = users.id) WHERE adverts.id=$1;", id.String())
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

	return advert.Advert{
		ID: adv.ID,
		Details: domain.AdvertDetails{
			Title:          adv.Title,
			Description:    adv.Description,
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
	advertDb := advertDB{
		ID:             advert.ID,
		UserID:         advert.User.ID,
		Title:          advert.Details.Title,
		Description:    advert.Details.Description,
		Type:           advert.Details.Type,
		ContactDetails: advert.Details.ContactDetails,
		CreatedAt:      advert.CreatedAt,
		DestroyedAt:    advert.DestroyedAt,
		UpdatedAt:      advert.UpdatedAt,
	}
	repo.db.WithContext(ctx)
	err := repo.db.Insert(&advertDb)
	if err != nil {
		return fmt.Errorf("adding advert failed while performing sql %w", err)
	}
	return nil
}

func (repo PostgresAdvertRepository) Delete(ctx context.Context, id uuid.UUID) error {
	panic("implement me")
}
