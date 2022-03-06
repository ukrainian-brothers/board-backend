package advert_repo

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/ukrainian-brothers/board-backend/domain"
	"github.com/ukrainian-brothers/board-backend/domain/advert"
	"github.com/ukrainian-brothers/board-backend/domain/user"
	"testing"
)

func getUser(t *testing.T) *user.User {
	usr, err := user.NewUser("Adam", "Małysz", "adam@wp.pl", "123", domain.NewContactDetails("adam@wp.pl", "123123123"))
	assert.Nil(t, err)
	return usr
}

func TestAdvertMemoryRepoAdd(t *testing.T) {
	usr := getUser(t)

	repo := NewMemoryAdvertRepository(map[uuid.UUID]*advert.Advert{})

	type testData struct {
		testName    string
		user        *user.User
		title       string
		description string
		advertType  domain.AdvertType
		expectedErr error
	}

	testCases := []testData{
		{
			testName:    "Correct data",
			user:        usr,
			title:       "Title",
			description: "Description",
			advertType:  domain.AdvertTypePlaceToStay,
		},
	}

	for _, test := range testCases {
		t.Run(test.testName, func(t *testing.T) {
			adv, err := advert.NewAdvert(test.user, test.title, test.description, test.advertType)
			assert.NoError(t, err)
			err = repo.Add(adv)
			assert.Equal(t, test.expectedErr, err)

			_, err = repo.Get(adv.ID)
			assert.NoError(t, err)
		})
	}
}

func TestAdvertMemoryRepoGet(t *testing.T) {
	usr := getUser(t)
	adv, err := advert.NewAdvert(usr, "title", "desc", domain.AdvertTypePlaceToStay)
	assert.NoError(t, err)

	repo := NewMemoryAdvertRepository(map[uuid.UUID]*advert.Advert{
		adv.ID: adv,
	})

	type testData struct {
		testName    string
		id          uuid.UUID
		expectedErr error
	}

	testCases := []testData{
		{
			testName: "Correct",
			id: adv.ID,
			expectedErr: nil,
		},
		{
			testName: "Not Found",
			id: uuid.MustParse("f47ac10b-58cc-0372-8567-0e02b2c3d479"),
			expectedErr: AdvertNotFound,
		},
	}

	for _, test := range testCases {
		t.Run(test.testName, func(t *testing.T) {
			_, err := repo.Get(test.id)
			assert.Equal(t, test.expectedErr, err)
		})
	}
}

func TestAdvertMemoryRepoDelete(t *testing.T) {
	usr := getUser(t)
	adv, err := advert.NewAdvert(usr, "title", "desc", domain.AdvertTypePlaceToStay)
	assert.NoError(t, err)

	repo := NewMemoryAdvertRepository(map[uuid.UUID]*advert.Advert{
		adv.ID: adv,
	})

	type testData struct {
		testName    string
		id          uuid.UUID
		expectedErr error
	}

	testCases := []testData{
		{
			testName: "Correct",
			id: adv.ID,
			expectedErr: nil,
		},
		{
			testName: "Not Found",
			id: uuid.MustParse("f47ac10b-58cc-0372-8567-0e02b2c3d479"),
			expectedErr: AdvertNotFound,
		},
	}

	for _, test := range testCases {
		t.Run(test.testName, func(t *testing.T) {
			err := repo.Delete(test.id)
			assert.Equal(t, test.expectedErr, err)
			if err == nil {
				_, err = repo.Get(test.id)
				assert.NotNil(t, err)
			}
			
		})
	}
}
