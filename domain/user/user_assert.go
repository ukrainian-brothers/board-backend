package user

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Assert(t *testing.T, expected *User,  actual *User) {
	assert.Equal(t, expected, actual)
	assert.Equal(t, expected.ID, actual.ID)
	assert.Equal(t, expected.Login, actual.Login)
	assert.Equal(t, expected.Password, actual.Password)
	assert.Equal(t, expected.Person.FirstName, actual.Person.FirstName)
	assert.Equal(t, expected.Person.Surname, actual.Person.Surname)
	assert.Equal(t, expected.ContactDetails.Mail, actual.ContactDetails.Mail)
	assert.Equal(t, expected.ContactDetails.PhoneNumber, actual.ContactDetails.PhoneNumber)
}