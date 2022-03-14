package domain

type ContactDetails struct {
	Mail        *string
	PhoneNumber *string
}

// TODO: Create ContactDetails aggregate which will have logic for checking if the mail and phone numbers are correct according to domain policy
func NewContactDetails(mail string, phoneNumber string) ContactDetails {
	details := ContactDetails{
		Mail:        &mail,
		PhoneNumber: &phoneNumber,
	}
	if phoneNumber == "" {
		details.PhoneNumber = nil
	}

	if mail == "" {
		details.Mail = nil
	}

	// TODO: Validate length of mail and phone number
	return details
}

func (cd ContactDetails) IsEmpty() bool {
	return cd.Mail == nil && cd.PhoneNumber == nil
}
