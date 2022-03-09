package domain

import "github.com/google/uuid"



type Person struct {
	ID        uuid.UUID
	FirstName string
	Surname   string
	Login     string
	Password  *string
}
