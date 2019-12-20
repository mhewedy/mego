package attendess

import (
	"errors"
	"github.com/mhewedy/ews"
	"github.com/mhewedy/ews/ewsutil"
)

func getAttendeeDetails(c ews.Client, e string) (*AttendeeDetails, error) {

	if attendeesDB == nil {
		return nil, errors.New("attendees db is empty, build the index first by invoking the search api")
	}

	attendee := attendeesDB[email(e)]
	// check cache
	if attendee.details != nil {
		return attendee.details, nil
	}

	persona, err := ewsutil.GetPersona(c, attendee.PersonaId)
	if err != nil {
		return nil, err
	}

	base64, err := ewsutil.GetUserPhotoBase64(c, e)
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
	attendeesDB[email(e)] = attendee // cache

	return attendee.details, nil
}
