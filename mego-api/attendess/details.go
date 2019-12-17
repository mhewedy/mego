package attendess

import (
	"errors"
	"github.com/mhewedy/ews"
	"github.com/mhewedy/ews/ewsutil"
)

func getAttendeeDetails(c ews.Client, email string) (*AttendeeDetails, error) {

	if attendeesIndex == nil {
		return nil, errors.New("index is empty, should run the index first")
	}

	attendee := attendeesIndex[email]
	// check cache
	if attendee.details != nil {
		return attendee.details, nil
	}

	persona, err := ewsutil.GetPersona(c, attendee.PersonaId)
	if err != nil {
		return nil, err
	}

	base64, err := ewsutil.GetUserPhotoBase64(c, email)
	if err != nil {
		base64 = "NA"
	}

	attendee.Image = base64
	attendee.details = &AttendeeDetails{
		Attendee:            attendee,
		Department:          persona.Department,
		BusinessPhoneNumber: persona.BusinessPhoneNumbers.PhoneNumberAttributedValue.Value.Number,
		MobilePhone:         persona.MobilePhones.PhoneNumberAttributedValue.Value.Number,
		OfficeLocation:      persona.OfficeLocations.StringAttributedValue.Value,
	}
	attendeesIndex[email] = attendee // cache

	return attendee.details, nil
}
