package domain

type AdvertType string

const (
	AdvertTypeTransport   AdvertType = "transport"
	AdvertTypeLawyer                 = "lawyer"
	AdvertTypePlaceToStay            = "place_to_stay"
	AdvertTypeJob                    = "job"
)

type Advert struct {
	Title          string
	Description    string
	Type           AdvertType
	Views          int
	ContactDetails ContactDetails
}
