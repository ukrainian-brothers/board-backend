package user_test

import (
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/ukrainian-brothers/board-backend/domain"
	"github.com/ukrainian-brothers/board-backend/domain/user"
	"github.com/ukrainian-brothers/board-backend/internal/common"
	internalUser "github.com/ukrainian-brothers/board-backend/internal/user"
	"testing"
)

func createString(s string) *string {
	return &s
}

func TestUserPostgresAdd(t *testing.T) {
	cfg, err := common.NewConfigFromFile("../../config/integration_test_config.json") // TODO: Describe how to achive this config in README.MD
	assert.NoError(t, err)
	db := common.InitPostgres(&cfg.Postgres)
	repo := internalUser.NewPostgresUserRepository(db)

	type testCase struct {
		name        string
		user        *user.User
		cleanUp     func(t *testing.T, id string)
		expectedErr error
	}

	testCases := []testCase {
		{
			name: "Success",
			user: &user.User {
				ID:       uuid.MustParse("6858fe22-2c04-4a13-bc75-eafeeb3cf767"),
				Login:    "Adam",
				Password: createString("awddwaawd"),
				Person: domain.Person {
					FirstName: "Adam",
					Surname:   "Małysz",
				},
				ContactDetails: domain.ContactDetails{
					Mail:        createString("adam@wp.pl"),
					PhoneNumber: createString("111222333"),
				},
			},
			cleanUp: func(t *testing.T, id string) {
				_, err := db.Exec("DELETE FROM users WHERE id=$1", id)
				assert.NoError(t, err)
			},
			expectedErr: nil,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			err := repo.Add(context.Background(), test.user)
			assert.Equal(t, test.expectedErr, err)
			test.cleanUp(t, test.user.ID.String())
		})
	}
}
