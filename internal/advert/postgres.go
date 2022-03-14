package advert

import (
	"context"
	"github.com/go-gorp/gorp"
	"github.com/google/uuid"
	"github.com/ukrainian-brothers/board-backend/domain"
	"github.com/ukrainian-brothers/board-backend/domain/advert"
	"github.com/ukrainian-brothers/board-backend/domain/user"
	"time"
)

type PostgresRepository struct {
	db *gorp.DbMap
}

func NewPostgresAdvertRepository(db *gorp.DbMap) *PostgresRepository {
	db.AddTableWithName(advertDB{}, "adverts").SetKeys(false, "id")

	return &PostgresRepository{
		db: db,
	}
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

func (repo PostgresRepository) Get(ctx context.Context, id uuid.UUID) (advert.Advert, error) {
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

	err := repo.db.SelectOne(&adv, "SELECT * FROM adverts INNER JOIN users ON (adverts.user_id = users.id) WHERE adverts.id=$1;", id.String())
	if err != nil {
		return advert.Advert{}, err
	}

	usr, err := user.NewUser(adv.FirstName, adv.Surname, adv.Login, *adv.Password, domain.NewContactDetails(*adv.ContactDetails.Mail, *adv.ContactDetails.PhoneNumber)) // This is temporary solution, there HAVE TO BE JOIN in query which will get the user out of another table
	return advert.Advert{
		ID:          adv.ID,
		Advert:      domain.Advert{
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
	}, err
}

func (repo PostgresRepository) Add(ctx context.Context, advert *advert.Advert) error {
	advertDb := advertDB {
		ID:             advert.ID,
		UserID:         advert.User.ID,
		Title:          advert.Advert.Title,
		Description:    advert.Advert.Description,
		Type:           advert.Advert.Type,
		ContactDetails: advert.Advert.ContactDetails,
		CreatedAt:      time.Now(),
	}
	repo.db.WithContext(ctx)
	err := repo.db.Insert(&advertDb)
	if err != nil {
		return err
	}
	return nil
}

func (repo PostgresRepository) Delete(ctx context.Context, id uuid.UUID) error {
	panic("implement me")
}
