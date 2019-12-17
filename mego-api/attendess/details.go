package attendess

import (
	"github.com/mhewedy/ews"
	"github.com/mhewedy/ews/ewsutil"
)

// TODO do caching
func getAttendeeDetails(c ews.Client, email string) (*Attendee, error) {

	attendee := attendeesIndex[email]

	persona, err := ewsutil.GetPersona(c, attendee.PersonaId)

	if err != nil {
		return nil, err
	}

	base64, err := ewsutil.GetUserPhotoBase64(c, email)
	if err != nil {
		base64 = "NA"
	}

	if attendeesIndex != nil {
		attendee.Image = base64
		attendeesIndex[email] = attendee
	}

	return &Attendee{
		PersonaId:           attendee.PersonaId,
		DisplayName:         attendee.DisplayName,
		Title:               attendee.Title,
		EmailAddress:        attendee.EmailAddress,
		Image:               base64,
		Department:          persona.Department,
		BusinessPhoneNumber: persona.BusinessPhoneNumbers.PhoneNumberAttributedValue.Value.Number,
		MobilePhone:         persona.MobilePhones.PhoneNumberAttributedValue.Value.Number,
		OfficeLocation:      persona.OfficeLocations.StringAttributedValue.Value,
	}, nil
}
