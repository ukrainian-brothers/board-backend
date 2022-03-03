package domain

import "github.com/google/uuid"



type User struct {
	ID             uuid.UUID
	FirstName      string
	SureName       string
	Login          string
	Password       *string
	ContactDetails ContactDetails
}
