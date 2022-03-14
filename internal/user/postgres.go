package user

import (
	"context"
	"fmt"
	"github.com/go-gorp/gorp"
	"github.com/google/uuid"
	"github.com/ukrainian-brothers/board-backend/domain"
	"github.com/ukrainian-brothers/board-backend/domain/user"
)

//@TODO: Wrapped errors in methods
type PostgresRepository struct {
	db *gorp.DbMap
}
type userDB struct {
	ID          uuid.UUID `db:"id"`
	Login       string    `db:"login"`
	Password    *string   `db:"password"`
	FirstName   string    `db:"name"`
	Surname     string    `db:"surname"`
	Mail        *string   `db:"mail"`
	PhoneNumber *string   `db:"phone_number"`
}

func NewPostgresUserRepository(db *gorp.DbMap) *PostgresRepository {
	db.AddTableWithName(userDB{}, "users").SetKeys(false, "id")
	return &PostgresRepository{db: db}
}

func (repo PostgresRepository) GetByID(ctx context.Context, id uuid.UUID) (*user.User, error) {
	sqlExecutor := repo.db.WithContext(ctx)

	var usr userDB
	err := sqlExecutor.SelectOne(&usr, `
	SELECT  login, id, password, name, surname, mail, phone_number FROM users
	WHERE id=$1`, id)
	if err != nil {
		return nil, fmt.Errorf("GetByID failed while selecting user %w", err)
	}

	return &user.User{
		ID:       id,
		Login:    usr.Login,
		Password: usr.Password,
		Person: domain.Person{
			FirstName: usr.FirstName,
			Surname:   usr.Surname,
		},
		ContactDetails: domain.ContactDetails{
			Mail:        usr.Mail,
			PhoneNumber: usr.PhoneNumber,
		},
	}, err
}

func (repo PostgresRepository) GetByLogin(ctx context.Context, login string) (*user.User, error) {
	sqlExecutor := repo.db.WithContext(ctx)

	var usr userDB
	err := sqlExecutor.SelectOne(&usr, "select * from users where login=$1", login)
	if err != nil {
		return nil, err
	}

	return &user.User{
		ID:       usr.ID,
		Login:    usr.Login,
		Password: usr.Password,
		Person: domain.Person{
			FirstName: usr.FirstName,
			Surname:   usr.Surname,
		},
		ContactDetails: domain.ContactDetails{
			Mail:        usr.Mail,
			PhoneNumber: usr.PhoneNumber,
		},
	}, nil
}

func (repo PostgresRepository) Add(ctx context.Context, user *user.User) error {

	userDB := userDB{
		ID:          user.ID,
		Login:       user.Login,
		Password:    user.Password,
		FirstName:   user.Person.FirstName,
		Surname:     user.Person.Surname,
		Mail:        user.ContactDetails.Mail,
		PhoneNumber: user.ContactDetails.PhoneNumber,
	}
	repo.db.WithContext(ctx)
	err := repo.db.Insert(&userDB)
	if err != nil {
		return err
	}
	return nil
}

func (repo PostgresRepository) Delete(ctx context.Context, id uuid.UUID) error {
	panic("implement me")
}
