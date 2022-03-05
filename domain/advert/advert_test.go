package advert

import (
	"github.com/stretchr/testify/assert"
	"github.com/ukrainian-brothers/board-backend/domain"
	"github.com/ukrainian-brothers/board-backend/domain/user"
	"testing"
)

func TestAdvertNew(t *testing.T) {
	type testData struct {
		testName    string
		user        *user.User
		title       string
		description string
		advertType  domain.AdvertType
		opts        []AdvertOption
		expectedErr error
		expectedContactDetails domain.ContactDetails
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
			description: "",
			advertType:  domain.AdvertTypeTransport,
			expectedErr: nil,
			// The constructor got no ContactDetails - it means that it will use user contact details by default
			expectedContactDetails: domain.NewContactDetails("adam@wp.pl", "111222333"),
		},
		{
			testName:    "Correct data with different contact details",
			user:        usr,
			title:       "Relocation for refugees",
			description: "",
			advertType:  domain.AdvertTypeTransport,
			opts:        []AdvertOption{WithContactDetails(domain.NewContactDetails("refugee_help@wp.pl", "111222333"))},
			expectedErr: nil,
			expectedContactDetails: domain.NewContactDetails("refugee_help@wp.pl", "111222333"),
		},
		{
			testName:    "Correct data with different contact details",
			user:        usr,
			title:       "Relocation for refugees",
			description: "",
			advertType:  domain.AdvertTypeTransport,
			opts:        []AdvertOption{WithContactDetails(domain.ContactDetails{})},
			expectedErr: ContactEmptyErr,
		},
	}

	for _, test := range testCases {
		t.Run(test.testName, func(t *testing.T) {
			adv, err := NewAdvert(test.user, test.title, test.description, test.advertType, test.opts...)
			assert.Equal(t, test.expectedErr, err)
			if err == nil {
				assert.Equal(t, test.expectedContactDetails, adv.Advert.ContactDetails)
			}
		})

	}

}
