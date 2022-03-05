package user

import (
	"github.com/stretchr/testify/assert"
	"github.com/ukrainian-brothers/board-backend/domain"
	"testing"
)

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
			testName:       "Correct data",
			firstName:      "Adam",
			surname:        "Małysz",
			login:          "adam@wp.pl",
			password:       "abc",
			contactDetails: domain.NewContactDetails("refugees_help@wp.pl", "111222333"),
			expectedErr:    nil,
		},
		{
			testName:    "Correct data",
			firstName:   "Adam",
			surname:     "Małysz",
			login:       "a@wp.pl",
			password:    "abc",
			expectedErr: MissingContactDataErr,
		},
		{
			testName:    "Correct data",
			firstName:   "",
			surname:     "Małysz",
			login:       "a@wp.pl",
			password:    "abc",
			expectedErr: MissingPersonalDataErr,
		},
	}

	for _, test := range testCases {
		t.Run(test.testName, func(t *testing.T) {
			_, err := NewUser(test.firstName, test.surname, test.login, test.password, test.contactDetails)
			assert.Equal(t, test.expectedErr, err)
		})

	}
}
