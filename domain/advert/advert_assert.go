package advert

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Assert(t *testing.T, expected *Advert, actual *Advert) {
	assert.Equal(t, expected, actual)
	assert.Equal(t, expected.ID, actual.ID)
	assert.Equal(t, expected.CreatedAt, actual.CreatedAt)
	assert.Equal(t, expected.UpdatedAt, actual.UpdatedAt)
	assert.Equal(t, expected.DestroyedAt, actual.DestroyedAt)
	assert.Equal(t, expected.Advert, actual.Advert)
	assert.Equal(t, expected.Advert.Title, actual.Advert.Title)
	assert.Equal(t, expected.Advert.Description, actual.Advert.Description)
	assert.Equal(t, expected.Advert.Type, actual.Advert.Type)
	assert.Equal(t, expected.Advert.Views, actual.Advert.Views)
	assert.Equal(t, expected.Advert.ContactDetails, actual.Advert.ContactDetails)
	assert.Equal(t, expected.Advert.ContactDetails.Mail, actual.Advert.ContactDetails.Mail)
	assert.Equal(t, expected.Advert.ContactDetails.PhoneNumber, actual.Advert.ContactDetails.PhoneNumber)
}