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
	assert.Equal(t, expected.Details, actual.Details)
	assert.Equal(t, expected.Details.Title, actual.Details.Title)
	assert.Equal(t, expected.Details.Description, actual.Details.Description)
	assert.Equal(t, expected.Details.Type, actual.Details.Type)
	assert.Equal(t, expected.Details.Views, actual.Details.Views)
	assert.Equal(t, expected.Details.ContactDetails, actual.Details.ContactDetails)
	assert.Equal(t, expected.Details.ContactDetails.Mail, actual.Details.ContactDetails.Mail)
	assert.Equal(t, expected.Details.ContactDetails.PhoneNumber, actual.Details.ContactDetails.PhoneNumber)
}