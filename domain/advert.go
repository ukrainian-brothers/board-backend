package domain

type AdvertType string

const (
	Transport   AdvertType = "transport"
	Lawyer                 = "lawyer"
	PlaceToStay            = "place_to_stay"
	Job                    = "job"
)

type Advert struct {
	Title          string
	Description    string
	Type           AdvertType
	Views          int
	ContactDetails ContactDetails
}
