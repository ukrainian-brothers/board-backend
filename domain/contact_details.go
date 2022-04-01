package domain

import (
	"errors"
	"regexp"
)

type ContactDetails struct {
	Mail        *string
	PhoneNumber *string
}

var (
	InvalidDataErr    = errors.New("invalid contact data")
	emailAddressRegex = regexp.MustCompile("(?:[a-z0-9!#$%&'*+/=?^_ \\x60{|}~-]+(?:\\.[a-z0-9!#$%&'*+/=?^_ \\x60{|}~-]+)*|\"(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x21\\x23-\\x5b\\x5d-\\x7f]|\\\\[\\x01-\\x09\\x0b\\x0c\\x0e-\\x7f])*\")@(?:(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?|\\[(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?|[a-z0-9-]*[a-z0-9]:(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x21-\\x5a\\x53-\\x7f]|\\\\[\\x01-\\x09\\x0b\\x0c\\x0e-\\x7f])+)\\])") // RFC 5322
	phoneRegex        = regexp.MustCompile("^\\+[0-9]{2} [0-9]{3} [0-9]{3} [0-9]{3}$")
)

func NewContactDetails(mail, phoneNumber string) (ContactDetails, error) {
	details := ContactDetails{
		Mail:        &mail,
		PhoneNumber: &phoneNumber,
	}
	isMailValid := emailAddressRegex.MatchString(mail)
	isPhoneValid := phoneRegex.MatchString(phoneNumber)

	if mail == "" {
		details.Mail = nil
	}

	if phoneNumber == "" {
		details.PhoneNumber = nil
	}

	if !isMailValid && !isPhoneValid {
		return ContactDetails{}, InvalidDataErr
	}

	return details, nil
}

func (cd ContactDetails) IsEmpty() bool {
	return cd.Mail == nil && cd.PhoneNumber == nil
}
