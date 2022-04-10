package domain

import . "github.com/ukrainian-brothers/board-backend/pkg/translation"

type AdvertType string

const (
	AdvertTypeTransport   AdvertType = "transport"
	AdvertTypeLawyer      AdvertType = "lawyer"
	AdvertTypePlaceToStay AdvertType = "place_to_stay"
	AdvertTypeJob         AdvertType = "job"
)

type AdvertDetails struct {
	Title          MultilingualString
	Description    MultilingualString
	Type           AdvertType
	Views          int
	ContactDetails ContactDetails
}
