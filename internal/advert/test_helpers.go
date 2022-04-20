package advert

import (
	"github.com/google/uuid"
	"github.com/ukrainian-brothers/board-backend/domain"
	"github.com/ukrainian-brothers/board-backend/pkg/test_helpers"
	"github.com/ukrainian-brothers/board-backend/pkg/translation"
	"time"
)

func GenerateTestAdvertDB(advertID uuid.UUID, userID uuid.UUID) AdvertDB {
	return AdvertDB{
		ID:     advertID,
		UserID: userID,
		Type:   domain.AdvertTypeTransport,
		Views:  test_helpers.RandomNumberRange(0, 150),
		ContactDetails: domain.ContactDetails{
			Mail: test_helpers.RandomMail(),
		},
		CreatedAt: time.Now(),
	}
}
func GenerateTestAdvertDetailsDB(advertID uuid.UUID, language translation.LanguageTag) AdvertDetailsDB {
	return AdvertDetailsDB{
		ID:          uuid.New(),
		AdvertID:    advertID,
		Language:    language,
		Title:       test_helpers.RandomString(10),
		Description: test_helpers.RandomString(50),
	}
}
