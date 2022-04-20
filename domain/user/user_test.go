package user

import (
	"github.com/stretchr/testify/assert"
	"github.com/ukrainian-brothers/board-backend/domain"
	"github.com/ukrainian-brothers/board-backend/pkg/test_helpers"
	"testing"
)

func newStringPtr(s string) *string {
	return &s
}

func TestUserCreation(t *testing.T) {
	type testData struct {
		testName       string
		firstName      string
		surname        string
		login          string
		password       string
		contactDetails domain.ContactDetails
		expectedErr    error
	}

	testCases := []testData{
		{
			testName:  "correct data",
			firstName: "Adam",
			surname:   "Małysz",
			login:     *test_helpers.RandomMail(),
			password:  "abc",
			contactDetails: domain.ContactDetails{
				Mail:        newStringPtr("mail"),
				PhoneNumber: newStringPtr("phone"),
			},
			expectedErr: nil,
		},
		{
			testName:    "missing contact data",
			firstName:   "Adam",
			surname:     "Małysz",
			login:       *test_helpers.RandomMail(),
			password:    "abc",
			expectedErr: MissingContactDataErr,
		},
		{
			testName:    "missing personal data",
			firstName:   "",
			surname:     "Małysz",
			login:       *test_helpers.RandomMail(),
			password:    "abc",
			expectedErr: MissingPersonalDataErr,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.testName, func(t *testing.T) {
			_, err := NewUser(tC.firstName, tC.surname, tC.login, tC.password, tC.contactDetails)
			assert.Equal(t, tC.expectedErr, err)
		})

	}
}
