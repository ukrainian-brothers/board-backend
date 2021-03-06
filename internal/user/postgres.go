package user

import (
	"context"
	"fmt"
	"github.com/go-gorp/gorp"
	"github.com/google/uuid"
	"github.com/ukrainian-brothers/board-backend/domain"
	"github.com/ukrainian-brothers/board-backend/domain/user"
)

type PostgresUserRepository struct {
	db *gorp.DbMap
}

type UserDB struct {
	ID          uuid.UUID `db:"id"`
	Login       string    `db:"login"`
	Password    *string   `db:"password"`
	FirstName   string    `db:"name"`
	Surname     string    `db:"surname"`
	Mail        *string   `db:"mail"`
	PhoneNumber *string   `db:"phone_number"`
}

func (usrDB *UserDB) LoadUser(usr *user.User) {
	usrDB.ID = usr.ID
	usrDB.Login = usr.Login
	usrDB.Password = usr.Password
	usrDB.FirstName = usr.Person.FirstName
	usrDB.Surname = usr.Person.Surname
	usrDB.Mail = usr.ContactDetails.Mail
	usrDB.PhoneNumber = usr.ContactDetails.PhoneNumber
}

func NewPostgresUserRepository(db *gorp.DbMap) *PostgresUserRepository {
	db.AddTableWithName(UserDB{}, "users").SetKeys(false, "id")
	return &PostgresUserRepository{db: db}
}

func (repo PostgresUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*user.User, error) {
	sqlExecutor := repo.db.WithContext(ctx)

	var usr UserDB
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

func (repo PostgresUserRepository) GetByLogin(ctx context.Context, login string) (*user.User, error) {
	sqlExecutor := repo.db.WithContext(ctx)

	var usr UserDB
	err := sqlExecutor.SelectOne(&usr, "select * from users where login=$1", login)
	if err != nil {
		return nil, fmt.Errorf("GetByLogin failed while selecting user %w", err)
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

func (repo PostgresUserRepository) Add(ctx context.Context, user *user.User) error {
	userDB := UserDB{
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
		return fmt.Errorf("adding user failed %w", err)
	}
	return nil
}

func (repo PostgresUserRepository) Exists(ctx context.Context, login string) (bool, error) {
	sqlExecutor := repo.db.WithContext(ctx)
	exists, err := sqlExecutor.SelectStr(`select exists(select 1 from users where login=$1)`, login)
	if err != nil {
		return false, fmt.Errorf("checking if user exists failed %w", err)
	}
	return exists == "true", nil
}

func (repo PostgresUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	sqlExecutor := repo.db.WithContext(ctx)
	_, err := sqlExecutor.Exec(`DELETE FROM users WHERE id=$1`, id)
	if err != nil {
		return fmt.Errorf("deleting user failed %w", err)
	}
	return nil
}
