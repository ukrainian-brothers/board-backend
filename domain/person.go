package domain

import "github.com/google/uuid"



type Person struct {
	ID             uuid.UUID
	FirstName      string
	SureName       string
	Login          string
	Password       *string
}
