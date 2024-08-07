package advert

import (
	"github.com/stretchr/testify/assert"
	"github.com/ukrainian-brothers/board-backend/domain"
	"github.com/ukrainian-brothers/board-backend/domain/user"
	"github.com/ukrainian-brothers/board-backend/pkg/test_helpers"
	. "github.com/ukrainian-brothers/board-backend/pkg/translation"
	"testing"
)

func newStringPtr(s string) *string {
	return &s
}

func TestAdvertNew(t *testing.T) {
	type expectations struct {
		err            error
		contactDetails domain.ContactDetails
	}
	type testData struct {
		testName     string
		user         *user.User
		title        MultilingualString
		description  MultilingualString
		advertType   domain.AdvertType
		opts         []AdvertOption
		expectations expectations
	}

	usr, err := user.NewUser(
		"Adam",
		"Małysz",
		*test_helpers.RandomMail(),
		"abc",
		domain.ContactDetails{
			Mail:        newStringPtr("mail"),
			PhoneNumber: newStringPtr("phone"),
		},
	)
	assert.NoError(t, err)

	testCases := []testData{
		{
			testName:    "Correct data using user contact details",
			user:        usr,
			title:       MultilingualString{"en": "x"},
			description: MultilingualString{"en": "x"},
			advertType:  domain.AdvertTypeTransport,
			expectations: expectations{
				err: nil,
				contactDetails: domain.ContactDetails{
					Mail:        newStringPtr("mail"),
					PhoneNumber: newStringPtr("phone"),
				},
			},
		},
		{
			testName:    "Correct data with different contact details",
			user:        usr,
			title:       MultilingualString{"en": "x"},
			description: MultilingualString{"en": "x"},
			advertType:  domain.AdvertTypeTransport,
			opts: []AdvertOption{
				WithContactDetails(
					domain.ContactDetails{
						Mail:        newStringPtr("mail"),
						PhoneNumber: newStringPtr("phone"),
					})},
			expectations: expectations{
				err: nil,
				contactDetails: domain.ContactDetails{
					Mail:        newStringPtr("mail"),
					PhoneNumber: newStringPtr("phone"),
				},
			},
		},
		{
			testName:   "Correct data with different contact details",
			user:       usr,
			title:      MultilingualString{"en": "x"},
			advertType: domain.AdvertTypeTransport,
			opts:       []AdvertOption{WithContactDetails(domain.ContactDetails{})},
			expectations: expectations{
				err: ContactEmptyErr,
			},
		},
		{
			testName: "Missing basic info",
			user:     usr,
			expectations: expectations{
				err: MissingBasicInfoErr,
			},
		},
		{
			testName: "No user provided",
			user:     nil,
			expectations: expectations{
				err: NoUserProvidedErr,
			},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.testName, func(t *testing.T) {
			adv, err := NewAdvert(tC.user, tC.title, tC.description, tC.advertType, tC.opts...)
			assert.Equal(t, tC.expectations.err, err)
			if err == nil {
				assert.Equal(t, tC.expectations.contactDetails, adv.Details.ContactDetails)
			}
		})
	}
}
