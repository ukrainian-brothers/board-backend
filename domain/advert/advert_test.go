package advert

import (
	"github.com/stretchr/testify/assert"
	"github.com/ukrainian-brothers/board-backend/domain"
	"github.com/ukrainian-brothers/board-backend/domain/user"
	"testing"
)

func TestAdvertNew(t *testing.T) {
	type expectations struct {
		err            error
		contactDetails domain.ContactDetails
	}
	type testData struct {
		testName     string
		user         *user.User
		title        string
		description  string
		advertType   domain.AdvertType
		opts         []AdvertOption
		expectations expectations
	}

	usr, err := user.NewUser(
		"Adam",
		"Małysz",
		"adam@wp.pl",
		"abc",
		domain.NewContactDetails("adam@wp.pl", "111222333"),
	)
	assert.NoError(t, err)

	testCases := []testData{
		{
			testName:    "Correct data using user contact details",
			user:        usr,
			title:       "Relocation for refugees",
			description: "desc",
			advertType:  domain.AdvertTypeTransport,
			expectations: expectations{
				err:            nil,
				contactDetails: domain.NewContactDetails("adam@wp.pl", "111222333"),
			},
		},
		{
			testName:    "Correct data with different contact details",
			user:        usr,
			title:       "Relocation for refugees",
			description: "desc",
			advertType:  domain.AdvertTypeTransport,
			opts:        []AdvertOption{WithContactDetails(domain.NewContactDetails("refugee_help@wp.pl", "111222333"))},
			expectations: expectations{
				err:            nil,
				contactDetails: domain.NewContactDetails("refugee_help@wp.pl", "111222333"),
			},
		},
		{
			testName:   "Correct data with different contact details",
			user:       usr,
			title:      "Relocation for refugees",
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

	for _, test := range testCases {
		t.Run(test.testName, func(t *testing.T) {
			adv, err := NewAdvert(test.user, test.title, test.description, test.advertType, test.opts...)
			assert.Equal(t, test.expectations.err, err)
			if err == nil {
				assert.Equal(t, test.expectations.contactDetails, adv.Advert.ContactDetails)
			}
		})

	}

}
