package domain

type ContactDetails struct {
	Mail        *string
	PhoneNumber *string
}

// TODO: Create ContactDetails aggregate which will have logic for checking if the mail and phone numbers are correct according to domain policy
func NewContactDetails(mail string, phoneNumber string) ContactDetails {
	return ContactDetails{
		Mail:        &mail,
		PhoneNumber: &phoneNumber,
	}
}

func (cd ContactDetails) IsEmpty() bool {
	return cd.Mail == nil && cd.PhoneNumber == nil
}
